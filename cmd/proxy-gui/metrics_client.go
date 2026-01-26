// GUI 端用于拉取 metrics 的辅助方法。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type metricsSnapshot struct {
	VideoInBytes    uint64 `json:"video_in_bytes"`
	VideoOutBytes   uint64 `json:"video_out_bytes"`
	VideoDrops      uint64 `json:"video_drops"`
	ControlInBytes  uint64 `json:"control_in_bytes"`
	ControlOutBytes uint64 `json:"control_out_bytes"`
	AudioInBytes    uint64 `json:"audio_in_bytes"`
	AudioOutBytes   uint64 `json:"audio_out_bytes"`
	VideoQueueLen   int    `json:"video_queue_len"`
}

func (m metricsSnapshot) Text() string {
	return fmt.Sprintf(
		"video_out=%dB video_drop=%d queue=%d control_out=%dB audio_out=%dB",
		m.VideoOutBytes,
		m.VideoDrops,
		m.VideoQueueLen,
		m.ControlOutBytes,
		m.AudioOutBytes,
	)
}

type metricsClient struct {
	client *http.Client
}

func newMetricsClient() *metricsClient {
	return &metricsClient{
		client: &http.Client{
			Timeout: 800 * time.Millisecond,
		},
	}
}

func (c *metricsClient) Fetch(ctx context.Context, addr string) (metricsSnapshot, error) {
	if strings.TrimSpace(addr) == "" {
		return metricsSnapshot{}, errors.New("metrics 地址为空")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/metrics", nil)
	if err != nil {
		return metricsSnapshot{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return metricsSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return metricsSnapshot{}, fmt.Errorf("metrics 状态码异常: %d", resp.StatusCode)
	}
	var snapshot metricsSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
		return metricsSnapshot{}, err
	}
	return snapshot, nil
}

func waitMetricsAddr(path string, timeout time.Duration) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("metrics 文件路径为空")
	}
	deadline := time.Now().Add(timeout)
	for {
		data, err := os.ReadFile(path)
		if err == nil {
			addr := strings.TrimSpace(string(data))
			if addr != "" {
				return addr, nil
			}
		}
		if time.Now().After(deadline) {
			return "", errors.New("等待 metrics 地址超时")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
