# 配置模块（config）

## 目的
提供 YAML 配置加载、默认值填充与校验。

## 模块概述
- **职责:** 读取配置文件、默认值合并、字段校验、导出运行配置
- **状态:** ✅稳定
- **最后更新:** 2026-01-23

## 规范

### 需求: 单一 base_port 配置
**模块:** config
仅要求用户设置 Sunshine 的 base_port，其它参数提供合理默认。

#### 场景: 默认配置启动
用户只填写 Sunshine 的 base_port 即可运行代理。
- 预期结果1：internal_offset 自动为 1000，外部端口为 base_port+offset
- 预期结果2：video/control/audio 具备默认策略
- 预期结果3：internal_host 默认 127.0.0.1，可指向局域网 Sunshine

#### 场景: 端口族直通映射
用户显式设置 internal_offset=0 并指定 internal_host，将外部端口与 Sunshine 端口保持一致。
- 预期结果1：外部端口与 Sunshine 端口号一致
- 预期结果2：转发目标指向指定的局域网主机

### 需求: 连接日志开关
**模块:** config
提供 connection_log.enable 用于控制连接日志输出。

#### 场景: 开启连接日志
用户设置 connection_log.enable=true。
- 预期结果1：UDP/TCP 有连接时输出日志
- 预期结果2：日志包含端口、流类型、客户端与内部目标

## API接口
- 暂无对外 API

## 数据模型
- 配置结构体（base_port/internal_offset/video/control/audio/connection_log）

## 依赖
- core

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - 配置加载与默认值
