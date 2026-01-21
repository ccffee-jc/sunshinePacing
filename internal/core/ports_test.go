// 端口映射测试验证偏移规则。
package core

import "testing"

func TestBuildPortMap(t *testing.T) {
	ports, err := BuildPortMap(47989, 1000)
	if err != nil {
		t.Fatalf("构建端口失败: %v", err)
	}
	video := mustFindPort(t, ports.UDP, "video")
	if video.ExternalPort != 48998 || video.InternalPort != 47998 {
		t.Fatalf("video 端口不符合预期: %+v", video)
	}
	control := mustFindPort(t, ports.UDP, "control")
	if control.ExternalPort != 48999 || control.InternalPort != 47999 {
		t.Fatalf("control 端口不符合预期: %+v", control)
	}
	audio := mustFindPort(t, ports.UDP, "audio")
	if audio.ExternalPort != 49000 || audio.InternalPort != 48000 {
		t.Fatalf("audio 端口不符合预期: %+v", audio)
	}
	mic := mustFindPort(t, ports.UDP, "mic")
	if mic.ExternalPort != 49002 || mic.InternalPort != 48002 {
		t.Fatalf("mic 端口不符合预期: %+v", mic)
	}
	rtspUDP := mustFindPort(t, ports.UDP, "rtsp-udp")
	if rtspUDP.ExternalPort != 49010 || rtspUDP.InternalPort != 48010 {
		t.Fatalf("rtsp-udp 端口不符合预期: %+v", rtspUDP)
	}
	http := mustFindPort(t, ports.TCP, "http")
	if http.ExternalPort != 48989 || http.InternalPort != 47989 {
		t.Fatalf("http 端口不符合预期: %+v", http)
	}
	web := mustFindPort(t, ports.TCP, "web")
	if web.ExternalPort != 48990 || web.InternalPort != 47990 {
		t.Fatalf("web 端口不符合预期: %+v", web)
	}
	https := mustFindPort(t, ports.TCP, "https")
	if https.ExternalPort != 48984 || https.InternalPort != 47984 {
		t.Fatalf("https 端口不符合预期: %+v", https)
	}
	rtsp := mustFindPort(t, ports.TCP, "rtsp")
	if rtsp.ExternalPort != 49010 || rtsp.InternalPort != 48010 {
		t.Fatalf("rtsp 端口不符合预期: %+v", rtsp)
	}
}

func TestBuildPortMapInvalid(t *testing.T) {
	if _, err := BuildPortMap(0, 1000); err == nil {
		t.Fatal("期望 base_port 无效时报错")
	}
}

func mustFindPort(t *testing.T, entries []PortEntry, name string) PortEntry {
	t.Helper()
	for _, entry := range entries {
		if entry.Name == name {
			return entry
		}
	}
	t.Fatalf("未找到端口: %s", name)
	return PortEntry{}
}
