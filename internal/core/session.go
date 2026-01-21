// 会话管理用于保存单客户端的 UDP 端点映射。
package core

import (
	"net"
	"sync"
	"time"
)

// StreamType 表示端口用途。
type StreamType int

const (
	StreamVideo StreamType = iota
	StreamControl
	StreamAudio
)

// Session 保存单客户端地址与超时。
type Session struct {
	mu       sync.Mutex
	addrs    map[StreamType]*net.UDPAddr
	lastSeen time.Time
	timeout  time.Duration
}

// NewSession 创建会话管理器。
func NewSession(timeout time.Duration) *Session {
	return &Session{
		addrs:   make(map[StreamType]*net.UDPAddr),
		timeout: timeout,
	}
}

// SetClient 更新客户端地址。
func (s *Session) SetClient(stream StreamType, addr *net.UDPAddr) {
	if addr == nil {
		return
	}
	clone := cloneAddr(addr)
	s.mu.Lock()
	s.addrs[stream] = clone
	s.lastSeen = time.Now()
	s.mu.Unlock()
}

// GetClient 获取客户端地址，若超时则清空。
func (s *Session) GetClient(stream StreamType) *net.UDPAddr {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lastSeen.IsZero() && time.Since(s.lastSeen) > s.timeout {
		s.addrs = make(map[StreamType]*net.UDPAddr)
		s.lastSeen = time.Time{}
		return nil
	}
	addr := s.addrs[stream]
	if addr == nil {
		return nil
	}
	return cloneAddr(addr)
}

func cloneAddr(addr *net.UDPAddr) *net.UDPAddr {
	if addr == nil {
		return nil
	}
	clone := *addr
	if addr.IP != nil {
		clone.IP = append([]byte{}, addr.IP...)
	}
	return &clone
}
