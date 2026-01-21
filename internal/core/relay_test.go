// TCP 转发测试覆盖基础双向转发能力。
package core

import (
	"io"
	"net"
	"testing"
	"time"

	"sunshinePacing/internal/config"
)

func TestHandleTCPConn(t *testing.T) {
	internalListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("创建内部监听失败: %v", err)
	}
	defer internalListener.Close()
	internalPort := internalListener.Addr().(*net.TCPAddr).Port

	echoDone := make(chan struct{})
	go func() {
		defer close(echoDone)
		conn, err := internalListener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4)
		n, err := io.ReadFull(conn, buf)
		if err != nil {
			return
		}
		_, _ = conn.Write(buf[:n])
	}()

	externalListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("创建外部监听失败: %v", err)
	}
	defer externalListener.Close()

	proxy := &Proxy{cfg: config.Config{InternalHost: "127.0.0.1"}}
	entry := PortEntry{
		PortSpec: PortSpec{
			Name:     "rtsp",
			Protocol: ProtocolTCP,
		},
		InternalPort: internalPort,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := externalListener.Accept()
		if err != nil {
			return
		}
		proxy.handleTCPConn(entry, conn)
	}()

	client, err := net.Dial("tcp", externalListener.Addr().String())
	if err != nil {
		t.Fatalf("客户端连接失败: %v", err)
	}
	defer client.Close()

	_ = client.SetDeadline(time.Now().Add(2 * time.Second))
	payload := []byte("ping")
	if _, err := client.Write(payload); err != nil {
		t.Fatalf("发送失败: %v", err)
	}
	reply := make([]byte, len(payload))
	if _, err := io.ReadFull(client, reply); err != nil {
		t.Fatalf("读取失败: %v", err)
	}
	if string(reply) != string(payload) {
		t.Fatalf("回包不匹配: got=%s want=%s", reply, payload)
	}
	_ = client.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("TCP 转发未按时完成")
	}
	select {
	case <-echoDone:
	case <-time.After(2 * time.Second):
		t.Fatal("回环服务未按时完成")
	}
}
