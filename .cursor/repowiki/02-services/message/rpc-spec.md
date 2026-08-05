> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/message/rpc/message.proto`

---

# Message RPC Spec

## 服务名

`Message` — 站内信与通知微服务，通过 etcd 服务发现（key: `message.rpc`），监听 `0.0.0.0:8085`。

## RPC 方法总览

### 通知模板

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveNoticeTemplate` | `NoticeTemplateSaveReq { id, name, code, type, title, content, isSmsTemplate }` | `IdReply { id }` | 新增/更新通知模板 |
| `DeleteNoticeTemplate` | `IdReq { id }` | `Empty {}` | 按 ID 删除通知模板 |
| `GetNoticeTemplate` | `IdReq { id }` | `NoticeTemplateVO { id, name, code, type, status, title, content, isSmsTemplate, createTime }` | 按 ID 查询通知模板 |
| `ListNoticeTemplates` | `PageReq { pageNo, pageSize }` | `NoticeTemplateListReply { total, list }` | 分页查询通知模板 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 通知模板 ID |
| `name` | string | 通知模板名称 |
| `code` | string | 模板代号，例如：verify-code |
| `type` | int32 | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知 |
| `title` | string | 通知标题，短信模板可以不填 |
| `content` | string | 通知内容模板 |
| `isSmsTemplate` | bool | 是否是短信模板 |

> `NoticeTemplateVO.status`（0-草稿，1-使用中，2-停用）仅出现在响应中，`NoticeTemplateSaveReq` 未定义该字段。

---

### 通知任务

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveNoticeTask` | `NoticeTaskSaveReq { id, templateId, name, partial, pushTime, interval, expireTime, maxTimes }` | `IdReply { id }` | 新增/更新通知任务 |
| `DeleteNoticeTask` | `IdReq { id }` | `Empty {}` | 按 ID 删除通知任务 |
| `GetNoticeTask` | `IdReq { id }` | `NoticeTaskVO { id, templateId, name, partial, pushTime, interval, expireTime, maxTimes, finished, createTime }` | 按 ID 查询通知任务 |
| `ListNoticeTasks` | `PageReq { pageNo, pageSize }` | `NoticeTaskListReply { total, list }` | 分页查询通知任务 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `templateId` | int64 | 任务对应的通知模板 ID |
| `name` | string | 任务名称 |
| `partial` | bool | 是否是部分人的通告，默认 false |
| `pushTime` | string | 任务预期执行时间 |
| `interval` | int32 | 任务延迟执行时间间隔，单位是分钟 |
| `expireTime` | string | 任务失效时间 |
| `maxTimes` | int32 | 任务重复执行次数上限，1 则只发一次 |

> `NoticeTaskVO.finished`（任务是否完成）仅出现在响应中，由服务端维护。

---

### 短信模板

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveMessageTemplate` | `MessageTemplateSaveReq { id, name, platformCode, signName, thirdTemplateCode, content, templateId, status }` | `IdReply { id }` | 新增/更新第三方短信模板 |
| `DeleteMessageTemplate` | `IdReq { id }` | `Empty {}` | 按 ID 删除短信模板 |
| `ListMessageTemplates` | `PageReq { pageNo, pageSize }` | `MessageTemplateListReply { total, list }` | 分页查询短信模板 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 模板名称 |
| `platformCode` | string | 第三方短信平台代号 |
| `signName` | string | 签名 |
| `thirdTemplateCode` | string | 第三方短信模板 code |
| `content` | string | 第三方短信模板内容预览 |
| `templateId` | int64 | 关联的通知模板 ID |
| `status` | int32 | 模板状态：0-禁用，1-启用 |

> 短信模板未提供 `Get` 单查方法，只能通过 `ListMessageTemplates` 分页获取。

---

### 短信平台

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ListSmsPlatforms` | `Empty {}` | `SmsPlatformListReply { list }` | 查询全部第三方短信平台 |

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 短信平台 ID |
| `name` | string | 短信平台名称 |
| `code` | string | 短信平台代码，例如：ali |
| `priority` | int32 | 数字越小优先级越高，最小为 0 |
| `status` | int32 | 短信平台状态：0-禁用，1-启用 |

> `SmsPlatformListReply` 不含 `total` 字段，为全量列表返回。

---

### 用户收件箱

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SendNotice` | `SendNoticeReq { userId, type, title, content, publisher, expireTime }` | `IdReply { id }` | 向指定用户投递一条站内信 |
| `ListInbox` | `InboxPageReq { pageNo, pageSize, userId, type }` | `InboxListReply { total, list }` | 分页查询用户收件箱 |
| `MarkInboxRead` | `IdReq { id }` | `Empty {}` | 标记单条站内信已读 |
| `DeleteInbox` | `IdReq { id }` | `Empty {}` | 删除单条站内信 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | int64 | 收件用户 ID |
| `type` | int32 | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知，4-私信 |
| `title` | string | 通知标题 |
| `content` | string | 通知或私信内容 |
| `publisher` | int64 | 通知的发送者 ID，0 则代表是系统 |
| `expireTime` | string | 过期时间，一旦过期用户端不再展示 |

