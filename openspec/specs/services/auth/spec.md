# auth（认证授权）Specification

## Purpose

auth（认证授权）是认证域的基础微服务，职责为「RBAC 权限模型（角色 / 菜单 / 权限 / 角色-菜单、角色-权限、账户-角色绑定）+ JWT 令牌签发 + 登录记录审计」。HTTP API 层 `auth-api` 监听 `0.0.0.0:8802`，RPC 服务 `Auth` 监听 `0.0.0.0:8082`（gRPC），经 etcd 以 key `auth.rpc` 注册与发现。数据存储于 MySQL `tj_auth` 库（`role` / `menu` / `privilege` / `role_menu` / `role_privilege` / `account_role` / `login_record` 七张表），角色/菜单/权限数据由 Redis 缓存托管（node 单机模式）。登录链路上 `auth-api` 依赖 `user.rpc` 的 `LoginVerify` 核验用户凭证，核验通过后由自身 `Auth.SignToken` 签发 JWT；`learning-api` 推测使用 `authclient` 进行 token 校验。

## Requirements

### Requirement: 登录与令牌管理 HTTP 接口

> **状态说明**：下表来自聚合文档 `docs/tjxt.openapi.json`（api-spec.md，最后同步 2026-08-05）。源表「认证」列均标注「否」、「权限标签」列均为空，为原始聚合文档标注，本规格按源文件忠实转写。

auth 服务 SHALL 暴露以下登录与令牌管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /accounts/admin/login | 管理端登录并获取token | LoginFormDTO | R{data: R«string»} |
| POST | /accounts/login | 登录并获取token | LoginFormDTO | R{data: R«string»} |
| POST | /accounts/logout | 退出登录 | - | R{data: R} |
| GET | /accounts/refresh | 刷新token | - | R{data: R«string»} |

统一约定（引用全局规范）：响应格式为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页为 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码见共享错误码表（`specs/shared/error-codes`）。

#### Scenario: 登录并获取 token

- **WHEN** 客户端以 `LoginFormDTO` 调用 `POST /accounts/login`（或管理端调用 `POST /accounts/admin/login`）
- **THEN** 系统经 `user.rpc` 的 `LoginVerify` 核验凭证，成功后调 `Auth.SignToken` 签发 JWT 并返回 `R{data: R«string»}`
- **AND** 登录成功后异步调用 `SaveLoginRecord` 记录登录时间 / IP

#### Scenario: 刷新 token

- **WHEN** 客户端调用 `GET /accounts/refresh`
- **THEN** 系统基于刷新令牌重新签发访问令牌，返回 `R{data: R«string»}`

### Requirement: 菜单管理 HTTP 接口

> **状态说明**：同上，源自 `docs/tjxt.openapi.json`（api-spec.md）。

auth 服务 SHALL 暴露以下菜单管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /menus | 查询菜单，按照多级菜单组成树结构 | - | R{data: R«List«MenuOptionVO»»} |
| POST | /menus | 新增菜单 | MenuDTO | R{data: R} |
| GET | /menus/me | 查询我的菜单，按照多级菜单组成树结构 | - | R{data: R«List«MenuOptionVO»»} |
| GET | /menus/parent/{pid} | 根据父菜单id查询子菜单 | - | R{data: R«List«MenuOptionVO»»} |
| DELETE | /menus/role/{roleId} | 解除角色的菜单权限 | bindRoleMenusUsingPOSTMenuids | R{data: R} |
| POST | /menus/role/{roleId} | 绑定角色与菜单权限 | bindRoleMenusUsingPOSTMenuids | R{data: R} |
| DELETE | /menus/{id} | 根据id删除菜单 | - | R{data: R} |
| GET | /menus/{id} | 根据id查询菜单 | - | R{data: R«MenuOptionVO»} |
| PUT | /menus/{id} | 更新菜单 | MenuDTO | R{data: R} |

#### Scenario: 查询菜单树

- **WHEN** 管理端调用 `GET /menus`（或当前登录用户调用 `GET /menus/me`）
- **THEN** 系统按多级菜单组成树结构返回 `R{data: R«List«MenuOptionVO»»}`（树构建规则见「菜单树管理规则」）

#### Scenario: 绑定与解除角色的菜单权限

