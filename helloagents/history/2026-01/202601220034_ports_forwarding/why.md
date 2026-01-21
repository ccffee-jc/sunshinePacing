# 变更提案: Sunshine 全端口转发支持

## 需求背景
当前代理仅覆盖 Sunshine 的视频/控制/音频 3 个 UDP 端口，无法处理其他 TCP/UDP 端口（例如 RTSP、配对、发现、Web UI、UPnP 等）。Sunshine 的端口由 base_port 计算偏移，用户无法单独配置每个端口，因此代理需要对所有官方端口偏移进行转发，以保证功能完整与兼容性。

## 变更内容
1. 扩展端口映射，覆盖 Sunshine 官方文档的所有端口偏移，并补充 UDP 48010 兼容。
2. 新增 TCP 端口转发能力，按端口偏移与内部偏移规则双向转发。
3. 默认全开端口转发，保持现有 video/control/audio UDP pacing 逻辑不变，其他端口仅做纯转发。

## 影响范围
- **模块:** config, core, cli/gui
- **文件:** internal/config/config.go, internal/core/ports.go, internal/core/proxy.go, internal/core/relay.go, cmd/proxy*, cmd/proxy-gui*
- **API:** 配置结构体新增端口集合（内部使用）
- **数据:** 无

## 核心场景

### 需求: Sunshine 全端口转发
**模块:** core/config
扩展端口映射并实现 TCP/UDP 双栈转发，确保所有端口都可通过 base_port 与 internal_offset 推导。

#### 场景: 默认全开
代理启动后自动监听全部外部端口，并将流量转发到内部 Sunshine 对应端口。
- 预期结果: 既有 47998/47999/48000 行为不变
- 预期结果: 其他端口默认转发成功

#### 场景: 兼容 UDP 48010
除 RTSP TCP 48010 外，增加 UDP 48010 的转发以兼容 Moonlight 文档。
- 预期结果: TCP/UDP 48010 均可转发

## 风险评估
- **风险:** 新增 TCP 转发引入连接与资源管理复杂度
- **缓解:** 统一在代理生命周期内管理 listener/conn，限制并发与超时，失败日志清晰可追踪
