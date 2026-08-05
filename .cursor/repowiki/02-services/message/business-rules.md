> 版本：v1.2 | 更新：2026-08-06 | 来源：2026-08-06 复核（go build 全模块通过 + 逻辑文件清点）

---

# Message Business Rules

## ⚠️ 实现状态

message 服务的业务逻辑**已全部实现**：37 个 logic 文件（RPC 19 + API 18）均已落地，`go build` 通过，全仓 0 处 goctl 占位。下列各方法状态已校正为「已实现」。

### RPC 层（`apps/message/rpc/internal/logic/`）

| 分组 | Logic 文件 | 方法 | 状态 |
|------|-----------|------|------|
| 通知模板 | `savenoticetemplatelogic.go` | `SaveNoticeTemplate` | ✅ 已实现 |
| 通知模板 | `deletenoticetemplatelogic.go` | `DeleteNoticeTemplate` | ✅ 已实现 |
| 通知模板 | `getnoticetemplatelogic.go` | `GetNoticeTemplate` | ✅ 已实现 |
| 通知模板 | `listnoticetemplateslogic.go` | `ListNoticeTemplates` | ✅ 已实现 |
| 通知任务 | `savenoticetasklogic.go` | `SaveNoticeTask` | ✅ 已实现 |
| 通知任务 | `deletenoticetasklogic.go` | `DeleteNoticeTask` | ✅ 已实现 |
| 通知任务 | `getnoticetasklogic.go` | `GetNoticeTask` | ✅ 已实现 |
| 通知任务 | `listnoticetaskslogic.go` | `ListNoticeTasks` | ✅ 已实现 |
| 短信模板 | `savemessagetemplatelogic.go` | `SaveMessageTemplate` | ✅ 已实现 |
| 短信模板 | `deletemessagetemplatelogic.go` | `DeleteMessageTemplate` | ✅ 已实现 |
| 短信模板 | `listmessagetemplateslogic.go` | `ListMessageTemplates` | ✅ 已实现 |
| 短信平台 | `listsmsplatformslogic.go` | `ListSmsPlatforms` | ✅ 已实现 |
| 用户收件箱 | `sendnoticelogic.go` | `SendNotice` | ✅ 已实现 |
| 用户收件箱 | `listinboxlogic.go` | `ListInbox` | ✅ 已实现 |
| 用户收件箱 | `markinboxreadlogic.go` | `MarkInboxRead` | ✅ 已实现 |
| 用户收件箱 | `deleteinboxlogic.go` | `DeleteInbox` | ✅ 已实现 |
| 公告 | `savepublicnoticelogic.go` | `SavePublicNotice` | ✅ 已实现 |
| 公告 | `deletepublicnoticelogic.go` | `DeletePublicNotice` | ✅ 已实现 |
| 公告 | `listpublicnoticeslogic.go` | `ListPublicNotices` | ✅ 已实现 |

### API 层（`apps/message/api/internal/logic/`）