- **WHEN** 管理员以菜单 ID 列表调用 `POST /menus/role/{roleId}`
- **THEN** 系统以原子替换语义绑定角色与菜单权限，返回 `R{data: R}`（规则见「角色-菜单/权限分配规则」）
- **AND** 调用 `DELETE /menus/role/{roleId}` 解除该角色的菜单权限

### Requirement: 权限管理 HTTP 接口

> **状态说明**：同上，源自 `docs/tjxt.openapi.json`（api-spec.md）。

auth 服务 SHALL 暴露以下权限（API 访问权限）管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /privileges | 分页查询所有权限 | - | R{data: R«PageDTO«PrivilegeDTO»»} |
| POST | /privileges | 新增权限 | PrivilegeDTO | R{data: R«PrivilegeDTO»} |
| GET | /privileges/options/{menuId} | 查询菜单下的所有权限，作为下拉选框菜单 | - | R{data: R«List«PrivilegeOptionVO»»} |
| DELETE | /privileges/role/{roleId} | 解除角色的API权限 | bindRolePrivilegesUsingPOSTPrivilegeids | R{data: R} |
| POST | /privileges/role/{roleId} | 绑定角色与API权限 | bindRolePrivilegesUsingPOSTPrivilegeids | R{data: R} |
| GET | /privileges/roles/{roleId}/{menuId} | 查询菜单下的权限列表，某个角色的权限 | - | R{data: R«List«PrivilegeOptionVO»»} |
| DELETE | /privileges/{id} | 删除权限 | - | R{data: R} |
| PUT | /privileges/{id} | 修改权限 | PrivilegeDTO | R{data: R«PrivilegeDTO»} |

#### Scenario: 分页查询所有权限

- **WHEN** 管理端调用 `GET /privileges`
- **THEN** 系统按分页约定返回 `R{data: R«PageDTO«PrivilegeDTO»»}`

#### Scenario: 绑定角色与 API 权限

- **WHEN** 管理员以权限 ID 列表调用 `POST /privileges/role/{roleId}`
- **THEN** 系统以原子替换语义绑定角色与 API 权限，返回 `R{data: R}`
- **AND** 调用 `DELETE /privileges/role/{roleId}` 解除该角色的 API 权限

### Requirement: 角色管理 HTTP 接口

> **状态说明**：同上，源自 `docs/tjxt.openapi.json`（api-spec.md）。

auth 服务 SHALL 暴露以下角色管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /roles | 查询员工角色列表 | - | R{data: R«List«RoleDTO»»} |
| POST | /roles | 新增角色 | RoleDTO | R{data: R«RoleDTO»} |
| GET | /roles/list | 查询员工角色列表 | - | R{data: R«List«RoleDTO»»} |
| DELETE | /roles/{id} | 删除角色信息 | - | R{data: R} |
| GET | /roles/{id} | 根据id查询角色 | - | R{data: R«RoleDTO»} |
| PUT | /roles/{id} | 修改角色信息 | RoleDTO | R{data: R} |

#### Scenario: 新增角色

- **WHEN** 管理端以 `RoleDTO` 调用 `POST /roles`
- **THEN** 系统按「RBAC 角色管理规则」完成校验后新增角色，返回 `R{data: R«RoleDTO»}`

#### Scenario: 删除固定角色被拒绝

- **WHEN** 管理端对固定角色（type=0，如 admin）调用 `DELETE /roles/{id}`
- **THEN** 系统返回冲突错误，固定角色不可删除（规则见「RBAC 角色管理规则」）

### Requirement: RPC 服务契约（Auth 服务）

`Auth` 服务 SHALL 作为身份认证与权限管理微服务经 etcd 服务发现注册（key: `auth.rpc`），监听 `0.0.0.0:8082`（gRPC），方法按以下六个域提供（proto 来源 `apps/auth/rpc/auth.proto`）。

**角色管理**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveRole` | `RoleSaveReq { id, code, name, type }` | `IdReply { id }` | 新增/更新角色，`id <= 0` 为新增 |
| `DeleteRole` | `IdReq { id }` | `Empty {}` | 软删除角色，固定角色(type=0)不可删除 |
| `GetRole` | `IdReq { id }` | `RoleVO { id, code, name, type, createTime }` | 按 ID 查询角色 |
| `ListRoles` | `PageReq { pageNo, pageSize }` | `RoleListReply { total, list }` | 分页查询角色列表 |

`RoleSaveReq` 关键字段：`id`（int64，角色 ID，新增时省略）、`code`（string，角色代号，如 admin/teacher/student，唯一）、`name`（string，角色名称）、`type`（int32，0-固定角色（不可改代号、不可删），1-自定义角色）。

**菜单管理**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetMenuTree` | `Empty {}` | `MenuTreeReply { list }` | 返回完整树形菜单结构 |
| `SaveMenu` | `MenuSaveReq { id, parentId, label, path, icon, priority }` | `IdReply { id }` | 新增/更新菜单 |
| `DeleteMenu` | `IdReq { id }` | `Empty {}` | 删除菜单（须无子菜单） |

