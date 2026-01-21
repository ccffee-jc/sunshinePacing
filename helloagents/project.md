# 项目技术约定

---

## 技术栈
- **核心:** Go >= 1.18
- **GUI(Windows):** Fyne
- **配置:** YAML

---

## 开发约定
- **代码规范:** gofmt + go vet
- **命名约定:** Go 规范（CamelCase/驼峰）

---

## 错误与日志
- **策略:** 显式错误返回，关键路径增加上下文
- **日志:** 使用标准库 log 输出（含时间戳）

---

## 测试与流程
- **测试:** go test（单元为主，核心转发与pacer逻辑需覆盖）
- **提交:** Conventional Commits（可选）

---

## 构建与发布
- **脚本:** scripts/build.sh（默认输出到项目根目录 dist/，可用 OUT_DIR 覆盖）
- **跨平台:** 默认构建 Linux CLI（sunshine-proxy）与 Windows CLI（sunshine-proxy-cli.exe）；Windows GUI 输出 sunshine-proxy.exe，在 Windows 环境自动构建，或设置 BUILD_WINDOWS_GUI=1 强制交叉编译
