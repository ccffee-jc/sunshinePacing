# Sunshine UDP Pacing Proxy (方案 B：本机代理转发) 设计与落地文档

> 目标：在 **Windows Sunshine 主机本机** 部署一个 **UDP 代理转发（relay/proxy）+ pacing** 程序，  
> 通过“端口对齐 + 按端口语义差异化队列/限速”，把 Sunshine 的 **视频突发（burst）** 压平，降低运营商内部整形导致的卡顿。  
>  
> 本方案采用 **方案 B：本机代理转发**（不做透明抓包/重注入），通过“外部端口由 Proxy 占用、Sunshine 改到内部端口”实现。

---

## 1. 背景与问题描述

- Sunshine 在画面复杂或出现参考刷新/关键帧时，可能产生 **短时间内的大 UDP 发送突发**（微突发/帧级 burst）。
- 若上游（运营商内部）存在 **policer/token-bucket** 或小 buffer，对 UDP 突发敏感，可能表现为：
  - RTT/jitter 瞬时升高或丢包成坨
  - 客户端偶发整屏变糊、卡顿
- 你的需求明确：
  - ✅ 可以接受持续较高码率
  - ❌ 不能接受突发码率尖峰
- 串流端口语义（默认 `base_port=47989`）：
  - Video：UDP `base+9`（默认 47998）
  - Control（键鼠/手柄）：UDP `base+10`（默认 47999）
  - Audio：UDP `base+11`（默认 48000）

---

## 2. 方案概览（本机代理转发）

### 2.1 拓扑

```
Moonlight Client  <——UDP——>  Pacing Proxy (外部端口/默认端口)
                                 |
                                 |（本机 UDP 转发）
                                 v
                           Sunshine (内部端口)
```

- Proxy 监听 **外部端口族（与 Sunshine 默认一致）**，对外表现成“Sunshine”。
- Sunshine 改为监听 **内部端口族**（与外部端口固定偏移对齐）。
- Proxy 在两端之间转发数据，并对 **Video** 方向做 pacing。

### 2.2 为什么这种代理能工作（关键机制）

- Sunshine 并不要求“看到真实客户端 IP/端口”，它只需要有一个“对端”在说 GameStream 协议。
- Proxy 作为“对端”，与 Sunshine 形成一个 UDP 会话；同时 Proxy 与真实客户端也形成 UDP 会话。
- Proxy 维护两边的映射关系：
  - **client endpoint**（客户端公网/局域网 IP:port）
  - **sunshine endpoint**（本机 127.0.0.1:internal_port）
- Sunshine 只和 Proxy 通信；Proxy 再把数据转发给真实客户端。

> 注意：这是标准的 UDP relay/ALG 模式，不需要透明 spoof 源地址。

---

## 3. 端口对齐与单一配置（Base Port）

### 3.1 单一配置项

- `base_port`：唯一必填/可配配置项（Sunshine 仍使用该基准端口）  
  - 默认：`47989`

Proxy 内部通过差值表推导端口：

| 端口用途 | 外部端口（对客户端） | 差值 | 内部端口（对 Sunshine） |
|---|---:|---:|---:|
| Video | `base+offset+9` | +9 | `base+9` |
| Control | `base+offset+10` | +10 | `base+10` |
| Audio | `base+offset+11` | +11 | `base+11` |

### 3.2 external_base（外部基准端口）

为了保持“用户只配一个端口”，外部端口基准采用固定偏移：

- `external_base = base_port + 1000`（固定规则，默认不暴露给用户）
  - 例如 base=47989 → external_base=48989
  - 对外 Video=48998、Control=48999、Audio=49000

> 若你们未来需要支持非默认偏移，可在高级配置中开放 `internal_offset`（实际为外部端口偏移），但**默认只保留 base_port**。

---

## 4. Sunshine 配置要求（必须）

### 4.1 Sunshine 端口保持默认（无需修改）

Sunshine 仍使用 `base_port` 作为其 base port：
- 客户端访问 `base_port + offset` 推导出来的 48998/48999/49000
- Sunshine 实际监听 `base_port` 推导出来的 47998/47999/48000（可在 internal_host 指定其地址）

> **关键点**：Proxy 对外使用偏移端口族，因此 Sunshine 无需修改。

### 4.2 防火墙建议

- Windows 防火墙允许：
  - Proxy：对外 UDP 48998-49000（以及你们实际用到的其它控制端口）
  - Sunshine：只需本机回环访问内部端口（127.0.0.1:47998-48000），可不放行到公网
- 建议将 Sunshine 绑定到 `127.0.0.1` 或仅内网接口，减少被扫描风险。

---

## 5. Proxy 的功能模块

### 5.1 模块划分

1. **Port Mapper**
   - 输入：`base_port`
   - 输出：外部/内部端口族（video/control/audio）

2. **Session Manager（会话管理）**
   - 维护 “客户端 ↔ Sunshine” 的会话映射
   - 支持单客户端（第一期）与多客户端扩展（后述）

