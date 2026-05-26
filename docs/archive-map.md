# 开发分支归档映射

本文记录开发分支、归档分支、归档 HEAD 和合入 `main` 所需的 `Archive-Ref` 字段。feature 完成联调或人工验收并迁出归档分支后，必须回到 feature 分支更新本文，再提交归档记录 commit。

## 记录规则

- 只有已经推送到专用归档仓库的归档分支才能写入“已归档记录”。
- `Archive-Ref` 必须直接复制到合入 `main` 的 squash commit trailer 区域。
- 归档 HEAD 以 `git ls-remote --heads archive <archive-branch>` 返回的 SHA 为准。
- 新的开发 feature 统一使用 `archive/feature/<feature-name>` 归档分支。
- 历史维护分支若已使用其他归档前缀，按实际分支登记，并在说明中标记为历史记录。

## 已归档记录

### 历史设计文档

- 开发分支：`codex/origin_docs`
- 归档分支：`archive/feature/origin-docs`
- 归档 HEAD：`74e5076bc9373b5502cf9bf1152b6b0be7aa3781`
- 合入状态：未作为完整开发 feature squash 合入 `main`，作为项目原始设计文档归档保留。
- 归档字段：

```text
Archive-Ref: git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/feature/origin-docs@74e5076bc9373b5502cf9bf1152b6b0be7aa3781
```

### 初始化项目工程骨架

- 开发分支：`feature/init-project`
- 归档分支：`archive/feature/init-project`
- 归档 HEAD：`d369e7a1b761a0725d722646b309b20c0291718a`
- `main` 聚合提交：`8713cc10fd79f052be4d297d8d21c7fc9d629ccd`
- 归档字段：

```text
Archive-Ref: git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/feature/init-project@d369e7a1b761a0725d722646b309b20c0291718a
```

### 迁移归档流程到专用仓库

- 开发分支：`fix/archive-repository-workflow`
- 归档分支：`archive/fix/archive-repository-workflow`
- 归档 HEAD：`f1cc649cfbf5c90cd380f2787419f397c47e12e6`
- `main` 聚合提交：`79fffca8ee598d77053841371fd5e8e44ffdc236`
- 说明：历史维护分支归档记录；后续新增开发统一使用 `feature/*` 和 `archive/feature/*`。
- 归档字段：

```text
Archive-Ref: git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/fix/archive-repository-workflow@f1cc649cfbf5c90cd380f2787419f397c47e12e6
```

## 待归档记录

### 微服务开发阶段规划

- 开发分支：`feature/microservice-development-planning`
- 预期归档分支：`archive/feature/microservice-development-planning`
- 当前状态：预开发规划中，待文档验收通过后按归档流程生成归档 HEAD 和 `Archive-Ref`。
