# 开发流程约束

本文定义 `go-micro-foundry` 的主线、feature、服务子分支、归档和提交约束。所有参与者在开发前必须阅读并遵守本文。

## 基本原则

- `main` 是稳定基线分支，只接收完成后的 feature。
- 禁止直接在 `main` 上开发、提交或推送。
- 项目采用单仓库目录模式，不为每个微服务维护长期开发分支。
- 微服务代码放在 `services/<service>` 下，公共代码按稳定性放入 `internal` 或 `pkg`。
- 功能开发必须小步推进，避免 feature 分支长期偏离 `main`。
- 涉及微服务实现、接口、配置、UI、部署或联调的开发，必须先完整阅读历史设计文档归档，再进入预开发规划或编码。

## 项目上下文基线

历史完整设计文档已经归档到专用归档仓库：

```text
git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/feature/origin-docs@74e5076bc9373b5502cf9bf1152b6b0be7aa3781
```

开发人员在开始涉及微服务能力的 feature 前，必须完整阅读该归档中的以下文档，并以这些文档作为项目上下文基线：

1. `docs/architecture.md`
2. `docs/api-conventions.md`
3. `docs/ui-conventions.md`
4. `docs/config-conventions.md`
5. `docs/config-center.md`
6. `docs/service-registry.md`
7. `docs/api-gateway.md`
8. `docs/auth-service.md`
9. `docs/tracing-service.md`
10. `docs/deployment.md`
11. `docs/roadmap.md`

阅读要求：

- 预开发规划文档必须写明所依据的 `origin-docs` 归档分支和归档 HEAD。
- 如果当前 feature 需要调整历史设计，必须在规划文档中明确说明调整点、原因和影响范围。
- 禁止只根据当前仓库 README 或服务占位目录直接编码；当前仓库可能只保留基线代码，完整项目设计以归档文档为准。
- 本地可以临时添加 `archive` 远端读取归档文档，但读取或归档完成后必须移除该远端，不得长期保留与归档仓库的本地关联。

临时读取归档文档示例：

```bash
git remote add archive git@github.com:tiankongzhise/go-micro-foundry-archive.git
git fetch archive archive/feature/origin-docs
git show FETCH_HEAD:docs/architecture.md
git remote remove archive
```

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
- 归档分支使用 `archive/feature/<feature-name>`，最终只保存在专用归档仓库 `git@github.com:tiankongzhise/go-micro-foundry-archive.git` 中。
- `<feature-name>` 使用小写英文、数字和连字符，例如 `feature/bootstrap-config-center`。
- `<service-name>` 必须和服务目录名一致，例如 `auth-service`。

## 预开发规划阶段

当一个 feature 涉及多个微服务、公共基础包或多人并行开发时，必须先在 feature 集成分支完成预开发规划，再创建服务子分支。

预开发规划必须至少明确：

- 已完整阅读的 `origin-docs` 归档版本和参考文档清单。
- feature 的目标、非目标、影响范围和验收方式。
- 各微服务在基础契约、核心能力、联调准备和联调验收阶段需要达到的进度。
- 服务子分支拆分方式、负责人或协作团队、依赖顺序和合并顺序。
- API 契约、配置项、错误码、日志字段、请求 ID 和 Trace ID 的统一约定。
- 联调用例、测试数据、运行说明和未完成事项登记方式。

预开发规划完成后，必须在 feature 分支提交规划 commit。服务子分支必须从该规划 commit 之后的 feature HEAD 创建，避免各团队基于不一致的目标开始开发。

微服务阶段规划模板见 [microservice-development-plan.md](microservice-development-plan.md)。

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
- 联调已通过，或该 feature 只有文档/流程变更且已经完成对应人工验收。
- 已按归档流程迁出归档分支、推送专用归档仓库，并在 [archive-map.md](archive-map.md) 中记录开发分支与归档分支的对应关系。

feature 合入 `main` 时使用 squash：

- `main` 保持一个 feature 一个聚合提交。
- squash commit 使用中文说明功能目标、关键技术变化、影响范围和验证方式。
- feature 的详细开发过程不进入 `main` 历史，而是保留在专用归档仓库的归档分支。
- squash commit 必须在提交信息末尾的 trailer 区域包含且只包含一个 `Archive-Ref` 字段，格式如下：

