# 变更提案: Sunshine UDP Pacing Proxy

## 需求背景
Sunshine 在画面复杂或关键帧时会出现 UDP 视频突发，导致运营商整形触发卡顿。需要一个本机代理转发并对视频做 pacing 的服务，Windows 提供 GUI，Linux 通过配置启动。第一期为单客户端 MVP。

## 变更内容
1. 新增 UDP relay + video pacing 核心逻辑（单客户端会话）。
2. 新增配置系统（YAML + 默认值）。
3. 新增 Windows Fyne GUI 与 Linux CLI 入口。

## 影响范围
- **模块:** core, config, cli, gui-win
- **文件:** 新增 Go 代码与配置模板
- **API:** 无对外 API
- **数据:** YAML 配置结构

## 核心场景

### 需求: 本机 UDP 代理转发
**模块:** core
实现外部端口监听与内部端口转发，并维护单客户端会话映射。

#### 场景: 单客户端基础转发
客户端通过外部端口连接代理，Sunshine 使用内部端口通信。
- 预期结果：Control/Audio 直通，Video 进入 pacing 队列

### 需求: 视频 pacing
**模块:** core
对 sunshine→client 的视频方向执行 token bucket pacing 并限制排队时间。

#### 场景: 突发压平
通过 burst 与 max_queue_delay 限制突发并丢弃超时包。
- 预期结果：外发视频包更均匀，延迟不显著增加

### 需求: Windows GUI 与 Linux CLI
**模块:** gui-win / cli
Windows 提供 GUI 控制与状态查看，Linux 通过配置启动。

#### 场景: GUI/CLI 启动
GUI 能加载/保存配置并启动代理；CLI 通过 -config 启动。
- 预期结果：Windows 有可视化入口，Linux 无 GUI 依赖

## 风险评估
- **风险:** UDP 端口占用、Windows 防火墙拦截、pacing 参数不当导致画质或延迟问题。
- **缓解:** 启动前端口占用检测；提供默认参数与日志提示；允许调整配置。
