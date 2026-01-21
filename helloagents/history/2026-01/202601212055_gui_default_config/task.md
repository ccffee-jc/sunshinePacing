# 任务清单: Windows GUI 默认配置与启动

目录: `helloagents/plan/202601212055_gui_default_config/`

---

## 1. GUI 启动行为
- [√] 1.1 在 `cmd/proxy-gui/main_windows.go` 中启动时自动定位可执行文件目录，默认使用 `sunshine-proxy.yml`，不存在则生成并加载。
- [√] 1.2 在 `cmd/proxy-gui/main_windows.go` 中加载默认配置到表单并提示状态。

## 2. 构建与输出
- [√] 2.1 在 `scripts/build.sh` 中将 Windows GUI 输出命名为 `sunshine-proxy.exe`，保留 Windows CLI 为 `sunshine-proxy-cli.exe`。

## 3. 文档更新
- [√] 3.1 更新 `helloagents/wiki/modules/gui-win.md` 说明默认配置文件逻辑。
- [√] 3.2 更新 `helloagents/project.md` 描述 Windows 构建产物命名。
- [√] 3.3 更新 `helloagents/CHANGELOG.md` 记录本次变更。

## 4. 安全检查
- [√] 4.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）
