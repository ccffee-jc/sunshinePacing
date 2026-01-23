# 变更提案: 连接日志与配置开关

## 需求背景
当前代理只输出错误或状态日志，缺少连接建立的可观测性。需要在 UDP/TCP 有连接时记录详细日志，并提供配置开关以便按需开启。

## 变更内容
1. 新增连接日志配置开关，默认关闭。
2. UDP：首次看到客户端时记录连接日志。
3. TCP：接受连接时记录连接日志。
4. 日志内容尽量详细（端口、流类型、客户端、内部目标等）。

## 影响范围
- **模块:** config / core / cli
- **文件:** internal/config/config.go, internal/core/session.go, internal/core/relay.go, proxy.yml, helloagents/wiki/modules/config.md, helloagents/wiki/modules/core.md, helloagents/wiki/modules/cli.md
- **API:** 无
- **数据:** 无

## 核心场景

### 需求: 连接日志可配置
**模块:** config
新增日志开关用于控制连接日志输出。

#### 场景: 开启连接日志
用户设置 connection_log.enable=true。
- 预期结果1：UDP/TCP 有连接时输出日志
- 预期结果2：日志包含端口、流类型、客户端与内部目标

### 需求: UDP 连接日志
**模块:** core
首次看到客户端即记录一次连接日志。

#### 场景: UDP 首次接入
收到某端口客户端首包。
- 预期结果1：仅首次记录该客户端
- 预期结果2：记录端口、流类型、客户端地址

### 需求: TCP 连接日志
**模块:** core
accept 连接时记录一次连接日志。

#### 场景: TCP 接入
监听端口 accept 新连接。
- 预期结果1：记录一次连接日志
- 预期结果2：记录端口、客户端、目标地址

## 风险评估
- **风险:** 日志过多影响可读性与性能。
- **缓解:** 提供配置开关，默认关闭，仅在需要时开启。
