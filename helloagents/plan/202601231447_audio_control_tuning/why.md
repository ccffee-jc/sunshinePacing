# 变更提案: 配置调参与音频轻量限制

## 需求背景
用户希望按新的视频 pacing 参数运行，并明确要求 control 直接放行、audio 进行轻量限制，以提升连接稳定性与音视频体验。

## 变更内容
1. 更新示例配置：应用新的 video 参数，设置 control bypass，audio 轻量限制。
2. 在核心转发中实现 audio 轻量限制逻辑，保持 control 直通。
3. 补充配置与核心模块文档说明。

## 影响范围
- **模块:** config / core / cli
- **文件:** internal/config/config.go, internal/core/pacer.go, internal/core/proxy.go, internal/core/relay.go, internal/core/pacer_test.go, proxy.yml, helloagents/wiki/modules/config.md, helloagents/wiki/modules/core.md, helloagents/wiki/modules/cli.md
- **API:** 无
- **数据:** 无

## 核心场景

### 需求: 视频参数更新
**模块:** config
将 video pacing 参数调整为用户提供的值。

#### 场景: 新视频参数启用
配置 rate_mbps=14, burst_kb=8, max_queue_delay_ms=5, tick_ms=1。
- 预期结果1：视频 pacing 参数以配置为准
- 预期结果2：保持其它默认行为不变

### 需求: Control 直通
**模块:** core
Control 流量不做排队或限速。

#### 场景: Control 连接
Control 包直接转发。
- 预期结果1：无额外延迟
- 预期结果2：不受 pacing 影响

### 需求: Audio 轻量限制
**模块:** core
Audio 流量在保持低延迟的前提下做轻量 pacing。

#### 场景: Audio 连接
启用轻量限制模式。
- 预期结果1：音频包仍然低延迟
- 预期结果2：微突发被轻量平滑

## 风险评估
- **风险:** 轻量限制参数不合适可能增加音频延迟。
- **缓解:** 采用高于音频常规码率的限速默认值，保留可配置项。
