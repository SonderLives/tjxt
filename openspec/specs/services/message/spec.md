# message（消息通知）Specification

## Purpose

message（消息通知）是站内信与通知微服务（gRPC 服务名 `Message`，etcd key `message.rpc`，监听 `0.0.0.0:8087`），管理通知模板、通知任务、第三方短信模板与短信平台、用户收件箱（站内信）、全站公告五类实体。API 层 `message-api`（HTTP :8807）为纯转发层，18 个接口全部转发到自身 RPC（19 个方法，`SendNotice` 仅暴露为 RPC）；数据落在 MySQL `tj_message` 库的 7 张表（无数据库外键），Redis 缓存主键查询，etcd（`127.0.0.1:2379`）服务注册与发现。当前除自身 API 层外无其它服务消费 `message.rpc`——架构文档所述 pay/trade/learning/promotion → message 的通知链路在代码中尚无实现（见「RPC 服务契约与服务发现」）。

## Requirements

### Requirement: HTTP API 用户收件箱接口

`message-api` SHALL 提供当前用户收件箱（站内信）的查询、标记已读与删除接口（三个接口均需 JWT 认证，权限标签 `inbox`；`message.api` 未声明任何 RBAC 权限注解）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /inbox | 分页查询当前用户收件箱 | InboxListReq (form) | InboxListVO |
| DELETE | /inbox/:id | 删除一条站内信 | IdPathReq | OkVO |
| PUT | /inbox/:id/read | 标记站内信为已读 | IdPathReq | OkVO |

- `InboxListReq`：`PageNo`/`PageSize`/`Type` 均为 optional；HTTP 侧无 `userId` 字段，由 API 层从 JWT 提取当前用户后填入 RPC `InboxPageReq.userId`
- 列表项 `UserInboxVO`：`Id, UserId, Type, Title, Content, IsRead, Publisher, PushTime, ExpireTime`；`isRead`（是否已读）与 `pushTime`（推送时间）由服务端写入
- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；错误码遵循共享错误码表（见 `openspec/specs/shared/error-codes`）

#### Scenario: 查询收件箱并处理单条消息

- **WHEN** 已登录用户 GET /inbox 携带可选的 `Type` 过滤（未传时不按类型过滤）
- **THEN** 以 JWT 提取的 `userId` 作为归属条件分页返回 `InboxListVO{Total, List}`
- **AND** 用户对单条消息可 PUT /inbox/:id/read 置已读、DELETE /inbox/:id 删除，均返回 `OkVO{Success}`

### Requirement: HTTP API 通知模板管理接口

`message-api` SHALL 提供通知模板的管理与查询接口（均需 JWT 认证，权限标签 `noticetemplate`）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /notice-templates | 分页查询通知模板 | PageReq (form) | NoticeTemplateListVO |
| POST | /notice-templates | 新增/更新通知模板 | NoticeTemplateSaveReq | IdVO |
| DELETE | /notice-templates/:id | 根据 id 删除通知模板 | IdPathReq | OkVO |
| GET | /notice-templates/:id | 根据 id 查询通知模板 | IdPathReq | NoticeTemplateVO |

- `NoticeTemplateSaveReq`：`Id`(optional), `Name`, `Code`, `Type`, `Title`(optional), `Content`, `IsSmsTemplate`(optional)
- `NoticeTemplateVO.status`（0-草稿，1-使用中，2-停用）仅出现在响应中，保存请求未定义该字段（状态不可由请求直接设置）

#### Scenario: 管理端维护通知模板

- **WHEN** 管理端 POST /notice-templates 提交 `code`/`content` 模板
- **THEN** 落库并返回 `IdVO{Id}`，后续可 GET /notice-templates 分页展示、GET /notice-templates/:id 单查
- **AND** 删除走 DELETE /notice-templates/:id（物理删除），返回 `OkVO`

### Requirement: HTTP API 通知任务管理接口

`message-api` SHALL 提供通知任务（定时通告）的管理与查询接口（均需 JWT 认证，权限标签 `noticetask`）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /notice-tasks | 分页查询通知任务 | PageReq (form) | NoticeTaskListVO |
| POST | /notice-tasks | 新增/更新通知任务 | NoticeTaskSaveReq | IdVO |
| DELETE | /notice-tasks/:id | 根据 id 删除通知任务 | IdPathReq | OkVO |
| GET | /notice-tasks/:id | 根据 id 查询通知任务 | IdPathReq | NoticeTaskVO |

