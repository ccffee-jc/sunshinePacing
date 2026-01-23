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
	StreamOther
)

// Session 保存单客户端地址与超时。
type Session struct {
	mu       sync.Mutex
	addrs    map[int]*net.UDPAddr
	lastSeen time.Time
	timeout  time.Duration
}

// NewSession 创建会话管理器。
func NewSession(timeout time.Duration) *Session {
	return &Session{
		addrs:   make(map[int]*net.UDPAddr),
		timeout: timeout,
	}
}

// SetClient 更新客户端地址，返回是否为新客户端。
func (s *Session) SetClient(port int, addr *net.UDPAddr) bool {
	if addr == nil {
		return false
	}
	clone := cloneAddr(addr)
	s.mu.Lock()
	changed := !sameAddr(s.addrs[port], addr)
	s.addrs[port] = clone
	s.lastSeen = time.Now()
	s.mu.Unlock()
	return changed
}

// GetClient 获取客户端地址，若超时则清空。
func (s *Session) GetClient(port int) *net.UDPAddr {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.lastSeen.IsZero() && time.Since(s.lastSeen) > s.timeout {
		s.addrs = make(map[int]*net.UDPAddr)
		s.lastSeen = time.Time{}
		return nil
	}
	addr := s.addrs[port]
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

func sameAddr(a *net.UDPAddr, b *net.UDPAddr) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Port != b.Port || a.Zone != b.Zone {
		return false
	}
	if a.IP == nil || b.IP == nil {
		return a.IP == nil && b.IP == nil
	}
	return a.IP.Equal(b.IP)
}
