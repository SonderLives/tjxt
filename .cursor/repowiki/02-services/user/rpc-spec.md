> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/user/rpc/user.proto`

---

# User RPC Spec

## 服务名

`User` — 用户中心微服务，通过 etcd 服务发现（key: `user.rpc`）。

职责边界：只负责「身份核验 + 用户信息管理」，登录入口与 JWT 签发统一由 auth 服务负责。

## RPC 方法总览

### 身份核验

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `LoginVerify` | `LoginVerifyRequest { cell_phone, username, password, type }` | `LoginVerifyResponse { user_id, username, cell_phone, type, status, name }` | 供 auth 服务调用，只核验凭证不签发令牌 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `cell_phone` | string | 学员/老师用手机号，优先按手机号 + type 定位 |
| `username` | string | 员工用用户名，手机号为空时才使用 |
| `password` | string | 明文密码，服务端用 bcrypt 比对 |
| `type` | int32 | 1-员工 2-学员 3-老师 |

---

### 学员账户

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `RegisterStudent` | `StudentFormRequest { cell_phone, code, password }` | `IdResponse { id }` | 学员注册，写 user + user_detail |
| `UpdateStudentPassword` | `StudentFormRequest { cell_phone, code, password }` | `EmptyResponse {}` | 学员改密，按手机号 + type=2 定位 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `cell_phone` | string | 手机号，必填 |
| `code` | string | 短信验证码（proto 已定义，logic 中未校验） |
| `password` | string | 密码，必填，落库前 bcrypt 加密 |

---

### 用户信息

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

**`UserDTO` 字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 用户 ID |
| `username` | string | 用户名 |
| `cell_phone` | string | 手机号 |
| `type` | int32 | 1-其他员工, 2-普通学员, 3-老师 |
| `name` | string | 姓名 |
| `gender` | int32 | 0-男性, 1-女性 |
| `icon` | string | 头像地址 |
| `email` | string | 邮箱 |
| `qq` | string | QQ 号码 |
| `job` | string | 岗位 |
| `province` / `city` / `district` | string | 省 / 市 / 区 |
| `intro` | string | 个人介绍 |
| `photo` | string | 形象照地址 |
| `role_id` | int64 | 角色 ID，老师和学生不用填 |

**`UserFormRequest` 相对 `UserDTO` 的增量字段**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `old_password` | string | 原始密码，携带新密码时必填 |
| `password` | string | 新密码 |

**`UpdateStatusRequest` 字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `user_id` | int64 | 目标用户 ID |
| `status` | int32 | 0-禁用 1-正常，其它取值拒绝 |
| `operator` | int64 | 操作人 ID，写入 `updater` |

---

### 管理后台分页查询

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `PageQueryStudents` | `UserPageRequest` | `StudentPageResponse { total, pages, list }` | 学员分页（固定 type=2） |
| `PageQueryTeachers` | `UserPageRequest` | `TeacherPageResponse { total, pages, list }` | 老师分页（固定 type=3） |
| `PageQueryStaffs` | `UserPageRequest` | `StaffPageResponse { total, pages, list }` | 员工分页（固定 type=1） |

**`UserPageRequest` 字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `page_no` | int64 | 页码，< 1 归一化为 1 |
| `page_size` | int64 | 每页条数，缺省 10，上限 100 |
| `sort_by` | string | 排序字段，白名单：name / status / id，其它落回 create_time |
| `is_asc` | bool | true 升序，false 降序 |
| `name` | string | 姓名模糊匹配 `d.name like %?%` |
| `phone` | string | 手机号模糊匹配 `u.cell_phone like %?%` |
| `status` | int32 | 账户状态，负值表示不过滤 |

**响应 VO 差异**：

| VO | 独有字段 |
|----|---------|
| `StudentPageVO` | `gender`, `course_amount` |
| `TeacherPageVO` | `photo`, `job`, `intro`, `course_amount`, `exam_question_amount` |
| `StaffVO` | `role_id`, `role_name` |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `user-api` (自身 API 层) | `apps/user/api/internal/svc/servicecontext.go` import `userclient "tjxt/apps/user/rpc/client/user"` | 用户管理全部 HTTP 接口最终走自身 RPC |
| `auth-api` | `apps/auth/api/internal/svc/servicecontext.go` import `userclient "tjxt/apps/user/rpc/client/user"` | 登录接口 `apps/auth/api/internal/logic/account/loginlogic.go:32` 调用 `LoginVerify` 核验凭证 |

（注：user 服务是身份域的基础设施，`user.rpc` 只被上述两个 API 服务在 `svc/servicecontext.go` 中显式装配。）

---

## 调用典型场景

1. **用户登录** → `auth-api` 登录 handler 调 `User.LoginVerify` 核验手机号/用户名 + 密码 → 成功后由 auth 服务签发 JWT
2. **学员注册** → 前端提交手机号 + 验证码 + 密码 → `RegisterStudent` 查重后写 `user` 与 `user_detail`
3. **个人中心** → 前端调 `GetUserDetail` 拉取当前用户资料 → 编辑后调 `UpdateCurrentUser` 合并保存
4. **管理端用户列表** → 后台调 `PageQueryStudents` / `PageQueryTeachers` / `PageQueryStaffs` 分类分页展示
5. **账户封禁** → 管理员调 `UpdateUserStatus(status=0)` → 模型层同步失效 id / cellPhone+type / username 三个缓存键
6. **忘记密码** → 管理员调 `ResetPassword` 重置为默认初始口令，或学员自助调 `UpdateStudentPassword`
7. **跨服务用户信息聚合** → 其它域拿到 userId 列表后调 `GetUsersByIds` 一次性回填昵称头像

---

## 自定义 Model 方法

`usermodel.go` 扩展了：
- `FindPageByType(ctx, userType, name, phone, status, offset, limit, sortCol, sortDir)` — 按类型 + 条件分页，`LEFT JOIN user_detail` 返回 `UserWithDetail` 联合视图；`status < 0` 不过滤；排序列由调用方白名单映射后传入
- `FindByIdsWithDetail(ctx, ids)` — 批量查询用户并附带资料，`IN (?)` 占位符拼接，按 `u.id DESC` 排序
- `ExistsByCellPhone(ctx, cellPhone)` — 手机号占用判断（不区分用户类型）
- `UpdateStatus(ctx, id, status, updater)` — 更新账户状态并失效 `cache:tjUser:user:id:`、`cache:tjUser:user:cellPhone:type:`、`cache:tjUser:user:username:` 三个缓存键

`userdetailmodel.go` 未新增扩展方法，仅内嵌 goctl 生成的 `userDetailModel` 接口。
