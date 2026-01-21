// 端口映射测试验证偏移规则。
package core

import "testing"

func TestBuildPortMap(t *testing.T) {
	ports, err := BuildPortMap(47989, 1000)
	if err != nil {
		t.Fatalf("构建端口失败: %v", err)
	}
	if ports.External.Video != 48998 || ports.External.Control != 48999 || ports.External.Audio != 49000 {
		t.Fatalf("外部端口不符合预期: %+v", ports.External)
	}
	if ports.Internal.Video != 47998 || ports.Internal.Control != 47999 || ports.Internal.Audio != 48000 {
		t.Fatalf("内部端口不符合预期: %+v", ports.Internal)
	}
}

func TestBuildPortMapInvalid(t *testing.T) {
	if _, err := BuildPortMap(0, 1000); err == nil {
		t.Fatal("期望 base_port 无效时报错")
	}
}
