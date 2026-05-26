# 微服务开发阶段规划

本文用于 `feature/microservice-development-planning` 预开发阶段，目标是在正式编码前统一五个基础微服务的任务边界、阶段目标、联调顺序和验收口径。当前分支只做规划和流程约束，不实现业务代码。

## 规划依据

本规划基于已归档的 `origin_docs` 设计文档整理，归档位置为：

```text
git@github.com:tiankongzhise/go-micro-foundry-archive.git archive/feature/origin-docs@74e5076bc9373b5502cf9bf1152b6b0be7aa3781
```

阅读顺序和规划依据：

- `docs/architecture.md`：五个基础服务的职责边界、调用关系、部署拓扑和统一请求链路。
- `docs/api-conventions.md`：只使用 `GET` 和 `POST`，状态变更通过 payload 的 `action` 表达，响应、错误码、审计和健康检查保持一致。
- `docs/config-conventions.md`：所有服务使用 `TOML + .env`，样例配置必须逐项中文注释，敏感项只出现在 `.env` 或真实环境变量中。
- `docs/ui-conventions.md`：每个服务内嵌 `/ui`，由 Go `go:embed` 打包，UI 只调用当前服务 API。
- 五个服务详细设计：配置中心、服务注册发现、网关、鉴权和链路追踪的 API、数据模型、UI、权限、指标和失败场景。
- `docs/deployment.md` 和 `docs/roadmap.md`：默认单机多进程、systemd 管理，推荐先实现公共基础库，再实现鉴权、追踪、配置中心、注册发现和网关。

## 开发目标

- 按归档设计文档还原服务边界，不把业务编排、复杂平台能力或重型依赖提前带入 MVP。
- 统一各阶段每个微服务必须达到的完成度，避免后续团队并行开发后进入联调时能力错位。
- 明确服务子分支、公共基础包和联调分支的协作方式。
- 为每个后续 feature 保留清晰的验收、归档和合入 `main` 依据。

## 阶段定义

### P0 预开发对齐

P0 只在 feature 集成分支完成，不创建服务子分支。

- 阅读归档设计文档，确认本次 feature 是否沿用原设计或需要显式调整。
- 写清 feature 目标、非目标、影响范围和验收方式。
- 拆分服务子分支、公共基础包任务和联调任务。
- 对齐 API、配置、错误码、审计字段、`X-Request-Id`、`X-Trace-Id` 和服务间凭证约定。
- 定义联调用例、依赖启动顺序和验收数据。

### P1 公共契约与服务壳

P1 结束时，所有服务应能独立启动，并暴露统一的基础接口。

- 每个服务提供 `GET /health`、`GET /metrics` 和 `/ui` 静态入口。
- 每个服务具备 `config.example.toml`、`.env.example` 和配置校验。
- 统一响应格式、错误码、请求 ID、Trace ID、日志字段和审计字段。
- 公共能力优先沉淀在 `internal/foundation`，只有确认稳定且需要对外复用时才进入 `pkg`。
- 每个服务至少具备启动测试、配置加载测试和基础响应测试。

### P2 核心 API 与数据闭环

P2 结束时，各服务完成 MVP 范围内的核心 API、数据模型和最小业务闭环。

- 所有业务接口仍只使用 `GET` 和 `POST`。
- 状态变更必须通过 JSON payload 的 `action` 字段表达。
- 管理类 `POST` 接口必须记录审计信息。
- 对外错误响应不得泄露密钥、堆栈、数据库路径或内部拓扑。
- 存储层先保留接口边界，MVP 可使用 SQLite 或等价本地轻量存储，并预留 PostgreSQL 适配。

### P3 UI、运维与服务间接入

P3 结束时，各服务应具备管理员可操作、开发者可接入、运维可排障的最小能力。

- `/ui` 能完成本服务核心资源的列表、详情、创建、编辑、启用、停用、发布或回滚操作。
- 危险操作必须确认并填写 `reason`。
- 监控页能展示 `GET /metrics` 的关键指标。
- 使用说明页能离线打开，并说明部署、接入和常见排障。
- 服务间调用使用鉴权服务发放的服务应用凭证。

### P4 联调验收与归档准备

P4 结束后才允许进入 feature 归档和合入 `main` 的准备流程。

- 所有服务子分支已经 merge 回 feature 集成分支。
- feature 分支已同步最新 `main`，并解决冲突。
- 完成端到端联调用例并记录验证方式。
- 未完成事项必须拆成后续 feature，不能隐藏在归档记录或 squash commit 中。
- 按 `docs/development-workflow.md` 迁出归档分支、推送归档仓库、更新 `docs/archive-map.md` 并生成 `Archive-Ref`。

## 推荐实现顺序

后续编码阶段按以下顺序推进，符合归档路线图对依赖关系的建议：

1. 公共基础能力：配置加载、统一响应、错误码、日志、请求上下文、Trace Header、健康检查和指标骨架。
2. 鉴权服务：先解决管理员登录、RBAC、服务应用凭证和 token 校验，支撑其他服务 UI 和服务间调用。
3. 链路追踪服务：尽早提供 Span 写入和 Trace 查询，让后续服务开发时即可接入排障链路。
4. 配置中心：提供配置发布、回滚和客户端拉取，支撑业务服务配置接入。
5. 服务注册发现：提供实例注册、心跳和发现查询，支撑网关代理。
6. 网关服务：最后接入鉴权、注册发现和追踪，形成外部请求端到端闭环。

