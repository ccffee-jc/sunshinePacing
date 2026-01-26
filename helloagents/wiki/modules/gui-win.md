# Windows GUI 模块（gui-win）

## 目的
提供 Windows 平台的 Fyne GUI，用于配置、启动/停止代理与查看状态。

## 模块概述
- **职责:** 配置编辑/保存、运行控制、日志/指标展示
- **状态:** ✅稳定
- **最后更新:** 2026-01-21

## 规范

### 需求: Windows GUI
**模块:** gui-win
使用 Fyne 提供基础界面与运行管理。

#### 场景: 基础 GUI 操作
用户通过 GUI 修改配置并启动代理。
- 预期结果1：提供 base_port 与视频参数编辑
- 预期结果2：可查看运行状态与基础统计
- 预期结果3：启动时默认加载同目录 `sunshine-proxy.yml`，不存在则生成并填充表单
- 预期结果4：启动与保存以配置文件为基准，表单仅覆盖 base_port/internal_host/video 参数

### 需求: 实时突发图表
**模块:** gui-win
在 GUI 内展示实时突发曲线。

#### 场景: 运行中观测突发
代理运行中观察 video 曲线变化。
- 预期结果1：500ms 刷新曲线
- 预期结果2：图表仅显示 video 队列长度变化
- 预期结果3：图表渲染在后台完成，GUI 主线程仅更新图像引用

### 需求: 进程外指标拉取
**模块:** gui-win
GUI 启动代理子进程，通过本地 HTTP 拉取指标。

#### 场景: GUI 启动与停止代理
GUI 启动代理并展示指标，停止时终止子进程。
- 预期结果1：GUI 启动代理子进程并读取 metrics 地址
- 预期结果2：GUI 每 500ms 轮询 `/metrics`
- 预期结果3：停止时终止子进程并停止轮询

## API接口
- 暂无对外 API

## 数据模型
- GUI 表单状态与运行时状态

## 依赖
- config
- core

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - GUI 运行管理
- [202601212055_gui_default_config](../../history/2026-01/202601212055_gui_default_config/) - Windows GUI 默认配置与启动
