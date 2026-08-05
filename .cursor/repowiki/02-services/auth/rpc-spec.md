> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/auth/rpc/auth.proto`

---

# Auth RPC Spec

## 服务名

`Auth` — 身份认证与权限管理微服务，通过 etcd 服务发现（key: `auth.rpc`）。

## RPC 方法总览

### 角色管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveRole` | `RoleSaveReq { id, code, name, type }` | `IdReply { id }` | 新增/更新角色，`id <= 0` 为新增 |
| `DeleteRole` | `IdReq { id }` | `Empty {}` | 软删除角色，固定角色(type=0)不可删除 |
| `GetRole` | `IdReq { id }` | `RoleVO { id, code, name, type, createTime }` | 按 ID 查询角色 |
| `ListRoles` | `PageReq { pageNo, pageSize }` | `RoleListReply { total, list }` | 分页查询角色列表 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 角色 ID，新增时省略 |
| `code` | string | 角色代号（如 admin/teacher/student），唯一 |
| `name` | string | 角色名称 |
| `type` | int32 | 0-固定角色（不可改代号、不可删），1-自定义角色 |

---

### 菜单管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetMenuTree` | `Empty {}` | `MenuTreeReply { list }` | 返回完整树形菜单结构 |
| `SaveMenu` | `MenuSaveReq { id, parentId, label, path, icon, priority }` | `IdReply { id }` | 新增/更新菜单 |
| `DeleteMenu` | `IdReq { id }` | `Empty {}` | 删除菜单（须无子菜单） |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 菜单 ID，新增时省略 |
| `parentId` | int64 | 父菜单 ID，0 表示一级菜单 |
| `label` | string | 菜单显示文本 |
| `path` | string | 菜单路径 |
| `icon` | string | 菜单图标 |
| `priority` | int32 | 排序优先级，默认 127，值越小越前 |

---

### 权限管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetPrivilegesByMenu` | `IdReq { id }` | `PrivilegeListReply { list }` | 按菜单查权限列表 |
| `SavePrivilege` | `PrivilegeSaveReq { id, menuId, intro, method, uri, internal }` | `IdReply { id }` | 新增/更新权限 |
| `DeletePrivilege` | `IdReq { id }` | `Empty {}` | 删除权限并解绑角色 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `menuId` | int64 | 所属菜单 ID |
| `intro` | string | 权限说明 |
| `method` | string | HTTP 请求方法 (GET/POST/PUT/DELETE) |
| `uri` | string | API 请求路径 |
| `internal` | bool | 是否内部接口 |

---

### 角色-菜单/权限分配

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveRoleMenus` | `RoleMenuReq { roleId, menuIds }` | `Empty {}` | 替换角色菜单分配（原子替换） |
| `SaveRolePrivileges` | `RolePrivilegeReq { roleId, privilegeIds }` | `Empty {}` | 替换角色权限分配 |
| `GetRoleMenus` | `IdReq { id }` | `IdListReply { ids }` | 查询角色的菜单 ID 列表 |
| `GetRolePrivileges` | `IdReq { id }` | `IdListReply { ids }` | 查询角色的权限 ID 列表 |

---

### 账户-角色分配

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveAccountRoles` | `AccountRoleReq { accountId, roleIds }` | `Empty {}` | 替换账户角色绑定 |
| `GetAccountRoles` | `IdReq { id }` | `IdListReply { ids }` | 查询账户的角色 ID 列表 |

---

### 登录记录

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveLoginRecord` | `LoginRecordReq { userId, cellPhone, ipv4 }` | `Empty {}` | 记录登录日志（失败不抛异常） |
| `ListLoginRecords` | `LoginRecordPageReq { pageNo, pageSize, userId }` | `LoginRecordListReply { total, list }` | 分页查询登录记录 |

---

### 令牌签发

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SignToken` | `SignTokenReq { userId, accountId, roleCode, expireSec }` | `SignTokenReply { accessToken, refreshToken, expiresAt }` | 签发 JWT 访问/刷新令牌 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | int64 | 用户 ID |
| `accountId` | int64 | 员工账号 ID，学员登录为 0 |
| `roleCode` | string | 角色代码 (USER/STUDENT/TEACHER/ADMIN) |
| `expireSec` | int64 | 访问令牌有效期秒数，0 用配置默认 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `auth-api` (自身 API 层) | HTTP Handler → `authclient.Auth` RPC | 所有 RBAC 管理接口最终走自身 RPC |
| `learning-api` | 推测使用 `authclient` 进行 token 校验 | 学习服务需要 JWT 鉴权 |

（注：auth 服务是认证域的基础设施，多数其他服务的 handler 中通过 middleware/jwt 使用 auth rpc 做鉴权，但具体 import 路径需看 `apps/*/api/internal/middleware/jwt`。）

---

## 调用典型场景

1. **用户登录** → 前端登录接口调用 `user` RPC 验证凭证 → 成功后调 `Auth.SignToken` 签发 JWT → 返回 `accessToken` / `refreshToken`
2. **菜单加载** → 管理端前端调 `GetMenuTree` → 组装左侧导航树
3. **角色权限分配** → 管理员在后台调 `SaveRoleMenus` / `SaveRolePrivileges` 替换角色菜单/权限
4. **账户角色绑定** → 管理员调 `SaveAccountRoles` 给员工分配角色
5. **登录审计** → 登录成功后异步调 `SaveLoginRecord` 记录登录时间/IP

---

## 自定义 Model 方法

`rolemodel.go` 扩展了：
- `FindPage(ctx, offset, limit)` — 分页查询（已删除角色排除）
- `ExistsByCode(ctx, code, excludeId)` — 角色 code 唯一性校验
- `SoftDelete(ctx, id, updater)` — 逻辑删除角色