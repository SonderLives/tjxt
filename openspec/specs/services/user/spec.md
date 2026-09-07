# user（用户中心）Specification

## Purpose

user（用户中心）是身份域的基础微服务，职责边界为「身份核验 + 用户信息管理」：登录入口与 JWT 签发统一由 auth 服务负责，本服务只找用户、比对密码、检查状态。HTTP API 层 `user-api` 监听 `0.0.0.0:8801`，RPC 服务 `User` 监听 `0.0.0.0:8081`（gRPC），经 etcd 以 key `user.rpc` 注册与发现。数据存储于 MySQL `tj_user` 库（`user` / `user_detail` 两张表，共享主键 1:1），主键/手机号/用户名查询由 Redis 缓存托管。RPC 消费方为 `user-api` 自身（用户管理 HTTP 接口全部落地自身 RPC）与 `auth-api`（登录时调 `LoginVerify` 核验凭证）。

## Requirements

### Requirement: 用户信息管理 HTTP 接口

> **状态说明**：下表来自聚合文档 `docs/tjxt.openapi.json`（api-spec.md，最后同步 2026-08-05）。源表「认证」列均标注「否」，为原始聚合文档标注，本规格按源文件忠实转写；所有 HTTP 接口最终经 `user.rpc` 落地。

user 服务 SHALL 暴露以下用户信息管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /users | 新增用户，一般是员工或教师 | UserDTO | R{data: R«long»} |
| PUT | /users | 更新当前登录用户信息，可修改密码 | UserFormDTO | R{data: R} |
| GET | /users/checkCellphone | 检查用户手机号是否存在 | - | R{data: R«boolean»} |
| GET | /users/me | 获取当前登录用户信息 | - | R{data: R«UserDetailVO»} |
| GET | /users/{id} | 根据id查询用户信息 | - | R{data: R«UserDTO»} |
| PUT | /users/{id} | 更新用户信息 | UserDTO | R{data: R} |

统一约定（引用全局规范）：响应格式为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页为 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码见共享错误码表（`specs/shared/error-codes`）。

#### Scenario: 新增用户

- **WHEN** 管理端以 `UserDTO` 调用 `POST /users`
- **THEN** 系统新增用户（一般是员工或教师），密码固定为默认初始口令（规则见「管理端新增用户规则」）
- **AND** 返回 `R{data: R«long»}`（新用户 ID）

#### Scenario: 获取当前登录用户信息

- **WHEN** 当前登录用户调用 `GET /users/me`
- **THEN** 系统返回当前登录用户详情 `R{data: R«UserDetailVO»}`

#### Scenario: 检查手机号是否已注册

- **WHEN** 客户端调用 `GET /users/checkCellphone` 并携带手机号
- **THEN** 系统返回该手机号是否已被注册 `R{data: R«boolean»}`

### Requirement: 密码与状态管理 HTTP 接口

> **状态说明**：同上，源自 `docs/tjxt.openapi.json`（api-spec.md）。

user 服务 SHALL 暴露以下密码与状态管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| PUT | /users/{id}/password/default | 重置密码 | - | R{data: R} |
| PUT | /users/{id}/status/{status} | 修改用户状态, status=0为禁用，status=1为正常 | - | R{data: R} |

#### Scenario: 重置密码为默认初始口令

- **WHEN** 管理端调用 `PUT /users/{id}/password/default`
- **THEN** 系统将该用户密码重置为默认初始口令（`defaultInitialPassword`，常量 `"123456"`，规则见「密码管理规则」）
- **AND** 返回 `R{data: R}`

#### Scenario: 修改用户状态

- **WHEN** 管理端调用 `PUT /users/{id}/status/{status}`
- **THEN** 系统按路径参数修改账户状态：`status=0` 为禁用，`status=1` 为正常（白名单约束见「账户状态管理规则」）
- **AND** 账户状态变更后同步失效对应缓存键，避免登录核验读到旧状态

