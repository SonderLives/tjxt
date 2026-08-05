# Git 工作流与版本发布规范

> 版本：v1.0 | 更新：2026-08-05

## 分支策略

```
main (保护分支，仅合并 PR，禁止直推)
  │
  ├── develop (集成分支，日常开发合入)
  │     │
  │     ├── feature/<svc>-<short-desc>  (功能分支，从 develop 切出)
  │     ├── bugfix/<svc>-<issue-id>     (修复分支)
  │     └── refactor/<scope>            (重构分支)
  │
  ├── release/v<major>.<minor>          (发布分支，从 develop 切出，仅修版本号/修复阻塞 Bug)
  │
  └── hotfix/<issue-id>                 (热修复，从 main 切出，合回 main + develop)
```

## 提交信息规范（Conventional Commits）

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

**Scope 建议**：服务名（`auth`、`course`、`trade`...）或共享库（`pkg`、`spec`）。

## PR 流程

1. **从 `develop` 切分支** → 功能开发
2. **本地跑通**：`make fmt && make test && make verify`
3. **提交 PR** → 目标分支 `develop`
4. **PR 描述模板**：
   ```markdown
   ## 变更摘要
   - 简述做了什么

   ## 变更类型
   - [ ] feat  [ ] fix  [ ] refactor  [ ] docs  [ ] test  [ ] chore

   ## 影响范围
   - 服务：auth / course / trade ...
   - 接口变更：是/否（如是，需同步更新 .spec/02-services/<svc>/api-spec.md）
   - 数据库变更：是/否（如是，需同步 sql/ddl/ 与 sql/migration/）

   ## 测试验证
   - 单测：`make test` 通过
   - 手工验证步骤/截图

   ## 关联 Issue
   - Closes #XXX
   ```
5. **Code Review** → 至少 1 人 Approve
6. **合并** → Squash merge 到 `develop`，删除源分支

## 版本发布流程

### 语义化版本
- `MAJOR`：不兼容 API 变更（如响应结构调整、认证方式变更）
- `MINOR`：向下兼容的新功能（新增接口、新增字段）
- `PATCH`：向下兼容的 Bug 修复

### 发布步骤
```bash
# 1. 从 develop 切发布分支
git checkout develop
git pull
git checkout -b release/v1.2.0

# 2. 更新版本号（如有 VERSION 文件或 go.mod 版本标签）
# 3. 仅修复阻塞 Bug，禁止加新功能
# 4. 打 Tag 合并回 main
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

## 变更日志
- 自动从 Conventional Commits 生成 `CHANGELOG.md`
- 发布时附在 GitHub Release Notes

## 保护分支规则（GitHub/GitLab 设置）
- `main`、`develop`：Require PR review、Require status checks、Dismiss stale reviews
- `main`：Require linear history、Include administrators