## 服务阶段矩阵

| 服务 | P1 公共契约与服务壳 | P2 核心 API 与数据闭环 | P3 UI、运维与服务间接入 | P4 联调验收 |
| --- | --- | --- | --- | --- |
| `auth-service` | 端口 `18030`、健康检查、指标、JWT 和启动密钥配置校验 | `POST /tokens`、`GET/POST /users`、`GET/POST /roles`、`GET/POST /service-apps`，完成登录、刷新、校验、RBAC 和服务凭证 | 登录页、用户管理、角色权限、服务应用、令牌管理、审计日志；密钥只展示一次并脱敏 | 网关可调用 `POST /tokens` 校验 Bearer Token，其他基础服务可用服务凭证校验管理请求 |
| `tracing-service` | 端口 `18040`、Trace Header 约定、Span 数据结构、存储和保留配置 | `POST /spans`、`GET /traces`、`GET /spans`、`GET/POST /sampling-rules`、`POST /trace-storage`，完成写入、查询、采样和清理 | Trace 查询、Span 详情、慢请求、错误链路、采样规则和存储管理页面 | 网关和至少两个基础服务能上报 Span，按 Trace ID 可查到完整调用链 |
| `config-center` | 端口 `18010`、应用/环境/配置项/版本模型、配置拉取默认策略 | `GET/POST /apps`、`GET/POST /configs`、`GET /config-versions`、`GET /client-configs`，完成配置草稿、发布、回滚和客户端拉取 | 应用管理、配置管理、发布管理、版本历史、审计日志和接入说明 | 测试服务可拉取已发布配置，发布和回滚后能观察到版本变化 |
| `service-registry` | 端口 `18020`、服务/实例/心跳模型、心跳间隔和超时配置 | `GET/POST /services`、`GET/POST /instances`、`GET /discoveries`，完成注册、心跳、健康状态、发现和下线保护 | 服务列表、实例列表、心跳记录、下线管理、审计日志和注册接入说明 | 测试服务可注册并心跳，网关可从 `GET /discoveries` 获取 healthy 实例 |
| `api-gateway` | 端口 `18080`、路由/上游缓存/访问日志模型、代理超时和限流配置 | `GET/POST /routes`、`GET/POST /upstreams`、`GET /access-logs`，完成路由代理、鉴权接入、服务发现缓存、固定窗口限流和失败实例跳过 | 路由管理、上游服务、限流规则、访问日志、监控和审计日志 | 外部请求经过网关完成鉴权、发现上游、代理调用、访问日志和 Span 上报 |
| `internal/foundation` | 配置加载、响应、错误、日志、请求上下文和 Trace Header 基础能力 | 抽象存储、审计、指标和服务间调用辅助能力，不绑定具体业务服务 | 为各服务提供一致的配置校验、健康检查、指标和 UI 静态资源挂载辅助 | 支撑五个服务在联调中保持响应、日志、Trace 和配置行为一致 |

## 服务子分支拆分

后续涉及完整 MVP 的 feature 可以按服务创建子分支：

```text
feature/<feature-name>/foundation
feature/<feature-name>/auth-service
feature/<feature-name>/tracing-service
feature/<feature-name>/config-center
feature/<feature-name>/service-registry
feature/<feature-name>/api-gateway
feature/<feature-name>/integration
```

协作要求：

- `foundation` 子分支优先合并，避免各服务重复实现公共能力。
- 服务子分支只实现本服务职责和必要公共能力，不抢占其他服务边界。
- `integration` 子分支只补联调脚本、测试桩、运行说明和验证记录，不承载单个服务的大量业务实现。
- 服务子分支合并回 feature 时必须使用 `merge --no-ff`，保留分支历史和集成顺序。

## 联调用例

P4 至少完成以下端到端链路：

1. 启动 `tracing-service` 和 `auth-service`。
2. 在 `auth-service` 创建管理员、`api-gateway` 服务应用和基础服务间凭证。
3. 启动 `config-center`，创建应用和环境，发布一份测试配置。
4. 启动 `service-registry`，注册一个测试业务服务实例并持续心跳。
5. 启动 `api-gateway`，创建指向测试业务服务的路由并发布。
6. 外部请求携带 Bearer Token 访问网关。
7. 网关调用鉴权服务校验 token，从注册中心发现 healthy 实例，代理到测试业务服务。
8. 测试业务服务拉取配置中心已发布配置，并返回配置版本摘要。
9. 网关、测试业务服务和相关基础服务向链路追踪服务写入 Span。
10. 在链路追踪 UI 或 API 中通过 Trace ID 查到完整链路。

## 验收清单

feature 结束前必须确认：

- 当前 feature 范围内的服务阶段矩阵均达到约定完成度。
- 所有 API 示例遵守 `GET`、`POST` 和 payload `action` 约定。
- 所有样例配置逐项中文注释，真实 `config.toml`、`.env` 和密钥未提交。
- 所有服务保留 `/health`、`/metrics`、`/ui` 和使用说明入口。
- 联调记录包含启动顺序、测试数据、请求样例、验证结果和已知限制。
- 已完成归档映射记录，并生成合入 `main` 所需的 `Archive-Ref`。
