# 数据模型

## 概述
本项目的主要数据为配置文件与运行时统计数据（内存态）。

---

## 数据表/集合

### 配置文件（YAML）
**描述:** 代理运行参数与默认值配置。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| base_port | int | 必填/默认47989 | Sunshine 基准端口 |
| internal_offset | int | 可选/默认1000 | 对外端口偏移（客户端使用 base_port+offset） |
| internal_host | string | 可选/默认127.0.0.1 | Sunshine 内部目标地址 |
| video.enable | bool | 默认true | 是否启用视频 pacing |
| video.rate_mbps | int | 默认16 | 放行速率 |
| video.burst_kb | int | 默认16 | 突发上限 |
| video.max_queue_delay_ms | int | 默认8 | 排队驻留上限 |
| video.tick_ms | int | 默认1 | 调度步长 |
| control.enable | bool | 默认true | 控制流开关 |
| control.mode | string | 默认bypass | 控制流策略 |
| audio.enable | bool | 默认true | 音频流开关 |
| audio.mode | string | 默认priority | 音频策略 |
| audio.max_queue_delay_ms | int | 默认10 | 音频队列上限 |

**关联关系:**
- 无。
