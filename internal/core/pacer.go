// 视频 pacing 使用令牌桶与队列驻留上限控制突发。
package core

import (
	"context"
	"net"
	"sync"
	"time"
)

type queuedPacket struct {
	data      []byte
	addr      *net.UDPAddr
	enqueued  time.Time
	byteCount int
}

// Pacer 控制视频发送速率。
type Pacer struct {
	rateBps       float64
	burstBits     float64
	maxQueueDelay time.Duration
	tickInterval  time.Duration
	send          func([]byte, *net.UDPAddr) error
	stats         *Stats

	mu    sync.Mutex
	queue []queuedPacket
	head  int
}

// NewPacer 创建视频 pacing 实例。
func NewPacer(rateMbps int, burstKB int, maxQueueDelay time.Duration, tick time.Duration, send func([]byte, *net.UDPAddr) error, stats *Stats) *Pacer {
	if tick <= 0 {
		tick = time.Millisecond
	}
	return &Pacer{
		rateBps:       float64(rateMbps) * 1_000_000,
		burstBits:     float64(burstKB) * 1024 * 8,
		maxQueueDelay: maxQueueDelay,
		tickInterval:  tick,
		send:          send,
		stats:         stats,
		queue:         make([]queuedPacket, 0, 256),
	}
}

// Enqueue 将视频包加入队列。
func (p *Pacer) Enqueue(data []byte, addr *net.UDPAddr) {
	if p == nil || addr == nil {
		if p != nil && p.stats != nil {
			p.stats.AddVideoDrop()
		}
		return
	}
	copyBuf := make([]byte, len(data))
	copy(copyBuf, data)
	pkt := queuedPacket{
		data:      copyBuf,
		addr:      cloneAddr(addr),
		enqueued:  time.Now(),
		byteCount: len(copyBuf),
	}
	p.mu.Lock()
	p.queue = append(p.queue, pkt)
	p.mu.Unlock()
}

// QueueLen 返回当前队列长度。
func (p *Pacer) QueueLen() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.queue) - p.head
}

// Run 启动 pacing 循环。
func (p *Pacer) Run(ctx context.Context) {
	ticker := time.NewTicker(p.tickInterval)
	defer ticker.Stop()
	last := time.Now()
	var tokens float64

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			dt := now.Sub(last).Seconds()
			last = now
			tokens += p.rateBps * dt
			if tokens > p.burstBits {
				tokens = p.burstBits
			}
			p.flush(now, &tokens)
		}
	}
}

func (p *Pacer) flush(now time.Time, tokens *float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.head >= len(p.queue) {
		p.queue = p.queue[:0]
		p.head = 0
		return
	}
	maxDelay := p.maxQueueDelay
	for p.head < len(p.queue) {
		pkt := p.queue[p.head]
		if maxDelay > 0 && now.Sub(pkt.enqueued) > maxDelay {
			p.head++
			if p.stats != nil {
				p.stats.AddVideoDrop()
			}
			continue
		}
		cost := float64(pkt.byteCount * 8)
		if *tokens < cost {
			break
		}
		if p.send != nil && pkt.addr != nil {
			if err := p.send(pkt.data, pkt.addr); err == nil {
				if p.stats != nil {
					p.stats.AddVideoOut(pkt.byteCount)
				}
			} else if p.stats != nil {
				p.stats.AddVideoDrop()
			}
		}
		*tokens -= cost
		p.head++
	}
	if p.head > 0 && p.head >= len(p.queue) {
		p.queue = p.queue[:0]
		p.head = 0
		return
	}
	if p.head > 0 && p.head > 256 {
		p.queue = append([]queuedPacket{}, p.queue[p.head:]...)
		p.head = 0
	}
}
