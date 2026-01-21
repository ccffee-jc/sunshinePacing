# 配置模块（config）

## 目的
提供 YAML 配置加载、默认值填充与校验。

## 模块概述
- **职责:** 读取配置文件、默认值合并、字段校验、导出运行配置
- **状态:** ✅稳定
- **最后更新:** 2026-01-21

## 规范

### 需求: 单一 base_port 配置
**模块:** config
仅要求用户设置 Sunshine 的 base_port，其它参数提供合理默认。

#### 场景: 默认配置启动
用户只填写 Sunshine 的 base_port 即可运行代理。
- 预期结果1：internal_offset 自动为 1000，外部端口为 base_port+offset
- 预期结果2：video/control/audio 具备默认策略
- 预期结果3：internal_host 默认 127.0.0.1，可指向局域网 Sunshine

## API接口
- 暂无对外 API

## 数据模型
- 配置结构体（base_port/internal_offset/video/control/audio）

## 依赖
- core

## 变更历史
- [202601211643_sunshine_pacing_proxy](../../history/2026-01/202601211643_sunshine_pacing_proxy/) - 配置加载与默认值
