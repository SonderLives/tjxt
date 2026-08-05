> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_message.sql`, `sql/ddl/tj_message_model.sql`

---

# Message Data Model

## 两份 DDL 的关系

message 是全仓唯一拥有**两份 DDL** 的服务，二者建表语句同源但用途不同：

| DDL 文件 | 表数量 | 用途 | 是否用于 goctl 生成 |
|---------|--------|------|-------------------|
| `sql/ddl/tj_message.sql` | 7 | 完整库结构，文件头注明 `extracted from migration/tj_message.sql`，字段注释最完整 | 否 |
| `sql/ddl/tj_message_model.sql` | 6 | 精简版，文件头注明 `notice_task_target excluded due to composite PK` | **是** |

**差异点**：

1. `tj_message_model.sql` **剔除了 `notice_task_target` 表**——该表使用 `(task_id, target_id)` 复合主键，goctl model 生成器不支持，故排除。
2. `tj_message_model.sql` 的字段 `COMMENT` 被截短（例如 `user_inbox.type` 由 `'通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知，4-私信'` 简化为 `'通知类型'`）。生成的 `*_gen.go` 结构体行尾注释与精简版一致，可据此反向确认模型确由 `tj_message_model.sql` 生成。
3. `tj_message.sql` 保留了 `-- Records of xxx` 数据段占位注释，`tj_message_model.sql` 没有。

> **结论**：字段语义以 `tj_message.sql` 为准；模型代码以 `tj_message_model.sql` 为准。`notice_task_target` **无对应 Model**，需手写 SQL 访问。

---

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `messagetemplatemodel.go` | `message_template` | 无（仅 goctl 空壳） |
| `noticetaskmodel.go` | `notice_task` | 无（仅 goctl 空壳） |
| `noticetemplatemodel.go` | `notice_template` | 无（仅 goctl 空壳） |
| `publicnoticemodel.go` | `public_notice` | 无（仅 goctl 空壳） |
| `smsthirdplatformmodel.go` | `sms_third_platform` | 无（仅 goctl 空壳） |
| `userinboxmodel.go` | `user_inbox` | 无（仅 goctl 空壳） |
| （缺失） | `notice_task_target` | 未生成 Model，复合主键不受 goctl 支持 |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `userinboxmodel.go` 扩展 `UserInboxModel` 接口）——**当前 6 个自定义文件全部为空壳，尚未追加任何方法**。

6 个 Model 均在 `apps/message/rpc/internal/svc/servicecontext.go` 中通过 `sqlx.NewMysql(c.DataSource)` + `c.Cache` 注入，走带缓存的 `sqlc.CachedConn`。

**缓存 key 前缀**（`*_gen.go` 中定义）：

| Model | 缓存前缀 |
|-------|---------|
| `MessageTemplate` | `cache:messageTemplate:id:` |
| `NoticeTask` | `cache:noticeTask:id:` |
| `NoticeTemplate` | `cache:noticeTemplate:id:` |
| `PublicNotice` | `cache:publicNotice:id:` |
| `SmsThirdPlatform` | `cache:smsThirdPlatform:id:` |
| `UserInbox` | `cache:userInbox:id:` |

---

## 表清单与字段说明

### 1. `notice_template` — 通知模板表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 通知模板 id | PK |
| `name` | varchar(50) | 通知模板名称 | - |
| `code` | varchar(50) | 模板代号，例如：verify-code | `idx_code` |
| `type` | tinyint | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知 | - |
| `status` | tinyint | 模板状态：0-草稿，1-使用中，2-停用，默认 0 | - |
| `title` | varchar(50) | 通知标题，短信模板可以不填（可空） | - |
| `content` | varchar(255) | 通知内容模板 | - |
| `is_sms_template` | bit(1) | 是否是短信模板，默认 b'0' | - |
| `creater` / `updater` | bigint | 创建人 / 更新人 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |

> Go 结构体 `NoticeTemplate` 中 `title` 映射为 `sql.NullString`，`is_sms_template` 映射为 `byte`。

---

### 2. `notice_task` — 通知任务表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 公告任务 id | PK |
| `template_id` | bigint | 任务对应的通知模板 id | - |
| `name` | varchar(128) | 任务名称 | - |
| `partial` | bit(1) | 是否是部分人的通告，默认 false | - |
| `push_time` | datetime | 任务预期执行时间（可空） | - |
| `interval` | int | 任务延迟执行时间间隔，单位是分钟（可空） | - |
| `expire_time` | datetime | 任务失效时间（可空） | - |
| `max_times` | int | 任务重复执行次数上限，1 则只发一次，默认 1 | - |
| `finished` | bit(1) | 任务是否完成，默认 false | - |
| `creater` / `updater` | bigint | 创建人 / 更新人 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |

> Go 结构体 `NoticeTask` 中 `push_time`/`expire_time` 为 `sql.NullTime`，`interval` 为 `sql.NullInt64`，`partial`/`finished` 为 `byte`。`template_id` 在 DDL 中**未建索引**。

---

### 3. `notice_task_target` — 通知任务目标用户表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `task_id` | bigint | 任务 id | 复合 PK |
| `target_id` | bigint | 目标用户 id | 复合 PK |

> 仅存在于 `sql/ddl/tj_message.sql`。使用 `(task_id, target_id)` 复合主键，**未生成 Model**，当 `notice_task.partial = 1` 时用于圈定目标人群。字符集为 `utf8mb4`，与其余 6 张 `utf8` 表不同。

---

### 4. `message_template` — 第三方短信模板表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 短信发送模板 id | PK |
| `name` | varchar(50) | 模板名称 | - |
| `platform_code` | varchar(50) | 第三方短信平台代号 | - |
| `sign_name` | varchar(50) | 签名 | - |
| `third_template_code` | varchar(50) | 第三方短信模板 code | - |
| `content` | varchar(255) | 第三方短信模板内容预览 | - |
| `template_id` | bigint | 通知模板 id | `idx_template_id` |
| `status` | tinyint | 模板状态：0-禁用，1-启用，默认 0 | - |
| `creater` / `updater` | bigint | 创建者 / 更新者 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |

> 表注释为「第三方短信平台签名和模板信息」。`platform_code` 逻辑上对应 `sms_third_platform.code`，DDL 中**无外键约束**。

---

### 5. `sms_third_platform` — 第三方云通讯平台表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 短信平台 id，`AUTO_INCREMENT` | PK |
| `name` | varchar(50) | 短信平台名称 | - |
| `code` | varchar(50) | 短信平台代码，例如：ali | - |
| `priority` | int unsigned | 数字越小优先级越高，最小为 0，默认 0 | - |
| `status` | int | 短信平台状态：0-禁用，1-启用，默认 1 | - |
| `creater` / `updater` | bigint | 创建人 / 更新人 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |

> 本库唯一使用 `AUTO_INCREMENT` 的表（`AUTO_INCREMENT = 4`），其余表主键由业务侧生成。Go 结构体中 `priority` 为 `uint64`。

---

### 6. `user_inbox` — 用户通知记录表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 用户通知 id | PK |
| `user_id` | bigint | 用户 id | `user_id` 索引 |
| `type` | tinyint | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知，4-私信，默认 4 | - |
| `title` | varchar(64) | 通知标题，默认 `''`（可空） | - |
| `content` | varchar(255) | 通知或私信内容 | - |
| `is_read` | bit(1) | 公告是否已读，默认 b'0' | - |
| `publisher` | bigint | 通知的发送者 id，0 则代表是系统，默认 0 | - |
| `push_time` | datetime | 创建时间 | `push_time` 索引 |
| `expire_time` | datetime | 过期时间，一旦过期用户端不再展示 | - |

> `push_time` 单独建索引，支撑收件箱按时间倒序分页；`user_id` 索引支撑按用户过滤。二者为**单列索引**，非联合索引。

---

### 7. `public_notice` — 公告消息模板表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 公告 id | PK |
| `title` | varchar(50) | 公告标题 | - |
| `content` | varchar(255) | 公告通知内容，可以存放公告消息模板 | - |
| `type` | tinyint | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知 | - |
| `push_time` | datetime | 通知发布时间 | - |
| `expire_time` | datetime | 通知失效时间 | - |

> 全库唯一**不含** `creater`/`updater`/`create_time`/`update_time` 审计字段的表。

---

## 关系图

```
notice_template (1) ──┬── (N) message_template   （template_id，无外键约束）
                      │
                      └── (N) notice_task        （template_id，无外键约束）
                                  │
                                  └── (N) notice_task_target  （task_id，复合主键，无 Model）

