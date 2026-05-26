# go-micro-foundry

`go-micro-foundry` 是一套面向小团队的 Go 轻量微服务基础设施项目。项目目标是用低依赖、易部署、可运维的方式沉淀配置中心、服务注册发现、网关、鉴权和链路追踪等基础能力，让后续业务服务可以按统一规范接入和演进。

## 项目定位

- 使用 Go 构建轻量基础服务，默认面向单仓库协作。
- 优先保持部署简单、依赖清晰、服务边界明确。
- 每个微服务独立实现、独立启动，并共享统一的 API、配置、日志、监控和追踪约定。
- 主分支 `main` 只保存稳定基线，功能开发通过 feature 分支完成。

## 微服务目录模式

项目采用单仓库目录模式，不为每个微服务维护长期开发分支。后续代码按以下结构组织：

```text
.
├── services
│   ├── config-center
│   ├── service-registry
│   ├── api-gateway
│   ├── auth-service
│   └── tracing-service
├── internal
├── pkg
└── docs
```

目录说明：

- `services/config-center`：配置中心，负责配置管理、版本、发布和回滚。
- `services/service-registry`：服务注册发现，负责实例注册、心跳、健康状态和服务发现。
- `services/api-gateway`：网关服务，负责统一入口、路由代理、鉴权接入、限流和访问日志。
- `services/auth-service`：鉴权服务，负责用户、角色、权限、令牌和服务凭证。
- `services/tracing-service`：链路追踪服务，负责 Trace/Span 接收、查询、采样和日志关联。
- `internal`：仓库内部共享实现，禁止作为外部公共 API 使用。
- `pkg`：确需对外复用的公共包，新增前必须确认稳定性和复用价值。

## 开发流程

完整开发约束见 [docs/development-workflow.md](docs/development-workflow.md)。
微服务阶段拆分和联调口径见 [docs/microservice-development-plan.md](docs/microservice-development-plan.md)。
开发分支与归档分支的对应关系见 [docs/archive-map.md](docs/archive-map.md)。
历史设计与开发文档已归档到 `git@github.com:tiankongzhise/go-micro-foundry-archive.git` 的 `archive/feature/origin-docs` 分支。

核心规则：

- 禁止直接在 `main` 上开发或直接推送。
- 每个功能从 `main` 新建 `feature/<feature-name>`。
- 涉及微服务实现、接口、配置、UI、部署或联调的开发，必须先完整阅读归档的 `origin_docs` 设计文档。
- 涉及多个微服务的 feature 必须先完成预开发规划，统一阶段目标、服务边界和联调验收口径。
- 涉及多个微服务时，从 feature 分支创建服务子分支。
- 服务子分支开发完成后必须用 merge 合并回对应 feature 分支，保留子分支开发记录。
- feature 完成后 squash 合入 `main`。
- 联调通过后，先迁出临时归档分支并推送到专用归档仓库，再回到 feature 更新归档映射文档。
- squash commit 必须在 trailer 区域记录 `Archive-Ref`，指向归档仓库、归档分支和归档 HEAD SHA。
- 每实现一个功能点提交一次中文 commit，说明变化、技术细节、影响范围和验证方式。

## 当前状态

当前 `main` 是项目基准分支，包含仓库说明、Go `.gitignore` 和开发流程约束。服务代码与详细设计将在后续 feature 分支中按目录模式逐步补充。
