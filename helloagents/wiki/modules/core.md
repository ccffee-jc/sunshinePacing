# 核心代理模块（core）

## 目的
实现 UDP relay、会话映射与视频 pacing 的核心逻辑。

## 模块概述
- **职责:** 端口映射、UDP/TCP 转发、视频队列与令牌桶、丢弃策略、统计采集
- **状态:** ✅稳定
- **最后更新:** 2026-01-23

## 规范

### 需求: 本机 UDP 代理转发
**模块:** core
在外部端口（base_port+offset）监听并转发到 Sunshine 端口（base_port），支持 internal_host 指向局域网目标。

#### 场景: 单客户端基础转发
满足 Sunshine 与客户端之间的 UDP 会话转发，Control/Audio 直通，Video 进入 pacing。
- 预期结果1：Sunshine 仅看到本机代理作为对端
- 预期结果2：客户端通信端口保持原有对外端口

#### 场景: 全端口转发
按 Sunshine 官方端口偏移表默认全开，支持 TCP/UDP 双栈转发。
- 预期结果1：HTTPS/HTTP/Web/RTSP 等 TCP 端口可用
- 预期结果2：Video/Control/Audio/Mic 等 UDP 端口可用
- 预期结果3：UDP 48010 兼容转发

### 需求: 视频 pacing
**模块:** core
对 sunshine→client 的视频方向做 token bucket pacing 并限制排队延迟。

#### 场景: 突发压平
限制 burst 并在驻留超过阈值时丢弃视频包。
- 预期结果1：外发包间隔更均匀
- 预期结果2：控制/音频不被视频排队影响

### 需求: 纯转发模式
**模块:** core
当 pacing.enable=false 时，video/control/audio 全部直通。

#### 场景: 关闭 pacing
禁用 pacing 后内部→外部直接发送。
- 预期结果1：video 不进入队列
- 预期结果2：control/audio 行为不变

### 需求: 连接日志
**模块:** core
在 UDP/TCP 有连接时输出一次日志，供排障使用。

#### 场景: UDP 首次接入
首次看到客户端地址时记录连接日志。
- 预期结果1：同一端口重复包不会反复输出
- 预期结果2：日志包含端口、协议、客户端与内部目标

#### 场景: TCP 接入
accept 新连接时记录连接日志。
- 预期结果1：每条连接记录一次
- 预期结果2：日志包含端口、协议、客户端与内部目标

### 需求: 回包日志
**模块:** core
在内部→外部转发路径输出逐包回包日志。

#### 场景: 内部→外部逐包记录
Sunshine 回包转发给客户端时输出日志。
- 预期结果1：每个包都记录日志
- 预期结果2：日志包含端口、协议、包大小、客户端与内部目标

## API接口
- 暂无对外 API

## 数据模型
- 运行时结构：session map、video queue、统计计数器

## 依赖
- config

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - 初始代理实现
- [202601220034_ports_forwarding](../../history/2026-01/202601220034_ports_forwarding/) - Sunshine 全端口转发支持
- [202601231353_connection_logging](../../history/2026-01/202601231353_connection_logging/) - 连接日志与配置开关
- [202601231513_response_logging](../../history/2026-01/202601231513_response_logging/) - 回包日志与配置开关
- [202601231529_pacing_bypass](../../history/2026-01/202601231529_pacing_bypass/) - 纯转发模式开关
