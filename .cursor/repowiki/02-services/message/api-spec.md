# message 服务 HTTP API 接口规格

> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/message/api/message.api`（`docs/tjxt.openapi.json` 未收录 message 服务，路径与类型均取自 `.api` 源文件，并经 `apps/message/api/internal/handler/routes.go` 核对）

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| GET | /inbox | 分页查询当前用户收件箱 | 是 | InboxListReq (form) | InboxListVO | inbox |
| DELETE | /inbox/:id | 删除一条站内信 | 是 | IdPathReq | OkVO | inbox |
| PUT | /inbox/:id/read | 标记站内信为已读 | 是 | IdPathReq | OkVO | inbox |
| GET | /message-templates | 分页查询短信模板 | 是 | PageReq (form) | MessageTemplateListVO | messagetemplate |
| POST | /message-templates | 新增/更新短信模板 | 是 | MessageTemplateSaveReq | IdVO | messagetemplate |
| DELETE | /message-templates/:id | 根据 id 删除短信模板 | 是 | IdPathReq | OkVO | messagetemplate |
| GET | /notice-tasks | 分页查询通知任务 | 是 | PageReq (form) | NoticeTaskListVO | noticetask |
| POST | /notice-tasks | 新增/更新通知任务 | 是 | NoticeTaskSaveReq | IdVO | noticetask |
| DELETE | /notice-tasks/:id | 根据 id 删除通知任务 | 是 | IdPathReq | OkVO | noticetask |
| GET | /notice-tasks/:id | 根据 id 查询通知任务 | 是 | IdPathReq | NoticeTaskVO | noticetask |
| GET | /notice-templates | 分页查询通知模板 | 是 | PageReq (form) | NoticeTemplateListVO | noticetemplate |
| POST | /notice-templates | 新增/更新通知模板 | 是 | NoticeTemplateSaveReq | IdVO | noticetemplate |
| DELETE | /notice-templates/:id | 根据 id 删除通知模板 | 是 | IdPathReq | OkVO | noticetemplate |
| GET | /notice-templates/:id | 根据 id 查询通知模板 | 是 | IdPathReq | NoticeTemplateVO | noticetemplate |
| GET | /notices | 分页查询公告 | 是 | PageReq (form) | PublicNoticeListVO | notice |
| POST | /notices | 新增/更新公告 | 是 | PublicNoticeSaveReq | IdVO | notice |
| DELETE | /notices/:id | 根据 id 删除公告 | 是 | IdPathReq | OkVO | notice |
| GET | /sms-platforms | 查询全部第三方短信平台 | 是 | - | SmsPlatformListVO | smsplatform |

> 「权限标签」列填写的是 `.api` 中 `@server` 块的 `group` 值（同时也是 handler / logic 的包名），`message.api` 未声明任何 RBAC 权限注解。
> 「认证」列全部为「是」：六个 `@server` 块均声明 `jwt: Auth`，对应 `routes.go` 中 6 处 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)`。

## 请求/响应类型速查

| 类型 | 字段 |
|------|------|
| `IdPathReq` | `Id int64` (path:"id") |
| `IdVO` | `Id int64` |
| `OkVO` | `Success bool` |
| `PageReq` | `PageNo int64`(optional), `PageSize int64`(optional) |
| `InboxListReq` | `PageNo int64`(optional), `PageSize int64`(optional), `Type int32`(optional) |
| `NoticeTemplateSaveReq` | `Id`(optional), `Name`, `Code`, `Type`, `Title`(optional), `Content`, `IsSmsTemplate`(optional) |
| `NoticeTaskSaveReq` | `Id`(optional), `TemplateId`, `Name`, `Partial`(optional), `PushTime`(optional), `Interval`(optional), `ExpireTime`(optional), `MaxTimes`(optional) |
| `MessageTemplateSaveReq` | `Id`(optional), `Name`, `PlatformCode`, `SignName`, `ThirdTemplateCode`, `Content`, `TemplateId`, `Status`(optional) |
| `PublicNoticeSaveReq` | `Id`(optional), `Title`, `Content`, `Type`, `PushTime`, `ExpireTime` |
| `UserInboxVO` | `Id`, `UserId`, `Type`, `Title`, `Content`, `IsRead`, `Publisher`, `PushTime`, `ExpireTime` |
| `SmsPlatformVO` | `Id`, `Name`, `Code`, `Priority`, `Status` |

> 列表类响应 `XxxListVO` 统一为 `{ Total int64, List []XxxVO }`；`SmsPlatformListVO` 例外，仅有 `{ List []SmsPlatformVO }`，无 `Total`。

## 与 RPC 层的差异

| 项 | 说明 |
|----|------|
| 接口数量 | HTTP 18 个 vs RPC 19 个，`SendNotice` 仅存在于 RPC，不对外暴露 HTTP 入口 |
| 用户身份 | HTTP 侧 `InboxListReq` 无 `userId`，需由 API 层从 JWT 提取后填入 RPC 的 `InboxPageReq.userId` |
| 分页字段类型 | HTTP 侧 `PageNo`/`PageSize` 为 `int64`，RPC 侧为 `int32`，转换时需收窄 |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../../03-shared/error-codes.md)
