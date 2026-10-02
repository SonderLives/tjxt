# AGENTS.md

## 项目概况

- 本仓库是基于 go-zero 的 Go 微服务项目，使用 Go 1.26 和 `go.work` 管理多个模块。
- 服务分为 API 与 RPC 层；`apps/data` 的 API/RPC 位于额外一层 `data/` 子目录。
- 修改前先查看根目录 `README.md`、相关服务代码和 `Makefile`，沿用项目现有结构与约定。

## 修改约定

- API 契约以 `.api` 文件为准，RPC 契约以 `.proto` 文件为准。契约变更时使用 `make api` 或 `make rpc` 生成代码，并检查生成差异。
- 业务逻辑放在对应服务的 `internal/logic/`；不要无故手改 goctl 生成的 handler、types、pb、server 或 model 代码。
- 数据库 DDL 和迁移分别维护在 `sql/ddl/` 与 `sql/migration/`；不要把密钥、真实凭据或本地环境配置写入仓库。
- 保持改动范围与用户请求一致；不要覆盖已有的用户修改。
- 面向用户的说明及 Git 提交标题使用中文。

## 校验

- 变更 Go 代码后，按改动范围运行对应检查；完整本地门禁为 `make ci`，包含结构校验、Lint、测试和构建。
- CI 会逐个处理 `go.work` 中的模块；只在根目录运行一次 `go test ./...` 或 `go build ./...` 不代表已覆盖全部服务。
- 若环境缺少 Go、goctl、golangci-lint 或服务依赖，说明具体未运行的检查及原因。

## GitHub 与 Codex

- `@codex` 的 PR 审查依赖用户将 ChatGPT/Codex 账号连接到 GitHub，并授权访问本仓库；不要将 Codex 伪装成提交作者或共同作者。
- 除非用户明确要求，不要自行提交、推送、评论或合并 PR。
