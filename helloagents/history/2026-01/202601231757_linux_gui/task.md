# 任务清单: Linux GUI 构建支持

目录: `helloagents/plan/202601231757_linux_gui/`

---

## 1. GUI 入口
- [√] 1.1 在 `cmd/proxy-gui/main_linux.go` 中新增 Linux GUI 入口，复用 Fyne 界面与运行控制逻辑

## 2. 构建脚本
- [√] 2.1 在 `scripts/build.sh` 中增加 Linux GUI 构建产物 `sunshine-proxy-gui`

## 3. 文档更新
- [√] 3.1 更新 `helloagents/wiki/overview.md`，补充 Linux GUI 范围与模块索引
- [√] 3.2 新增 `helloagents/wiki/modules/gui-linux.md`，记录 Linux GUI 模块说明
- [√] 3.3 更新 `helloagents/CHANGELOG.md`，记录新增 Linux GUI 支持

## 4. 测试
- [√] 4.1 构建 Linux GUI 产物：`go build -o dist/sunshine-proxy-gui ./cmd/proxy-gui`
> 备注: 已在升级 Go 与安装 X11/GL 依赖后完成构建。