- `NoticeTaskSaveReq`：`Id`(optional), `TemplateId`, `Name`, `Partial`(optional), `PushTime`(optional), `Interval`(optional), `ExpireTime`(optional), `MaxTimes`(optional)
- `NoticeTaskVO.finished`（任务是否完成）仅出现在响应中，由服务端维护，客户端不可直接改写

#### Scenario: 配置定时通告

- **WHEN** 管理员 POST /notice-tasks 绑定 `TemplateId` 并设置 `PushTime`/`Interval`/`MaxTimes`
- **THEN** 通知任务落库并返回 `IdVO{Id}`，由调度侧按 `finished` 标记推进
- **AND** GET /notice-tasks/:id 单查可见服务端维护的 `finished` 状态

### Requirement: HTTP API 短信模板与短信平台接口

`message-api` SHALL 提供第三方短信模板的管理接口与短信平台查询接口（均需 JWT 认证，权限标签分别为 `messagetemplate` / `smsplatform`）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /message-templates | 分页查询短信模板 | PageReq (form) | MessageTemplateListVO |
| POST | /message-templates | 新增/更新短信模板 | MessageTemplateSaveReq | IdVO |
| DELETE | /message-templates/:id | 根据 id 删除短信模板 | IdPathReq | OkVO |
| GET | /sms-platforms | 查询全部第三方短信平台 | - | SmsPlatformListVO |

- `MessageTemplateSaveReq`：`Id`(optional), `Name`, `PlatformCode`, `SignName`, `ThirdTemplateCode`, `Content`, `TemplateId`, `Status`(optional)
- `SmsPlatformVO`：`Id, Name, Code, Priority, Status`；`SmsPlatformListVO` 为例外结构，仅有 `{ List }`、无 `Total`（全量列表返回）；其余列表类响应统一为 `{ Total, List }`
- 短信模板未提供 HTTP 单查接口，详情仅能通过分页列表获取

#### Scenario: 对接第三方短信平台

- **WHEN** 管理端 GET /sms-platforms
- **THEN** 返回全部第三方短信平台（`SmsPlatformListVO`，无 Total）
- **AND** POST /message-templates 绑定平台 `PlatformCode` 与第三方模板 `ThirdTemplateCode` 后返回 `IdVO`

### Requirement: HTTP API 公告管理接口

`message-api` SHALL 提供全站公告的管理与查询接口（均需 JWT 认证，权限标签 `notice`）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /notices | 分页查询公告 | PageReq (form) | PublicNoticeListVO |
| POST | /notices | 新增/更新公告 | PublicNoticeSaveReq | IdVO |
| DELETE | /notices/:id | 根据 id 删除公告 | IdPathReq | OkVO |

- `PublicNoticeSaveReq`：`Id`(optional), `Title`, `Content`, `Type`, `PushTime`, `ExpireTime`
- 公告未提供单查接口，仅有分页列表

#### Scenario: 发布全站公告

- **WHEN** 运营 POST /notices 写入 `Title`/`Content` 并设置 `PushTime`/`ExpireTime`
- **THEN** 公告落库并返回 `IdVO{Id}`
- **AND** 端上 GET /notices 分页拉取，按 `pushTime`/`expireTime` 窗口展示

### Requirement: RPC 服务契约与服务发现

`Message` 服务 SHALL 以 gRPC 服务名 `Message` 注册于 etcd（key `message.rpc`，监听 `0.0.0.0:8087`），提供 19 个 RPC 方法；HTTP 侧 18 个接口为其子集，`SendNotice` 仅存在于 RPC、不对外暴露 HTTP 入口。消费方：

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `message-api`（自身 API 层） | HTTP Handler → `messageclient.Message` RPC（`apps/message/api/internal/svc/servicecontext.go` 注入 `MessageRpc`） | 18 个 HTTP 接口全部指向自身 RPC |

> 已知缺口（跨域通知链路未落地）：全仓 Grep `apps/message/rpc/message` / `messageclient` / `MessageRpc`，除 `apps/message/` 自身外无其它服务引用。`specs/shared/mq-events` 与 `00-architecture/service-topology.md` 中描述的 pay-rpc / trade-rpc / learning-rpc / promotion-rpc → message-rpc 通知链路，目前在代码中尚无对应 import 或 MQ 消费者实现——即 message 并未消费任何其它域事件（订单、课程学习等业务事件触发通知）来驱动 `SendNotice`，该写入口当前仅由契约定义。

#### Scenario: message-api handler 转发

