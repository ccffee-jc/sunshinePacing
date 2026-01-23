# 技术设计: 本机47989端口族直转到内网Sunshine

## 技术方案
### 核心技术
- Go 标准库 net（TCP/UDP 监听与转发）
- YAML 配置加载（gopkg.in/yaml.v3）

### 实现要点
- 放宽 internal_offset 的校验范围，允许 0 以实现外部端口=base_port。
- 生成 proxy.yml，设置 base_port=47989、internal_offset=0、internal_host=192.168.2.110。
- Linux 编译 CLI 到项目根目录，作为与配置同目录的可执行文件。

## 架构决策 ADR
### ADR-002: 允许 internal_offset=0
**上下文:** 用户需要外部端口与 Sunshine 原始端口一致，当前校验要求 internal_offset>=1。
**决策:** 将 internal_offset 最小值放宽到 0。
**理由:** 保持配置简洁，满足端口直通需求。
**替代方案:** 保持 internal_offset>=1 → 拒绝原因: 无法实现外部端口与原始端口一致。
**影响:** 可能与本机已占用端口冲突，需要用户自行确认。

## 安全与性能
- **安全:** 不记录敏感数据；仅按配置转发指定目标。
- **性能:** 端口映射逻辑无新增开销。

## 测试与部署
- **测试:** 执行 go test ./... 验证端口映射与转发逻辑。
- **部署:** 使用 proxy.yml 启动 CLI。 
