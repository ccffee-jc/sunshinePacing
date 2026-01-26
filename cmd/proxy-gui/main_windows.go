//go:build windows

// Windows GUI 入口使用 Fyne 启动代理并显示状态。
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v3"

	"sunshinePacing/internal/config"
)

func main() {
	setupLogger()

	fyneApp := app.New()
	window := fyneApp.NewWindow("Sunshine Pacing Proxy")

	cfgPathEntry := widget.NewEntry()
	cfgPathEntry.SetPlaceHolder("配置文件路径，如 C:\\proxy.yml")

	baseEntry := widget.NewEntry()
	hostEntry := widget.NewEntry()
	rateEntry := widget.NewEntry()
	burstEntry := widget.NewEntry()
	queueEntry := widget.NewEntry()
	tickEntry := widget.NewEntry()

	setFormDefaults(baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)

	statusBind := binding.NewString()
	metricsBind := binding.NewString()
	_ = statusBind.Set("未启动")
	_ = metricsBind.Set("暂无数据")

	statusLabel := widget.NewLabelWithData(statusBind)
	metricsLabel := widget.NewLabelWithData(metricsBind)
	burstChart := NewBurstChart(120)
	burstLegend := newBurstLegend()

	if defaultPath, err := defaultConfigPath(); err != nil {
		dialog.ShowError(err, window)
	} else {
		cfgPathEntry.SetText(defaultPath)
		loadOrCreateDefaultConfig(defaultPath, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry, statusBind, window)
	}

	runtimeState := &proxyRuntime{}
	metricsClient := newMetricsClient()

	loadBtn := widget.NewButton("加载配置", func() {
		path := cfgPathEntry.Text
		if path == "" {
			dialog.ShowError(fmt.Errorf("配置路径为空"), window)
			return
		}
		cfg, err := config.Load(path)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		fillForm(cfg, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
	})

	saveBtn := widget.NewButton("保存配置", func() {
		path := cfgPathEntry.Text
		if path == "" {
			dialog.ShowError(fmt.Errorf("配置路径为空"), window)
			return
		}
		cfg, err := configFromFormWithBase(path, true, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		data, err := yaml.Marshal(cfg)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			dialog.ShowError(err, window)
			return
		}
		_ = statusBind.Set("配置已保存")
	})

	startBtn := widget.NewButton("启动", func() {
		if cmd, _ := runtimeState.Snapshot(); cmd != nil {
			dialog.ShowInformation("提示", "代理已在运行", window)
			return
		}
		cfgPath := cfgPathEntry.Text
		cfg, err := configFromFormWithBase(cfgPath, false, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		data, err := yaml.Marshal(cfg)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
			dialog.ShowError(err, window)
			return
		}
		execPath, err := os.Executable()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		metricsFile, err := os.CreateTemp(filepath.Dir(cfgPath), "sunshine-metrics-*.txt")
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		metricsPath := metricsFile.Name()
		_ = metricsFile.Close()
		cliPath := filepath.Join(filepath.Dir(execPath), "sunshine-proxy-cli.exe")
		cmd := exec.Command(cliPath, "-config", cfgPath, "-metrics-addr", "127.0.0.1:0", "-metrics-file", metricsPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			_ = os.Remove(metricsPath)
			dialog.ShowError(err, window)
			return
		}
		pollCtx, pollCancel := context.WithCancel(context.Background())
		runtimeState.Store(cmd, metricsPath, pollCancel)
		_ = statusBind.Set("启动中")
		_ = metricsBind.Set("正在连接指标")

		go func(cmd *exec.Cmd, metricsPath string) {
			_ = cmd.Wait()
			if metricsFile, pollCancel, cleared := runtimeState.ClearIf(cmd); cleared {
				if pollCancel != nil {
					pollCancel()
				}
				if metricsFile != "" {
					_ = os.Remove(metricsFile)
				}
				fyne.Do(func() {
					_ = statusBind.Set("已停止")
					_ = metricsBind.Set("暂无数据")
				})
			}
		}(cmd, metricsPath)

		go func(cmd *exec.Cmd, metricsPath string) {
			addr, err := waitMetricsAddr(metricsPath, 3*time.Second)
			if err != nil {
				if metricsFile, pollCancel, cleared := runtimeState.ClearIf(cmd); cleared {
					if pollCancel != nil {
						pollCancel()
					}
					if cmd.Process != nil {
						_ = cmd.Process.Kill()
					}
					if metricsFile != "" {
						_ = os.Remove(metricsFile)
					}
					fyne.Do(func() {
						_ = statusBind.Set("启动失败")
						_ = metricsBind.Set("暂无数据")
						dialog.ShowError(err, window)
					})
				}
				return
			}
			runtimeState.SetMetricsAddr(addr)
			fyne.Do(func() {
				_ = statusBind.Set("运行中")
			})
		}(cmd, metricsPath)

		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-pollCtx.Done():
					return
				case <-ticker.C:
					cmd, addr := runtimeState.Snapshot()
					if cmd == nil || addr == "" {
						continue
					}
					snapshot, err := metricsClient.Fetch(pollCtx, addr)
					if err != nil {
						fyne.Do(func() {
							_ = metricsBind.Set("指标连接失败")
						})
						continue
					}
					text := snapshot.Text()
					queueLen := snapshot.VideoQueueLen
					fyne.Do(func() {
						burstChart.Push(queueLen)
						_ = metricsBind.Set(text)
					})
				}
			}
		}()
	})

	stopBtn := widget.NewButton("停止", func() {
		cmd, metricsPath, pollCancel := runtimeState.Take()
		if cmd == nil {
			dialog.ShowInformation("提示", "代理未运行", window)
			return
		}
		if pollCancel != nil {
			pollCancel()
		}
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		if metricsPath != "" {
			_ = os.Remove(metricsPath)
		}
		_ = statusBind.Set("已停止")
		_ = metricsBind.Set("暂无数据")
	})

	form := container.NewVBox(
		widget.NewLabel("配置文件"),
		cfgPathEntry,
		container.NewHBox(loadBtn, saveBtn),
		widget.NewSeparator(),
		widget.NewLabel("运行参数"),
		container.NewGridWithColumns(2,
			widget.NewLabel("sunshine_base_port"), baseEntry,
			widget.NewLabel("internal_host"), hostEntry,
			widget.NewLabel("video.rate_mbps"), rateEntry,
			widget.NewLabel("video.burst_kb"), burstEntry,
			widget.NewLabel("video.max_queue_delay_ms"), queueEntry,
			widget.NewLabel("video.tick_ms"), tickEntry,
		),
		container.NewHBox(startBtn, stopBtn),
		widget.NewSeparator(),
		widget.NewLabel("状态"),
		statusLabel,
		widget.NewLabel("统计"),
		metricsLabel,
		widget.NewSeparator(),
		widget.NewLabel("实时突发图表"),
		burstLegend,
		burstChart,
	)

	window.SetContent(form)
	window.Resize(fyne.NewSize(520, 520))

	window.ShowAndRun()
}

