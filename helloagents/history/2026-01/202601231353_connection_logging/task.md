# 任务清单: 连接日志与配置开关

目录: `helloagents/plan/202601231353_connection_logging/`

---

## 1. 配置与模型
- [√] 1.1 在 `internal/config/config.go` 中新增 connection_log.enable 配置并设置默认值，验证 why.md#需求-连接日志可配置-场景-开启连接日志
- [√] 1.2 在项目根目录更新 `proxy.yml` 示例配置，验证 why.md#需求-连接日志可配置-场景-开启连接日志

## 2. 核心日志逻辑
- [√] 2.1 在 `internal/core/session.go` 中新增 UDP 连接日志去重能力，验证 why.md#需求-udp-连接日志-场景-udp-首次接入
- [√] 2.2 在 `internal/core/relay.go` 中实现 UDP/TCP 连接日志输出，验证 why.md#需求-udp-连接日志-场景-udp-首次接入、why.md#需求-tcp-连接日志-场景-tcp-接入

## 3. 文档同步
- [√] 3.1 更新 `helloagents/wiki/modules/config.md` 说明 connection_log.enable
- [√] 3.2 更新 `helloagents/wiki/modules/core.md` 说明连接日志行为
- [√] 3.3 更新 `helloagents/wiki/modules/cli.md` 说明如何开启连接日志

## 4. 安全检查
- [√] 4.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 5. 测试
- [√] 5.1 执行 `go test ./...`，验证核心逻辑