### Requirement: RPC 服务契约（User 服务）

`User` 服务 SHALL 作为用户中心微服务经 etcd 服务发现注册（key: `user.rpc`），监听 `0.0.0.0:8081`（gRPC），只负责「身份核验 + 用户信息管理」，登录入口与 JWT 签发统一由 auth 服务负责。

**身份核验**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `LoginVerify` | `LoginVerifyRequest { cell_phone, username, password, type }` | `LoginVerifyResponse { user_id, username, cell_phone, type, status, name }` | 供 auth 服务调用，只核验凭证不签发令牌 |

`LoginVerifyRequest` 关键字段：`cell_phone`（学员/老师用手机号，优先按手机号 + type 定位）、`username`（员工用用户名，手机号为空时才使用）、`password`（明文密码，服务端用 bcrypt 比对）、`type`（1-员工 2-学员 3-老师）。

**学员账户**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `RegisterStudent` | `StudentFormRequest { cell_phone, code, password }` | `IdResponse { id }` | 学员注册，写 user + user_detail |
| `UpdateStudentPassword` | `StudentFormRequest { cell_phone, code, password }` | `EmptyResponse {}` | 学员改密，按手机号 + type=2 定位 |

`StudentFormRequest` 关键字段：`cell_phone`（必填）、`code`（短信验证码，proto 已定义，logic 中未校验）、`password`（必填，落库前 bcrypt 加密）。

**用户信息**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `AddUser` | `UserDTO` | `IdResponse { id }` | 管理端新增用户，密码固定为默认初始口令 |
| `GetUserById` | `UserIdRequest { user_id }` | `UserDTO` | 按 ID 查用户（资料缺失时退化为基础信息） |
| `GetUsersByIds` | `UserIdsRequest { user_ids }` | `UserListResponse { list }` | 批量查用户，联合视图避免 N+1 |
| `GetUserDetail` | `UserIdRequest { user_id }` | `UserDetailVO` | 当前登录用户详情 |
| `UpdateUserById` | `UserDTO` | `EmptyResponse {}` | 管理端更新指定用户，仅覆盖非零字段 |
| `UpdateCurrentUser` | `UserFormRequest` | `EmptyResponse {}` | 更新当前登录用户，可同时改密码 |
| `CheckCellPhone` | `CheckCellPhoneRequest { cell_phone }` | `BoolResponse { result }` | 检查手机号是否已被注册 |
| `ResetPassword` | `UserIdRequest { user_id }` | `EmptyResponse {}` | 重置为默认初始口令 |
| `UpdateUserStatus` | `UpdateStatusRequest { user_id, status, operator }` | `EmptyResponse {}` | 启用/禁用账户并失效缓存 |

`UserDTO` 关键字段：`id`（int64）、`username`、`cell_phone`、`type`（1-其他员工, 2-普通学员, 3-老师）、`name`、`gender`（0-男性, 1-女性）、`icon`、`email`、`qq`、`job`、`province`/`city`/`district`、`intro`、`photo`、`role_id`（角色 ID，老师和学生不用填）。`UserFormRequest` 相对 `UserDTO` 增量字段：`old_password`（原始密码，携带新密码时必填）、`password`（新密码）。`UpdateStatusRequest` 关键字段：`user_id`、`status`（0-禁用 1-正常，其它取值拒绝）、`operator`（操作人 ID，写入 `updater`）。

**管理后台分页查询**：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `PageQueryStudents` | `UserPageRequest` | `StudentPageResponse { total, pages, list }` | 学员分页（固定 type=2） |
| `PageQueryTeachers` | `UserPageRequest` | `TeacherPageResponse { total, pages, list }` | 老师分页（固定 type=3） |
| `PageQueryStaffs` | `UserPageRequest` | `StaffPageResponse { total, pages, list }` | 员工分页（固定 type=1） |

