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

	<-ctx.Done()
	proxy.Stop()
}

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
