// 运行状态管理，避免 GUI 并发访问冲突。
package main

import (
	"context"
	"os/exec"
	"sync"
)

type proxyRuntime struct {
	mu          sync.Mutex
	cmd         *exec.Cmd
	metricsAddr string
	metricsFile string
	pollCancel  context.CancelFunc
}

func (r *proxyRuntime) Store(cmd *exec.Cmd, metricsFile string, pollCancel context.CancelFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cmd = cmd
	r.metricsAddr = ""
	r.metricsFile = metricsFile
	r.pollCancel = pollCancel
}

func (r *proxyRuntime) SetMetricsAddr(addr string) {
	r.mu.Lock()
	r.metricsAddr = addr
	r.mu.Unlock()
}

func (r *proxyRuntime) Snapshot() (*exec.Cmd, string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cmd, r.metricsAddr
}

func (r *proxyRuntime) Take() (*exec.Cmd, string, context.CancelFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cmd := r.cmd
	metricsFile := r.metricsFile
	pollCancel := r.pollCancel
	r.cmd = nil
	r.metricsAddr = ""
	r.metricsFile = ""
	r.pollCancel = nil
	return cmd, metricsFile, pollCancel
}

func (r *proxyRuntime) ClearIf(cmd *exec.Cmd) (string, context.CancelFunc, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cmd != cmd {
		return "", nil, false
	}
	metricsFile := r.metricsFile
	pollCancel := r.pollCancel
	r.cmd = nil
	r.metricsAddr = ""
	r.metricsFile = ""
	r.pollCancel = nil
	return metricsFile, pollCancel, true
}
