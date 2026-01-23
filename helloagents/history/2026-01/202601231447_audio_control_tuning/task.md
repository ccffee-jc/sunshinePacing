# 任务清单: 配置调参与音频轻量限制

目录: `helloagents/plan/202601231447_audio_control_tuning/`

---

> **状态:** 未执行（用户清理）

## 1. 配置与 pacing 模型
- [-] 1.1 在 `internal/config/config.go` 中扩展 audio 配置字段并设置默认值，验证 why.md#需求-视频参数更新-场景-新视频参数启用
- [-] 1.2 在 `internal/core/pacer.go` 中支持 bps 速率与统计回调，验证 why.md#需求-audio-轻量限制-场景-audio-连接

## 2. 代理逻辑调整
- [-] 2.1 在 `internal/core/proxy.go` 中初始化 audio pacer，控制 mode=limit/light 才启用
- [-] 2.2 在 `internal/core/relay.go` 中接入 audio pacer，control 保持直通

## 3. 配置与文档更新
- [-] 3.1 更新 `proxy.yml` 应用新视频参数与 control/audio 配置
- [-] 3.2 更新 `helloagents/wiki/modules/config.md` 说明 audio 轻量限制字段
- [-] 3.3 更新 `helloagents/wiki/modules/core.md` 说明 control bypass 与 audio 轻量限制
- [-] 3.4 更新 `helloagents/wiki/modules/cli.md` 说明配置启用方式
- [-] 3.5 更新 `helloagents/wiki/arch.md` 增加 ADR-004
- [-] 3.6 更新 `helloagents/CHANGELOG.md` 记录变更

## 4. 安全检查
- [-] 4.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 5. 测试
- [-] 5.1 执行 `go test ./...`，验证核心逻辑