`MenuSaveReq` 关键字段：`id`（int64，新增时省略）、`parentId`（int64，父菜单 ID，0 表示一级菜单）、`label`（菜单显示文本）、`path`（菜单路径）、`icon`（菜单图标）、`priority`（int32，排序优先级，默认 127，值越小越前）。

**权限管理**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetPrivilegesByMenu` | `IdReq { id }` | `PrivilegeListReply { list }` | 按菜单查权限列表 |
| `SavePrivilege` | `PrivilegeSaveReq { id, menuId, intro, method, uri, internal }` | `IdReply { id }` | 新增/更新权限 |
| `DeletePrivilege` | `IdReq { id }` | `Empty {}` | 删除权限并解绑角色 |

`PrivilegeSaveReq` 关键字段：`menuId`（所属菜单 ID）、`intro`（权限说明）、`method`（HTTP 请求方法 GET/POST/PUT/DELETE）、`uri`（API 请求路径）、`internal`（bool，是否内部接口）。

**角色-菜单/权限分配**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveRoleMenus` | `RoleMenuReq { roleId, menuIds }` | `Empty {}` | 替换角色菜单分配（原子替换） |
| `SaveRolePrivileges` | `RolePrivilegeReq { roleId, privilegeIds }` | `Empty {}` | 替换角色权限分配 |
| `GetRoleMenus` | `IdReq { id }` | `IdListReply { ids }` | 查询角色的菜单 ID 列表 |
| `GetRolePrivileges` | `IdReq { id }` | `IdListReply { ids }` | 查询角色的权限 ID 列表 |

**账户-角色分配**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveAccountRoles` | `AccountRoleReq { accountId, roleIds }` | `Empty {}` | 替换账户角色绑定 |
| `GetAccountRoles` | `IdReq { id }` | `IdListReply { ids }` | 查询账户的角色 ID 列表 |

**登录记录与令牌签发**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveLoginRecord` | `LoginRecordReq { userId, cellPhone, ipv4 }` | `Empty {}` | 记录登录日志（失败不抛异常） |
| `ListLoginRecords` | `LoginRecordPageReq { pageNo, pageSize, userId }` | `LoginRecordListReply { total, list }` | 分页查询登录记录 |
| `SignToken` | `SignTokenReq { userId, accountId, roleCode, expireSec }` | `SignTokenReply { accessToken, refreshToken, expiresAt }` | 签发 JWT 访问/刷新令牌 |

`SignTokenReq` 关键字段：`userId`（int64，用户 ID）、`accountId`（int64，员工账号 ID，学员登录为 0）、`roleCode`（string，角色代码 USER/STUDENT/TEACHER/ADMIN）、`expireSec`（int64，访问令牌有效期秒数，0 用配置默认）。

消费方 SHALL 为：`auth-api`（自身 API 层，HTTP Handler → `authclient.Auth` RPC，所有 RBAC 管理接口最终走自身 RPC）与 `learning-api`（推测使用 `authclient` 进行 token 校验，学习服务需要 JWT 鉴权）。

#### Scenario: 用户登录签发令牌

- **WHEN** 前端登录接口调用 `user` RPC 验证凭证成功后调 `Auth.SignToken`
- **THEN** 系统返回 `accessToken` / `refreshToken`（及 `expiresAt`），学员登录时 `accountId` 传 0

#### Scenario: 管理端加载菜单与分配权限

- **WHEN** 管理端前端调 `GetMenuTree` 组装左侧导航树，管理员在后台调 `SaveRoleMenus` / `SaveRolePrivileges`
- **THEN** 系统返回完整树形菜单结构，并以原子替换语义更新角色的菜单/权限分配

#### Scenario: 登录审计

- **WHEN** 登录成功后异步调用 `SaveLoginRecord`
- **THEN** 系统记录登录时间 / IP，插入失败不抛异常（规则见「登录记录规则」）