- **WHEN** message-api 任一 handler 收到请求
- **THEN** logic 层调用 `l.svcCtx.MessageRpc` 对应 RPC 方法（API 层为纯转发层，不直连数据库）
- **AND** `MessageRpc` 客户端按 `MessageRpc.Etcd`（`127.0.0.1:2379`，key `message.rpc`）经 etcd 发现 `message.rpc`

#### Scenario: 业务事件触发通知（设计意图）

- **WHEN** 业务方（设计上为 trade / pay / learning / promotion）产生需要通知用户的事件
- **THEN** 设计上应调 `SendNotice` 写入 `user_inbox`，用户端经 `ListInbox` 拉取
- **AND** 当前代码中无任何业务服务装配 `MessageRpc` 客户端，也无 MQ 消费者实现，该链路尚未接线

### Requirement: RPC 通知模板与通知任务方法

`Message` 服务 SHALL 提供通知模板 4 个、通知任务 4 个共 8 个管理方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveNoticeTemplate` | `NoticeTemplateSaveReq { id, name, code, type, title, content, isSmsTemplate }` | `IdReply { id }` | 新增/更新通知模板 |
| `DeleteNoticeTemplate` | `IdReq { id }` | `Empty {}` | 按 ID 删除通知模板 |
| `GetNoticeTemplate` | `IdReq { id }` | `NoticeTemplateVO { id, name, code, type, status, title, content, isSmsTemplate, createTime }` | 按 ID 查询通知模板 |
| `ListNoticeTemplates` | `PageReq { pageNo, pageSize }` | `NoticeTemplateListReply { total, list }` | 分页查询通知模板 |
| `SaveNoticeTask` | `NoticeTaskSaveReq { id, templateId, name, partial, pushTime, interval, expireTime, maxTimes }` | `IdReply { id }` | 新增/更新通知任务 |
| `DeleteNoticeTask` | `IdReq { id }` | `Empty {}` | 按 ID 删除通知任务 |
| `GetNoticeTask` | `IdReq { id }` | `NoticeTaskVO { id, templateId, name, partial, pushTime, interval, expireTime, maxTimes, finished, createTime }` | 按 ID 查询通知任务 |
| `ListNoticeTasks` | `PageReq { pageNo, pageSize }` | `NoticeTaskListReply { total, list }` | 分页查询通知任务 |

关键字段约束：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 模板代号，例如：verify-code |
| `type` | int32 | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知 |
| `title` | string | 通知标题，短信模板可以不填 |
| `isSmsTemplate` | bool | 是否是短信模板 |
| `partial` | bool | 是否是部分人的通告，默认 false |
| `pushTime` | string | 任务预期执行时间 |
| `interval` | int32 | 任务延迟执行时间间隔，单位是分钟（非秒、非毫秒） |
| `expireTime` | string | 任务失效时间 |
| `maxTimes` | int32 | 任务重复执行次数上限，1 则只发一次 |

- `NoticeTemplateVO.status`（0-草稿，1-使用中，2-停用）与 `NoticeTaskVO.finished` 仅出现在响应中，由服务端维护

#### Scenario: 管理端维护通知模板

- **WHEN** 后台调 `SaveNoticeTemplate` 录入 `code`/`content` 模板
- **THEN** 返回 `IdReply{id}`，可经 `ListNoticeTemplates` 分页展示、`GetNoticeTemplate` 单查
- **AND** 管理员调 `SaveNoticeTask` 绑定 `templateId` 并设置 `pushTime`/`interval`/`maxTimes`，由调度侧按 `finished` 标记推进

### Requirement: RPC 短信模板与短信平台方法

`Message` 服务 SHALL 提供短信模板 3 个与短信平台 1 个共 4 个方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveMessageTemplate` | `MessageTemplateSaveReq { id, name, platformCode, signName, thirdTemplateCode, content, templateId, status }` | `IdReply { id }` | 新增/更新第三方短信模板 |
| `DeleteMessageTemplate` | `IdReq { id }` | `Empty {}` | 按 ID 删除短信模板 |
| `ListMessageTemplates` | `PageReq { pageNo, pageSize }` | `MessageTemplateListReply { total, list }` | 分页查询短信模板 |
| `ListSmsPlatforms` | `Empty {}` | `SmsPlatformListReply { list }` | 查询全部第三方短信平台（不含 `total`，全量返回） |

关键字段约束：

| 字段 | 类型 | 说明 |
|------|------|------|
| `platformCode` | string | 第三方短信平台代号（逻辑上对应 `sms_third_platform.code`，DDL 无外键约束） |
| `signName` | string | 签名 |
| `thirdTemplateCode` | string | 第三方短信模板 code |
| `templateId` | int64 | 关联的通知模板 ID |
| `status`（模板） | int32 | 模板状态：0-禁用，1-启用 |
| `priority`（平台响应） | int32 | 数字越小优先级越高，最小为 0 |
| `status`（平台响应） | int32 | 短信平台状态：0-禁用，1-启用 |