`UserPageRequest` 关键字段：`page_no`（< 1 归一化为 1）、`page_size`（缺省 10，上限 100）、`sort_by`（排序白名单：name / status / id，其它落回 create_time）、`is_asc`（true 升序，false 降序）、`name`（姓名模糊匹配 `d.name like %?%`）、`phone`（手机号模糊匹配 `u.cell_phone like %?%`）、`status`（负值表示不过滤）。响应 VO 独有字段：`StudentPageVO`（gender, course_amount）、`TeacherPageVO`（photo, job, intro, course_amount, exam_question_amount）、`StaffVO`（role_id, role_name）。

消费方 SHALL 为 `user-api` 自身 API 层（`apps/user/api/internal/svc/servicecontext.go` import `userclient "tjxt/apps/user/rpc/client/user"`，用户管理全部 HTTP 接口最终走自身 RPC）与 `auth-api`（`apps/auth/api/internal/svc/servicecontext.go` import 同一 client，登录接口 `apps/auth/api/internal/logic/account/loginlogic.go:32` 调用 `LoginVerify` 核验凭证）。user 服务是身份域的基础设施，`user.rpc` 只被上述两个 API 服务在 `svc/servicecontext.go` 中显式装配。

#### Scenario: 用户登录核验

- **WHEN** `auth-api` 登录接口调用 `User.LoginVerify` 传入手机号或用户名 + 密码
- **THEN** user 服务只核验凭证不签发令牌，核验通过后返回 userId / username / cellPhone / type / status，由 auth 服务签发 JWT

#### Scenario: 跨服务用户信息聚合

- **WHEN** 其它域拿到 userId 列表后调用 `GetUsersByIds`
- **THEN** 系统以单条 SQL `LEFT JOIN user_detail` 一次性回填昵称头像等资料，避免 N+1

#### Scenario: 账户封禁

- **WHEN** 管理员调用 `UpdateUserStatus(status=0)`
- **THEN** 模型层同步失效 id / cellPhone+type / username 三个缓存键

> 已知缺口：`LoginVerifyResponse.name` 字段在 proto 中已定义，但 logic 未回填（未查 `user_detail`），实际恒为空串。

### Requirement: 身份核验规则（LoginVerify）

user 服务 SHALL 在 `LoginVerify` 中只做「找用户 + 比对密码 + 状态检查」，**不签发任何令牌**；角色归属由 auth 服务通过 `account_role` 自查。具体约束：

- 定位优先级：`cell_phone` 非空优先按「手机号 + type」定位（`FindOneByCellPhoneType`），否则按 `username` 定位（`FindOneByUsername`）
- 两者皆空：直接返回 `ErrUserNotFound`，不查库
- 状态前置：`status != 1` 返回 `ErrUserDisabled`，**先于密码比对**，禁用账户不泄漏密码正确性
- 密码比对：`bcrypt.CompareHashAndPassword`，失败返回 `ErrBadCredential`
- 错误类型：使用包级哨兵错误（`ErrUserNotFound` / `ErrBadCredential` / `ErrUserDisabled`）而非 `xerr`，由调用方 auth 服务翻译为业务错误码

#### Scenario: 登录核验成功路径

- **WHEN** `LoginVerify` 收到 `cell_phone` 非空的请求
- **THEN** 系统按 `FindOneByCellPhoneType(cell_phone, type)` 定位用户，记录不存在时返回 `ErrUserNotFound`
- **AND** 用户 `status=1` 且 bcrypt 比对通过后，返回 userId / username / cellPhone / type / status；密码比对失败返回 `ErrBadCredential`

#### Scenario: 禁用账户先于密码比对

- **WHEN** 定位到的用户 `status != 1`
- **THEN** 系统返回 `ErrUserDisabled`，不再进行密码比对
- **AND** 禁用账户不泄漏密码正确性

#### Scenario: 标识皆空短路

- **WHEN** `cell_phone` 与 `username` 均为空
- **THEN** 系统直接返回 `ErrUserNotFound`，不查库

