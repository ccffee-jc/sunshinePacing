# 任务清单: 代理进程分离与本地 Metrics HTTP

目录: `helloagents/plan/202601261134_proxy_metrics_http/`

---

## 1. Metrics HTTP 服务
- [√] 1.1 在 `internal/core` 中新增 Metrics HTTP 服务与监听逻辑，验证 why.md#需求-进程外-metrics-服务-场景-gui-启动代理拉取数据
- [√] 1.2 在 `internal/core` 中补充 `/metrics` JSON 输出与超时设置，验证 why.md#需求-进程外-metrics-服务-场景-gui-启动代理拉取数据

## 2. CLI 启动与动态端口
- [√] 2.1 在 `cmd/proxy/main.go` 中新增 `--metrics-addr` 与 `--metrics-file` 参数，并将实际监听地址写入文件，验证 why.md#需求-进程外-metrics-服务-场景-gui-启动代理拉取数据
- [√] 2.2 为 CLI 增加 Metrics 读取/启动失败的错误提示，验证 why.md#需求-进程外-metrics-服务-场景-gui-启动代理拉取数据

## 3. GUI 进程管理与轮询
- [√] 3.1 在 `cmd/proxy-gui/main_linux.go` 中改为启动代理子进程并轮询 `/metrics`，验证 why.md#需求-进程外-metrics-服务-场景-gui-启动代理拉取数据
- [√] 3.2 在 `cmd/proxy-gui/main_windows.go` 中改为启动代理子进程并轮询 `/metrics`，验证 why.md#需求-进程外-metrics-服务-场景-gui-启动代理拉取数据
- [√] 3.3 在 GUI 停止逻辑中终止子进程并清理轮询，验证 why.md#需求-停止代理进程-场景-用户点击停止

## 4. 异常处理
- [√] 4.1 在 GUI 中处理 metrics 连接失败/超时并提示状态，验证 why.md#需求-metrics-异常处理-场景-metrics-服务不可达

## 5. 安全检查
- [√] 5.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 6. 文档更新
- [√] 6.1 更新 `helloagents/wiki/modules/gui-linux.md` 说明进程拆分与 metrics 拉取
- [√] 6.2 更新 `helloagents/wiki/modules/gui-win.md` 说明进程拆分与 metrics 拉取
- [√] 6.3 更新 `helloagents/wiki/arch.md` 增加进程拆分说明与 ADR 索引
- [√] 6.4 更新 `helloagents/CHANGELOG.md` 记录新增功能
- [√] 6.5 更新 `helloagents/history/index.md` 索引新方案包

## 7. 测试
- [√] 7.1 增加 Metrics HTTP handler 单元测试，验证 JSON 字段与状态码
- [√] 7.2 本地构建并验证 GUI 启动代理与指标刷新：`./scripts/build.sh`