> 短信模板未提供 `Get` 单查方法，只能通过 `ListMessageTemplates` 分页获取。

#### Scenario: 选择短信平台并绑定模板

- **WHEN** 管理端调 `ListSmsPlatforms`
- **THEN** 取出可用短信平台列表（`SmsPlatformListReply{list}`，按 `priority` 择优见「短信模板与平台规则」）
- **AND** 调 `SaveMessageTemplate` 绑定平台 `platformCode` 与第三方模板 `thirdTemplateCode`，一个通知模板可对应多个平台的短信模板

### Requirement: RPC 用户收件箱与公告方法

`Message` 服务 SHALL 提供用户收件箱 4 个与公告 3 个共 7 个方法，其中 `SendNotice` 是本服务对外的核心写入口：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SendNotice` | `SendNoticeReq { userId, type, title, content, publisher, expireTime }` | `IdReply { id }` | 向指定用户投递一条站内信（仅 RPC，无 HTTP 入口） |
| `ListInbox` | `InboxPageReq { pageNo, pageSize, userId, type }` | `InboxListReply { total, list }` | 分页查询用户收件箱 |
| `MarkInboxRead` | `IdReq { id }` | `Empty {}` | 标记单条站内信已读 |
| `DeleteInbox` | `IdReq { id }` | `Empty {}` | 删除单条站内信 |
| `SavePublicNotice` | `PublicNoticeSaveReq { id, title, content, type, pushTime, expireTime }` | `IdReply { id }` | 新增/更新公告 |
| `DeletePublicNotice` | `IdReq { id }` | `Empty {}` | 按 ID 删除公告 |
| `ListPublicNotices` | `PageReq { pageNo, pageSize }` | `PublicNoticeListReply { total, list }` | 分页查询公告 |

`SendNoticeReq` 关键字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | int64 | 收件用户 ID |
| `type` | int32 | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知，4-私信 |
| `title` | string | 通知标题 |
| `content` | string | 通知或私信内容 |
| `publisher` | int64 | 通知的发送者 ID，0 则代表是系统 |
| `expireTime` | string | 过期时间，一旦过期用户端不再展示 |

- `InboxPageReq.type` 为可选过滤条件，未传时不按类型过滤；`UserInboxVO` 额外返回 `isRead` 与 `pushTime`，二者由服务端写入
- 公告 `content` 可以存放公告消息模板；`type` 取值：0-系统通知，1-笔记通知，2-问答通知，3-其它通知

#### Scenario: 用户读消息全流程

- **WHEN** 业务侧（设计上）调 `SendNotice` 向用户投递站内信
- **THEN** `user_inbox` 写入一行，返回 `IdReply{id}`
- **AND** 前端 `ListInbox` 分页拉取 → 点开后 `MarkInboxRead` 置已读 → 不需要时 `DeleteInbox` 删除

### Requirement: 通知模板管理规则

通知模板是所有通知内容的来源，SHALL 按 `code` 引用、按 `status` 控制是否可用：

- `code` 为业务引用键：DDL 建 `idx_code` 索引，按 code 高频查询，应校验唯一性
- `type` 四类通知：0-系统通知，1-笔记通知，2-问答通知，3-其它通知
- `status` 三态流转：0-草稿（默认）→ 1-使用中 → 2-停用；`status` 不可由请求直接设置（`NoticeTemplateSaveReq` 无该字段），新建落 DDL 默认值 0（草稿），状态迁移需另行提供入口
- 短信模板可不填 title：`is_sms_template = 1` 时放宽 title 必填校验（DDL `title` 可空）
- 物理删除：`DeleteNoticeTemplate` 为真删除，删除前应校验是否被 `notice_task` / `message_template` 引用

> 注：本组规则为依据 proto/DDL/.api 契约推导的设计意图（business-rules.md）；37 个 logic 已全部实现并编译通过，建议对照源码最终确认。

#### Scenario: 保存通知模板流程（SaveNoticeTemplate）

- **WHEN** `SaveNoticeTemplate` 收到请求
- **THEN** 校验 `name` / `code` / `content` 非空，并校验 `code` 唯一性（当前 model 缺 `ExistsByCode`，需补）；有 id 则 Update、无 id 则生成 ID 后 Insert
- **AND** `is_sms_template = false` 时校验 `title` 非空（短信模板放宽）

### Requirement: 通知任务调度规则

通知任务 SHALL 把「通知模板」按时间策略投递出去，支持延迟、周期与部分人群：

- 必须绑定模板：`SaveNoticeTask` 应校验 `templateId > 0` 且对应模板存在、`status = 1`（使用中）（DDL `template_id NOT NULL`）
- `partial` 决定人群范围：`partial=1` 时目标人群写入 `notice_task_target`；`partial=0` 为全站
- `interval` 单位为分钟：非秒、非毫秒，调度器换算需注意
- `max_times` 控制重复次数：DDL 默认 1，注释「1 则只发一次」，达到上限后置 `finished = 1`
- `finished` 由服务端维护：客户端不可直接改写完成状态
- `expire_time` 后不再投递：过期任务应跳过并终结
- 时间字段可空：`push_time` / `expire_time` / `interval` 均 NULL，三者皆空时语义为「立即一次性投递」（推导）

> 注：本组规则为契约推导的设计意图（business-rules.md）；仓库内无通知任务调度器/定时轮询的配置与实现描述。

#### Scenario: 保存通知任务流程（SaveNoticeTask）

- **WHEN** `SaveNoticeTask` 收到请求
- **THEN** 校验 `templateId > 0` 且模板存在、status = 1（使用中），校验 `name` 非空；`pushTime` / `expireTime` 解析为 `time.Time` 写入 `sql.NullTime`
- **AND** `maxTimes <= 0` 时落 DDL 默认值 1；新建时 `finished = 0`

#### Scenario: 周期任务终结

- **WHEN** 任务重复执行次数达到 `max_times` 上限，或任务已过 `expire_time`
- **THEN** 置 `finished = 1`（或跳过并终结），不再调度投递
- **AND** `finished` 由服务端维护，客户端不可经 `SaveNoticeTask` 直接改写

### Requirement: 短信模板与平台规则

短信 SHALL 走第三方平台：`message_template` 保存平台侧的签名与模板 code，`sms_third_platform` 保存平台本身：

- 平台按 priority 择优：数字越小优先级越高，最小为 0；`ListSmsPlatforms` 应按 `priority ASC` 排序返回
- 平台需启用才可用：发送时应过滤 `status = 0` 的平台（平台 `status` 默认 1）
- 短信模板挂靠平台：`message_template.platform_code` ↔ `sms_third_platform.code` 无外键约束，需在 logic 层校验平台存在
- 短信模板关联通知模板：`message_template.template_id` + `idx_template_id`，一个通知模板可对应多个平台的短信模板
- 短信模板默认禁用：`message_template.status` 默认 0（0-禁用，1-启用），与平台表默认值相反，新建后需显式启用
- 无单查接口：详情通过列表返回，无 `GetMessageTemplate`

> 注：本组规则为契约推导的设计意图（business-rules.md）；短信平台是否真正对接第三方发送通道建议对照源码复核。

#### Scenario: 平台择优排序

- **WHEN** 调用 `ListSmsPlatforms`
- **THEN** 返回平台列表按 `priority ASC` 排序（数字越小优先级越高）
- **AND** 发送短信时应过滤 `status = 0` 的禁用平台

#### Scenario: 短信模板挂靠校验

- **WHEN** `SaveMessageTemplate` 提交 `platformCode`
- **THEN** 因 DDL 无外键约束，应在 logic 层校验该 `platformCode` 对应的平台存在
- **AND** 新建短信模板默认 `status = 0`（禁用），需显式启用后才可用

### Requirement: 用户收件箱投递与查询规则

`SendNotice` 是本服务对外的核心写入口，`user_inbox` 一行即一条用户可见消息，SHALL 满足：

- `publisher = 0` 代表系统：系统通知与用户私信共用一张表
- `type = 4` 为私信：0/1/2/3 为通知，4-私信；DDL 默认值即 4
- 过期后端上不展示：`ListInbox` 应带 `expire_time > now()` 过滤
- 按用户 + 时间分页：依赖 `user_id`、`push_time` 两个单列索引，收件箱按 `push_time DESC` 倒序
- `type` 为可选过滤条件：未传时不按类型过滤
- `userId` 由 JWT 注入：API 层从 JWT 取当前用户 ID 填入 RPC 请求
- 读写归属校验：`MarkInboxRead` / `DeleteInbox` 仅传 `id`，必须在 logic 层校验该 `id` 的 `user_id` 等于当前登录用户，否则存在越权风险
- 物理删除：`DeleteInbox` 为真删除，用户删除后不可恢复

> 注：本组规则为契约推导的设计意图（business-rules.md）；`ListInbox` 越权校验是否落地建议对照源码复核。

#### Scenario: 投递站内信流程（SendNotice）

- **WHEN** `SendNotice` 收到请求
- **THEN** 校验 `userId > 0`、`content` 非空；`type` 缺省落 DDL 默认值 4（私信）；`push_time = now()`、`is_read = 0`
- **AND** `expire_time` 解析请求值（缺省策略未定义），Insert `user_inbox` 后返回 `IdReply{id}`

#### Scenario: 已读与删除的归属校验

- **WHEN** 用户对非本人 `user_id` 的收件箱记录调 `MarkInboxRead` / `DeleteInbox`
- **THEN** logic 层应校验该 `id` 的 `user_id` 等于当前登录用户，校验不通过则拒绝
- **AND** 删除为物理 DELETE，不可恢复

### Requirement: 公告展示规则

`public_notice` 是全站公告，SHALL 与用户收件箱解耦：

- 生效窗口：`push_time` / `expire_time` 均 NOT NULL，端上展示需满足 `push_time <= now() < expire_time`
- 内容可存模板：`content` 可以存放公告消息模板，支持占位符渲染（推导）
- 与收件箱解耦：两表无关联字段，公告不产生 `user_inbox` 记录，已读状态无处记录
- 无审计字段：无 `creater`/`updater`/`create_time`/`update_time`，公告变更无法追溯操作人
- 物理删除：`DeletePublicNotice` 为真删除

> 注：本组规则为契约推导的设计意图（business-rules.md）。

#### Scenario: 公告生效窗口展示

- **WHEN** 端上经 `ListPublicNotices` 拉取公告
- **THEN** 仅 `push_time <= now() < expire_time` 窗口内的公告对用户展示
- **AND** 公告不写入任何 `user_inbox` 记录，已读状态无处记录

### Requirement: API 层与 RPC 层关系

message 服务的 API 层 SHALL 为纯转发层：

- API 全部转发 RPC：`apps/message/api/internal/svc/servicecontext.go` 注入 `MessageRpc`，API 层不直连数据库，仅做 DTO 转换
- 全部接口需登录：`message.api` 六个 `@server` 块均声明 `jwt: Auth`，对应 `routes.go` 六处 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)`，无匿名接口
- API 比 RPC 少 1 个方法：HTTP 18 个 vs RPC 19 个，`SendNotice` 仅暴露为 RPC，不对外提供 HTTP 入口
- 分页参数类型不一致：api `PageReq{PageNo,PageSize int64}` vs proto `PageReq{pageNo,pageSize int32}`，API → RPC 转换时需做 int64 → int32 收窄

