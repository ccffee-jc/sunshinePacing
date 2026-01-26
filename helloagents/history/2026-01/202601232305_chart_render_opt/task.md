# 任务清单: 图表渲染性能优化

目录: `helloagents/plan/202601232305_chart_render_opt/`

---

## 1. 渲染优化
- [√] 1.1 在 `cmd/proxy-gui/chart.go` 中将图表渲染移到后台，主线程仅更新图像引用
- [√] 1.2 限制高频刷新时的重复渲染（合并脏数据）

## 2. 文档更新
- [√] 2.1 更新 `helloagents/wiki/modules/gui-linux.md` 说明图表渲染优化
- [√] 2.2 更新 `helloagents/wiki/modules/gui-win.md` 说明图表渲染优化
- [√] 2.3 更新 `helloagents/CHANGELOG.md` 记录性能优化

## 3. 测试
- [√] 3.1 构建 GUI 验证：`go build -o dist/sunshine-proxy-gui ./cmd/proxy-gui`
