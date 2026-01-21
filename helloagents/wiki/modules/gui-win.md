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

## API接口
- 暂无对外 API

## 数据模型
- GUI 表单状态与运行时状态

## 依赖
- config
- core

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - GUI 运行管理