> 注：auth 服务是认证域的基础设施，多数其他服务的 handler 中通过 middleware/jwt 使用 auth rpc 做鉴权，但具体 import 路径需看 `apps/*/api/internal/middleware/jwt`；`learning-api` 的消费方式在源文件中为推测表述。

### Requirement: RBAC 角色管理规则

角色 SHALL 分为两种类型：固定角色（type=0，如 admin）与自定义角色（type=1）。角色管理 SHALL 满足以下约束：

| 规则 | 约束 |
|------|------|
| 角色 code 唯一性 | 新增时校验 code 不重复，更新时可排除自身（`ExistsByCode`） |
| 固定角色不可改代号 | type=0 的角色修改 code 时拒绝 |
| 固定角色不可删除 | type=0 的角色 `DeleteRole` 直接返回冲突错误 |
| 已分配的固定角色不可删除 | 需先解除所有账户绑定 |
| 软删除 | 删除时 `deleted=1` 并失效缓存 |

`SaveRole` 流程：校验 code 非空、name 非空 → 校验 code 唯一性（`ExistsByCode`）→ 有 id 走更新、无 id 走新增（生成 ID）→ 更新时 type=0 且 code 变更则拒绝。

#### Scenario: 保存角色（SaveRole）

- **WHEN** 客户端调用 `SaveRole` 新增或更新角色
- **THEN** 系统先校验 code 与 name 非空，再经 `ExistsByCode` 校验 code 唯一性；请求带 id 时走更新，无 id 时新增并生成 ID
- **AND** `id <= 0` 视为新增

#### Scenario: 角色 code 唯一性校验

- **WHEN** 新增角色时 code 与既有角色重复
- **THEN** 系统拒绝该请求；更新场景下 `ExistsByCode` 排除自身 ID 后判定唯一性

#### Scenario: 固定角色不可修改代号

- **WHEN** 更新 type=0 的固定角色且请求变更其 code
- **THEN** 系统拒绝该修改

#### Scenario: 固定角色不可删除（软删除语义）

- **WHEN** 对 type=0 的固定角色调用 `DeleteRole`
- **THEN** 系统直接返回冲突错误；对自定义角色删除时置 `deleted=1` 并失效缓存
- **AND** 已分配的固定角色需先解除所有账户绑定方可处理删除

### Requirement: 菜单树管理规则

菜单 SHALL 通过 `parent_id` 自引用构建树形结构。菜单管理 SHALL 满足以下约束：

| 规则 | 约束 |
|------|------|
| 父菜单不能是自身 | `id == parentId` 时拒绝 |
| 子菜单存在不可删 | 删除菜单前检查 `CountChildren > 0` 则拒绝 |
| 删除级联清理 | 删除菜单时同时删除其下所有权限、角色关联 |
| has_children 同步 | 增/删/移动菜单后调用 `SyncHasChildren` 更新父节点和自身 |
| 排序优先级 | priority 值越小越前，默认 127 |

`GetMenuTree` 流程：查所有菜单（FindAll）→ 构建 id→节点 map → 按 parentId 拼接子节点 → 提取 root（parentId <= 0）→ 按 priority 排序。

#### Scenario: 构建完整菜单树（GetMenuTree）

- **WHEN** 客户端调用 `GetMenuTree`
- **THEN** 系统 FindAll 查所有菜单，构建 id→节点 map，按 parentId 拼接子节点，提取 parentId <= 0 的根节点并按 priority 排序后返回树形结构

#### Scenario: 父菜单不能是自身

- **WHEN** 保存菜单时 `id == parentId`
- **THEN** 系统拒绝该请求

#### Scenario: 有子菜单不可删除

- **WHEN** 对存在子菜单的菜单（`CountChildren > 0`）调用 `DeleteMenu`
- **THEN** 系统拒绝删除

#### Scenario: 删除菜单级联清理并同步 has_children

- **WHEN** 删除一个无子菜单的菜单
- **THEN** 系统同时删除其下所有权限及角色关联
- **AND** 增/删/移动菜单后调用 `SyncHasChildren` 更新父节点和自身的 `has_children` 冗余字段

### Requirement: 权限管理规则

权限 SHALL 依附于菜单，描述具体的 API 访问规则（method + uri）。权限管理 SHALL 满足：

