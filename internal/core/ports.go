// 端口映射负责根据 base_port 推导外部与内部端口。
package core

import "fmt"

const (
	videoOffset   = 9
	controlOffset = 10
	audioOffset   = 11
)

// Ports 表示一组端口。
type Ports struct {
	Video   int
	Control int
	Audio   int
}

// PortMap 表示外部与内部端口映射。
type PortMap struct {
	External Ports
	Internal Ports
}

// BuildPortMap 构建端口映射。
func BuildPortMap(basePort int, internalOffset int) (PortMap, error) {
	if basePort <= 0 {
		return PortMap{}, fmt.Errorf("base_port 无效: %d", basePort)
	}
	return PortMap{
		External: Ports{
			Video:   basePort + internalOffset + videoOffset,
			Control: basePort + internalOffset + controlOffset,
			Audio:   basePort + internalOffset + audioOffset,
		},
		Internal: Ports{
			Video:   basePort + videoOffset,
			Control: basePort + controlOffset,
			Audio:   basePort + audioOffset,
		},
	}, nil
}
