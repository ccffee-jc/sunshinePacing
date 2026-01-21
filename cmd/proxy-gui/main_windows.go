//go:build windows

// Windows GUI 入口使用 Fyne 启动代理并显示状态。
package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
	"sunshinePacing/internal/core"
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

	var runningProxy *core.Proxy
	var runningCancel context.CancelFunc

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
		cfg, err := configFromForm(baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
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
		if runningProxy != nil {
			dialog.ShowInformation("提示", "代理已在运行", window)
			return
		}
		cfg, err := configFromForm(baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		proxy, err := core.NewProxy(cfg)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		if err := proxy.Start(ctx); err != nil {
			cancel()
			dialog.ShowError(err, window)
			return
		}
		runningProxy = proxy
		runningCancel = cancel
		_ = statusBind.Set("运行中")
	})

	stopBtn := widget.NewButton("停止", func() {
		if runningProxy == nil {
			dialog.ShowInformation("提示", "代理未运行", window)
			return
		}
		if runningCancel != nil {
			runningCancel()
		}
		runningProxy.Stop()
		runningProxy = nil
		runningCancel = nil
		_ = statusBind.Set("已停止")
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
	)

	window.SetContent(form)
	window.Resize(fyne.NewSize(520, 520))

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if runningProxy == nil {
				continue
			}
			m := runningProxy.Metrics()
			_ = metricsBind.Set(fmt.Sprintf(
				"video_out=%dB video_drop=%d queue=%d control_out=%dB audio_out=%dB",
				m.VideoOutBytes,
				m.VideoDrops,
				m.VideoQueueLen,
				m.ControlOutBytes,
				m.AudioOutBytes,
			))
		}
	}()

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

func fillForm(cfg config.Config, baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry *widget.Entry) {
	baseEntry.SetText(strconv.Itoa(cfg.BasePort))
	hostEntry.SetText(cfg.InternalHost)
	rateEntry.SetText(strconv.Itoa(cfg.Video.RateMbps))
	burstEntry.SetText(strconv.Itoa(cfg.Video.BurstKB))
	queueEntry.SetText(strconv.Itoa(cfg.Video.MaxQueueDelayMs))
	tickEntry.SetText(strconv.Itoa(cfg.Video.TickMs))
}

func configFromForm(baseEntry, hostEntry, rateEntry, burstEntry, queueEntry, tickEntry *widget.Entry) (config.Config, error) {
	cfg := config.DefaultConfig()
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