| 分组（package） | Logic 文件 | 方法 | 状态 |
|----------------|-----------|------|------|
| `noticetemplate` | `savenoticetemplatelogic.go` | `SaveNoticeTemplate` | ✅ 已实现 |
| `noticetemplate` | `deletenoticetemplatelogic.go` | `DeleteNoticeTemplate` | ✅ 已实现 |
| `noticetemplate` | `getnoticetemplatelogic.go` | `GetNoticeTemplate` | ✅ 已实现 |
| `noticetemplate` | `listnoticetemplateslogic.go` | `ListNoticeTemplates` | ✅ 已实现 |
| `noticetask` | `savenoticetasklogic.go` | `SaveNoticeTask` | ✅ 已实现 |
| `noticetask` | `deletenoticetasklogic.go` | `DeleteNoticeTask` | ✅ 已实现 |
| `noticetask` | `getnoticetasklogic.go` | `GetNoticeTask` | ✅ 已实现 |
| `noticetask` | `listnoticetaskslogic.go` | `ListNoticeTasks` | ✅ 已实现 |
| `messagetemplate` | `savemessagetemplatelogic.go` | `SaveMessageTemplate` | ✅ 已实现 |
| `messagetemplate` | `deletemessagetemplatelogic.go` | `DeleteMessageTemplate` | ✅ 已实现 |
| `messagetemplate` | `listmessagetemplateslogic.go` | `ListMessageTemplates` | ✅ 已实现 |
| `smsplatform` | `listsmsplatformslogic.go` | `ListSmsPlatforms` | ✅ 已实现 |
| `inbox` | `listinboxlogic.go` | `ListInbox` | ✅ 已实现 |
| `inbox` | `markinboxreadlogic.go` | `MarkInboxRead` | ✅ 已实现 |
| `inbox` | `deleteinboxlogic.go` | `DeleteInbox` | ✅ 已实现 |
| `notice` | `savepublicnoticelogic.go` | `SavePublicNotice` | ✅ 已实现 |
| `notice` | `deletepublicnoticelogic.go` | `DeletePublicNotice` | ✅ 已实现 |
| `notice` | `listpublicnoticeslogic.go` | `ListPublicNotices` | ✅ 已实现 |

### 统计

| 层 | 已实现 | 总计 | 比例 |
|----|--------|------|------|
| RPC (`apps/message/rpc/internal/logic/`) | 19 | 19 | 100% |
| API (`apps/message/api/internal/logic/`) | 18 | 18 | 100% |
| **合计** | **37** | **37** | **100%** |

> 以下各节均为**📋 设计意图（契约推导）**，依据 `apps/message/rpc/message.proto` 注释、`apps/message/api/message.api` 类型定义与 `sql/ddl/tj_message.sql` 表结构推导，**2026-08-06 复核：logic 已全部实现并编译通过；以下规则为依据 proto/DDL/.api 契约推导，建议对照源码最终确认**。`docs/tjxt.openapi.json` 中**未收录** message 服务任何路径，无法作为交叉参照。


## 已知缺口

- 已实现并编译通过；与原设计意图的主要差异建议对照源码复核（如 `SendNotice` 是否真正投递、短信平台是否对接第三方、`ListInbox` 越权校验是否落地）。

---

## 1. 通知模板管理 📋 设计意图（契约推导）

**核心规则**：通知模板是所有通知内容的来源，按 `code` 引用，按 `status` 控制是否可用。

| 规则 | 依据 | 说明 |
|------|------|------|
| `code` 为业务引用键 | DDL `idx_code` 索引 + 注释「模板代号，例如：verify-code」 | 建索引说明按 code 高频查询，应校验唯一性 |
| `type` 四类通知 | DDL 注释 | 0-系统通知，1-笔记通知，2-问答通知，3-其它通知 |
| `status` 三态流转 | DDL 注释，默认 0 | 0-草稿 → 1-使用中 → 2-停用 |
| `status` 不可由请求直接设置 | proto `NoticeTemplateSaveReq` 无 `status` 字段，`NoticeTemplateVO` 有 | 新建落 DDL 默认值 0（草稿），状态迁移需另行提供入口 |
| 短信模板可不填 title | DDL `title` 可空 + 注释「短信模板可以不填」 | `is_sms_template = 1` 时放宽 title 必填校验 |
| 物理删除 | DDL 无 `deleted` 字段 | `DeleteNoticeTemplate` 为真删除，删除前应校验是否被 `notice_task` / `message_template` 引用 |

```
流程（SaveNoticeTemplate，设计意图）:
  1. 校验 name / code / content 非空
  2. 校验 code 唯一性（当前 model 缺 ExistsByCode，需补）
  3. 有 id → Update；无 id → 生成 ID 后 Insert
  4. is_sms_template = false 时校验 title 非空
```

## 2. 通知任务调度 📋 设计意图（契约推导）