#### Scenario: 用户身份经 JWT 注入

- **WHEN** 已登录用户调用 GET /inbox
- **THEN** HTTP 侧 `InboxListReq` 无 `userId` 字段，API 层从 JWT 提取当前用户 ID 填入 RPC `InboxPageReq.userId`
- **AND** 分页字段由 `int64` 收窄为 RPC 侧 `int32`

### Requirement: 数据模型表清单与关系

message 服务 SHALL 将数据持久化于 MySQL `tj_message` 库的 7 张表。message 是全仓唯一拥有两份 DDL 的服务：`sql/ddl/tj_message.sql`（7 表，完整库结构，字段注释最完整，字段语义以此为准）与 `sql/ddl/tj_message_model.sql`（6 表，精简版，goctl 模型生成源，剔除了复合主键的 `notice_task_target`）：

| 表 | 职责 | 关键字段与索引 |
|----|------|--------------|
| `notice_template` | 通知模板 | `code`（idx_code，业务引用键）、`type`（0-3）、`status`（0-草稿默认/1-使用中/2-停用）、`title`（可空）、`is_sms_template`（默认 b'0'） |
| `notice_task` | 通知任务 | `template_id`（未建索引）、`partial`、`push_time`/`expire_time`（可空）、`interval`（分钟，可空）、`max_times`（默认 1）、`finished`（默认 false） |
| `notice_task_target` | 通知任务目标用户 | `(task_id, target_id)` 复合主键；仅存在于 `tj_message.sql`，无对应 Model；`partial=1` 时圈定目标人群 |
| `message_template` | 第三方短信模板 | `platform_code`、`sign_name`、`third_template_code`、`template_id`（idx_template_id）、`status`（0-禁用默认/1-启用） |
| `sms_third_platform` | 第三方云通讯平台 | `code`（如 ali）、`priority`（越小越优先，默认 0）、`status`（0-禁用/1-启用默认） |
| `user_inbox` | 用户通知记录（站内信） | `user_id`（索引）、`type`（默认 4）、`is_read`（默认 b'0'）、`publisher`（默认 0）、`push_time`（索引）、`expire_time` |
| `public_notice` | 公告消息模板 | `title`、`content`、`type`（0-3）、`push_time`/`expire_time`（均 NOT NULL）；全库唯一无审计字段表 |