| 规则 | 约束 |
|------|------|
| method 统一大写 | `strings.ToUpper` 规范化 |
| method 和 uri 不可空 | 基础校验 |
| 删除权限时解绑所有角色 | `DeleteByPrivilegeId` 清理关联 |
| 软删除 | 逻辑删除 `deleted=1` |

#### Scenario: 保存权限的规范化与基础校验

- **WHEN** 客户端调用 `SavePrivilege` 保存权限
- **THEN** 系统对 method 做 `strings.ToUpper` 统一大写规范化，并校验 method 与 uri 不可空

#### Scenario: 删除权限并解绑角色

- **WHEN** 客户端调用 `DeletePrivilege`
- **THEN** 系统逻辑删除该权限（`deleted=1`），并经 `DeleteByPrivilegeId` 解绑其与所有角色的关联

### Requirement: 角色-菜单/权限分配规则（原子替换）

角色-菜单与角色-权限分配 SHALL 使用 Replace 模式：先清除旧关联再插入新关联，保证原子性。分配操作 SHALL 满足：

| 规则 | 约束 |
|------|------|
| 角色存在性校验 | 分配前校验角色存在且未删除 |
| 资源存在性校验 | 批量分配菜单/权限时，逐 ID 校验是否存在 |
| Replace 语义 | `ReplaceByRoleId` 先 delete 再 batch insert |

`SaveRoleMenus` 流程：校验 roleId 有效 → 查角色，确认存在且未删除 → 若 menuIds 非空，逐 ID 校验存在性 → `ReplaceByRoleId(roleId, menuIds)` 原子替换。

#### Scenario: 原子替换角色的菜单分配

- **WHEN** 客户端以非空 menuIds 调用 `SaveRoleMenus`
- **THEN** 系统校验 roleId 有效、角色存在且未删除、逐 ID 校验菜单存在性后，先 delete 旧关联再 batch insert 新关联完成原子替换
- **AND** `SaveRolePrivileges` 对权限分配采用同一 Replace 语义

#### Scenario: 资源存在性校验失败

- **WHEN** 批量分配的菜单/权限 ID 列表中任一 ID 不存在
- **THEN** 系统拒绝本次分配请求

### Requirement: 账户-角色绑定规则

账户 SHALL 可绑定多角色，绑定关系 SHALL 使用 Replace 语义（`ReplaceByAccountId`）。绑定操作 SHALL 满足：

- 账户 ID 必须 > 0（基础校验）
- 每个角色 ID 必须 > 0 且存在（逐角色校验）
- 已删除角色不可绑定：`role.deleted = 1` 时拒绝
- Replace 替换：同一账户旧角色全部清除再插入新角色

#### Scenario: 替换账户角色绑定

- **WHEN** 客户端调用 `SaveAccountRoles` 提交 accountId 与 roleIds
- **THEN** 系统校验 accountId > 0、每个角色 ID > 0 且存在后，同一账户旧角色全部清除再插入新角色（`ReplaceByAccountId`）

#### Scenario: 已删除角色不可绑定

- **WHEN** 绑定请求中包含 `deleted = 1` 的角色
- **THEN** 系统拒绝该绑定

### Requirement: JWT 令牌签发规则

令牌签发 SHALL 使用 `pkg/auth` 工具包的 `Sign` 方法签发 JWT，配置默认值如下：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `AccessExpire` | 7200 秒（2 小时） | 访问令牌有效期 |
| `RefreshExpire` | 604800 秒（7 天） | 刷新令牌有效期 |
| `AccessSecret` | 配置项 `AccessSecret` | JWT 签名密钥 |
| `RefreshSecret` | 配置项 `RefreshSecret` | 刷新令牌独立密钥 |

`SignToken` 流程：`expireSec <= 0` 时取 `Jwt.AccessExpire`、再缺省 7200 → `authutil.Sign(userId, roleCode, accessSecret, expire)` 生成 accessToken → `authutil.Sign(userId, roleCode, refreshSecret, refreshExpire)` 生成 refreshToken → 返回 accessToken、refreshToken、expiresAt（unix 秒）。

#### Scenario: 有效期缺省回退

- **WHEN** `SignToken` 请求的 `expireSec <= 0`
- **THEN** 系统取配置 `Jwt.AccessExpire` 作为访问令牌有效期，配置再缺省时回退 7200 秒

#### Scenario: 双密钥签发访问/刷新令牌

