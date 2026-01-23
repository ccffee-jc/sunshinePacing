# 命令行模块（cli）

## 目的
提供 Linux/通用 CLI 入口，基于配置启动代理。

## 模块概述
- **职责:** 解析参数、加载配置、启动/停止服务、输出日志与状态
- **状态:** ✅稳定
- **最后更新:** 2026-01-23

## 规范

### 需求: Linux 配置启动
**模块:** cli
支持 `-config` 参数加载 YAML 并启动代理。

#### 场景: 无 GUI 启动
在 Linux 环境使用 CLI 启动代理服务。
- 预期结果1：无 GUI 依赖
- 预期结果2：可输出运行日志

#### 场景: 本地编译产物与配置同目录
在项目根目录构建 Linux CLI，并与 proxy.yml 放在同目录。
- 预期结果1：产物位于项目根目录 `sunshine-proxy`
- 预期结果2：可直接使用 `./sunshine-proxy -config proxy.yml` 启动

#### 场景: 开启连接日志启动
在 proxy.yml 中设置 connection_log.enable=true 后启动。
- 预期结果1：UDP/TCP 有连接时输出日志
- 预期结果2：日志包含端口、协议、客户端与内部目标

## API接口
- 暂无对外 API

## 数据模型
- CLI 参数结构体

## 依赖
- config
- core

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - CLI 启动与日志
