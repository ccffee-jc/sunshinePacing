# API 手册

## 概述
当前项目提供本地 metrics HTTP 接口用于 GUI 拉取运行指标，核心仍为 UDP relay/pacing 服务与本地 GUI/CLI 控制。

## 认证方式
无。

---

## 接口列表

### 本地指标

#### GET /metrics
**描述:** 返回代理运行指标快照（仅本地回环地址可用）。

**响应:**
```json
{
  "video_in_bytes": 0,
  "video_out_bytes": 0,
  "video_drops": 0,
  "control_in_bytes": 0,
  "control_out_bytes": 0,
  "audio_in_bytes": 0,
  "audio_out_bytes": 0,
  "video_queue_len": 0
}
```
