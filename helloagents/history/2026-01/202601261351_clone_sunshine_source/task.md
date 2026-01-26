# 任务清单: 引入 Sunshine 官方源码镜像用于检索

目录: `helloagents/plan/202601261351_clone_sunshine_source/`

---

## 1. 仓库镜像
- [√] 1.1 在 `third_party/sunshine` 克隆 Sunshine 官方仓库并递归拉取子模块，验证 why.md#需求-本地检索-sunshine-源码-场景-克隆官方仓库与依赖并忽略跟踪
- [√] 1.2 在 `.gitignore` 添加 `third_party/sunshine/` 忽略项，验证 why.md#需求-本地检索-sunshine-源码-场景-克隆官方仓库与依赖并忽略跟踪

## 2. 文档更新
- [√] 2.1 更新 `helloagents/wiki/overview.md` 记录本地源码镜像用途与位置

## 3. 安全检查
- [√] 3.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 4. 测试
- [√] 4.1 验证 `git -C third_party/sunshine submodule status` 可用，且 `git status -s` 无未忽略文件
