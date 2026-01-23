# 任务清单: GUI 启动配置一致性修复

目录: `helloagents/plan/202601232204_gui_config_fix/`

---

## 1. GUI 配置读取
- [√] 1.1 在 `cmd/proxy-gui/main_linux.go` 中启动/保存时以配置文件为基准，避免丢失字段
- [√] 1.2 在 `cmd/proxy-gui/main_windows.go` 中启动/保存时以配置文件为基准，避免丢失字段

## 2. 文档更新
- [√] 2.1 更新 `helloagents/wiki/modules/gui-linux.md` 说明 GUI 启动的配置合并逻辑
- [√] 2.2 更新 `helloagents/wiki/modules/gui-win.md` 说明 GUI 启动的配置合并逻辑
- [√] 2.3 更新 `helloagents/CHANGELOG.md` 记录修复