### Requirement: 学员注册规则（RegisterStudent）

user 服务 SHALL 按以下规则处理学员注册：

- 必填校验：`cell_phone` 与 `password` 均不可为空，否则 `BadRequest`
- 手机号查重：`ExistsByCellPhone` 全局查重（不区分 type），已存在返回 `Conflict("该手机号已注册")`
- username 兜底：注册时 `username` 直接取 `cell_phone` 值
- 密码加密：`bcrypt.GenerateFromPassword(..., bcrypt.DefaultCost)`
- ID 生成：`idgen.NextID()` 雪花 ID，`user` 与 `user_detail` 共用同一 ID
- 类型固定：`type = 2`（普通学员），`status = 1`（正常）
- 审计字段：`creater` / `updater` 均写入自身 ID（自注册）
- 验证码：`code` 字段接收但 **logic 未做任何校验**，短信验证留待接入

#### Scenario: 学员注册成功

- **WHEN** 学员以未注册的手机号与密码调用 `RegisterStudent`
- **THEN** 系统通过必填校验与全局手机号查重后，以雪花 ID 先写 `user`（username=cell_phone, type=2, status=1，密码 bcrypt 加密），再写 `user_detail`（仅骨架，资料字段留空），两表共用同一 ID
- **AND** `creater` / `updater` 均写入学员自身 ID

#### Scenario: 手机号已注册

- **WHEN** 提交的手机号在任意用户类型中已存在（`ExistsByCellPhone` 不区分 type）
- **THEN** 系统返回 `Conflict("该手机号已注册")`

#### Scenario: 验证码字段不校验

- **WHEN** 注册请求携带 `code` 字段
- **THEN** logic 层接收该字段但不做任何校验，短信验证留待接入
- **AND** 该接口必须由上层 API 完成验证码校验后才可暴露

> 已知缺口（一致性风险）：写 `user` 与写 `user_detail` 两次 Insert 未包裹事务，`user` 写成功而 `user_detail` 写失败会产生孤儿用户；`GetUserById` 对此做了退化兜底（资料缺失时只返回基础信息）。

### Requirement: 管理端新增用户规则（AddUser）

user 服务 SHALL 按以下规则处理管理端新增用户：

- 至少一个标识：`cell_phone` 与 `username` 不能同时为空（否则 `BadRequest`）
- 手机号查重：`cell_phone` 非空时走 `ExistsByCellPhone`，冲突返回 `Conflict("该手机号已存在")`
- 用户名查重：`username` 非空时走 `FindOneByUsername`，查到即 `Conflict("该用户名已存在")`；非 `ErrNotFound` 的错误按内部错误处理
- 类型缺省：`type == 0` 时兜底为 `2`（普通学员）
- 密码策略：**密码不随请求传入**，统一使用 `defaultInitialPassword`（常量 `"123456"`）
- 资料落库：空串字段经 `nullStr` 转为 SQL NULL，不写空串

#### Scenario: 管理端新增用户成功

- **WHEN** 管理端提交至少含 `cell_phone` 或 `username` 之一的用户信息
- **THEN** 系统完成手机号 / 用户名查重，`type` 为 0 时兜底为 2，以 bcrypt 加密默认初始口令后写 `user` 与 `user_detail`（共用雪花 ID）
- **AND** 空串资料字段经 `nullStr` 转为 SQL NULL 落库，不写空串

#### Scenario: 标识皆空拒绝

- **WHEN** `cell_phone` 与 `username` 都为空
- **THEN** 系统返回 `BadRequest`

> 已知缺口（安全提示）：`defaultInitialPassword = "123456"` 是硬编码常量，源码注释明确「正式环境应改为短信或随机下发」。

### Requirement: 用户信息更新规则（UpdateUserById / UpdateCurrentUser）

`UpdateUserById` 与 `UpdateCurrentUser` SHALL 共用同一套「非零字段合并」语义，仅覆盖请求中的非零值，避免误清零未传字段。具体约束：