> `UserInboxVO` 额外返回 `isRead`（是否已读）与 `pushTime`（推送时间），二者由服务端写入。

---

### 公告

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SavePublicNotice` | `PublicNoticeSaveReq { id, title, content, type, pushTime, expireTime }` | `IdReply { id }` | 新增/更新公告 |
| `DeletePublicNotice` | `IdReq { id }` | `Empty {}` | 按 ID 删除公告 |
| `ListPublicNotices` | `PageReq { pageNo, pageSize }` | `PublicNoticeListReply { total, list }` | 分页查询公告 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `title` | string | 公告标题 |
| `content` | string | 公告通知内容，可以存放公告消息模板 |
| `type` | int32 | 通知类型：0-系统通知，1-笔记通知，2-问答通知，3-其它通知 |
| `pushTime` | string | 通知发布时间 |
| `expireTime` | string | 通知失效时间 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `message-api` (自身 API 层) | HTTP Handler → `messageclient.Message` RPC | `apps/message/api/internal/svc/servicecontext.go` 注入 `MessageRpc`，18 个 HTTP 接口全部指向自身 RPC |

（注：全仓 Grep `apps/message/rpc/message` / `messageclient` / `MessageRpc`，除 `apps/message/` 自身外**无其它服务引用**。`.cursor/repowiki/03-shared/mq-events.md` 与 `00-architecture/service-topology.md` 中描述的 pay-rpc / trade-rpc / learning-rpc / promotion-rpc → message-rpc 通知链路，目前在代码中尚无对应 import 或 MQ 消费者实现。）

---

## 调用典型场景

1. **管理端维护通知模板** → 后台调 `SaveNoticeTemplate` 录入 `code`/`content` 模板 → `ListNoticeTemplates` 分页展示
2. **配置定时通告** → 管理员调 `SaveNoticeTask` 绑定 `templateId`，设置 `pushTime`/`interval`/`maxTimes` → 由调度侧按 `finished` 标记推进
3. **对接第三方短信** → `ListSmsPlatforms` 取出按 `priority` 排序的可用平台 → `SaveMessageTemplate` 绑定平台 `platformCode` 与 `thirdTemplateCode`
4. **业务事件发通知** → 业务方（设计上为 trade / pay / learning / promotion）调 `SendNotice` 写入 `user_inbox` → 用户端 `ListInbox` 拉取
5. **用户读消息** → 前端调 `ListInbox` 分页 → 点开后 `MarkInboxRead` 置已读 → 不需要时 `DeleteInbox` 删除
6. **发布全站公告** → 运营调 `SavePublicNotice` 写入公告 → 端上 `ListPublicNotices` 按 `pushTime`/`expireTime` 展示

---

## 自定义 Model 方法

`apps/message/rpc/internal/model/` 下 6 个非 `_gen.go` 文件（`messagetemplatemodel.go`、`noticetaskmodel.go`、`noticetemplatemodel.go`、`publicnoticemodel.go`、`smsthirdplatformmodel.go`、`userinboxmodel.go`）**当前均为 goctl 生成的空壳**，接口体内仅内嵌同名小写接口，未追加任何扩展方法：

```go
MessageTemplateModel interface {
    messageTemplateModel
}
```

因此可用方法仅为 goctl 生成的 4 个基础方法（每个 model 一致）：

- `Insert(ctx, data)` — 插入
- `FindOne(ctx, id)` — 按主键查询（带缓存）
- `Update(ctx, data)` — 全字段更新
- `Delete(ctx, id)` — 按主键删除

> proto 中的 `ListNoticeTemplates` / `ListNoticeTasks` / `ListMessageTemplates` / `ListPublicNotices` / `ListInbox` / `ListSmsPlatforms` 所需的分页与条件查询方法（如 `FindPage`、`FindByUserId`）**尚未定义**，实现这些 RPC 前需先在自定义 model 文件中补齐。
