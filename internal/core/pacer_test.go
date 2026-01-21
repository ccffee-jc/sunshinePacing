// 视频 pacing 测试覆盖丢弃与发送逻辑。
package core

import (
	"net"
	"testing"
	"time"
)

func TestPacerFlushDropsExpired(t *testing.T) {
	stats := &Stats{}
	p := NewPacer(10, 16, 5*time.Millisecond, time.Millisecond, func([]byte, *net.UDPAddr) error {
		return nil
	}, stats)
	p.queue = []queuedPacket{{
		data:      []byte{1, 2, 3},
		addr:      &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234},
		enqueued:  time.Now().Add(-10 * time.Millisecond),
		byteCount: 3,
	}}
	p.head = 0
	tokens := 1_000_000.0
	p.flush(time.Now(), &tokens)
	if stats.Snapshot().VideoDrops != 1 {
		t.Fatalf("期望丢弃 1 个包，实际为 %d", stats.Snapshot().VideoDrops)
	}
}

func TestPacerFlushSends(t *testing.T) {
	stats := &Stats{}
	sent := 0
	p := NewPacer(10, 16, 10*time.Millisecond, time.Millisecond, func([]byte, *net.UDPAddr) error {
		sent++
		return nil
	}, stats)
	p.queue = []queuedPacket{{
		data:      []byte{1, 2, 3, 4},
		addr:      &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234},
		enqueued:  time.Now(),
		byteCount: 4,
	}}
	p.head = 0
	tokens := 1_000_000.0
	p.flush(time.Now(), &tokens)
	if sent != 1 {
		t.Fatalf("期望发送 1 个包，实际为 %d", sent)
	}
	if stats.Snapshot().VideoOutBytes != 4 {
		t.Fatalf("期望发送字节 4，实际为 %d", stats.Snapshot().VideoOutBytes)
	}
}
