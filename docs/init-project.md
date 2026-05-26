# init-project 初始化边界

`feature/init-project` 是从 `main` 迁出的第一个开发分支，目标是为后续五个基础服务准备最小可编译工程骨架。

## 范围

- 初始化 Go module：`go-micro-foundry`。
- 建立 `cmd`、`services`、`internal` 和 `pkg` 目录。
- 预置五个服务目录：配置中心、服务注册发现、网关、鉴权和链路追踪。
- 在 `internal/foundation` 下沉淀公共基础能力：
  - 配置加载顺序骨架：默认值、TOML、`.env`、系统环境变量、命令行参数。
  - 统一 HTTP 响应格式。
  - `X-Request-Id` 和 `X-Trace-Id` 生成与透传。
  - 基础日志字段约定。

## 非目标

- 不实现具体业务接口。
- 不引入数据库或迁移工具。
- 不提交真实 `config.toml`、`.env` 或密钥。
- 不把历史设计文档分支全量合入当前 feature；历史文档已归档到 `git@github.com:tiankongzhise/go-micro-foundry-archive.git` 的 `archive/feature/origin-docs` 分支。

## 验证

当前初始化完成后必须至少执行：

```bash
go test ./...
```
