# Sunshine UDP Pacing Proxy

> 本文件包含项目级别的核心信息。详细的模块文档见 `modules/` 目录。

---

## 1. 项目概述

### 目标与背景
为 Sunshine 提供本机 UDP 代理转发 + pacing，压平视频突发，降低运营商整形导致的卡顿；Windows 提供 GUI，Linux 使用配置启动。

### 范围
- **范围内:** UDP relay、端口映射、视频 pacing、单客户端会话、基础观测指标、Windows GUI、Linux CLI
- **范围外:** 多客户端支持、动态调参、DSCP/路由策略、透明抓包/重注入

### 干系人
- **负责人:** 需求方/维护者

---

## 2. 模块索引

| 模块名称 | 职责 | 状态 | 文档 |
|---------|------|------|------|
| core | UDP relay/session/pacer 核心逻辑 | ✅稳定 | [core](modules/core.md) |
| config | 配置加载与默认值 | ✅稳定 | [config](modules/config.md) |
| cli | Linux 命令行入口与运行管理 | ✅稳定 | [cli](modules/cli.md) |
| gui-win | Windows GUI（Fyne）与运行管理 | ✅稳定 | [gui-win](modules/gui-win.md) |

---

## 3. 快速链接
- [技术约定](../project.md)
- [架构设计](arch.md)
- [API 手册](api.md)
- [数据模型](data.md)
- [变更历史](../history/index.md)