- 用户必须存在：`FindOne` 失败且为 `ErrNotFound` → `NotFound("用户不存在")`
- 资料可缺失：`user_detail` 查不到时**就地构造**空壳 `{Id, Type: u.Type, Creater: in.Id}` 后继续
- 字符串合并：走 `applyStr(cur, next)`：`next` 为空则保留原值，否则覆盖
- 数值合并：`Name` / `Gender` / `Type` / `RoleId` 判 `!= 0` / `!= ""` 才覆盖
- 审计字段：`u.Updater` 与 `d.Updater` 均写入 `in.Id`
- 写入顺序：先 `UserModel.Update`，后 `UserDetailModel.Update`，任一失败即返回

`UpdateCurrentUser` SHALL 独有改密规则：

- 原密码必填：`password != "" && old_password == ""` → `BadRequest("修改密码需提供原密码")`
- 原密码校验：bcrypt 比对失败 → `Unauthorized("原密码错误")`
- 校验时机：在两次 `Update` **之前**完成，校验不通过不落任何库

#### Scenario: 非零字段合并

- **WHEN** 更新请求中某字符串字段为空串或数值字段为 0
- **THEN** 系统保留该字段原值不覆盖，避免误清零未传字段

#### Scenario: 当前用户改密

- **WHEN** `UpdateCurrentUser` 携带新密码（`password` 非空）
- **THEN** 系统要求 `old_password` 非空（否则 `BadRequest("修改密码需提供原密码")`），并以 bcrypt 比对原密码（失败返回 `Unauthorized("原密码错误")`）
- **AND** 校验在两次 `Update` 之前完成，校验不通过不落任何库；通过后重新 bcrypt 生成 `u.Password`，再依次更新 user 与 user_detail

#### Scenario: 性别无法从女性改回男性

- **WHEN** 试图通过更新接口把性别从「女性(1)」改回「男性(0)」
- **THEN** 合并逻辑 `if in.Gender != 0` 将 0 视为未传，性别保持女性——该修改无法通过更新接口完成

> 已知缺口（一致性风险）：两次 `Update` 未包裹事务；且 `UpdateUserById` 用 `in.Id` 而非操作人 ID 填 `updater`，实际记录的是被改用户自身 ID。

### Requirement: 密码管理规则

user 服务 SHALL 按以下三种场景管理密码：

| 场景 | 方法 | 规则 |
|------|------|------|
| 学员自助改密 | `UpdateStudentPassword` | 按「手机号 + type=2」定位；**不校验原密码、不校验短信验证码**；`updater` 写用户自身 ID |
| 管理员重置 | `ResetPassword` | 重置为 `defaultInitialPassword`；`updater` 写 `in.UserId`（即被重置者自身） |
| 当前用户改密 | `UpdateCurrentUser` | 必须提供并通过 `old_password` 校验（见「用户信息更新规则」） |

#### Scenario: 学员自助改密

- **WHEN** 学员以手机号 + 新密码调用 `UpdateStudentPassword`
- **THEN** 系统校验 `cell_phone` / `password` 非空，按 `FindOneByCellPhoneType(cell_phone, 2)` 定位用户（不存在则 `NotFound`），bcrypt 加密新密码后更新 `user`
- **AND** 全程不校验原密码、不校验短信验证码，`updater` 写用户自身 ID

#### Scenario: 管理员重置密码

- **WHEN** 管理员调用 `ResetPassword`
- **THEN** 系统将该用户密码重置为默认初始口令 `defaultInitialPassword`（`"123456"`）

> 已知缺口（安全提示）：`UpdateStudentPassword` 只凭手机号即可重设密码，`StudentFormRequest.code`（短信验证码）在 logic 中**完全未被读取**。该接口必须由上层 API 完成验证码校验后才可暴露。

### Requirement: 账户状态管理规则（UpdateUserStatus）

