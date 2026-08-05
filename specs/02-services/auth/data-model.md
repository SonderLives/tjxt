> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_auth.sql`

---

# Auth Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `rolemodel.go` | `role` | FindPage, ExistsByCode, SoftDelete |
| `menumodel.go` | `menu` | SyncHasChildren, FindByIds, CountChildren |
| `roleprivilegemodel.go` | `role_privilege` | DeleteByRoleId, DeleteByPrivilegeId |
| `rolemenumodel.go` | `role_menu` | DeleteByRoleId, DeleteByMenuId, FindMenuIdsByRoleId |
| `accountrolemodel.go` | `account_role` | ReplaceByAccountId, FindRoleIdsByAccountId, CountByRoleId |
| `privilegemodel.go` | `privilege` | FindByMenuId, FindByIds, SoftDelete |
| `loginrecordmodel.go` | `login_record` | FindPage (按 userId 分页) |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `rolemodel.go` 扩展 `RoleModel` 接口）。

---

## 表清单与字段说明

### 1. `role` — 角色表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键 | PK |
| `code` | varchar(64) | 角色代号（admin/teacher/student 等），唯一业务键 | 唯一约束 |
| `name` | varchar(255) | 角色名称 | - |
| `type` | tinyint | 0=固定角色(不可改), 1=自定义角色 | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建者 ID | - |
| `updater` | bigint | 更新者 ID | - |
| `dep_id` | bigint | 部门 ID | - |
| `deleted` | tinyint | 0=正常, 1=逻辑删除 | 过滤条件 |

---

### 2. `menu` — 菜单表（树形结构）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键 | PK |
| `parent_id` | bigint | 父菜单 ID，0=一级菜单 | 树形查询 |
| `has_children` | tinyint | 是否有子菜单 | 缓存用 |
| `label` | varchar(16) | 菜单文本 | - |
| `path` | varchar(255) | 菜单路径 | - |
| `icon` | varchar(32) | 菜单图标 | - |
| `priority` | tinyint | 排序值，默认 127 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |
| `creater` / `updater` | bigint | 创建/更新人 | - |
| `dep_id` | bigint | 部门 ID | - |
| `deleted` | tinyint | 逻辑删除 | 过滤条件 |

> **关系**：通过 `parent_id` 自引用构建树形结构，`has_children` 冗余字段便于前端展示。

---

### 3. `privilege` — 权限表（API 访问权限）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键 | PK |
| `menu_id` | bigint | 所属菜单 ID | 外键关联 menu |
| `intro` | varchar(255) | 权限说明 | - |
| `method` | varchar(16) | HTTP 方法 (GET/POST...) | - |
| `uri` | varchar(255) | API 路径 | - |
| `internal` | tinyint | 是否内部接口 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |
| `creater` / `updater` | bigint | 创建/更新人 | - |
| `dep_id` | bigint | 部门 ID | - |
| `deleted` | tinyint | 逻辑删除 | 过滤条件 |

---

### 4. `role_menu` — 角色-菜单关联表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 自增主键 | PK |
| `role_id` | bigint | 角色 ID | FK → role |
| `menu_id` | bigint | 菜单 ID | FK → menu |

> 一对多关系：一个角色可以绑定多个菜单。

---

### 5. `role_privilege` — 角色-权限关联表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 自增主键 | PK |
| `role_id` | bigint | 角色 ID | FK → role |
| `privilege_id` | bigint | 权限 ID | FK → privilege |

> 一对多关系：一个角色可以拥有多个权限。

---

### 6. `account_role` — 账户-角色关联表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 自增主键 | PK |
| `account_id` | bigint | 账户 ID | FK → user.account |
| `role_id` | bigint | 角色 ID | FK → role |

> 多对多关系：一个账户可有多角色，一个角色可分配给多账户。

---

### 7. `login_record` — 登录记录表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 自增主键 | PK |
| `user_id` | bigint | 用户 ID | 过滤条件 |
| `cell_phone` | varchar(11) | 手机号码 | - |
| `login_time` | datetime | 登录时间 | 自动填充 |
| `logout_time` | datetime | 登出时间 | 可空 |
| `login_date` | date | 登录日期 | 索引 |
| `duration` | bigint | 登录时长(秒) | - |
| `ipv4` | varchar(15) | IP 地址 | - |

---

## 关系图

```
role (1) ──┬── (N) role_menu ─── (N) menu
           │
           └── (N) role_privilege ─── (N) privilege

account (1) ─── (N) account_role ─── (N) role

login_record → user (外键，在 user 域)
```

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
rolemodel_gen.go          ← goctl 生成，只读
rolemodel.go              ← 手写扩展 FindPage / ExistsByCode / SoftDelete
```

当前项目自定义 Model 模式：
- `rolemodel.go` — FindPage, ExistsByCode, SoftDelete
- `menumodel.go` — SyncHasChildren, FindByIds, CountChildren
- `accountrolemodel.go` — ReplaceByAccountId, FindRoleIdsByAccountId, CountByRoleId
- `roleprivilegemodel.go` — DeleteByRoleId, DeleteByPrivilegeId
- `privilegemodel.go` — FindByMenuId, FindByIds, SoftDelete
- `rolemenumodel.go` — DeleteByRoleId, DeleteByMenuId, FindMenuIdsByRoleId
- `loginrecordmodel.go` — FindPage

`sqlhelper.go` 提供通用 SQL 工具函数：`inPlaceholders`, `pairValuePlaceholders`, `dedupeIds`。