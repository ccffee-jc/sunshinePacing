// 转发循环负责 UDP 收发与分类处理。
package core

import (
	"errors"
	"log"
	"net"
	"time"
)

const readTimeout = 500 * time.Millisecond

func (p *Proxy) startLoops() {
	p.startExternal(StreamVideo, p.extVideo, p.intVideo)
	if p.cfg.Control.Enable {
		p.startExternal(StreamControl, p.extControl, p.intControl)
	} else {
		log.Printf("控制流已禁用，外部端口不转发")
	}
	if p.cfg.Audio.Enable {
		p.startExternal(StreamAudio, p.extAudio, p.intAudio)
	} else {
		log.Printf("音频流已禁用，外部端口不转发")
	}

	p.startInternal(StreamVideo, p.intVideo, p.extVideo)
	if p.cfg.Control.Enable {
		p.startInternal(StreamControl, p.intControl, p.extControl)
	}
	if p.cfg.Audio.Enable {
		p.startInternal(StreamAudio, p.intAudio, p.extAudio)
	}
}

func (p *Proxy) startExternal(stream StreamType, ext *net.UDPConn, internal *net.UDPConn) {
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
				log.Printf("外部端口读取失败: stream=%v err=%v", stream, err)
				continue
			}
			p.session.SetClient(stream, addr)
			p.addStreamIn(stream, n)
			if _, err := internal.Write(buf[:n]); err != nil {
				log.Printf("转发到 Sunshine 失败: stream=%v err=%v", stream, err)
			}
		}
	}()
}

func (p *Proxy) startInternal(stream StreamType, internal *net.UDPConn, ext *net.UDPConn) {
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
				log.Printf("内部端口读取失败: stream=%v err=%v", stream, err)
				continue
			}
			addr := p.session.GetClient(stream)
			if addr == nil {
				if stream == StreamVideo {
					p.stats.AddVideoDrop()
				}
				continue
			}

			switch stream {
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
			}
		}
	}()
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
