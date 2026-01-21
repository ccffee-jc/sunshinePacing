# 任务清单: Sunshine UDP Pacing Proxy

目录: `helloagents/plan/202601211643_sunshine_pacing_proxy/`

---

## 1. 配置与端口映射
- [√] 1.1 在 `internal/config/config.go` 中实现配置结构体与默认值合并，验证 why.md#需求-本机-udp-代理转发-场景-单客户端基础转发
- [√] 1.2 在 `internal/core/ports.go` 中实现 base_port → 外部/内部端口映射，验证 why.md#需求-本机-udp-代理转发-场景-单客户端基础转发，依赖任务1.1

## 2. 会话与转发核心
- [√] 2.1 在 `internal/core/session.go` 中实现单客户端映射与超时清理，验证 why.md#需求-本机-udp-代理转发-场景-单客户端基础转发
- [√] 2.2 在 `internal/core/relay.go` 中实现 UDP 双向转发（control/audio 直通，video 进入队列），验证 why.md#需求-本机-udp-代理转发-场景-单客户端基础转发，依赖任务2.1

## 3. Video Pacer
- [√] 3.1 在 `internal/core/pacer.go` 中实现 token bucket + max_queue_delay 丢弃，验证 why.md#需求-视频-pacing-场景-突发压平

## 4. CLI 入口
- [√] 4.1 在 `cmd/proxy/main.go` 中实现 CLI 参数与启动逻辑，验证 why.md#需求-windows-gui-与-linux-cli-场景-gui-cli-启动

## 5. Windows GUI
- [√] 5.1 在 `cmd/proxy-gui/main.go` 中实现 Fyne GUI（配置编辑、启动/停止、状态展示），验证 why.md#需求-windows-gui-与-linux-cli-场景-gui-cli-启动

## 6. 安全检查
- [√] 6.1 执行安全检查（按G9: 端口占用检测、日志脱敏、避免敏感信息输出）

## 7. 文档更新
- [√] 7.1 更新 `helloagents/wiki/modules/core.md`
- [√] 7.2 更新 `helloagents/wiki/modules/config.md`
- [√] 7.3 更新 `helloagents/wiki/modules/cli.md`
- [√] 7.4 更新 `helloagents/wiki/modules/gui-win.md`
- [√] 7.5 更新 `helloagents/wiki/arch.md`（如与实现不一致）

## 8. 测试
- [√] 8.1 在 `internal/core/pacer_test.go` 中覆盖 pacing 速率与丢弃逻辑
- [√] 8.2 在 `internal/core/ports_test.go` 中覆盖端口映射逻辑
