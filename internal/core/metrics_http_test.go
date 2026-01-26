// Metrics HTTP 服务的基础测试。
package core

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"sunshinePacing/internal/config"
)

func TestMetricsServer(t *testing.T) {
	proxy, err := NewProxy(config.DefaultConfig())
	if err != nil {
		t.Fatalf("创建代理失败: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server, err := StartMetricsServer(ctx, proxy, "127.0.0.1:0")
	if err != nil {
		t.Fatalf("启动 metrics 服务失败: %v", err)
	}
	defer func() {
		_ = server.Stop(context.Background())
	}()

	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get("http://" + server.Addr() + "/metrics")
	if err != nil {
		t.Fatalf("请求 metrics 失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("metrics 状态码异常: %d", resp.StatusCode)
	}
	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("解析 metrics 失败: %v", err)
	}
	if _, ok := payload["video_out_bytes"]; !ok {
		t.Fatalf("缺少字段: video_out_bytes")
	}
	if _, ok := payload["video_queue_len"]; !ok {
		t.Fatalf("缺少字段: video_queue_len")
	}
}
