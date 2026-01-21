# 任务清单: 内部转发目标主机可配置

目录: `helloagents/plan/202601211732_internal_host/`

---

## 1. 配置与默认值
- [√] 1.1 在 `internal/config/config.go` 中增加 `internal_host` 配置，默认值为 127.0.0.1

## 2. 转发目标
- [√] 2.1 在 `internal/core/proxy.go` 中使用 `internal_host` 作为 Sunshine 目标地址

## 3. GUI 支持
- [√] 3.1 在 `cmd/proxy-gui/main_windows.go` 中新增 `internal_host` 表单字段并保存

## 4. 文档同步
- [√] 4.1 更新 `helloagents/wiki/data.md`、`helloagents/wiki/modules/config.md`、`helloagents/wiki/modules/core.md`
- [√] 4.2 更新 `helloagents/wiki/arch.md` 与 `sunshine_pacing_proxy_plan.md` 的转发说明

## 5. 变更记录
- [√] 5.1 更新 `helloagents/CHANGELOG.md`

## 6. 测试
- [√] 6.1 运行 `go test ./...`