表间关系（无数据库外键，逻辑关联）：

- `notice_template (1) ──(template_id)── (N) message_template`（无外键约束）
- `notice_template (1) ──(template_id)── (N) notice_task`（无外键约束）→ `notice_task (1) ──(task_id)── (N) notice_task_target`（复合主键，无 Model）
- `sms_third_platform (1) ──(code ↔ platform_code)── (N) message_template`（逻辑关联）
- `public_notice`（独立表，全站公告）；`user_inbox`（独立表，`user_id` 指向 user 域，跨库无外键）

Model 层与缓存：

| Model | 对应表 | 缓存前缀 | 扩展方法 |
|-------|--------|---------|---------|
| `NoticeTemplate` | `notice_template` | `cache:noticeTemplate:id:` | 无（空壳） |
| `NoticeTask` | `notice_task` | `cache:noticeTask:id:` | 无（空壳） |
| `MessageTemplate` | `message_template` | `cache:messageTemplate:id:` | 无（空壳） |
| `SmsThirdPlatform` | `sms_third_platform` | `cache:smsThirdPlatform:id:` | 无（空壳） |
| `UserInbox` | `user_inbox` | `cache:userInbox:id:` | 无（空壳） |
| `PublicNotice` | `public_notice` | `cache:publicNotice:id:` | 无（空壳） |
| （缺失） | `notice_task_target` | - | 未生成 Model，复合主键不受 goctl 支持，需手写 SQL 访问 |

