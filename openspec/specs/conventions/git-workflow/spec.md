# Git 工作流与版本发布规范 Specification

## Purpose

规定 tjxt 的分支模型、提交信息格式、PR 流程与语义化版本发布流程，使 13 个服务并行开发时的代码合入路径可预测、变更历史可由 Conventional Commits 自动生成 CHANGELOG、发布动作可回滚。本规格适用于所有仓库贡献者，是 PR 合并与打 Tag 的判定依据。

## Requirements

### Requirement: 分支策略

仓库 SHALL 采用「main + develop + 三类临时分支」的模型：

```
main (保护分支，仅合并 PR，禁止直推)
  ├── develop (集成分支，日常开发合入)
  │     ├── feature/<svc>-<short-desc>  (功能分支，从 develop 切出)
  │     ├── bugfix/<svc>-<issue-id>     (修复分支)
  │     └── refactor/<scope>            (重构分支)
  ├── release/v<major>.<minor>          (发布分支，从 develop 切出，仅修版本号/阻塞 Bug)
  └── hotfix/<issue-id>                 (热修复，从 main 切出，合回 main + develop)
```

#### Scenario: 日常功能开发

- **WHEN** 开发者开始一项新功能
- **THEN** 从 `develop` 切出 `feature/<svc>-<short-desc>` 分支
- **AND** 完成后以 `develop` 为目标分支发起 PR，Squash merge 后删除源分支

#### Scenario: 线上紧急修复

- **WHEN** 生产环境出现需立即修复的问题
- **THEN** 从 `main` 切出 `hotfix/<issue-id>`
- **AND** 修复后同时合回 `main` 与 `develop`，避免修复在下一版本丢失

### Requirement: 提交信息规范

提交信息 SHALL 遵循 Conventional Commits 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

| type | 说明 | 示例 |
|------|------|------|
| feat | 新功能 | `feat(auth): add OAuth2 login` |
| fix | 修复 Bug | `fix(trade): correct refund amount calculation` |
| refactor | 重构（无功能变更） | `refactor(pkg): extract idgen to utils` |
| docs | 文档更新 | `docs(spec): add api-contracts.md` |
| style | 格式调整（不影响逻辑） | `style: gofmt all files` |
| test | 测试相关 | `test(learning): add progress logic test` |
| chore | 构建/工具/依赖 | `chore: upgrade go-zero to v1.10.3` |
| perf | 性能优化 | `perf(course): add redis cache for course list` |

`scope` 建议使用服务名（`auth`、`course`、`trade`…）或共享库（`pkg`、`spec`）。

#### Scenario: 一次典型提交

- **WHEN** 完成课程服务的缓存优化
- **THEN** 提交信息写作 `perf(course): add redis cache for course list`
- **AND** 该提交可被工具识别并归入 CHANGELOG 的性能优化条目

### Requirement: PR 流程

PR SHALL 按「本地校验 → 模板化描述 → Review → Squash merge」的流程合入 `develop`。

1. 从 `develop` 切分支并完成开发
2. 本地跑通 `make fmt && make test && make verify`
3. 提交 PR，目标分支 `develop`
4. PR 描述使用模板：

```markdown
## 变更摘要
- 简述做了什么

## 变更类型
- [ ] feat  [ ] fix  [ ] refactor  [ ] docs  [ ] test  [ ] chore

## 影响范围
- 服务：auth / course / trade ...
- 接口变更：是/否（如是，需同步更新 openspec/specs/services/<svc>/spec.md）
- 数据库变更：是/否（如是，需同步 sql/ddl/ 与 sql/migration/）

## 测试验证
- 单测：`make test` 通过
- 手工验证步骤/截图

## 关联 Issue
- Closes #XXX
```

5. Code Review：至少 1 人 Approve
6. Squash merge 到 `develop`，删除源分支

#### Scenario: 接口或数据库变更

- **WHEN** PR 修改了 `.api` / `.proto` 或表结构
- **THEN** 必须在同一 PR 内同步更新对应 openspec 规格与 `sql/ddl/`、`sql/migration/`
- **AND** 「影响范围」中显式勾选「接口变更：是」或「数据库变更：是」

#### Scenario: 合并前门禁

- **WHEN** PR 等待合并
- **THEN** 需满足：本地 `make fmt && make test && make verify` 通过、至少 1 人 Approve、目标分支为 `develop`
- **AND** 采用 Squash merge，合并后删除源分支

### Requirement: 语义化版本

版本号 SHALL 遵循 SemVer（`MAJOR.MINOR.PATCH`）：

- `MAJOR`：不兼容 API 变更（响应结构调整、认证方式变更）
- `MINOR`：向下兼容的新功能（新增接口、新增字段）
- `PATCH`：向下兼容的 Bug 修复

#### Scenario: 判断版本号增量

- **WHEN** 一次发布同时包含新增接口与响应结构调整
- **THEN** 以最高影响级别为准，递增 `MAJOR`
- **AND** 纯新增接口（字段只增不减）递增 `MINOR`

### Requirement: 发布流程

发布 SHALL 经 release 分支执行，并在合入 `main` 后打 Tag、回合 `develop`、删除发布分支。

```bash
# 1. 从 develop 切发布分支
git checkout develop && git pull
git checkout -b release/v1.2.0

# 2. 更新版本号
# 3. 仅修复阻塞 Bug，禁止加新功能
# 4. 打 Tag 并合并回 main
git tag -a v1.2.0 -m "Release v1.2.0"
git checkout main
git merge --no-ff release/v1.2.0
git push origin main --tags

# 5. 合回 develop
git checkout develop
git merge --no-ff release/v1.2.0
git push origin develop

# 6. 删除发布分支
git branch -d release/v1.2.0
```

#### Scenario: 标准发布

- **WHEN** 准备发布 v1.2.0
- **THEN** 从 `develop` 切 `release/v1.2.0`，仅修版本号与阻塞 Bug
- **AND** 依次完成打 Tag、`--no-ff` 合入 `main` 并推送 tags、回合 `develop`、删除发布分支

### Requirement: 变更日志与分支保护

CHANGELOG SHALL 由 Conventional Commits 自动生成并附于 Release Notes；保护分支 SHALL 开启强制校验。

- CHANGELOG：自动生成 `CHANGELOG.md`，发布时附在 GitHub Release Notes
- `main`、`develop`：Require PR review、Require status checks、Dismiss stale reviews
- `main`：额外 Require linear history、Include administrators

#### Scenario: 保护分支被直推

- **WHEN** 有人尝试直接 push 到 `main`
- **THEN** 被分支保护规则拒绝（Include administrators 亦生效）
- **AND** 变更必须经 PR 与 status checks 后合入
