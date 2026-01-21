# 核心代理模块（core）

## 目的
实现 UDP relay、会话映射与视频 pacing 的核心逻辑。

## 模块概述
- **职责:** 端口映射、收发转发、视频队列与令牌桶、丢弃策略、统计采集
- **状态:** ✅稳定
- **最后更新:** 2026-01-21

## 规范

### 需求: 本机 UDP 代理转发
**模块:** core
在外部端口（base_port+offset）监听并转发到 Sunshine 端口（base_port），支持 internal_host 指向局域网目标。

#### 场景: 单客户端基础转发
满足 Sunshine 与客户端之间的 UDP 会话转发，Control/Audio 直通，Video 进入 pacing。
- 预期结果1：Sunshine 仅看到本机代理作为对端
- 预期结果2：客户端通信端口保持原有对外端口

### 需求: 视频 pacing
**模块:** core
对 sunshine→client 的视频方向做 token bucket pacing 并限制排队延迟。

#### 场景: 突发压平
限制 burst 并在驻留超过阈值时丢弃视频包。
- 预期结果1：外发包间隔更均匀
- 预期结果2：控制/音频不被视频排队影响

## API接口
- 暂无对外 API

## 数据模型
- 运行时结构：session map、video queue、统计计数器

## 依赖
- config

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - 初始代理实现
