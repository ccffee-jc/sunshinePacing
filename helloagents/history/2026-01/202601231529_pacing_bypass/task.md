# 任务清单: 纯转发模式开关

目录: `helloagents/plan/202601231529_pacing_bypass/`

---

## 1. 配置与逻辑
- [√] 1.1 在 `internal/config/config.go` 中新增 pacing.enable 配置并设置默认值，验证 why.md#需求-纯转发模式开关-场景-禁用-pacing
- [√] 1.2 在 `internal/core/proxy.go` 中接入 pacing.enable，关闭时全部直通

## 2. 配置与文档更新
- [√] 2.1 更新 `proxy.yml` 示例配置
- [√] 2.2 更新 `helloagents/wiki/modules/config.md` 说明 pacing.enable
- [√] 2.3 更新 `helloagents/wiki/modules/core.md` 说明纯转发模式
- [√] 2.4 更新 `helloagents/wiki/modules/cli.md` 说明启用方式
- [√] 2.5 更新 `helloagents/wiki/arch.md` 增加 ADR-006
- [√] 2.6 更新 `helloagents/CHANGELOG.md` 记录变更

## 3. 安全检查
- [√] 3.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 4. 测试
- [√] 4.1 执行 `go test ./...`，验证核心逻辑
