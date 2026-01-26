# 技术设计: 引入 Sunshine 官方源码镜像用于检索

## 技术方案
### 核心技术
- Git 克隆 + 子模块递归拉取

### 实现要点
- 目标路径固定为 `third_party/sunshine`，便于统一检索路径。
- 使用 `git clone --recurse-submodules` 拉取官方仓库及其子模块。
- 将 `third_party/sunshine/` 加入主项目 `.gitignore`，防止误提交。
- 在知识库 `wiki/overview.md` 增补本地源码镜像说明。

## 安全与性能
- **安全:** 仅用于源码检索，不执行或修改第三方代码。
- **性能:** 增加磁盘占用；不参与构建流程。

## 测试与部署
- **测试:** 确认 `third_party/sunshine` 存在且子模块拉取完成；`git status -s` 不出现未忽略文件。
- **部署:** 无需部署。