**核心规则**：通知任务把「通知模板」按时间策略投递出去，支持延迟、周期与部分人群。

| 规则 | 依据 | 说明 |
|------|------|------|
| 必须绑定模板 | DDL `template_id NOT NULL` | `SaveNoticeTask` 应校验 `templateId` 对应模板存在 |
| `partial` 决定人群范围 | DDL 注释「是否是部分人的通告，默认 false」 | `partial=1` 时目标人群写入 `notice_task_target`；`partial=0` 为全站 |
| `interval` 单位为分钟 | DDL 注释 | 非秒、非毫秒，调度器换算需注意 |
| `max_times` 控制重复次数 | DDL 默认 1，注释「1 则只发一次」 | 达到上限后置 `finished = 1` |
| `finished` 由服务端维护 | proto `NoticeTaskSaveReq` 无该字段，`NoticeTaskVO` 有 | 客户端不可直接改写完成状态 |
| `expire_time` 后不再投递 | DDL 注释「任务失效时间」 | 过期任务应跳过并终结 |
| 时间字段可空 | DDL `push_time` / `expire_time` / `interval` 均 NULL | 三者皆空时语义为「立即一次性投递」（推导） |

```
流程（SaveNoticeTask，设计意图）:
  1. 校验 templateId > 0 且模板存在、status = 1（使用中）
  2. 校验 name 非空
  3. pushTime / expireTime 解析为 time.Time，写入 sql.NullTime
  4. maxTimes <= 0 → 落 DDL 默认值 1
  5. 新建时 finished = 0
```

## 3. 短信模板与平台 📋 设计意图（契约推导）

**核心规则**：短信走第三方平台，`message_template` 保存平台侧的签名与模板 code，`sms_third_platform` 保存平台本身。

| 规则 | 依据 | 说明 |
|------|------|------|
| 平台按 priority 择优 | DDL 注释「数字越小优先级越高，最小为 0」 | `ListSmsPlatforms` 应按 `priority ASC` 排序返回 |
| 平台需启用才可用 | DDL `status` 默认 1，注释「0-禁用，1-启用」 | 发送时应过滤 `status = 0` 的平台 |
| 短信模板挂靠平台 | `message_template.platform_code` ↔ `sms_third_platform.code` | 无外键约束，需在 logic 层校验平台存在 |
| 短信模板关联通知模板 | `message_template.template_id` + `idx_template_id` | 一个通知模板可对应多个平台的短信模板 |
| 短信模板默认禁用 | DDL `status` 默认 0，注释「0-禁用，1-启用」 | 与平台表默认值相反，新建后需显式启用 |
| 无单查接口 | proto 仅有 Save / Delete / List | 详情通过列表返回，无 `GetMessageTemplate` |

## 4. 用户收件箱（站内信） 📋 设计意图（契约推导）

**核心规则**：`SendNotice` 是本服务对外的核心写入口，`user_inbox` 一行即一条用户可见消息。

| 规则 | 依据 | 说明 |
|------|------|------|
| `publisher = 0` 代表系统 | DDL 注释「通知的发送者 id，0 则代表是系统」 | 系统通知与用户私信共用一张表 |
| `type = 4` 为私信 | DDL 注释与默认值 | 0/1/2/3 为通知，4-私信；DDL 默认值即 4 |
| 过期后端上不展示 | DDL 注释「一旦过期用户端不再展示」 | `ListInbox` 应带 `expire_time > now()` 过滤 |
| 按用户 + 时间分页 | `user_id`、`push_time` 两个单列索引 | 收件箱按 `push_time DESC` 倒序 |
| `type` 为可选过滤条件 | proto `InboxPageReq.type`、api `InboxListReq.Type` 带 `optional` | 未传时不按类型过滤 |
| userId 由 JWT 注入 | api `InboxListReq` 无 userId 字段，proto `InboxPageReq` 有 | API 层从 JWT 取当前用户 ID 填入 RPC 请求 |
| 读写归属校验 | `MarkInboxRead` / `DeleteInbox` 仅传 `id` | 必须在 logic 层校验该 `id` 的 `user_id` 等于当前登录用户，否则存在越权风险 |
| 物理删除 | DDL 无 `deleted` 字段 | `DeleteInbox` 为真删除，用户删除后不可恢复 |