- **WHEN** 客户端调用 `SignToken`
- **THEN** 系统以 `AccessSecret` 签发 accessToken（有效期 expireSec 或默认 7200 秒），以 `RefreshSecret` 签发 refreshToken（有效期 604800 秒）
- **AND** 返回 `expiresAt` 为 unix 秒

### Requirement: 登录记录规则

登录记录 SHALL 满足以下约束：

| 规则 | 约束 |
|------|------|
| 容错插入 | 插入失败只 log 错误，不抛异常（`l.Errorf`） |
| 按 userId 分页 | 登录记录按用户 ID 分组查询（`FindPage`） |

#### Scenario: 登录记录容错插入

- **WHEN** `SaveLoginRecord` 写入 `login_record` 失败
- **THEN** 系统仅记录错误日志（`l.Errorf`），不向调用方抛异常，登录主流程不受影响

#### Scenario: 按用户分页查询登录记录

- **WHEN** 客户端以 `LoginRecordPageReq { pageNo, pageSize, userId }` 调用 `ListLoginRecords`
- **THEN** 系统按 userId 分组分页返回 `LoginRecordListReply { total, list }`

### Requirement: 数据模型与存储不变量

auth 服务 SHALL 将数据存储于 MySQL `tj_auth` 库（DDL 来源 `sql/ddl/tj_auth.sql`）的七张表：

| 表 | 用途 | 关键字段 |
|----|------|---------|
| `role` | 角色表 | `id`（主键）、`code`（varchar(64)，角色代号 admin/teacher/student 等，唯一约束）、`name`、`type`（tinyint，0=固定角色(不可改), 1=自定义角色）、`create_time`/`update_time`（自动填充/更新）、`creater`/`updater`、`dep_id`、`deleted`（0=正常, 1=逻辑删除，过滤条件） |
| `menu` | 菜单表（树形结构） | `id`（主键）、`parent_id`（父菜单 ID，0=一级菜单，树形查询）、`has_children`（tinyint，是否有子菜单，缓存用冗余字段）、`label`（varchar(16)）、`path`、`icon`（varchar(32)）、`priority`（tinyint，排序值默认 127）、`deleted`（逻辑删除） |
| `privilege` | 权限表（API 访问权限） | `id`（主键）、`menu_id`（外键关联 menu）、`intro`（权限说明）、`method`（varchar(16)，HTTP 方法）、`uri`（API 路径）、`internal`（是否内部接口）、`deleted`（逻辑删除） |
| `role_menu` | 角色-菜单关联表 | `id`（自增主键）、`role_id`（FK → role）、`menu_id`（FK → menu）；一对多：一个角色可绑定多个菜单 |
| `role_privilege` | 角色-权限关联表 | `id`（自增主键）、`role_id`（FK → role）、`privilege_id`（FK → privilege）；一对多：一个角色可拥有多个权限 |
| `account_role` | 账户-角色关联表 | `id`（自增主键）、`account_id`（FK → user.account）、`role_id`（FK → role）；多对多：一个账户可有多角色，一个角色可分配给多账户 |
| `login_record` | 登录记录表 | `id`（自增主键）、`user_id`（过滤条件）、`cell_phone`（varchar(11)）、`login_time`（自动填充）、`logout_time`（可空）、`login_date`（date，索引）、`duration`（登录时长秒）、`ipv4`（varchar(15)） |

关系结构：`role (1) ─ (N) role_menu ─ (N) menu`；`role (1) ─ (N) role_privilege ─ (N) privilege`；`account (1) ─ (N) account_role ─ (N) role`；`login_record → user`（外键在 user 域）。

存储与模型不变量：

- `role.code` 具有唯一约束，是唯一业务键；角色删除为软删除（`deleted=1`），所有查询以 `deleted` 为过滤条件
- `menu` 通过 `parent_id` 自引用构建树，`has_children` 为便于前端展示的冗余字段，由 `SyncHasChildren` 维护一致性
- 关联表（`role_menu` / `role_privilege` / `account_role`）无逻辑删除位，一致性由 Replace 语义与级联清理（`DeleteByRoleId` / `DeleteByMenuId` / `DeleteByPrivilegeId`）维护
- 模型扩展模式：所有 `*_gen.go` 由 goctl 自动生成**禁止修改**，扩展方法统一放同名自定义 `.go` 文件——`rolemodel.go`（FindPage, ExistsByCode, SoftDelete）、`menumodel.go`（SyncHasChildren, FindByIds, CountChildren）、`roleprivilegemodel.go`（DeleteByRoleId, DeleteByPrivilegeId）、`rolemenumodel.go`（DeleteByRoleId, DeleteByMenuId, FindMenuIdsByRoleId）、`accountrolemodel.go`（ReplaceByAccountId, FindRoleIdsByAccountId, CountByRoleId）、`privilegemodel.go`（FindByMenuId, FindByIds, SoftDelete）、`loginrecordmodel.go`（FindPage，按 userId 分页）
- `sqlhelper.go` 提供通用 SQL 工具函数：`inPlaceholders`、`pairValuePlaceholders`、`dedupeIds`

