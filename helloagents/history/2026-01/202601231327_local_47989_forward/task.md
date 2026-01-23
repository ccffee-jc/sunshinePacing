# 任务清单: 本机47989端口族直转到内网Sunshine

目录: `helloagents/plan/202601231327_local_47989_forward/`

---

## 1. 配置与端口校验
- [√] 1.1 在 `internal/config/config.go` 中允许 internal_offset=0，验证 why.md#需求-端口族直通映射-场景-internal_offset--0
- [√] 1.2 在项目根目录新增 `proxy.yml` 示例配置，验证 why.md#需求-端口族直通映射-场景-internal_offset--0

## 2. 运行与编译说明
- [√] 2.1 在 `helloagents/wiki/modules/config.md` 中补充 internal_offset=0 的说明与示例
- [√] 2.2 在 `helloagents/wiki/modules/cli.md` 中补充 Linux 编译产物位置与启动说明

## 3. 安全检查
- [√] 3.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 4. 测试
- [√] 4.1 执行 `go test ./...`，验证端口映射与核心逻辑稳定
