// 统计模块用于采集基础的流量与丢包指标。
package core

import "sync/atomic"

// Stats 保存运行时计数。
type Stats struct {
	videoInBytes    uint64
	videoOutBytes   uint64
	videoDrops      uint64
	controlInBytes  uint64
	controlOutBytes uint64
	audioInBytes    uint64
	audioOutBytes   uint64
}

// StatsSnapshot 为统计快照。
type StatsSnapshot struct {
	VideoInBytes    uint64 `json:"video_in_bytes"`
	VideoOutBytes   uint64 `json:"video_out_bytes"`
	VideoDrops      uint64 `json:"video_drops"`
	ControlInBytes  uint64 `json:"control_in_bytes"`
	ControlOutBytes uint64 `json:"control_out_bytes"`
	AudioInBytes    uint64 `json:"audio_in_bytes"`
	AudioOutBytes   uint64 `json:"audio_out_bytes"`
}

// Snapshot 生成当前统计快照。
func (s *Stats) Snapshot() StatsSnapshot {
	return StatsSnapshot{
		VideoInBytes:    atomic.LoadUint64(&s.videoInBytes),
		VideoOutBytes:   atomic.LoadUint64(&s.videoOutBytes),
		VideoDrops:      atomic.LoadUint64(&s.videoDrops),
		ControlInBytes:  atomic.LoadUint64(&s.controlInBytes),
		ControlOutBytes: atomic.LoadUint64(&s.controlOutBytes),
		AudioInBytes:    atomic.LoadUint64(&s.audioInBytes),
		AudioOutBytes:   atomic.LoadUint64(&s.audioOutBytes),
	}
}

func (s *Stats) AddVideoIn(n int) {
	if n > 0 {
		atomic.AddUint64(&s.videoInBytes, uint64(n))
	}
}

func (s *Stats) AddVideoOut(n int) {
	if n > 0 {
		atomic.AddUint64(&s.videoOutBytes, uint64(n))
	}
}

func (s *Stats) AddVideoDrop() {
	atomic.AddUint64(&s.videoDrops, 1)
}

func (s *Stats) AddControlIn(n int) {
	if n > 0 {
		atomic.AddUint64(&s.controlInBytes, uint64(n))
	}
}

func (s *Stats) AddControlOut(n int) {
	if n > 0 {
		atomic.AddUint64(&s.controlOutBytes, uint64(n))
	}
}

func (s *Stats) AddAudioIn(n int) {
	if n > 0 {
		atomic.AddUint64(&s.audioInBytes, uint64(n))
	}
}

func (s *Stats) AddAudioOut(n int) {
	if n > 0 {
		atomic.AddUint64(&s.audioOutBytes, uint64(n))
	}
}