user 服务 SHALL 按以下规则管理账户状态：

- 状态白名单：仅接受 0（禁用）与 1（正常），其它取值 → `BadRequest("状态取值非法")`
- 存在性校验：由模型层 `UpdateStatus` 内部先 `FindOne`，不存在则冒泡 `ErrNotFound` → `NotFound`
- 操作人记录：`in.Operator` 写入 `updater` 列（**本服务唯一正确记录操作人的接口**）
- 更新方式：`ExecNoCacheCtx` 直接 UPDATE，同时刷新 `update_time = now()`

`UpdateStatus` 绕过了 goctl 的缓存写通道，SHALL 手动清理三个缓存键（`cache:tjUser:user:id:{id}`、`cache:tjUser:user:cellPhone:type:{cellPhone}:{type}`、`cache:tjUser:user:username:{username}`），否则 `LoginVerify` 会读到旧的 `status`。

#### Scenario: 禁用账户并失效缓存

- **WHEN** 管理员调用 `UpdateUserStatus` 且 `status=0`（合法取值）
- **THEN** 模型层先 `FindOne(id)` 取旧记录拿到 CellPhone / Type / Username 用于拼键，再 `UPDATE user SET status, updater, update_time`（`ExecNoCacheCtx` 直改 SQL）
- **AND** 依次删除 `cache:tjUser:user:id:{id}`、`cache:tjUser:user:cellPhone:type:{cellPhone}:{type}`、`cache:tjUser:user:username:{username}` 三个缓存键

#### Scenario: 非法状态值拒绝

- **WHEN** `status` 取值不是 0 或 1
- **THEN** 系统返回 `BadRequest("状态取值非法")`

#### Scenario: 缓存键漏删导致旧状态

- **WHEN** 三个缓存键中任一个未被删除
- **THEN** 对应查询路径（按 ID / 按手机号登录 / 按用户名登录）将继续命中旧状态
- **AND** 因此三个键必须**全部**删除

### Requirement: 查询与分页规则

**单体与批量查询**，user 服务 SHALL 满足：

- 资料退化：`GetUserById` / `GetUserDetail` 中 `user_detail` 缺失不报错，`toUserDTO(u, nil)` 只填基础字段
- 空入参短路：`GetUsersByIds` 收到空 ID 列表直接返回空 `UserListResponse`，不查库
- 批量走视图：`FindByIdsWithDetail` 单条 SQL `LEFT JOIN` 完成，杜绝 N+1
- NULL 降级：`strVal` 将 `sql.NullString` 无效值统一转空串
- 时间格式：`create_time` 统一格式化为 `2006-01-02 15:04:05`

**分页查询**：三个分页方法共用 `FindPageByType`，仅 `userType` 常量不同（`PageQueryStaffs` 固定 type=1、`PageQueryStudents` 固定 type=2、`PageQueryTeachers` 固定 type=3），SHALL 满足：

- 参数归一化：`page.Normalize`：`pageNo < 1` → 1；`pageSize < 1` → 10；`pageSize > 100` → 100
- 总页数：`page.CalcPages(total, limit)`，`total <= 0` 时返回 0
- 条件拼接：`name` / `phone` 走 `like %?%`；`status >= 0` 才追加状态过滤
- 先查列表后查总数：两条 SQL 复用同一组 `args`，列表额外追加 `offset, limit`
- SQL 注入防护：排序字段**不允许透传**，`sortClause` 做白名单映射后才交给模型层：`"name"` → `d.name`、`"status"` → `u.status`、`"id"` → `u.id`，其它/空 → `u.create_time`（默认按注册时间）；方向 `isAsc ? "ASC" : "DESC"`（默认 DESC）

#### Scenario: 分页参数归一化

- **WHEN** 客户端传 `pageNo=0` 或 `pageSize=200`
- **THEN** 系统归一化为 `pageNo=1`、`pageSize=100`（`pageSize < 1` 归一化为 10）

