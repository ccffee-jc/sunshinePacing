// 端口映射负责根据 base_port 推导外部与内部端口。
package core

import "fmt"

// Protocol 表示传输协议。
type Protocol string

const (
	ProtocolUDP Protocol = "udp"
	ProtocolTCP Protocol = "tcp"
)

// PortSpec 描述端口偏移与用途。
type PortSpec struct {
	Name     string
	Offset   int
	Protocol Protocol
	Stream   StreamType
}

// PortEntry 表示具备实际端口号的映射。
type PortEntry struct {
	PortSpec
	ExternalPort int
	InternalPort int
}

// PortMap 汇总 UDP/TCP 端口映射。
type PortMap struct {
	UDP []PortEntry
	TCP []PortEntry
}

var defaultPortSpecs = []PortSpec{
	{Name: "https", Offset: -5, Protocol: ProtocolTCP, Stream: StreamOther},
	{Name: "http", Offset: 0, Protocol: ProtocolTCP, Stream: StreamOther},
	{Name: "web", Offset: 1, Protocol: ProtocolTCP, Stream: StreamOther},
	{Name: "rtsp", Offset: 21, Protocol: ProtocolTCP, Stream: StreamOther},
	{Name: "video", Offset: 9, Protocol: ProtocolUDP, Stream: StreamVideo},
	{Name: "control", Offset: 10, Protocol: ProtocolUDP, Stream: StreamControl},
	{Name: "audio", Offset: 11, Protocol: ProtocolUDP, Stream: StreamAudio},
	{Name: "mic", Offset: 13, Protocol: ProtocolUDP, Stream: StreamOther},
	{Name: "rtsp-udp", Offset: 21, Protocol: ProtocolUDP, Stream: StreamOther},
}

// BuildPortMap 构建端口映射。
func BuildPortMap(basePort int, internalOffset int) (PortMap, error) {
	if basePort <= 0 {
		return PortMap{}, fmt.Errorf("base_port 无效: %d", basePort)
	}
	var udpEntries []PortEntry
	var tcpEntries []PortEntry
	required := map[StreamType]bool{
		StreamVideo:   false,
		StreamControl: false,
		StreamAudio:   false,
	}
	for _, spec := range defaultPortSpecs {
		entry, err := buildEntry(spec, basePort, internalOffset)
		if err != nil {
			return PortMap{}, err
		}
		if spec.Protocol == ProtocolUDP {
			udpEntries = append(udpEntries, entry)
		} else {
			tcpEntries = append(tcpEntries, entry)
		}
		if spec.Stream != StreamOther {
			required[spec.Stream] = true
		}
	}
	for stream, ok := range required {
		if !ok {
			return PortMap{}, fmt.Errorf("缺少必需端口: %v", stream)
		}
	}
	return PortMap{
		UDP: udpEntries,
		TCP: tcpEntries,
	}, nil
}

func buildEntry(spec PortSpec, basePort int, internalOffset int) (PortEntry, error) {
	externalPort := basePort + internalOffset + spec.Offset
	internalPort := basePort + spec.Offset
	if err := validatePort(externalPort); err != nil {
		return PortEntry{}, fmt.Errorf("外部端口无效(%s/%s): %w", spec.Name, spec.Protocol, err)
	}
	if err := validatePort(internalPort); err != nil {
		return PortEntry{}, fmt.Errorf("内部端口无效(%s/%s): %w", spec.Name, spec.Protocol, err)
	}
	return PortEntry{
		PortSpec:     spec,
		ExternalPort: externalPort,
		InternalPort: internalPort,
	}, nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("端口超出范围: %d", port)
	}
	return nil
}
