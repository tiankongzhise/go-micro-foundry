# 开发流程约束

本文定义 `go-micro-foundry` 的主线、feature、服务子分支、归档和提交约束。所有参与者在开发前必须阅读并遵守本文。

## 基本原则

- `main` 是稳定基线分支，只接收完成后的 feature。
- 禁止直接在 `main` 上开发、提交或推送。
- 项目采用单仓库目录模式，不为每个微服务维护长期开发分支。
- 微服务代码放在 `services/<service>` 下，公共代码按稳定性放入 `internal` 或 `pkg`。
- 功能开发必须小步推进，避免 feature 分支长期偏离 `main`。

## 分支模型

每个功能从 `main` 创建一个 feature 集成分支：

```bash
git switch main
git pull --ff-only origin main
git switch -c feature/<feature-name>
```

分支命名要求：

- feature 分支使用 `feature/<feature-name>`。
- 服务子分支使用 `feature/<feature-name>/<service-name>`。
- 归档分支使用 `archive/feature/<feature-name>`。
- `<feature-name>` 使用小写英文、数字和连字符，例如 `feature/bootstrap-config-center`。
- `<service-name>` 必须和服务目录名一致，例如 `auth-service`。

## 服务子分支

当一个 feature 涉及多个微服务，或需要多人并行开发时，必须从对应 feature 分支创建服务子分支：

```bash
git switch feature/<feature-name>
git switch -c feature/<feature-name>/<service-name>
```

服务子分支开发完成后，必须使用 merge 合并回对应 feature 分支：

```bash
git switch feature/<feature-name>
git merge --no-ff feature/<feature-name>/<service-name>
```

服务子分支合并规则：

- 必须使用 merge 合并回 feature 分支。
- 不允许使用 squash 合并服务子分支。
- 不允许通过 rebase 改写服务子分支历史后再合并。
- 合并提交必须保留，以便在归档分支中看到服务子分支的完整开发记录、提交顺序和合并关系。
- 服务子分支合并回 feature 后，可以删除远端活动分支；完整记录已经进入 feature 历史。

## feature 合入 main

feature 完成后必须满足以下条件：

- 相关服务子分支均已 merge 回 feature。
- 文档、配置示例和迁移说明已同步更新。
- 本地测试、构建或人工验证已完成并记录在 PR 或 commit 中。
- feature 分支已同步最新 `main`，并解决冲突。

feature 合入 `main` 时使用 squash：

- `main` 保持一个 feature 一个聚合提交。
- squash commit 使用中文说明功能目标、关键技术变化、影响范围和验证方式。
- feature 的详细开发过程不进入 `main` 历史，而是保留在归档分支。

## feature 归档

feature squash 合入 `main` 后，必须归档 feature 分支：

```bash
git push origin feature/<feature-name>:archive/feature/<feature-name>
git push origin --delete feature/<feature-name>
```

归档约束：

- `archive/feature/<feature-name>` 只用于追溯历史。
- 归档分支不可更改。
- 禁止向归档分支继续提交。
- 禁止 force push 归档分支。
- 禁止 rebase、reset、cherry-pick 后覆盖归档分支。
- 禁止把归档分支作为新开发基线。
- 如需修复问题，必须从当前 `main` 新建新的 `feature/*` 分支。
- 仓库管理员应在 GitHub 上为 `archive/feature/*` 配置保护规则，禁止删除和直接推送。

可选为归档点打 tag：

```bash
git tag archive/<feature-name>/merged-YYYYMMDD archive/feature/<feature-name>
git push origin archive/<feature-name>/merged-YYYYMMDD
```

## commit 约束

每实现一个功能点必须提交一次 commit。commit 应保持原子性，一个 commit 对应一个可解释、可验证的功能、修复或文档变更。

commit 信息必须使用中文，并尽量包含：

- 变更内容：做了什么。
- 技术细节：采用了什么实现方式，为什么这样处理。
- 影响范围：影响哪些服务、接口、配置、数据或文档。
- 验证方式：执行了哪些测试、构建或人工验证。

推荐格式：

```text
类型: 简要说明

变更内容：
- ...

技术细节：
- ...

影响范围：
- ...

验证方式：
- ...
```

示例：

```text
feat: 增加配置中心配置项保存能力

变更内容：
- 新增配置项保存接口和本地存储实现。

技术细节：
- 使用标准库 http.Handler 暴露 POST 接口。
- 存储层先定义接口，再提供文件存储实现，便于后续替换 PostgreSQL。

影响范围：
- 影响 services/config-center。
- 新增配置文件路径配置项。

验证方式：
- 执行 go test ./services/config-center/...。
- 手动验证重复保存同一配置项会覆盖旧版本。
```

## 大型 feature 控制

大型 feature 必须拆分为可独立验证的小 feature。确需提前合入但暂不启用的能力，应使用 feature flag 或配置开关控制默认行为。

长期 feature 必须定期从 `main` 同步：

```bash
git fetch origin
git switch feature/<feature-name>
git merge --no-ff origin/main
```

同步 `main` 时必须保留 merge 记录，方便归档分支还原集成过程。

## GitHub 保护规则

`main` 必须启用分支保护：

- 要求通过 Pull Request 合入。
- 要求至少一次 review。
- 禁止直接 push 到 `main`。
- 有 CI 后启用必需状态检查。

`archive/feature/*` 必须启用保护规则：

- 禁止 force push。
- 禁止删除。
- 禁止直接 push。
- 仅仓库管理员可调整保护规则。

如果当前 GitHub 权限或规则能力不足，必须至少在仓库管理记录中登记待办，并在具备权限后补齐保护规则。