3. **UDP Relay（双向转发）**
   - client → proxy → sunshine
   - sunshine → proxy → client

4. **Traffic Classifier（按端口分类）**
   - base+offset+9：Video
   - base+offset+10：Control（键鼠/手柄输入）
   - base+offset+11：Audio

5. **Pacer / Scheduler（pacing 调度器）**
   - **只对 Video 的 sunshine→client 方向** 做 pacing（最主要的 burst 来源）
   - Control/Auido 走高优先级、尽量不排队

6. **Stats & Observability**
   - bps/pps、队列长度、排队时间分位数、丢包计数
   - 便于调参：rate、burst、queue_delay

---

## 6. 按端口语义的 pacing 策略

### 6.1 Control（UDP base+offset+10，默认 48999）
- 目标：输入手感稳定
- 策略：**直通**（bypass）或极小队列
- 不做 pacing、不做限速（或仅做极低保护阈值）

### 6.2 Audio（UDP base+offset+11，默认 49000）
- 目标：音频不卡顿、低抖动
- 策略：高优先级小队列，可选轻量 pacing
- 推荐：
  - queue_delay 上限 10ms
  - 不与视频共享队列

### 6.3 Video（UDP base+offset+9，默认 48998）
- 目标：压平突发；宁可偶发丢包/画质波动，也不要排队造成延迟跳变
- 策略：**Token Bucket / Leaky Bucket pacing + 严格队列驻留上限**
- 核心参数：
  - `rate_bps`：放行的持续速率（建议 16–18 Mbps；需结合你们上行 20Mbps 与 FEC）
  - `burst_bytes`：允许突发字节（建议 8–16 KB 起步，越小越平）
  - `max_queue_delay_ms`：队列驻留上限（建议 5–10 ms）
  - 超限策略：**deadline drop**（超过驻留上限直接丢）

---

## 7. 推荐默认参数（你的场景：1080p60、视频10Mbps、FEC开、上行约20Mbps）

> 说明：10 Mbps 对 1080p60 非常紧，任何参考刷新/关键帧/复杂纹理都会制造“更大的帧”。  
> 本 Proxy 目标是把 burst 压平；但如果平均预算太紧，仍可能出现“整屏糊一下”（编码器提高 QP）——这是用画质换稳定的正常代价。

**建议 Proxy 默认值：**
- Video `rate_mbps`: **16**
- Video `burst_kb`: **16**
- Video `max_queue_delay_ms`: **8**
- Control：bypass
- Audio：priority，小队列（10ms）

> 若仍出现运营商触发（丢包成坨/RTT尖峰），优先把 `burst_kb` 降到 8KB；  
> 若出现明显“排队导致手感肉”，把 `max_queue_delay_ms` 降到 5ms，并允许更多丢包（更糊但更稳）。

---

## 8. 代理转发实现细节（工程落地）

### 8.1 UDP socket 设计

为每个端口用途分别创建两个 socket：

- 外部监听 socket（对客户端）：
  - `ext_video_sock` 绑定 `0.0.0.0:base+9`
  - `ext_control_sock` 绑定 `0.0.0.0:base+10`
  - `ext_audio_sock` 绑定 `0.0.0.0:base+11`

- 内部通信 socket（对 Sunshine）：
  - `int_video_sock` 发送到 `internal_host:base_port+9`
  - `int_control_sock` 发送到 `internal_host:base_port+10`
  - `int_audio_sock` 发送到 `internal_host:base_port+11`
  - 同时也要接收 Sunshine 回包（可以 bind 到对应端口或用同一 socket 发送/接收）

> 推荐用“同端口绑定”方式：Proxy 内部端口也用固定端口监听，便于 Sunshine 回包定位。

### 8.2 会话映射（单客户端 MVP）

第一期建议只支持单客户端（足以验证效果）：

- 记录最近一次从外部端口收到包的 `client_endpoint`（IP:port）
- 将 Sunshine 回包全部转发到该 endpoint
- 若 client 超时（比如 10s 无包），清空映射

多客户端扩展见 8.4。

### 8.3 转发路径

**client → sunshine（通常不需要 pacing）**
- 收到外部包：
  - base+offset+10/11/9：直接转发到内部对应端口
- 目的：让 Sunshine 收到“来自 Proxy 的对端数据”

**sunshine → client（Video 做 pacing）**
- 收到内部包：
  - Control/Audio：立即转发到外部 client_endpoint
  - Video：进入 video 队列，由 pacer 按 rate 放行

### 8.4 多客户端扩展（可选）

若要支持多个客户端，需要以“会话 key”区分：

- 最简单：按 `client_ip` 区分（同一 IP 可能多个设备，仍有冲突风险）
- 更可靠：按 `(client_ip, client_port, dst_port)` 区分
- Sunshine 端可能同时存在多个 session；Proxy 需要在握手阶段建立映射关系（通常由 Control 流开始）

第一期建议先交付单客户端版本。

