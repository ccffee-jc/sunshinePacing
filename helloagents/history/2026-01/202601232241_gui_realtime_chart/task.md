# 任务清单: GUI 实时突发图表

目录: `helloagents/plan/202601232241_gui_realtime_chart/`

---

## 1. 图表组件
- [√] 1.1 在 `cmd/proxy-gui` 中新增折线图组件与环形缓冲逻辑，验证 why.md#需求-实时突发图表-场景-运行中观测突发
- [√] 1.2 接入 500ms 定时采样，分别维护 video/control/audio 三条序列（video=队列长度，control/audio=Δbytes）

## 2. GUI 集成
- [√] 2.1 在 `cmd/proxy-gui/main_linux.go` 中嵌入图表区域并显示三色曲线
- [√] 2.2 在 `cmd/proxy-gui/main_windows.go` 中嵌入图表区域并显示三色曲线

## 3. 文档更新
- [√] 3.1 更新 `helloagents/wiki/modules/gui-linux.md` 说明实时图表与指标含义
- [√] 3.2 更新 `helloagents/wiki/modules/gui-win.md` 说明实时图表与指标含义
- [√] 3.3 更新 `helloagents/CHANGELOG.md` 记录新增功能

## 4. 安全检查
- [√] 4.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 5. 测试
- [√] 5.1 构建 GUI 并验证图表刷新：`go build -o dist/sunshine-proxy-gui ./cmd/proxy-gui`