func setFormDefaults(baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry *widget.Entry) {
	cfg := config.DefaultConfig()
	baseEntry.SetText(strconv.Itoa(cfg.BasePort))
	hostEntry.SetText(cfg.InternalHost)
	rateEntry.SetText(strconv.Itoa(cfg.Video.RateMbps))
	burstEntry.SetText(strconv.Itoa(cfg.Video.BurstKB))
	queueEntry.SetText(strconv.Itoa(cfg.Video.MaxQueueDelayMs))
	tickEntry.SetText(strconv.Itoa(cfg.Video.TickMs))
}

func defaultConfigPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	return filepath.Join(filepath.Dir(execPath), "sunshine-proxy.yml"), nil
}

func loadOrCreateDefaultConfig(path string, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry *widget.Entry, statusBind binding.String, window fyne.Window) {
	if path == "" {
		return
	}
	_, err := os.Stat(path)
	if err == nil {
		cfg, err := config.Load(path)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		fillForm(cfg, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
		_ = statusBind.Set("已加载默认配置")
		return
	}
	if !errors.Is(err, os.ErrNotExist) {
		dialog.ShowError(fmt.Errorf("检查配置失败: %w", err), window)
		return
	}
	cfg := config.DefaultConfig()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		dialog.ShowError(err, window)
		return
	}
	fillForm(cfg, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
	_ = statusBind.Set("已生成默认配置")
}

func fillForm(cfg config.Config, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry *widget.Entry) {
	baseEntry.SetText(strconv.Itoa(cfg.BasePort))
	hostEntry.SetText(cfg.InternalHost)
	rateEntry.SetText(strconv.Itoa(cfg.Video.RateMbps))
	burstEntry.SetText(strconv.Itoa(cfg.Video.BurstKB))
	queueEntry.SetText(strconv.Itoa(cfg.Video.MaxQueueDelayMs))
	tickEntry.SetText(strconv.Itoa(cfg.Video.TickMs))
}

func configFromFormWithBase(path string, allowMissing bool, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry *widget.Entry) (config.Config, error) {
	if path == "" {
		return config.Config{}, fmt.Errorf("配置路径为空")
	}
	cfg := config.DefaultConfig()
	loaded, err := config.Load(path)
	if err != nil {
		if !(allowMissing && errors.Is(err, os.ErrNotExist)) {
			return cfg, err
		}
	} else {
		cfg = loaded
	}
	base, err := strconv.Atoi(baseEntry.Text)
	if err != nil {
		return cfg, fmt.Errorf("base_port 无效")
	}
	host := hostEntry.Text
	if host == "" {
		return cfg, fmt.Errorf("internal_host 不能为空")
	}
	rate, err := strconv.Atoi(rateEntry.Text)
	if err != nil {
		return cfg, fmt.Errorf("rate_mbps 无效")
	}
	burst, err := strconv.Atoi(burstEntry.Text)
	if err != nil {
		return cfg, fmt.Errorf("burst_kb 无效")
	}
	queueMs, err := strconv.Atoi(queueEntry.Text)
	if err != nil {
		return cfg, fmt.Errorf("max_queue_delay_ms 无效")
	}
	tickMs, err := strconv.Atoi(tickEntry.Text)
	if err != nil {
		return cfg, fmt.Errorf("tick_ms 无效")
	}

	cfg.BasePort = base
	cfg.InternalHost = host
	cfg.Video.RateMbps = rate
	cfg.Video.BurstKB = burst
	cfg.Video.MaxQueueDelayMs = queueMs
	cfg.Video.TickMs = tickMs

	if err := cfg.NormalizeAndValidate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
