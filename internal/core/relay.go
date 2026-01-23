// 转发循环负责 UDP 收发与分类处理。
package core

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const readTimeout = 500 * time.Millisecond
const tcpDialTimeout = 5 * time.Second

func (p *Proxy) startLoops() {
	for _, entry := range p.ports.UDP {
		if !p.shouldEnableUDP(entry) {
			continue
		}
		ext := p.udpExternal[entry.ExternalPort]
		internal := p.udpInternal[entry.ExternalPort]
		if ext == nil || internal == nil {
			log.Printf("UDP 端口未就绪: name=%s protocol=%s port=%d", entry.Name, entry.Protocol, entry.ExternalPort)
			continue
		}
		p.startExternal(entry, ext, internal)
		p.startInternal(entry, internal, ext)
	}
	p.startTCPForwarders()
}

func (p *Proxy) shouldEnableUDP(entry PortEntry) bool {
	switch entry.Stream {
	case StreamControl:
		if !p.cfg.Control.Enable {
			log.Printf("控制流已禁用，外部端口不转发: port=%d", entry.ExternalPort)
			return false
		}
	case StreamAudio:
		if !p.cfg.Audio.Enable {
			log.Printf("音频流已禁用，外部端口不转发: port=%d", entry.ExternalPort)
			return false
		}
	}
	return true
}

func (p *Proxy) startExternal(entry PortEntry, ext *net.UDPConn, internal *net.UDPConn) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		buf := make([]byte, 64*1024)
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
			}
			_ = ext.SetReadDeadline(time.Now().Add(readTimeout))
			n, addr, err := ext.ReadFromUDP(buf)
			if err != nil {
				if isTimeout(err) {
					continue
				}
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("外部端口读取失败: name=%s port=%d err=%v", entry.Name, entry.ExternalPort, err)
				continue
			}
			changed := p.session.SetClient(entry.ExternalPort, addr)
			if changed && p.cfg.ConnectionLog.Enable {
				internalTarget := fmt.Sprintf("%s:%d", p.cfg.InternalHost, entry.InternalPort)
				log.Printf("UDP 连接: name=%s protocol=%s external=%d internal=%s client=%s", entry.Name, entry.Protocol, entry.ExternalPort, internalTarget, addr.String())
			}
			p.addStreamIn(entry.Stream, n)
			if _, err := internal.Write(buf[:n]); err != nil {
				log.Printf("转发到 Sunshine 失败: name=%s port=%d err=%v", entry.Name, entry.InternalPort, err)
			}
		}
	}()
}

func (p *Proxy) startInternal(entry PortEntry, internal *net.UDPConn, ext *net.UDPConn) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		buf := make([]byte, 64*1024)
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
			}
			_ = internal.SetReadDeadline(time.Now().Add(readTimeout))
			n, err := internal.Read(buf)
			if err != nil {
				if isTimeout(err) {
					continue
				}
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("内部端口读取失败: name=%s port=%d err=%v", entry.Name, entry.InternalPort, err)
				continue
			}
			addr := p.session.GetClient(entry.ExternalPort)
			if addr == nil {
				if entry.Stream == StreamVideo {
					p.stats.AddVideoDrop()
				}
				continue
			}

			switch entry.Stream {
			case StreamVideo:
				if p.pacer != nil {
					p.stats.AddVideoIn(n)
					p.pacer.Enqueue(buf[:n], addr)
					continue
				}
				if err := sendDirect(ext, addr, buf[:n]); err == nil {
					p.stats.AddVideoOut(n)
				}
			case StreamControl:
				p.stats.AddControlIn(n)
				if err := sendDirect(ext, addr, buf[:n]); err == nil {
					p.stats.AddControlOut(n)
				}
			case StreamAudio:
				p.stats.AddAudioIn(n)
				if err := sendDirect(ext, addr, buf[:n]); err == nil {
					p.stats.AddAudioOut(n)
				}
			default:
				if err := sendDirect(ext, addr, buf[:n]); err != nil {
					log.Printf("内部端口转发失败: name=%s port=%d err=%v", entry.Name, entry.InternalPort, err)
				}
			}
		}
	}()
}

func (p *Proxy) startTCPForwarders() {
	for _, item := range p.tcpListeners {
		listener := item.listener
		entry := item.entry
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				conn, err := listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}
					log.Printf("TCP 接收失败: name=%s port=%d err=%v", entry.Name, entry.ExternalPort, err)
					continue
				}
				if p.cfg.ConnectionLog.Enable {
					target := fmt.Sprintf("%s:%d", p.cfg.InternalHost, entry.InternalPort)
					log.Printf("TCP 连接: name=%s protocol=%s external=%d internal=%s client=%s", entry.Name, entry.Protocol, entry.ExternalPort, target, conn.RemoteAddr().String())
				}
				p.wg.Add(1)
				go func(c net.Conn) {
					defer p.wg.Done()
					p.handleTCPConn(entry, c)
				}(conn)
			}
		}()
	}
}

func (p *Proxy) handleTCPConn(entry PortEntry, client net.Conn) {
	defer client.Close()
	target := fmt.Sprintf("%s:%d", p.cfg.InternalHost, entry.InternalPort)
	server, err := net.DialTimeout("tcp", target, tcpDialTimeout)
	if err != nil {
		log.Printf("TCP 连接 Sunshine 失败: name=%s port=%d err=%v", entry.Name, entry.InternalPort, err)
		return
	}
	defer server.Close()
	if p.ctx != nil && p.ctx.Done() != nil {
		go func() {
			<-p.ctx.Done()
			_ = client.Close()
			_ = server.Close()
		}()
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go proxyTCPCopy(server, client, &wg)
	go proxyTCPCopy(client, server, &wg)
	wg.Wait()
}

func proxyTCPCopy(dst net.Conn, src net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, _ = io.Copy(dst, src)
	closeWrite(dst)
	closeRead(src)
}

func closeWrite(conn net.Conn) {
	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		_ = tcpConn.CloseWrite()
		return
	}
	_ = conn.Close()
}

func closeRead(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.CloseRead()
	}
}

func sendDirect(conn *net.UDPConn, addr *net.UDPAddr, data []byte) error {
	if conn == nil || addr == nil {
		return errors.New("连接或地址为空")
	}
	_, err := conn.WriteToUDP(data, addr)
	return err
}

func (p *Proxy) addStreamIn(stream StreamType, n int) {
	switch stream {
	case StreamVideo:
		p.stats.AddVideoIn(n)
	case StreamControl:
		p.stats.AddControlIn(n)
	case StreamAudio:
		p.stats.AddAudioIn(n)
	}
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}
