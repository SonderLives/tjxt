> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_user.sql`

---

# User Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `usermodel.go` | `user` | FindPageByType, FindByIdsWithDetail, ExistsByCellPhone, UpdateStatus |
| `userdetailmodel.go` | `user_detail` | 无（仅内嵌 goctl 生成接口） |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `usermodel.go` 扩展 `UserModel` 接口）。

---

## 表清单与字段说明

### 1. `user` — 学员用户表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键（雪花 ID，非自增） | PK |
| `username` | varchar(32) | 用户名 | 唯一索引 `username_idx` |
| `cell_phone` | varchar(11) | 手机号 | 联合唯一索引 `cell_idx` (cell_phone, type) |
| `password` | varchar(255) | 密码，bcrypt 密文 | - |
| `type` | tinyint | 用户类型：1-员工, 2-普通学员, 3-老师，默认 0 | 普通索引 `type_idx` |
| `status` | tinyint | 账户状态：0-禁用 1-正常，默认 1 | 过滤条件 |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建者 ID，可空 | - |
| `updater` | bigint | 更新者 ID，默认 0 | - |

> **注意**：手机号唯一性由 `(cell_phone, type)` 联合约束保证，即同一手机号可分别注册为学员与老师；但 `ExistsByCellPhone` 在应用层做的是不分类型的全局查重，比 DDL 约束更严格。

---

### 2. `user_detail` — 教师详情表（实为全类型用户资料表）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 关联用户 ID，与 `user.id` 同值 | PK |
| `type` | tinyint | 用户类型：1-员工, 2-普通学员, 3-老师，默认 2 | - |
| `name` | varchar(64) | 名字，默认空串 | 全文索引 `name_idx` |
| `gender` | tinyint | 性别：0-男性，1-女性，默认 0 | - |
| `icon` | varchar(255) | 头像地址，可空 | - |
| `email` | varchar(255) | 邮箱，可空 | - |
| `qq` | varchar(18) | QQ 号码，可空 | - |
| `birthday` | date | 生日，可空 | - |
| `job` | varchar(32) | 岗位，可空 | - |
| `province` | varchar(32) | 省，可空 | - |
| `city` | varchar(32) | 市，可空 | - |
| `district` | varchar(32) | 区，可空 | - |
| `intro` | varchar(255) | 个人介绍，可空 | - |
| `photo` | varchar(255) | 形象照地址，可空 | - |
| `role_id` | bigint | 角色 ID，非空 | 关联 auth 域 role 表 |
| `course_amount` | smallint | 购买课程数量，学生才有该字段信息，默认 0 | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建者 ID，可空 | - |
| `updater` | bigint | 更新者 ID，默认 0 | - |
| `dep_id` | bigint | 部门 ID，默认 0 | - |

> **关系**：`user_detail.id` 与 `user.id` 为 1:1 共享主键，无外键约束，由应用层在 `AddUser` / `RegisterStudent` 中同事务外顺序写入两张表。
>
> **注意**：`birthday` 与 `dep_id` 在 proto 与 logic 中均未被读写，属于表结构预留字段。

---

## 联合视图 `UserWithDetail`

`usermodel.go` 中定义的非表结构投影，用于分页与批量查询，避免 N+1：

```
SELECT u.id, u.cell_phone, u.type, u.status, u.create_time,
       d.name, d.icon, d.gender, d.photo, d.job, d.intro,
       d.course_amount, d.role_id, d.email, d.qq,
       d.province, d.city, d.district
FROM `user` u LEFT JOIN `user_detail` d ON u.id = d.id
```

| 字段 | Go 类型 | 来源 |
|------|---------|------|
| `Id` / `CellPhone` / `Type` / `Status` / `CreateTime` | int64 / string / int64 / int64 / time.Time | `user` |
| `Name` / `Icon` / `Photo` / `Job` / `Intro` | sql.NullString | `user_detail` |
| `Gender` / `CourseAmount` / `RoleId` | int64 | `user_detail` |
| `Email` / `Qq` / `Province` / `City` / `District` | sql.NullString | `user_detail` |

> `LEFT JOIN` 保证资料缺失的用户仍能被查出，NULL 列由 `strVal` 统一降级为空串。

---

## 关系图

```
user (1) ─── (1) user_detail        共享主键 id

user_detail.role_id ──→ role.id     跨库引用（auth 域 tj_auth.role）

user.id ←── login_record.user_id    跨库引用（auth 域 tj_auth.login_record）
```

---

## 缓存键

goctl 为 `user` 表生成了三组缓存键（见 `usermodel_gen.go`）：

| 缓存键前缀 | 对应查询 |
|-----------|---------|
| `cache:tjUser:user:id:` | `FindOne(id)` |
| `cache:tjUser:user:cellPhone:type:` | `FindOneByCellPhoneType(cellPhone, type)` |
| `cache:tjUser:user:username:` | `FindOneByUsername(username)` |

`UpdateStatus` 因走 `ExecNoCacheCtx` 直改 SQL，需手动 `DelCacheCtx` 依次失效上述三个键。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
usermodel_gen.go          ← goctl 生成，只读
usermodel.go              ← 手写扩展 FindPageByType / FindByIdsWithDetail / ExistsByCellPhone / UpdateStatus
```

当前项目自定义 Model 模式：
- `usermodel.go` — FindPageByType, FindByIdsWithDetail, ExistsByCellPhone, UpdateStatus
- `userdetailmodel.go` — 无扩展方法，仅保留 `customUserDetailModel` 包装壳，便于后续按需追加

`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`，供 logic 层 `errors.Is` 判定「记录不存在」。
