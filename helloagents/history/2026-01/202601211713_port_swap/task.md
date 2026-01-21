# 任务清单: 端口基准对调

目录: `helloagents/plan/202601211713_port_swap/`

---

## 1. 端口映射逻辑
- [√] 1.1 在 `internal/core/ports.go` 中调整端口映射：外部端口使用 base+offset，内部端口使用 base
- [√] 1.2 在 `internal/core/ports_test.go` 中更新端口映射测试用例

## 2. 文档同步
- [√] 2.1 更新 `helloagents/wiki/data.md` 中 internal_offset 含义
- [√] 2.2 更新 `helloagents/wiki/modules/config.md` 的说明
- [√] 2.3 更新 `helloagents/wiki/modules/core.md` 的端口转发描述
- [√] 2.4 更新 `helloagents/wiki/arch.md` 的端口流向说明

## 3. UI/提示调整
- [√] 3.1 在 `cmd/proxy-gui/main_windows.go` 中明确 base_port 为 Sunshine 基准

## 4. 变更记录
- [√] 4.1 更新 `helloagents/CHANGELOG.md` 版本记录

## 5. 测试
- [√] 5.1 运行 `go test ./...`