#### Scenario: 排序字段白名单映射

- **WHEN** `sortBy` 传入白名单（name / status / id）之外的值或为空
- **THEN** 排序落回 `u.create_time`（默认按注册时间），方向默认 DESC，排序字段不允许透传以防止 SQL 注入

#### Scenario: 批量查询空入参短路

- **WHEN** `GetUsersByIds` 收到空 ID 列表
- **THEN** 系统直接返回空 `UserListResponse`，不查库

> 已知缺口（未实现的跨服务聚合，源码注释明确标注，当前返回本地值或空值）：`UserDetailVO.role_name`、`StaffVO.role_name` — 归属 auth 服务 `role` 表，暂留空；`TeacherPageVO.exam_question_amount` — 归属 exam 服务，暂留 0；`TeacherPageVO.course_amount` — 当前取 `user_detail.course_amount`（学员购课数），权威值应由 course 服务计算。

### Requirement: 数据模型与存储不变量

user 服务 SHALL 将用户数据存储于 MySQL `tj_user` 库（DDL 来源 `sql/ddl/tj_user.sql`）的两张表：

| 表 | 用途 | 关键字段 |
|----|------|---------|
| `user` | 用户账号表（学员/老师/员工共用） | `id`（主键，雪花 ID 非自增）、`username`（唯一索引 `username_idx`）、`cell_phone`（联合唯一索引 `cell_idx`(cell_phone, type)）、`password`（bcrypt 密文）、`type`（1-员工, 2-普通学员, 3-老师，默认 0，普通索引 `type_idx`）、`status`（0-禁用 1-正常，默认 1）、`create_time`/`update_time`（自动填充）、`creater`（可空）/`updater`（默认 0） |
| `user_detail` | 全类型用户资料表（DDL 注释为教师详情表，实为全类型资料） | `id`（与 `user.id` 同值）、`type`（默认 2）、`name`（默认空串，全文索引 `name_idx`）、`gender`（0-男性, 1-女性，默认 0）、`icon`/`email`/`qq`/`birthday`/`job`/`province`/`city`/`district`/`intro`/`photo`（可空）、`role_id`（非空，跨库引用 auth 域 `tj_auth.role`）、`course_amount`（购买课程数量，学生才有该字段信息，默认 0）、`dep_id`（默认 0） |

存储不变量：

- 手机号唯一性由 `(cell_phone, type)` 联合约束保证，即同一手机号可分别注册为学员与老师；但 `ExistsByCellPhone` 在应用层做的是不分类型的全局查重，比 DDL 约束更严格
- `user_detail.id` 与 `user.id` 为 1:1 共享主键，无外键约束，由应用层在 `AddUser` / `RegisterStudent` 中同事务外顺序写入两张表
- `birthday` 与 `dep_id` 在 proto 与 logic 中均未被读写，属于表结构预留字段
- 分页与批量查询使用非表结构投影 `UserWithDetail`（`user` LEFT JOIN `user_detail`）：`LEFT JOIN` 保证资料缺失的用户仍能被查出，NULL 列由 `strVal` 统一降级为空串
- 跨库引用：`user_detail.role_id` → auth 域 `tj_auth.role`；`user.id` ← auth 域 `tj_auth.login_record.user_id`
- 缓存键三组（goctl 生成）：`cache:tjUser:user:id:`（`FindOne`）、`cache:tjUser:user:cellPhone:type:`（`FindOneByCellPhoneType`）、`cache:tjUser:user:username:`（`FindOneByUsername`）；`UpdateStatus` 走 `ExecNoCacheCtx` 直改 SQL，需手动 `DelCacheCtx` 依次失效三个键
- 模型扩展模式：`*_gen.go` 由 goctl 生成**禁止修改**，扩展方法统一放同名自定义 `.go` 文件——`usermodel.go` 扩展 `FindPageByType` / `FindByIdsWithDetail` / `ExistsByCellPhone` / `UpdateStatus`；`userdetailmodel.go` 未新增扩展方法，仅内嵌 goctl 生成接口；`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound` 供 logic 层 `errors.Is` 判定「记录不存在」

