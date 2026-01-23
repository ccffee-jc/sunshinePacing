# Linux GUI 模块（gui-linux）

## 目的
提供 Linux 平台的 Fyne GUI，用于配置、启动/停止代理与查看状态。

## 模块概述
- **职责:** 配置编辑/保存、运行控制、日志/指标展示
- **状态:** ✅稳定
- **最后更新:** 2026-01-23

## 规范

### 需求: Linux GUI
**模块:** gui-linux
使用 Fyne 提供基础界面与运行管理。

#### 场景: 基础 GUI 操作
用户通过 GUI 修改配置并启动代理。
- 预期结果1：提供 base_port 与视频参数编辑
- 预期结果2：可查看运行状态与基础统计
- 预期结果3：启动时默认加载同目录 `sunshine-proxy.yml`，不存在则生成并填充表单
- 预期结果4：启动与保存以配置文件为基准，表单仅覆盖 base_port/internal_host/video 参数

## API接口
- 暂无对外 API

## 数据模型
- GUI 表单状态与运行时状态

## 依赖
- config
- core

## 变更历史
- [202601231757_linux_gui](../../history/2026-01/202601231757_linux_gui/) - Linux GUI 入口与构建支持