#### Scenario: 角色 code 唯一约束与软删除

- **WHEN** 新增角色的 code 与库中既有角色重复
- **THEN** 系统经 `ExistsByCode` 在应用层拒绝（唯一约束兜底于 DDL）；删除角色时仅置 `deleted=1`，后续查询自动过滤已删除记录

#### Scenario: 菜单树自引用与 has_children 冗余

- **WHEN** 菜单发生增/删/移动
- **THEN** `parent_id` 自引用结构随之更新，`SyncHasChildren` 同步父节点与自身的 `has_children`（0=无子菜单, 1=有子菜单）冗余字段

#### Scenario: goctl 生成文件禁止修改

- **WHEN** 需要为任一模型新增数据访问方法
- **THEN** 扩展方法必须写入同名自定义 `.go` 文件（如 `rolemodel.go` 扩展 `RoleModel` 接口），`*_gen.go` 保持 goctl 生成原样

### Requirement: 服务配置与端口

auth 服务 SHALL 按以下默认配置部署（API：`apps/auth/api/etc/auth-api.yaml`；RPC：`apps/auth/rpc/etc/auth.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| API `Name` / `Host` / `Port` | `auth-api` / `0.0.0.0` / `8802` | HTTP 监听 |
| API `Auth.AccessSecret` | `change-me-in-production` | JWT 签名密钥 |
| API `Auth.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| API `AuthRpc.Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `auth.rpc` | 自身 RPC 服务发现 |
| API `UserRpc.Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `user.rpc` | user RPC 服务发现（登录验证凭证） |
| RPC `Name` / `ListenOn` | `auth.rpc` / `0.0.0.0:8082` | gRPC 监听 |
| RPC `Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `auth.rpc` | RPC 注册 |
| RPC `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_auth?charset=utf8mb4&parseTime=true&loc=Local` | MySQL `tj_auth` 库（utf8mb4，时区 Local） |
| RPC `Cache[0]` | `127.0.0.1:6379`，Type `node`，Pass 空 | Redis 单机缓存（缓存角色/菜单/权限数据） |
| RPC `Jwt.AccessSecret` | `change-me-in-production` | JWT 签名密钥 |
| RPC `Jwt.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| RPC `Jwt.RefreshSecret` | `change-me-refresh-in-production` | 刷新令牌密钥 |
| RPC `Jwt.RefreshExpire` | `604800` | 刷新令牌有效期（秒）= 7 天 |

约束：

- 生产环境**必须修改**默认 JWT 密钥，否则 JWT 可被伪造
- `AccessSecret` 与 `RefreshSecret` 两个密钥应不同，实现令牌的独立轮换：`AccessSecret` 签发和校验 accessToken，`RefreshSecret` 签发和校验 refreshToken
- 依赖的外部服务：MySQL `tj_auth` 库（DataSource 配置，自建存储）、Redis（Cache 配置，缓存角色/菜单/权限数据）、`user.rpc`（RpcClient 配置，登录时验证用户凭证 `LoginVerify`）

#### Scenario: 登录链路跨服务协作

- **WHEN** 客户端调用 `auth-api` 的登录接口（`POST /accounts/login` 等）
- **THEN** HTTP handler 按 `UserRpc` 配置（etcd `127.0.0.1:2379`，key `user.rpc`）调用 user 服务 `LoginVerify` 核验凭证，成功后经 `AuthRpc`（key `auth.rpc`）调自身 `SignToken` 签发令牌

#### Scenario: 生产环境必须替换默认 JWT 密钥

- **WHEN** auth 服务以默认 `AccessSecret = change-me-in-production` / `RefreshSecret = change-me-refresh-in-production` 上线生产
- **THEN** JWT 可被伪造，生产环境必须替换两个默认密钥
- **AND** 两个密钥应设置为不同值，实现访问/刷新令牌的独立轮换