> 已知缺口（模型层）：以 data-model.md 为准——6 个自定义 Model 文件（`apps/message/rpc/internal/model/`）当前均为 goctl 生成的空壳，可用方法仅为 `Insert` / `FindOne`（带缓存）/ `Update` / `Delete` 四个基础方法；实现 RPC 分页/列表能力前需补写 `FindPage`（notice_template / notice_task / message_template / public_notice）、`FindPageByUserId`（user_inbox）、`FindAll`（sms_third_platform，按 `priority` 升序）以及 `notice_task_target` 的手写访问层。加注：business-rules.md v1.2（2026-08-06 复核）称 37 个 logic 文件（RPC 19 + API 18）已全部实现且 `go build` 通过，与上述「列表方法尚未定义」的 v1.0 表述存在冲突，以 data-model.md 表述为准、建议对照源码复核。
>
> 已知缺口：`docs/tjxt.openapi.json` 未收录 message 服务任何路径，无法作为交叉参照。

#### Scenario: 收件箱索引支撑分页

- **WHEN** `ListInbox` 按用户分页查询
- **THEN** 依赖 `user_inbox` 的 `user_id` 单列索引过滤归属、`push_time` 单列索引按时间倒序分页（二者非联合索引）
- **AND** 主键查询走 Redis 缓存（前缀 `cache:userInbox:id:`，`sqlc.CachedConn`）

#### Scenario: 部分人群任务的目标表访问

- **WHEN** `notice_task.partial = 1` 需要圈定目标人群
- **THEN** 目标用户写入 `notice_task_target`（`(task_id, target_id)` 复合主键）
- **AND** 该表无对应 goctl Model，需手写 SQL 访问

### Requirement: 数据模型关键不变量

message 服务的存储层 SHALL 满足以下不变量：

- 本库不使用逻辑删除：7 张表均无 `deleted` 字段，删除操作为物理 `DELETE`（与 auth 域的 `deleted=1` 软删除模式不同）
- `sms_third_platform` 是本库唯一使用 `AUTO_INCREMENT` 的表（`AUTO_INCREMENT = 4`），其余表主键由业务侧生成；Go 结构体中 `priority` 为 `uint64`
- `public_notice` 是全库唯一不含 `creater`/`updater`/`create_time`/`update_time` 审计字段的表
- `notice_task_target` 字符集为 `utf8mb4`，与其余 6 张 `utf8` 表不同
- `notice_task.template_id` 在 DDL 中未建索引
- Go 类型映射：`NoticeTemplate.title` 为 `sql.NullString`、`is_sms_template` 为 `byte`；`NoticeTask` 的 `push_time`/`expire_time` 为 `sql.NullTime`、`interval` 为 `sql.NullInt64`、`partial`/`finished` 为 `byte`
- `vars.go` 仅定义 `ErrNotFound = sqlx.ErrNotFound`，无 auth 域那样的 `sqlhelper.go` 通用 SQL 工具函数
- 字段语义以 `tj_message.sql` 为准；模型代码以 `tj_message_model.sql` 为准（其字段 COMMENT 被截短，`*_gen.go` 由其生成且禁止修改）