```text
Archive-Ref: git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/feature/<feature-name>@<40位HEAD_SHA>
```

CI 可使用以下正则检查 `Archive-Ref`：

```text
^Archive-Ref: git@github\.com:tiankongzhise/go-micro-foundry-archive\.git archive/feature/[a-z0-9-]+@[0-9a-f]{40}$
```

## feature 归档

feature squash 合入 `main` 前，必须先把已通过联调或人工验收的 feature 迁出为临时归档分支，再推送到专用归档仓库。归档完成后，必须回到 feature 分支更新归档映射文档并提交归档记录 commit。

```bash
git switch feature/<feature-name>
git switch -c archive/feature/<feature-name>
git remote add archive git@github.com:tiankongzhise/go-micro-foundry-archive.git
git push archive archive/feature/<feature-name>:archive/feature/<feature-name>
git ls-remote --heads archive archive/feature/<feature-name>
```

根据 `git ls-remote` 返回的 40 位 SHA 生成合入 `main` 需要的归档字段：

```text
Archive-Ref: git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/feature/<feature-name>@<40位HEAD_SHA>
```

然后回到 feature 分支，更新 [archive-map.md](archive-map.md)，说明开发分支、归档分支、归档 HEAD、归档字段和必要的归档过程。该更新必须单独提交一次 commit，commit 信息应详述归档过程并包含归档字段。

```bash
git switch feature/<feature-name>
git add docs/archive-map.md
git commit
git branch -d archive/feature/<feature-name>
git remote remove archive
```

归档约束：

- 本地 `archive/feature/<feature-name>` 只是迁出归档时使用的临时分支，归档和记录完成后必须删除。
- 本地仓库不得长期保留 `archive` 远端；归档完成后必须执行 `git remote remove archive`，解除本地与专用归档仓库的关联。
- `archive/feature/<feature-name>` 只用于追溯历史，最终必须位于专用归档仓库。
- 主仓库 `git@github.com:tiankongzhise/go-micro-foundry.git` 不保留 `archive/feature/*` 分支。
- 归档分支不可更改。
- 禁止向归档分支继续提交。
- 禁止 force push 归档分支。
- 禁止 rebase、reset、cherry-pick 后覆盖归档分支。
- 禁止把归档分支作为新开发基线。
- 如需修复问题，必须从当前 `main` 新建新的 `feature/*` 分支。
- 归档仓库管理员应在 GitHub 上为 `archive/feature/*` 配置保护规则，禁止删除和直接推送。

可选为归档点打 tag：

```bash
git tag archive/<feature-name>/merged-YYYYMMDD feature/<feature-name>
git push archive archive/<feature-name>/merged-YYYYMMDD
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

同步 `main` 时必须保留 merge 记录，方便归档仓库中的归档分支还原集成过程。

## GitHub 保护规则

`main` 必须启用分支保护：

- 要求通过 Pull Request 合入。
- 要求至少一次 review。
- 禁止直接 push 到 `main`。
- 有 CI 后启用必需状态检查。

专用归档仓库的 `archive/feature/*` 必须启用保护规则：

- 禁止 force push。
- 禁止删除。
- 禁止直接 push。
- 仅仓库管理员可调整保护规则。

如果当前 GitHub 权限或规则能力不足，必须至少在仓库管理记录中登记待办，并在具备权限后补齐保护规则。

## 仓库管理待办

当前基线文档已经写入 `main`。仓库管理员需要在 GitHub 仓库设置中补齐以下保护规则：

- 为 `main` 配置 branch protection rule。
- 为 `main` 启用 Require a pull request before merging。
- 为 `main` 启用 Require approvals，并要求至少 1 次 review。
- 禁止直接 push 到 `main`。
- CI 建立后，为 `main` 启用 Require status checks to pass before merging。
- 在专用归档仓库为 `archive/feature/*` 配置 branch protection rule。
- 在专用归档仓库为 `archive/feature/*` 禁止 force push。
- 在专用归档仓库为 `archive/feature/*` 禁止删除分支。
- 在专用归档仓库为 `archive/feature/*` 禁止直接 push，归档分支只读。

完成以上配置后，应在后续管理提交中更新本文，记录配置完成时间和执行人。