sms_third_platform (1) ─── (N) message_template  （code ↔ platform_code，逻辑关联）

public_notice   （独立表，全站公告）
user_inbox      （独立表，user_id 指向 user 域，跨库无外键）
```

> 本库**不使用逻辑删除**：7 张表均无 `deleted` 字段，删除操作为物理 `DELETE`（与 auth 域的 `deleted=1` 软删除模式不同）。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
userinboxmodel_gen.go     ← goctl 生成，只读（Insert / FindOne / Update / Delete）
userinboxmodel.go         ← 手写扩展位（当前为空壳，无任何扩展方法）
```

当前项目自定义 Model 模式：

- `messagetemplatemodel.go` — 无扩展
- `noticetaskmodel.go` — 无扩展
- `noticetemplatemodel.go` — 无扩展
- `publicnoticemodel.go` — 无扩展
- `smsthirdplatformmodel.go` — 无扩展
- `userinboxmodel.go` — 无扩展

`vars.go` 仅定义 `ErrNotFound = sqlx.ErrNotFound`，**无** auth 域那样的 `sqlhelper.go` 通用 SQL 工具函数。

📋 **待补齐（设计意图）**：实现 RPC 层分页/列表能力前，至少需要补写 `FindPage`（notice_template / notice_task / message_template / public_notice）、`FindPageByUserId`（user_inbox）、`FindAll`（sms_third_platform，按 `priority` 升序）以及 `notice_task_target` 的手写访问层。
