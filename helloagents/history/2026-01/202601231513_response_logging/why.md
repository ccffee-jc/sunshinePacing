# 变更提案: 回包日志与配置开关

## 需求背景
当前日志只覆盖连接与错误，用户希望对 Sunshine→客户端（内部→外部）的所有回包进行逐包日志，以便排障与分析。

## 变更内容
1. 新增回包日志配置开关，默认关闭。
2. 内部→外部方向的每个 UDP 回包都记录日志。
3. 日志内容尽可能详细（端口、协议、流类型、来源/目标、包大小等）。

## 影响范围
- **模块:** config / core / cli
- **文件:** internal/config/config.go, internal/core/relay.go, proxy.yml, helloagents/wiki/modules/config.md, helloagents/wiki/modules/core.md, helloagents/wiki/modules/cli.md
- **API:** 无
- **数据:** 无

## 核心场景

### 需求: 回包日志开关
**模块:** config
新增 response_log.enable 控制回包日志。

#### 场景: 开启回包日志
用户设置 response_log.enable=true。
- 预期结果1：内部→外部回包逐包输出日志
- 预期结果2：日志覆盖所有流

### 需求: 回包逐包记录
**模块:** core
每个内部回包都输出日志。

#### 场景: 内部→外部转发
Sunshine 回包到代理后转发给客户端。
- 预期结果1：每个包都记录
- 预期结果2：日志包含端口、流类型、包大小、源/目标地址

## 风险评估
- **风险:** 日志量极大影响性能与磁盘。
- **缓解:** 通过配置开关默认关闭，并提示仅排障时开启。