#### Scenario: 同一手机号分类型注册的约束差异

- **WHEN** 同一手机号尝试注册第二种用户类型（如已注册学员再注册老师）
- **THEN** DDL 层面 `(cell_phone, type)` 联合唯一索引允许两条不同 type 的记录
- **AND** 应用层 `ExistsByCellPhone` 不分类型全局查重，比 DDL 约束更严格，将拦截该注册并返回冲突

#### Scenario: 资料缺失用户仍可查出

- **WHEN** `user` 记录存在而 `user_detail` 缺失（如注册产生的孤儿用户）
- **THEN** `LEFT JOIN` 联合视图仍能查出该用户，NULL 资料列由 `strVal` 降级为空串
- **AND** `GetUserById` 退化为基础信息返回，不报错

### Requirement: 服务配置与端口

user 服务 SHALL 按以下默认配置部署（API：`apps/user/api/etc/user-api.yaml`；RPC：`apps/user/rpc/etc/user.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| API `Name` / `Host` / `Port` | `user-api` / `0.0.0.0` / `8801` | HTTP 监听 |
| API `Auth.AccessSecret` | `change-me-in-production` | JWT 签名密钥 |
| API `Auth.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| API `UserRpc.Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `user.rpc` | 自身 RPC 服务发现 |
| RPC `Name` / `ListenOn` | `user.rpc` / `0.0.0.0:8081` | gRPC 监听 |
| RPC `Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `user.rpc` | RPC 注册 |
| RPC `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_user?charset=utf8mb4&parseTime=true&loc=Local` | MySQL `tj_user` 库（utf8mb4，时区 Local） |
| RPC `Cache[0]` | `127.0.0.1:6379`，Type `node`，Pass 空 | Redis 单机缓存 |

约束：

- `parseTime=true` 是必需项：`UserWithDetail.CreateTime` 声明为 `time.Time`，需要驱动直接把 `datetime` 解析为时间类型
- RPC 配置结构体（`apps/user/rpc/internal/config/config.go`）内嵌 `zrpc.RpcServerConf`，外加 `DataSource` 与 `Cache`；RPC 层**不含 Jwt 配置**——令牌签发是 auth 服务的职责
- API 配置结构体（`apps/user/api/internal/config/config.go`）内嵌 `rest.RestConf`，`Auth` 为匿名结构体，`UserRpc` 为 `zrpc.RpcClientConf`
- 缓存由 `svc.NewServiceContext` 一次性传给两个模型；`UserDetailModel` 缓存键由 goctl 按 `user_detail` 主键生成
- `user.rpc` 作为被依赖方，在 `apps/user/api/etc/user-api.yaml`（配置节 `UserRpc`）与 `apps/auth/api/etc/auth-api.yaml`（配置节 `UserRpc`）中均声明为 key `user.rpc` 的客户端
- 可观测性注入配置（**已注入勿删**，configs.md 未记录、以实际 yaml 为准）：`Log`、`Prometheus`、`Telemetry` 已注入 `user-api.yaml` 与 `user.yaml`，本规格不展开

#### Scenario: API 层经 etcd 调用自身 RPC

- **WHEN** 客户端调用 `user-api` 的任一用户管理 HTTP 接口
- **THEN** HTTP handler 按 `UserRpc` 配置（etcd `127.0.0.1:2379`，key `user.rpc`）经服务发现调用自身 RPC 对应方法

#### Scenario: 生产环境必须替换默认 JWT 密钥

- **WHEN** `user-api` 以默认 `Auth.AccessSecret = change-me-in-production` 上线生产
- **THEN** JWT 可被伪造，生产环境必须替换该默认值
- **AND** `AccessExpire` 默认 7200 秒