```
流程（SendNotice，设计意图）:
  1. 校验 userId > 0、content 非空
  2. type 缺省落 DDL 默认值 4（私信）
  3. push_time = now()，is_read = 0
  4. expire_time 解析请求值；缺省策略未定义
  5. Insert user_inbox → 返回 IdReply{id}
```

## 5. 公告 📋 设计意图（契约推导）

**核心规则**：`public_notice` 是全站公告，不落到每个用户的收件箱。

| 规则 | 依据 | 说明 |
|------|------|------|
| 生效窗口 | DDL `push_time` / `expire_time` 均 `NOT NULL` | 端上展示需满足 `push_time <= now() < expire_time` |
| 内容可存模板 | DDL 表注释「公告消息模板」、字段注释「可以存放公告消息模板」 | 支持占位符渲染（推导） |
| 与收件箱解耦 | 两表无关联字段 | 公告不产生 `user_inbox` 记录，已读状态无处记录 |
| 无审计字段 | DDL 无 creater/updater/create_time/update_time | 公告变更无法追溯操作人 |
| 物理删除 | DDL 无 `deleted` 字段 | `DeletePublicNotice` 为真删除 |

## 6. API 层与 RPC 层的关系 📋 设计意图（契约推导）

| 规则 | 依据 | 说明 |
|------|------|------|
| API 全部转发 RPC | `apps/message/api/internal/svc/servicecontext.go` 注入 `MessageRpc` | API 层不直连数据库，仅做 DTO 转换 |
| 全部接口需登录 | `message.api` 六个 `@server` 块均声明 `jwt: Auth`，`routes.go` 六处 `rest.WithJwt` | 无匿名接口 |
| API 比 RPC 少 1 个方法 | API 18 个 vs RPC 19 个 | `SendNotice` 仅暴露为 RPC，不对外提供 HTTP 入口 |
| 分页参数类型不一致 | api `PageReq{PageNo,PageSize int64}` vs proto `PageReq{pageNo,pageSize int32}` | API → RPC 转换时需做 int64 → int32 收窄 |

---

## 状态说明

### 通知类型 `type`

| 值 | 含义 | 适用表 |
|----|------|--------|
| 0 | 系统通知 | notice_template / public_notice / user_inbox |
| 1 | 笔记通知 | notice_template / public_notice / user_inbox |
| 2 | 问答通知 | notice_template / public_notice / user_inbox |
| 3 | 其它通知 | notice_template / public_notice / user_inbox |
| 4 | 私信 | 仅 user_inbox（且为默认值） |

### 通知模板状态 `notice_template.status`

| 值 | 含义 |
|----|------|
| 0 | 草稿（默认） |
| 1 | 使用中 |
| 2 | 停用 |

### 短信模板状态 `message_template.status`

| 值 | 含义 |
|----|------|
| 0 | 禁用（默认） |
| 1 | 启用 |

### 短信平台状态 `sms_third_platform.status`

| 值 | 含义 |
|----|------|
| 0 | 禁用 |
| 1 | 启用（默认） |

### 通知任务标记

| 字段 | 值 | 含义 |
|------|----|------|
| `partial` | 0 | 全站通告（默认） |
| `partial` | 1 | 部分人通告，人群见 `notice_task_target` |
| `finished` | 0 | 未完成（默认） |
| `finished` | 1 | 已完成，不再调度 |

### 收件箱已读标记 `user_inbox.is_read`

| 值 | 含义 |
|----|------|
| 0 | 未读（默认） |
| 1 | 已读 |
