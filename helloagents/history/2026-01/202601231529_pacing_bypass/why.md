# 变更提案: 纯转发模式开关

## 需求背景
用户需要在排障时禁用所有 pacing，验证纯转发链路是否正常。

## 变更内容
1. 新增 pacing.enable 配置开关，默认开启。
2. 当 pacing.enable=false 时，video/control/audio 全部直通。
3. 更新示例配置与文档说明。

## 影响范围
- **模块:** config / core / cli
- **文件:** internal/config/config.go, internal/core/proxy.go, internal/core/relay.go, proxy.yml, helloagents/wiki/modules/config.md, helloagents/wiki/modules/core.md, helloagents/wiki/modules/cli.md
- **API:** 无
- **数据:** 无

## 核心场景

### 需求: 纯转发模式开关
**模块:** config
新增 pacing.enable。

#### 场景: 禁用 pacing
用户设置 pacing.enable=false。
- 预期结果1：video/control/audio 全部直通
- 预期结果2：不启用任何 pacing 逻辑

### 需求: 默认保持 pacing
**模块:** core
默认仍启用 video pacing。

#### 场景: 未配置开关
配置缺失时按默认启用 pacing。
- 预期结果1：行为与当前版本一致
- 预期结果2：仅在显式关闭时才纯转发

## 风险评估
- **风险:** 纯转发可能暴露突发导致卡顿。
- **缓解:** 默认保持启用，仅排障时关闭。
