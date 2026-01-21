// 代理核心负责初始化 UDP 端口并运行转发与 pacing。
package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"sunshinePacing/internal/config"
)

const defaultSessionTimeout = 10 * time.Second

// Metrics 为 GUI/CLI 提供的运行快照。
type Metrics struct {
	StatsSnapshot
	VideoQueueLen int
}

// Proxy 表示代理实例。
type Proxy struct {
	cfg     config.Config
	ports   PortMap
	session *Session
	pacer   *Pacer
	stats   *Stats

	extVideo   *net.UDPConn
	extControl *net.UDPConn
	extAudio   *net.UDPConn
	intVideo   *net.UDPConn
	intControl *net.UDPConn
	intAudio   *net.UDPConn

	cancel  context.CancelFunc
	ctx     context.Context
	wg      sync.WaitGroup
	running uint32
}

// NewProxy 创建代理实例。
func NewProxy(cfg config.Config) (*Proxy, error) {
	if err := cfg.NormalizeAndValidate(); err != nil {
		return nil, err
	}
	ports, err := BuildPortMap(cfg.BasePort, cfg.InternalOffset)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		cfg:     cfg,
		ports:   ports,
		session: NewSession(defaultSessionTimeout),
		stats:   &Stats{},
	}, nil
}

// Start 启动代理。
func (p *Proxy) Start(ctx context.Context) error {
	if atomic.LoadUint32(&p.running) == 1 {
		return errors.New("代理已在运行")
	}
	p.ctx, p.cancel = context.WithCancel(ctx)

	if err := p.openSockets(); err != nil {
		p.closeSockets()
		return err
	}

	if p.cfg.Video.Enable {
		p.pacer = NewPacer(
			p.cfg.Video.RateMbps,
			p.cfg.Video.BurstKB,
			time.Duration(p.cfg.Video.MaxQueueDelayMs)*time.Millisecond,
			time.Duration(p.cfg.Video.TickMs)*time.Millisecond,
			func(data []byte, addr *net.UDPAddr) error {
				_, err := p.extVideo.WriteToUDP(data, addr)
				return err
			},
			p.stats,
		)
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.pacer.Run(p.ctx)
		}()
	}

	p.startLoops()
	atomic.StoreUint32(&p.running, 1)
	log.Printf("代理启动: external=%+v internal=%+v", p.ports.External, p.ports.Internal)
	return nil
}

// Stop 停止代理。
func (p *Proxy) Stop() {
	if atomic.LoadUint32(&p.running) == 0 {
		return
	}
	if p.cancel != nil {
		p.cancel()
	}
	p.closeSockets()
	p.wg.Wait()
	atomic.StoreUint32(&p.running, 0)
	log.Printf("代理已停止")
}

// Metrics 返回运行统计快照。
func (p *Proxy) Metrics() Metrics {
	queueLen := 0
	if p.pacer != nil {
		queueLen = p.pacer.QueueLen()
	}
	return Metrics{
		StatsSnapshot: p.stats.Snapshot(),
		VideoQueueLen: queueLen,
	}
}

func (p *Proxy) openSockets() error {
	var err error
	p.extVideo, err = listenUDP("0.0.0.0", p.ports.External.Video)
	if err != nil {
		return err
	}
	p.extControl, err = listenUDP("0.0.0.0", p.ports.External.Control)
	if err != nil {
		return err
	}
	p.extAudio, err = listenUDP("0.0.0.0", p.ports.External.Audio)
	if err != nil {
		return err
	}
	p.intVideo, err = dialUDP(p.cfg.InternalHost, p.ports.Internal.Video)
	if err != nil {
		return err
	}
	p.intControl, err = dialUDP(p.cfg.InternalHost, p.ports.Internal.Control)
	if err != nil {
		return err
	}
	p.intAudio, err = dialUDP(p.cfg.InternalHost, p.ports.Internal.Audio)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) closeSockets() {
	closeConn(p.extVideo)
	closeConn(p.extControl)
	closeConn(p.extAudio)
	closeConn(p.intVideo)
	closeConn(p.intControl)
	closeConn(p.intAudio)
}

func listenUDP(host string, port int) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func dialUDP(host string, port int) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeConn(conn *net.UDPConn) {
	if conn != nil {
		_ = conn.Close()
	}
}
