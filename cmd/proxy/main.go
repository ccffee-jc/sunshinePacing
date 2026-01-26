// CLI 入口用于在 Linux 或无 GUI 环境启动代理。
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sunshinePacing/internal/config"
	"sunshinePacing/internal/core"
)

func main() {
	cfgPath := flag.String("config", "proxy.yml", "配置文件路径")
	metricsAddr := flag.String("metrics-addr", "", "metrics 监听地址(如 127.0.0.1:0)")
	metricsFile := flag.String("metrics-file", "", "metrics 监听地址输出文件")
	flag.Parse()

	setupLogger()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Printf("加载配置失败: %v", err)
		os.Exit(1)
	}

	proxy, err := core.NewProxy(cfg)
	if err != nil {
		log.Printf("初始化代理失败: %v", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := proxy.Start(ctx); err != nil {
		log.Printf("启动代理失败: %v", err)
		os.Exit(1)
	}

	var metricsServer *core.MetricsServer
	if *metricsFile != "" && *metricsAddr == "" {
		log.Printf("metrics-file 需要配合 metrics-addr 使用")
		os.Exit(1)
	}
	if *metricsAddr != "" {
		server, err := core.StartMetricsServer(ctx, proxy, *metricsAddr)
		if err != nil {
			log.Printf("启动 metrics 服务失败: %v", err)
			os.Exit(1)
		}
		metricsServer = server
		if *metricsFile != "" {
			if err := os.WriteFile(*metricsFile, []byte(server.Addr()), 0o644); err != nil {
				_ = server.Stop(context.Background())
				log.Printf("写入 metrics 地址失败: %v", err)
				os.Exit(1)
			}
		}
	}

	<-ctx.Done()
	if metricsServer != nil {
		_ = metricsServer.Stop(context.Background())
	}
	proxy.Stop()
}

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
