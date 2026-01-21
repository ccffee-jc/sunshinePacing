# 任务清单: Sunshine 全端口转发支持

目录: `helloagents/plan/202601220034_ports_forwarding/`

---

## 1. 端口映射扩展
- [√] 1.1 在 `internal/core/ports.go` 中扩展端口表结构（支持 TCP/UDP 与多端口），验证 why.md#需求:-sunshine-全端口转发-场景:-默认全开
- [√] 1.2 补充端口映射单测覆盖新增端口，验证 why.md#需求:-sunshine-全端口转发-场景:-默认全开

## 2. TCP 转发实现
- [√] 2.1 在 `internal/core/proxy.go` 中增加 TCP 监听与生命周期管理，验证 why.md#需求:-sunshine-全端口转发-场景:-默认全开
- [√] 2.2 在 `internal/core/relay.go` 中实现 TCP 双向转发通道，验证 why.md#需求:-sunshine-全端口转发-场景:-默认全开

## 3. UDP 额外端口转发
- [√] 3.1 在 `internal/core/relay.go` 中实现非视频/控制/音频 UDP 端口纯转发，验证 why.md#需求:-sunshine-全端口转发-场景:-默认全开
- [√] 3.2 支持 UDP 48010 转发，验证 why.md#需求:-sunshine-全端口转发-场景:-兼容-udp-48010

## 4. 配置与文档
- [√] 4.1 更新 `internal/config/config.go` 中端口表/默认值说明
- [√] 4.2 更新知识库模块文档 `helloagents/wiki/modules/core.md`

## 5. 安全检查
- [√] 5.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 6. 测试
- [√] 6.1 在 `internal/core/ports_test.go` 中新增端口表测试: 默认全开映射正确
- [√] 6.2 在 `internal/core/relay_test.go` 中新增 TCP 转发基础测试（本地回环）
