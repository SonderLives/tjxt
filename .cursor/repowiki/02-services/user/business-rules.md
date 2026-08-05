> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/user/rpc/internal/logic/user/*.go`

---

# User Business Rules

## 1. 身份核验（LoginVerify）

**核心规则**：user 服务只做「找用户 + 比对密码 + 状态检查」，**不签发任何令牌**；角色归属由 auth 服务通过 `account_role` 自查。

| 规则 | 说明 |
|------|------|
| 定位优先级 | `cell_phone` 非空优先按「手机号 + type」定位，否则按 `username` 定位 |
| 两者皆空 | 直接返回 `ErrUserNotFound`，不查库 |
| 状态前置 | `status != 1` 返回 `ErrUserDisabled`，**先于密码比对**，禁用账户不泄漏密码正确性 |
| 密码比对 | `bcrypt.CompareHashAndPassword`，失败返回 `ErrBadCredential` |
| 错误类型 | 使用包级哨兵错误而非 `xerr`，由调用方 auth 服务翻译为业务错误码 |

```
流程（LoginVerify）:
  1. cell_phone 非空 → FindOneByCellPhoneType(cell_phone, type)
     否则 username 非空 → FindOneByUsername(username)
     否则 → ErrUserNotFound
  2. sqlc.ErrNotFound → ErrUserNotFound
  3. user.Status != 1 → ErrUserDisabled
  4. bcrypt 比对失败 → ErrBadCredential
  5. 返回 userId / username / cellPhone / type / status
```

> **注意**：`LoginVerifyResponse.name` 字段在 proto 中已定义，但 logic 未回填（未查 `user_detail`），实际恒为空串。

---

## 2. 学员注册（RegisterStudent）

| 规则 | 说明 |
|------|------|
| 必填校验 | `cell_phone` 与 `password` 均不可为空，否则 `BadRequest` |
| 手机号查重 | `ExistsByCellPhone` 全局查重（不区分 type），已存在返回 `Conflict("该手机号已注册")` |
| username 兜底 | 注册时 `username` 直接取 `cell_phone` 值 |
| 密码加密 | `bcrypt.GenerateFromPassword(..., bcrypt.DefaultCost)` |
| ID 生成 | `idgen.NextID()` 雪花 ID，`user` 与 `user_detail` 共用同一 ID |
| 类型固定 | `type = 2`（普通学员），`status = 1`（正常） |
| 审计字段 | `creater` / `updater` 均写入自身 ID（自注册） |
| 验证码 | `code` 字段接收但 **logic 未做任何校验**，短信验证留待接入 |

```
流程（RegisterStudent）:
  1. 校验 cell_phone / password 非空
  2. ExistsByCellPhone → 存在则 Conflict
  3. bcrypt 加密密码
  4. id = idgen.NextID()
  5. Insert user   { username=cell_phone, type=2, status=1 }
  6. Insert user_detail { id, type=2 }   ← 仅骨架，资料字段留空
```

> **⚠️ 一致性风险**：步骤 5、6 未包裹事务，`user` 写成功而 `user_detail` 写失败会产生孤儿用户；`GetUserById` 对此做了退化兜底（资料缺失时只返回基础信息）。

---

## 3. 管理端新增用户（AddUser）

| 规则 | 说明 |
|------|------|
| 至少一个标识 | `cell_phone` 与 `username` 不能同时为空 |
| 手机号查重 | `cell_phone` 非空时走 `ExistsByCellPhone`，冲突返回 `Conflict("该手机号已存在")` |
| 用户名查重 | `username` 非空时走 `FindOneByUsername`，查到即 `Conflict("该用户名已存在")`；非 `ErrNotFound` 的错误按内部错误处理 |
| 类型缺省 | `type == 0` 时兜底为 `2`（普通学员） |
| 密码策略 | **密码不随请求传入**，统一使用 `defaultInitialPassword`（常量 `"123456"`） |
| 资料落库 | 空串字段经 `nullStr` 转为 SQL NULL，不写空串 |

```
流程（AddUser）:
  1. cell_phone 与 username 都空 → BadRequest
  2. cell_phone 非空 → ExistsByCellPhone 查重
  3. username 非空 → FindOneByUsername 查重
  4. userType = in.Type，为 0 则置 2
  5. bcrypt 加密 defaultInitialPassword
  6. Insert user + Insert user_detail（共用雪花 ID）
```

> **⚠️ 安全提示**：`defaultInitialPassword = "123456"` 是硬编码常量，源码注释明确「正式环境应改为短信或随机下发」。

---

## 4. 用户信息更新（UpdateUserById / UpdateCurrentUser）

**核心规则**：两者共用同一套「非零字段合并」语义，仅覆盖请求中的非零值，避免误清零未传字段。

| 规则 | 说明 |
|------|------|
| 用户必须存在 | `FindOne` 失败且为 `ErrNotFound` → `NotFound("用户不存在")` |
| 资料可缺失 | `user_detail` 查不到时**就地构造**空壳 `{Id, Type: u.Type, Creater: in.Id}` 后继续 |
| 字符串合并 | 走 `applyStr(cur, next)`：`next` 为空则保留原值，否则覆盖 |
| 数值合并 | `Name` / `Gender` / `Type` / `RoleId` 判 `!= 0` / `!= ""` 才覆盖 |
| 审计字段 | `u.Updater` 与 `d.Updater` 均写入 `in.Id` |
| 写入顺序 | 先 `UserModel.Update`，后 `UserDetailModel.Update`，任一失败即返回 |

**`UpdateCurrentUser` 独有的改密规则**：

| 规则 | 说明 |
|------|------|
| 原密码必填 | `password != "" && old_password == ""` → `BadRequest("修改密码需提供原密码")` |
| 原密码校验 | bcrypt 比对失败 → `Unauthorized("原密码错误")` |
| 校验时机 | 在两次 `Update` **之前**完成，校验不通过不落任何库 |

```
流程（UpdateCurrentUser）:
  1. FindOne(user) → 不存在则 NotFound
  2. FindOne(user_detail) → 缺失则构造空壳
  3. 合并 user 字段：username / cell_phone / type（非零覆盖）
  4. 合并 detail 字段：name / gender / role_id（非零覆盖）
                      icon / email / qq / job / province / city / district / intro / photo（applyStr）
  5. 若 password 非空：
       old_password 为空 → BadRequest
       bcrypt 比对失败   → Unauthorized
       通过 → 重新 bcrypt 生成 u.Password
  6. Update user → Update user_detail
```

> **⚠️ 一致性风险**：步骤 6 的两次 Update 未包裹事务；且 `UpdateUserById` 用 `in.Id` 而非操作人 ID 填 `updater`，实际记录的是被改用户自身 ID。

---

## 5. 密码管理

| 场景 | 方法 | 规则 |
|------|------|------|
| 学员自助改密 | `UpdateStudentPassword` | 按「手机号 + type=2」定位；**不校验原密码、不校验短信验证码**；`updater` 写用户自身 ID |
| 管理员重置 | `ResetPassword` | 重置为 `defaultInitialPassword`；`updater` 写 `in.UserId`（即被重置者自身） |
| 当前用户改密 | `UpdateCurrentUser` | 必须提供并通过 `old_password` 校验，见第 4 节 |

```
流程（UpdateStudentPassword）:
  1. 校验 cell_phone / password 非空
  2. FindOneByCellPhoneType(cell_phone, 2) → 不存在则 NotFound
  3. bcrypt 加密新密码
  4. Update user
```

> **⚠️ 安全提示**：`UpdateStudentPassword` 只凭手机号即可重设密码，`StudentFormRequest.code`（短信验证码）在 logic 中**完全未被读取**。该接口必须由上层 API 完成验证码校验后才可暴露。

---

## 6. 账户状态管理（UpdateUserStatus）

| 规则 | 说明 |
|------|------|
| 状态白名单 | 仅接受 0（禁用）与 1（正常），其它取值 → `BadRequest("状态取值非法")` |
| 存在性校验 | 由模型层 `UpdateStatus` 内部先 `FindOne`，不存在则冒泡 `ErrNotFound` → `NotFound` |
| 操作人记录 | `in.Operator` 写入 `updater` 列（**本服务唯一正确记录操作人的接口**） |
| 更新方式 | `ExecNoCacheCtx` 直接 UPDATE，同时刷新 `update_time = now()` |

**缓存失效策略**：`UpdateStatus` 绕过了 goctl 的缓存写通道，必须手动清理三个键，否则 `LoginVerify` 会读到旧的 `status`：

```
流程（UpdateStatus 缓存失效）:
  1. FindOne(id)  ← 先取旧记录，拿到 CellPhone / Type / Username 用于拼键
  2. UPDATE user SET status, updater, update_time
  3. DelCache  cache:tjUser:user:id:{id}
  4. DelCache  cache:tjUser:user:cellPhone:type:{cellPhone}:{type}
  5. DelCache  cache:tjUser:user:username:{username}
```

> 三个键必须**全部**删除：漏删任一个都会导致对应查询路径（按 ID / 按手机号登录 / 按用户名登录）继续命中旧状态。

---

## 7. 查询与分页

### 单体与批量查询

| 规则 | 说明 |
|------|------|
| 资料退化 | `GetUserById` / `GetUserDetail` 中 `user_detail` 缺失不报错，`toUserDTO(u, nil)` 只填基础字段 |
| 空入参短路 | `GetUsersByIds` 收到空 ID 列表直接返回空 `UserListResponse`，不查库 |
| 批量走视图 | `FindByIdsWithDetail` 单条 SQL `LEFT JOIN` 完成，杜绝 N+1 |
| NULL 降级 | `strVal` 将 `sql.NullString` 无效值统一转空串 |
| 时间格式 | `create_time` 统一格式化为 `2006-01-02 15:04:05` |

### 分页查询

三个分页方法共用 `FindPageByType`，仅 `userType` 常量不同：

| 方法 | 固定 type |
|------|----------|
| `PageQueryStaffs` | 1 |
| `PageQueryStudents` | 2 |
| `PageQueryTeachers` | 3 |

| 规则 | 说明 |
|------|------|
| 参数归一化 | `page.Normalize`：`pageNo < 1` → 1；`pageSize < 1` → 10；`pageSize > 100` → 100 |
| 总页数 | `page.CalcPages(total, limit)`，`total <= 0` 时返回 0 |
| 条件拼接 | `name` / `phone` 走 `like %?%`；`status >= 0` 才追加状态过滤 |
| 先查列表后查总数 | 两条 SQL 复用同一组 `args`，列表额外追加 `offset, limit` |

**SQL 注入防护**：排序字段**不允许透传**，`sortClause` 做白名单映射后才交给模型层：

```
sortClause(sortBy, isAsc):
  "name"   → d.`name`
  "status" → u.`status`
  "id"     → u.`id`
  其它/空  → u.`create_time`   ← 默认按注册时间
  isAsc ? "ASC" : "DESC"       ← 默认 DESC
```

> **未实现的跨服务聚合**（源码注释明确标注，当前返回本地值或空值）：
> - `UserDetailVO.role_name`、`StaffVO.role_name` — 归属 auth 服务 `role` 表，暂留空
> - `TeacherPageVO.exam_question_amount` — 归属 exam 服务，暂留 0
> - `TeacherPageVO.course_amount` — 当前取 `user_detail.course_amount`（学员购课数），权威值应由 course 服务计算

---

## 状态说明

### 用户类型 `type`

| 值 | 含义 | 备注 |
|----|------|------|
| 1 | 员工 | `PageQueryStaffs` 查询该类型 |
| 2 | 普通学员 | 注册默认值，`AddUser` 缺省兜底值 |
| 3 | 老师 | `PageQueryTeachers` 查询该类型 |

### 账户状态 `status`

| 值 | 含义 | 影响 |
|----|------|------|
| 0 | 禁用 | `LoginVerify` 返回 `ErrUserDisabled` |
| 1 | 正常 | 可正常登录，建表默认值 |

### 性别 `gender`

| 值 | 含义 |
|----|------|
| 0 | 男性 |
| 1 | 女性 |

> **注意**：`gender` 的合并逻辑为 `if in.Gender != 0`，因此**无法通过更新接口把性别从「女性」改回「男性」**（0 被当作未传）。

---

## 错误码约定

| 场景 | 错误构造 |
|------|---------|
| 参数非法 | `xerr.BadRequestf(...)` |
| 记录不存在 | `xerr.NotFound(...)` |
| 唯一性冲突 | `xerr.Conflict(...)` |
| 原密码错误 | `xerr.Unauthorized("原密码错误")` |
| 数据库/加密异常 | `xerr.Wrap(err, xerr.CodeInternal, ...)` |
| 身份核验失败 | 包级哨兵错误 `ErrUserNotFound` / `ErrBadCredential` / `ErrUserDisabled`（**不走 xerr**） |