---

## 9. Pacer（token bucket）实现建议

### 9.1 数据结构

- `video_queue`: FIFO 队列（存储 packet + enqueue_time + client_endpoint）
- `tokens_bits`: 当前可用 token（bit）
- `rate_bps`: 每秒补充 token（bit/s）
- `burst_bits`: token 上限（bit），对应 burst_bytes*8
- `max_queue_delay_ms`: 丢弃阈值

### 9.2 调度循环（建议 1ms tick）

伪代码：

```text
loop every 1ms:
  now = time()
  tokens += rate_bps * dt
  tokens = min(tokens, burst_bits)

  # control/audio：无队列或极小队列，尽量在接收线程直接发送

  # video：先丢弃超时包
  while video_queue not empty and now - pkt.enqueue_time > max_queue_delay:
      drop(pkt)

  # video：按 token 放行
  while video_queue not empty:
      pkt = peek()
      cost = pkt.size_bytes * 8
      if tokens < cost:
          break
      sendto(ext_video_sock, pkt.data, client_endpoint)
      tokens -= cost
      pop()
```

### 9.3 丢弃策略

- 采用 **deadline drop**：超过驻留时间就丢
- 这会造成：
  - 偶发画面块状/短暂更糊（可接受）
  - 但能有效避免“延迟突然增加”

---

## 10. 配置文件（YAML）示例

> 只需用户设置 `base_port`；其它有默认值。  
> 若你们需要快速迭代，可先保留可选的“高级参数”字段。

```yaml
base_port: 47989

internal_offset: 1000   # 外部端口偏移；可隐藏为固定规则
internal_host: 127.0.0.1 # Sunshine 目标地址，可指向局域网

video:
  enable: true
  rate_mbps: 16
  burst_kb: 16
  max_queue_delay_ms: 8
  tick_ms: 1

control:
  enable: true
  mode: bypass

audio:
  enable: true
  mode: priority
  max_queue_delay_ms: 10
```

---

## 11. 运行与部署步骤（Windows）

1. 安装/运行 Proxy（建议作为 Windows Service）
2. 配置 Proxy：确认 `base_port` 为 Sunshine 基准端口（默认 47989）
3. 客户端连接：使用 `base_port + internal_offset`（默认 48989）
4. 防火墙：
   - 对外放行 Proxy 的 UDP 48998-49000
   - Sunshine 内部端口仅本机
5. 启动顺序：
   - 先启动 Sunshine（内部端口）
   - 再启动 Proxy（外部端口）
6. 客户端连接到 Sunshine 主机（与以前一致，不需要修改）

---

## 12. 验收与排障

### 12.1 成功标志

- Wireshark 抓包看外发 48998：
  - 从“每 16.7ms 一坨”变成“更均匀铺开”
- Moonlight overlay：
  - RTT/jitter 尖峰减少
  - “画面一动就卡一下”明显改善
- 输入（47999）：
  - 手感不变“肉”

### 12.2 常见问题

1. **客户端连不上**
   - Proxy 未监听外部端口（被其他程序占用）
   - Sunshine 未改内部端口导致冲突
   - 防火墙未放行 Proxy 外部端口

2. **画面偶发更糊**
   - 这是编码器在低预算下提高 QP 的正常表现
   - 若糊得频繁：提升视频平均码率（例如 12–14 Mbps）或降低分辨率到 900p

3. **手感变肉**
   - video_queue 过长（max_queue_delay 太大）
   - 降低 max_queue_delay（如 8→5ms），允许更多丢包换延迟稳定

4. **音频爆音/断续**
   - 音频被视频队列挤压（队列没分开）
   - 确保 audio 使用独立高优先级路径

---

## 13. 安全与兼容性注意事项

- Proxy 对外暴露端口，应限制访问来源（如果可能）
- 建议支持“仅允许最近成功握手的 client_ip”
- 支持日志脱敏（不记录完整公网 IP 可选）

---

## 14. 后续增强（建议路线图）

- 多客户端会话管理（握手识别与映射）
- DSCP 标记（内网/自家路由支持时可用）：
  - Control：EF46
  - Audio：AF41/EF46
  - Video：AF41/CS0
- 动态调参（根据 queue_drop 率自动调整 rate/burst）
- 统计面板（本地 HTTP / Prometheus）

---

# 附录：默认端口速查

当 `base_port=47989` 且 `internal_offset=1000`：
- 对外 Video：UDP 48998
- 对外 Control：UDP 48999（键鼠/手柄输入）
- 对外 Audio：UDP 49000
- Sunshine Video：UDP 47998
- Sunshine Control：UDP 47999
- Sunshine Audio：UDP 48000

---

> 这份文档是“方案 B：本机代理转发”的工程落地版。  
> 若后续发现 Sunshine/协议握手阶段需要额外端口对齐（发现/配对/RTSP 等），可在 Port Mapper 中继续按差值表扩展，但 pacing 核心仍以 base+offset+9/10/11 为主。
