# 任务清单: 图表仅显示视频队列

目录: `helloagents/plan/202601232311_video_only_chart/`

---

## 1. 图表与采样
- [√] 1.1 在 `cmd/proxy-gui/chart.go` 中移除 control/audio 曲线，仅保留 video 队列
- [√] 1.2 在 `cmd/proxy-gui/main_linux.go` 与 `cmd/proxy-gui/main_windows.go` 中移除 control/audio 采样与图例

## 2. 文档更新
- [√] 2.1 更新 `helloagents/wiki/modules/gui-linux.md` 说明图表仅显示 video 队列
- [√] 2.2 更新 `helloagents/wiki/modules/gui-win.md` 说明图表仅显示 video 队列
- [√] 2.3 更新 `helloagents/CHANGELOG.md` 记录变更

## 3. 测试
- [√] 3.1 构建 GUI 验证：`go build -o dist/sunshine-proxy-gui ./cmd/proxy-gui`
