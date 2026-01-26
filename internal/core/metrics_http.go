// Metrics HTTP 服务用于向 GUI 提供本地运行指标快照。
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// MetricsServer 负责提供本地 metrics HTTP 服务。
type MetricsServer struct {
	addr   string
	server *http.Server
	ln     net.Listener
}

// StartMetricsServer 启动 metrics HTTP 服务。
func StartMetricsServer(ctx context.Context, proxy *Proxy, addr string) (*MetricsServer, error) {
	if proxy == nil {
		return nil, errors.New("proxy 不能为空")
	}
	normalized, err := normalizeMetricsAddr(addr)
	if err != nil {
		return nil, err
	}
	ln, err := net.Listen("tcp", normalized)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(proxy.Metrics())
	})
	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	metrics := &MetricsServer{
		addr:   ln.Addr().String(),
		server: server,
		ln:     ln,
	}
	go func() {
		_ = server.Serve(ln)
	}()
	if ctx != nil {
		go func() {
			<-ctx.Done()
			_ = metrics.Stop(context.Background())
		}()
	}
	return metrics, nil
}

// Addr 返回实际监听地址。
func (m *MetricsServer) Addr() string {
	if m == nil {
		return ""
	}
	return m.addr
}

// Stop 停止 metrics HTTP 服务。
func (m *MetricsServer) Stop(ctx context.Context) error {
	if m == nil || m.server == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return m.server.Shutdown(ctx)
}

func normalizeMetricsAddr(addr string) (string, error) {
	if strings.TrimSpace(addr) == "" {
		return "", errors.New("metrics 地址不能为空")
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	host = strings.TrimSpace(host)
	if host == "" {
		host = "127.0.0.1"
	}
	if !isLoopbackHost(host) {
		return "", fmt.Errorf("metrics 地址必须为本地回环地址: %s", host)
	}
	return net.JoinHostPort(host, port), nil
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
