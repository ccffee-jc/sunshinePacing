// 配置模块负责加载 YAML 并合并默认值。
package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示代理的运行配置。
type Config struct {
	BasePort       int          `yaml:"base_port"`
	InternalOffset int          `yaml:"internal_offset"`
	InternalHost   string       `yaml:"internal_host"`
	Video          VideoConfig  `yaml:"video"`
	Control        StreamConfig `yaml:"control"`
	Audio          AudioConfig  `yaml:"audio"`
}

// VideoConfig 为视频 pacing 参数。
type VideoConfig struct {
	Enable          bool `yaml:"enable"`
	RateMbps        int  `yaml:"rate_mbps"`
	BurstKB         int  `yaml:"burst_kb"`
	MaxQueueDelayMs int  `yaml:"max_queue_delay_ms"`
	TickMs          int  `yaml:"tick_ms"`
}

// StreamConfig 为控制流的配置。
type StreamConfig struct {
	Enable bool   `yaml:"enable"`
	Mode   string `yaml:"mode"`
}

// AudioConfig 为音频流的配置。
type AudioConfig struct {
	Enable          bool   `yaml:"enable"`
	Mode            string `yaml:"mode"`
	MaxQueueDelayMs int    `yaml:"max_queue_delay_ms"`
}

// DefaultConfig 返回包含默认值的配置。
func DefaultConfig() Config {
	return Config{
		BasePort:       47989,
		InternalOffset: 1000,
		InternalHost:   "127.0.0.1",
		Video: VideoConfig{
			Enable:          true,
			RateMbps:        16,
			BurstKB:         16,
			MaxQueueDelayMs: 8,
			TickMs:          1,
		},
		Control: StreamConfig{
			Enable: true,
			Mode:   "bypass",
		},
		Audio: AudioConfig{
			Enable:          true,
			Mode:            "priority",
			MaxQueueDelayMs: 10,
		},
	}
}

// Load 读取配置文件并合并默认值。
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, errors.New("配置文件路径为空")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("读取配置失败: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("解析配置失败: %w", err)
	}
	if err := cfg.NormalizeAndValidate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// NormalizeAndValidate 补齐默认值并校验配置。
func (c *Config) NormalizeAndValidate() error {
	if c.InternalOffset == 0 {
		c.InternalOffset = 1000
	}
	if c.InternalHost == "" {
		c.InternalHost = "127.0.0.1"
	}
	if c.Video.TickMs == 0 {
		c.Video.TickMs = 1
	}
	if c.Video.RateMbps == 0 {
		c.Video.RateMbps = 16
	}
	if c.Video.BurstKB == 0 {
		c.Video.BurstKB = 16
	}
	if c.Video.MaxQueueDelayMs == 0 {
		c.Video.MaxQueueDelayMs = 8
	}
	if c.Audio.MaxQueueDelayMs == 0 {
		c.Audio.MaxQueueDelayMs = 10
	}
	if c.Control.Mode == "" {
		c.Control.Mode = "bypass"
	}
	if c.Audio.Mode == "" {
		c.Audio.Mode = "priority"
	}
	if c.BasePort == 0 {
		return errors.New("base_port 必填")
	}
	if c.BasePort < 1024 || c.BasePort > 65535 {
		return fmt.Errorf("base_port 超出范围: %d", c.BasePort)
	}
	if c.InternalOffset < 1 || c.InternalOffset > 60000 {
		return fmt.Errorf("internal_offset 超出范围: %d", c.InternalOffset)
	}
	if c.InternalHost == "" {
		return errors.New("internal_host 不能为空")
	}
	return nil
}