#### Scenario: 物理删除不可恢复

- **WHEN** 执行 `DeleteInbox` / `DeleteNoticeTemplate` / `DeletePublicNotice` 等删除操作
- **THEN** 均为物理 `DELETE`（表无 `deleted` 字段），删除后不可恢复

#### Scenario: 平台表主键自增

- **WHEN** 插入 `sms_third_platform` 新平台记录
- **THEN** 主键由数据库 `AUTO_INCREMENT` 生成（本库唯一），其余 6 张表主键由业务侧生成

### Requirement: 服务配置与依赖

message 服务 SHALL 按以下配置运行（`apps/message/api/etc/message-api.yaml`、`apps/message/rpc/etc/message.yaml`）：

| 服务 | 配置项 | 默认值 |
|------|--------|--------|
| `message-api` | `Name` / `Host` / `Port` | `message-api` / `0.0.0.0` / `8807`（HTTP） |
| `message-api` | `Auth.AccessSecret` / `Auth.AccessExpire` | `change-me-in-production` / `7200`（秒） |
| `message-api` | `MessageRpc.Etcd.Hosts[0]` / `MessageRpc.Etcd.Key` | `127.0.0.1:2379` / `message.rpc` |
| `message.rpc` | `Name` / `ListenOn` | `message.rpc` / `0.0.0.0:8087`（gRPC） |
| `message.rpc` | `Etcd.Hosts[0]` / `Etcd.Key` | `127.0.0.1:2379` / `message.rpc` |
| `message.rpc` | `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_message?charset=utf8mb4&parseTime=true&loc=Local` |
| `message.rpc` | `Cache[0]` | `127.0.0.1:6379`，Pass 空，Type `node`（单机模式） |

- 依赖的外部服务：MySQL `tj_message` 库（6 个 Model 共用一个 `sqlx.NewMysql(c.DataSource)` 连接）、Redis 缓存（6 个 Model 均走 `sqlc.CachedConn`，前缀 `cache:<model>:id:`）、etcd（RPC 注册与发现）
- API 层不直连存储：`message-api` 未配置 MySQL / Redis，仅依赖自身 RPC `message.rpc` 与 etcd
- JWT 只校验不签发：`Auth.AccessSecret` 必须与 auth 服务的 `Jwt.AccessSecret` 保持一致，否则所有接口 401；`Auth.AccessExpire` 仅作为配置项被反序列化，代码内无签发逻辑引用；生产环境必须修改默认值 `change-me-in-production`
- RPC 层未配置 JWT 段：令牌校验在 API 层的 `rest.WithJwt` 完成，RPC 层不做鉴权（与 auth/user 不同，对照 auth：其 RPC 配置 JWT 段且签发令牌）
- 与 auth 服务无直接连接：仅共享 `AccessSecret` 完成离线 JWT 校验，不调用 auth RPC

> 注（字符集不一致）：连接串字符集为 `utf8mb4`，而 DDL 中 7 张表有 6 张建表字符集为 `utf8`（仅 `notice_task_target` 为 `utf8mb4`），两者不一致。
>
> 可观测性注入配置（Telemetry / Prometheus / Log 三段）：两份 yaml 中已注入勿删（约定见 `openspec/specs/infra/observability`）。

#### Scenario: API 层经 etcd 发现自身 RPC

- **WHEN** message-api 启动并处理首个请求
- **THEN** `MessageRpc` 客户端按 `MessageRpc.Etcd`（`127.0.0.1:2379`，key `message.rpc`）发现 `message.rpc` 并完成转发调用
- **AND** API 层不直连 MySQL / Redis，全部 18 个 logic 转发 RPC

#### Scenario: JWT 密钥不一致导致 401

- **WHEN** `message-api` 的 `Auth.AccessSecret` 与 auth 服务签发令牌所用密钥不一致
- **THEN** 六个 `@server` 块（`rest.WithJwt`）校验失败，所有 HTTP 接口返回 401
- **AND** 生产环境必须修改默认密钥 `change-me-in-production`

