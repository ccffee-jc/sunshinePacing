# Changelog

本文件记录项目所有重要变更。
格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/),
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 新增
- 支持 Sunshine 全端口 TCP/UDP 转发（包含 UDP 48010 兼容）。
- 支持连接日志开关（UDP/TCP 连接输出）。

### 变更
- Windows GUI 启动时自动生成/加载同目录 sunshine-proxy.yml。
- Windows 构建产物命名调整：GUI 为 sunshine-proxy.exe，CLI 为 sunshine-proxy-cli.exe。
- 允许 internal_offset=0 以保持外部端口与 base_port 一致，并提供示例配置。

## [0.2.3] - 2026-01-21

### 新增
- 一键构建脚本 scripts/build.sh，输出到 dist/ 并支持 Windows 目标。

### 变更
- .gitignore 忽略 dist/ 编译目录。

### 修复
- 无。

### 移除
- 无。

## [0.2.2] - 2026-01-21

### 新增
- 支持配置 internal_host，将转发目标指向局域网 Sunshine。

### 变更
- 无。

### 修复
- 无。

### 移除
- 无。

## [0.2.1] - 2026-01-21

### 新增
- 无。

### 变更
- 对调外部/内部端口基准：外部使用 base+offset，Sunshine 保持 base_port。

### 修复
- 无。

### 移除
- 无。

## [0.2.0] - 2026-01-21

### 新增
- UDP 代理转发与视频 pacing 核心实现。
- Windows Fyne GUI 与 Linux CLI 入口。
- YAML 配置加载与默认值校验。
- 端口映射与基础统计指标。

### 变更
- 无。

### 修复
- 无。

### 移除
- 无。

## [0.1.0] - 2026-01-21

### 新增
- 初始化知识库与方案设计前置文档。

### 变更
- 无。

### 修复
- 无。

### 移除
- 无。
