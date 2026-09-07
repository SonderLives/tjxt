# media（媒资文件）Specification

## Purpose

media（媒资文件）是媒资与文件管理微服务（RPC 服务名 `Media`），职责为「媒资视频管理 + 普通文件/图片管理 + 对象存储签名授权（上传/预览/播放）」。HTTP API 层 `media-api` 监听 `0.0.0.0:8806`，为自身 RPC 的转发层；RPC 服务 `Media` 监听 `0.0.0.0:8086`（gRPC），经 etcd 以 key `media.rpc` 注册与发现。数据存储于 MySQL `tj_media` 库（`media` / `file` 两张表，无数据库外键），主键查询由 Redis 缓存托管（node 单机模式）。course 服务与 media 为双向逻辑依赖（课程侧持有 `media_id` 引用媒资，media 侧 `MediaVO.useTimes` 需 course 回填引用计数），但跨服务 RPC 尚未接线；对象存储后端当前为 mock（指向本地 `http://127.0.0.1:9000`）。

> 实现状态（2026-08-06 复核，business-rules.md v1.2）：本服务业务 logic 已全部实现——20 个 logic 文件（RPC 10 + API 10）均已落地并编译通过（go build 全模块通过）；对象存储后端为 mock（见各处「已知缺口」）。业务规则各条为依据 `media.proto` / `tj_media.sql` / `media.api` / `tjxt.openapi.json` 契约推导的设计意图，建议对照源码最终确认。

## Requirements

### Requirement: 文件管理 HTTP 接口

> **状态说明**：来自聚合文档 `docs/tjxt.openapi.json`（api-spec.md，最后同步 2026-08-05）。源表「认证」列均标注「否」、「权限标签」列均为空，为原始聚合文档标注，本规格按源文件忠实转写；但 configs.md 指出 `media.api` 三组路由（`media` / `signature` / `file`）均声明了 `jwt: Auth`（两处冲突，见「服务配置与端口」）。

media 服务 SHALL 暴露以下文件管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /files | 上传文件 | JSON | R{data: R«FileDTO»} |
| GET | /files/{id} | 获取文件信息 | - | R{data: R«FileDTO»} |
| DELETE | /files/{id} | 删除文件 | - | R{data: R} |

统一约定（引用全局规范）：响应格式为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页为 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码见共享错误码表（`specs/shared/error-codes`）。

#### Scenario: 上传文件

- **WHEN** 客户端以 JSON 请求体调用 `POST /files`
- **THEN** 系统保存文件记录并返回 `R{data: R«FileDTO»}`（落库与状态流转规则见「文件保存、删除与状态流转规则」）

#### Scenario: 删除文件

- **WHEN** 客户端调用 `DELETE /files/{id}`
- **THEN** 系统按文件删除规则处理并返回 `R{data: R}`
- **AND** `status=3`（已使用）的文件应拒绝删除（设计意图，见「文件保存、删除与状态流转规则」）

### Requirement: 媒资管理 HTTP 接口

> **状态说明**：同上，源自 `docs/tjxt.openapi.json`（api-spec.md）。

media 服务 SHALL 暴露以下媒资管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /medias | 分页搜索已上传媒资信息 | - | R{data: R«PageDTO«MediaVO»»} |
| POST | /medias | 上传视频后保存媒资信息 | MediaUploadResultDTO | R{data: R«MediaDTO»} |
| DELETE | /medias | 批量删除媒资视频 | - | R{data: R} |
| DELETE | /medias/{mediaId} | 删除媒资视频 | - | R{data: R} |

#### Scenario: 分页搜索已上传媒资

- **WHEN** 管理端调用 `GET /medias` 并携带分页参数
- **THEN** 系统返回 `R{data: R«PageDTO«MediaVO»»}`（模糊搜索与排序规则见「媒资列表查询规则」）

#### Scenario: 上传视频后保存媒资信息

- **WHEN** 客户端以 `MediaUploadResultDTO`（视频上传结果）调用 `POST /medias`
- **THEN** 系统保存媒资记录并返回 `R{data: R«MediaDTO»}`（新增/更新语义与幂等缺口见「媒资保存与状态流转规则」）

#### Scenario: 删除媒资视频

- **WHEN** 客户端调用 `DELETE /medias/{mediaId}` 删除单个媒资，或 `DELETE /medias` 批量删除
- **THEN** 系统按媒资删除规则处理并返回 `R{data: R}`（逻辑删除 `deleted=1`；被课程引用的媒资应拒绝删除，见「服务消费与跨域引用计数」）

### Requirement: 签名授权 HTTP 接口

> **状态说明**：同上，源自 `docs/tjxt.openapi.json`（api-spec.md）。
>
> 已知缺口：对象存储后端为 mock——`config.go` / `etc/*.yaml` 无 `SecretId/SecretKey/Bucket` 配置，签名/上传/播放均指向本地 `http://127.0.0.1:9000`（`mockBaseURL`）；接入真实 COS/OSS 前媒资全链路仅为本地可跑的桩。

media 服务 SHALL 暴露以下签名授权接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /medias/signature/upload | 获取上传视频的授权签名 | - | R{data: R«string»} |
| GET | /medias/signature/preview | 管理端获取预览视频的授权签名 | - | R{data: R«VideoPlayVO»} |
| GET | /medias/signature/play | 获取播放视频的授权签名 | - | R{data: R«VideoPlayVO»} |

#### Scenario: 获取上传视频的授权签名

- **WHEN** 客户端在发起视频上传前调用 `GET /medias/signature/upload`
- **THEN** 系统返回 `R{data: R«string»}`（签名令牌，用于客户端直传对象存储；规则见「签名授权规则」）

#### Scenario: 获取播放/预览的授权签名

- **WHEN** 学员端调用 `GET /medias/signature/play` 播放视频，或管理端调用 `GET /medias/signature/preview` 预览视频
- **THEN** 系统返回 `R{data: R«VideoPlayVO»}`（带鉴权的播放地址，签名应限时）

### Requirement: 媒资管理 RPC 契约

media 服务 SHALL 以 gRPC 服务名 `Media` 注册于 etcd（key: `media.rpc`），提供以下媒资管理方法（`apps/media/rpc/media.proto`）：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `MediaGet` | `MediaIdRequest { mediaId }` | `MediaVO { id, filename, mediaUrl, coverUrl, duration, size, status, creater, createTime, useTimes }` | 按媒资 ID 查询单条媒资 |
| `MediaList` | `MediaListRequest { pageNo, pageSize, name, sortBy, isAsc }` | `MediaListReply { total, list }` | 分页搜索媒资列表 |
| `MediaSave` | `MediaSaveRequest { id, filename, duration, size, fileId }` | `MediaIdReply { id }` | 视频上传完成后保存媒资信息 |
| `MediaDelete` | `MediaIdRequest { mediaId }` | `Empty {}` | 删除媒资视频 |

请求/响应关键字段：

- `mediaId`（int64）媒资 ID；`id`（int64）媒资 ID，新增时省略
- `filename`（string）文件名称；`fileId`（string）文件在云端的唯一标示（对应 `media.file_id`）
- `duration`（double）视频时长，单位秒；`size`（int64）视频大小，单位字节
- `name`（string）列表按文件名模糊搜索；`sortBy`（string）排序字段；`isAsc`（string，**非 bool**）是否升序
- `MediaVO.status`（int32）：1-上传中，2-已上传；`MediaVO.useTimes`（int32）被引用次数（由课程域回填，不落表）

#### Scenario: 按媒资 ID 查询

- **WHEN** 消费方以 `MediaIdRequest{ mediaId }` 调用 `MediaGet`
- **THEN** 服务返回 `MediaVO`（含 `mediaUrl` / `coverUrl` 播放与封面地址、`status` 状态、`useTimes` 引用次数）

#### Scenario: 视频上传完成后保存媒资

- **WHEN** 视频直传对象存储完成后，消费方以 `MediaSaveRequest{ filename, duration, size, fileId }` 调用 `MediaSave`
- **THEN** 服务落库媒资记录并返回 `MediaIdReply{ id }`（新增/更新语义见「媒资保存与状态流转规则」）

### Requirement: 签名授权 RPC 契约

media 服务的 RPC 层 SHALL 提供三个签名方法，共用同一组 `SignatureRequest { mediaId, fileName, mediaType }` / `SignatureVO { token, url, uploadUrl, playUrl }`，由方法语义区分用途：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SignatureUpload` | `SignatureRequest` | `SignatureVO` | 获取上传视频的授权签名 |
| `SignaturePreview` | `SignatureRequest` | `SignatureVO` | 管理端获取预览视频的授权签名 |
| `SignaturePlay` | `SignatureRequest` | `SignatureVO` | 学员端获取播放视频的授权签名 |

字段说明：请求 `mediaId`（int64，上传场景可为 0）、`fileName`（string）、`mediaType`（string）；响应 `token`（云端签名令牌）、`url`（资源访问地址）、`uploadUrl`（上传地址，上传签名场景）、`playUrl`（播放地址，播放/预览签名场景）。

#### Scenario: 上传签名（mediaId 可为 0）

- **WHEN** 客户端以 `SignatureRequest{ fileName, mediaType }`（`mediaId=0`）调用 `SignatureUpload`
- **THEN** 服务返回 `SignatureVO`，关键字段为 `token` 与 `uploadUrl`（用于客户端直传对象存储）

#### Scenario: 播放与预览签名

- **WHEN** 学员端调用 `SignaturePlay` 播放已发布视频，或管理端调用 `SignaturePreview` 预览未发布视频
- **THEN** 服务返回 `SignatureVO`，关键字段为 `token` 与 `playUrl`（带鉴权的限时播放地址）

### Requirement: 文件管理 RPC 契约

media 服务的 RPC 层 SHALL 提供以下文件管理方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `FileGet` | `FileIdRequest { id }` | `FileVO { id, key, filename, path, status }` | 按 ID 获取文件信息 |
| `FileSave` | `FileSaveRequest { filename, key, size }` | `FileIdReply { id }` | 上传文件后保存文件记录 |
| `FileDelete` | `FileIdRequest { id }` | `Empty {}` | 删除文件 |

字段说明：请求 `id`（int64）文件 ID、`filename`（string）文件上传时的名称、`key`（string）文件在云端的唯一标示（如 `aaa.jpg`）、`size`（int64）文件大小；响应 `FileVO.path` 为文件访问路径、`FileVO.status`（int32）：1-待上传，2-已上传未使用，3-已使用。

> 已知缺口：proto 的 `FileVO.path` 在 DDL 中**无对应列**，需由 `key` 拼接对象存储访问前缀在应用层组装（当前无 CDN/对象存储域名配置，见「服务配置与端口」）。proto 中另声明了 `OkVO {}` 空消息，当前 `service Media` 未引用，删除类方法统一返回 `Empty`。

#### Scenario: 上传文件后保存文件记录

- **WHEN** 客户端直传文件完成后以 `FileSaveRequest{ filename, key, size }` 调用 `FileSave`
- **THEN** 服务落库文件记录并返回 `FileIdReply{ id }`，`status` 由 1（待上传）流转到 2（已上传未使用）

#### Scenario: 按 ID 获取文件信息

- **WHEN** 消费方以 `FileIdRequest{ id }` 调用 `FileGet`
- **THEN** 服务返回 `FileVO`，其中 `path` 由 `key` 拼接对象存储前缀在应用层组装（表中无 `path` 列）

### Requirement: 服务消费与跨域引用计数

media 服务的消费关系与跨域引用计数 SHALL 满足：

**已接线消费方**：`media-api`（自身 API 层）——`apps/media/api/internal/svc/servicecontext.go` 中 import `mediaclient "tjxt/apps/media/rpc/media"`，所有媒资/文件/签名 HTTP 接口最终走自身 RPC（HTTP Handler → `mediaclient.Media` RPC，经 etcd key `media.rpc` 发现）。

**跨域引用计数规则**（设计意图，契约推导）：

- `MediaVO.useTimes` 表示媒资被课程引用的次数，**不落 media 表**，需调 course 域 `Course.CourseMediaUseInfo(MediaIdsRequest)` 取 `MediaQuoteList{ media_id, quote_num }` 在应用层聚合
- 被引用的媒资不可删除：`quote_num > 0` 时 `MediaDelete` 应拒绝（与 `file.status=3` 语义一致）
- 绑定动作在 course 侧：`Course.CourseMediaSave` 写入 `media_id`，文件 `status` 应流转为 3（已使用）

> 已知缺口（未接线）：course 与 media 之间是双向逻辑依赖——`apps/course/api/internal/svc/servicecontext.go` **未** import `mediaclient`，`apps/course/api/etc/*.yaml` 亦无 `MediaRpc` 配置项；`apps/media/*/internal/svc/servicecontext.go` 也未持有 `CourseRpc`。跨服务调用尚未接线，两侧目前均以本地字段/占位方式表达（course 侧 `course.proto:52` 定义 `CourseMediaSave`、`course.proto:36` 定义 `CourseMediaUseInfo` 回传 `media_id` → `quote_num`，`sql/ddl/tj_course.sql:91,123` 两处 `media_id` 列）。

#### Scenario: 媒资列表回填引用计数（设计意图）

- **WHEN** 消费方调用 `MediaList` / `MediaGet` 返回 `MediaVO`
- **THEN** `useTimes` 字段应由 course 域 `CourseMediaUseInfo` 回传的 `quote_num` 聚合，而非读取 media 表
- **AND** 当前跨域 RPC 未接线，引用计数尚无法实际回填

#### Scenario: 被引用的媒资不可删除（设计意图）

- **WHEN** 对 `quote_num > 0`（被课程引用）的媒资调用 `MediaDelete`
- **THEN** 服务应拒绝删除（设计意图；当前因接线缺口该校验无法跨域取数）

### Requirement: 文件保存、删除与状态流转规则

`file.status` SHALL 按 DDL 注释三态流转：`1-待上传` → `2-已上传,未使用` → `3-已使用`；删除 SHALL 为逻辑删除（`deleted=1`）。规则组：

| 规则 | 说明 |
|------|------|
| 初始态为待上传 | 申请上传时先落库 `status=1` |
| 上传完成置为未使用 | 对象存储回调后 `status=2` |
| 被业务引用置为已使用 | 被课程/用户头像等引用后 `status=3` |
| 平台区分 | `platform` 1-腾讯 / 2-阿里，默认 1 |
| 逻辑删除 | 删除时 `deleted=1`，而非物理删除 |

状态流转（设计意图）：`1 待上传 ──上传成功──→ 2 已上传未使用 ──被引用──→ 3 已使用`；状态 1、2 可删除，状态 3（已使用）应拒绝删除，2 与 3 均可进入 `deleted=1`。

流程（FileSave，设计意图）：校验 `filename` / `key` 非空 → 生成 ID，落库 `status=1`（待上传）或 2（已上传未使用）→ `platform` 取默认 1（腾讯）→ 返回 `FileIdReply{ id }`。

流程（FileDelete，设计意图）：`FindOne` 校验存在且 `deleted=0` → `status=3`（已使用）时应拒绝删除 → 置 `deleted=1`，并清理对象存储中的 `key`。

> 已知缺口：生成的 `FileModel.Delete` 为**物理删除**，与 `deleted` 列的逻辑删除语义冲突，实现时需改用自定义 `SoftDelete`；当前自定义 model 为 goctl 空壳（见「数据模型」待补齐方法表）。

#### Scenario: 保存文件记录（FileSave）

- **WHEN** 客户端上传文件后以非空 `filename` / `key` 调用 `FileSave`
- **THEN** 系统生成 ID 落库，`status` 置 1（待上传）或 2（已上传未使用），`platform` 取默认 1（腾讯），返回 `FileIdReply{ id }`
- **AND** 文件上传成功后 `status` 流转为 2；被课程/用户头像等业务引用后流转为 3

#### Scenario: 删除文件（FileDelete）

- **WHEN** 对存在且 `deleted=0` 的文件调用 `FileDelete`
- **THEN** 若 `status=3`（已使用）系统应拒绝删除；否则置 `deleted=1` 并清理对象存储中的 `key`

### Requirement: 媒资保存与状态流转规则

`media.status` SHALL 按 DDL 注释两态流转：`1-上传中`（DDL DEFAULT 1）→ `2-已上传`。规则组：

| 规则 | 说明 |
|------|------|
| 默认上传中 | 建记录时 `status=1`（DDL DEFAULT 1） |
| 转码/上传完成置为已上传 | 回调后 `status=2`，回填 `media_url` / `cover_url` |
| `file_id` 关联云端 | `MediaSaveRequest.fileId` 写入 `media.file_id` |
| 时长与大小由客户端上报 | `duration`(秒) / `size`(字节) 随 `MediaSave` 传入 |
| 逻辑删除 | `MediaDelete` 应置 `deleted=1` |

`MediaSave` 流程（设计意图）：`MediaSaveRequest.id` 为 optional，沿用项目内 `id <= 0` 为新增的约定（参考 auth 域 `SaveRole`）——

1. 校验 `filename` / `fileId` 非空
2. `id <= 0` → 新增（生成雪花 ID，`status=1`，`creater` 取 JWT userId）；`id > 0` → 更新（校验记录存在且 `deleted=0`）
3. 写入 `filename` / `duration` / `size` / `file_id`
4. 返回 `MediaIdReply{ id }`

> 已知缺口（幂等）：同一 `fileId` 重复调用 `MediaSave` 会产生重复媒资记录，需补 `FindOneByFileId` 校验（见「数据模型」待补齐方法表）。

#### Scenario: 新增媒资（id <= 0）

- **WHEN** 客户端以非空 `filename` / `fileId` 且 `id <= 0`（省略）调用 `MediaSave`
- **THEN** 系统生成雪花 ID 新增记录，`status=1`（上传中），`creater` 取 JWT userId，返回 `MediaIdReply{ id }`

#### Scenario: 更新媒资（id > 0）

- **WHEN** 客户端以 `id > 0` 调用 `MediaSave`
- **THEN** 系统校验记录存在且 `deleted=0` 后更新 `filename` / `duration` / `size` / `file_id`，返回 `MediaIdReply{ id }`

#### Scenario: 转码完成后置为已上传（设计意图）

- **WHEN** 视频转码/上传完成回调到达
- **THEN** 系统置 `media.status=2`（已上传），并回填 `media_url` / `cover_url`

### Requirement: 媒资列表查询规则

`MediaList` SHALL 支持模糊搜索 + 动态排序 + 分页（`MediaListRequest`）：

| 字段 | 说明 |
|------|------|
| `pageNo` / `pageSize` | 分页参数，`.api` 中均为 optional，需在 logic 内兜底默认值 |
| `name` | 按 `media.filename` 模糊匹配 |
| `sortBy` | 排序字段名 |
| `isAsc` | 是否升序（proto 中为 **string** 类型，非 bool） |

> ⚠️ 安全约束：`sortBy` 直接来自客户端，拼接 `ORDER BY` 时**必须走白名单**，否则构成 SQL 注入。
>
> 已知缺口（能力）：`MediaModel` 无 `FindPage`，分页查询无法实现（自定义 model 为 goctl 空壳，见「数据模型」）。

#### Scenario: 模糊搜索与动态排序分页

- **WHEN** 管理端以 `name`（文件名模糊）、`sortBy`（白名单内的排序字段）、`isAsc`（升/降序）与 `pageNo` / `pageSize` 调用 `MediaList`
- **THEN** 系统按 `media.filename` 模糊匹配、按指定字段与方向排序后分页返回 `MediaListReply{ total, list }`
- **AND** `pageNo` / `pageSize` 缺省时由 logic 兜底默认值；`sortBy` 未命中白名单时应拒绝

### Requirement: 签名授权规则

三个签名方法 SHALL 共用 `SignatureRequest` / `SignatureVO`，按用途返回不同字段组合，并满足以下规则：

| 方法 | 用途 | 预期返回的关键字段 |
|------|------|------------------|
| `SignatureUpload` | 客户端直传对象存储 | `token`, `uploadUrl` |
| `SignaturePreview` | 管理端预览未发布视频 | `token`, `playUrl` |
| `SignaturePlay` | 学员端播放已发布视频 | `token`, `playUrl` |

规则组：

- 签名有时效：播放/预览签名应限时，避免地址外泄长期可用（通用对象存储实践）
- 上传签名不需要 mediaId：`SignatureRequest.mediaId` 在上传场景可为 0
- 平台适配：需按 `file.platform`（1-腾讯 / 2-阿里）选择不同签名算法

> 已知缺口（配置阻塞项）：`apps/media/rpc/internal/config/config.go` 仅有 `RpcServerConf` / `DataSource` / `Cache` 三项，**没有任何对象存储的 AppId / SecretId / SecretKey / Bucket / Region / 上传路径配置**，`apps/media/rpc/etc/media.yaml` 中亦无对应条目；当前对象存储为 mock（`mockBaseURL` 指向 `http://127.0.0.1:9000`）。建议补齐的配置结构见「服务配置与端口」。

#### Scenario: 客户端直传视频全链路（设计意图）

- **WHEN** 前端上传视频
- **THEN** 前端先调 `SignatureUpload` 获取上传签名（`token` + `uploadUrl`），直传对象存储，完成后回调 `MediaSave` 落库媒资记录（`status=2` 已上传）
- **AND** 当前对象存储后端为 mock，签名/上传均指向本地 `http://127.0.0.1:9000`

#### Scenario: 学员播放与管理端预览（设计意图）

- **WHEN** 学员在小节页播放视频（`SignaturePlay`），或管理员在媒资列表点击预览（`SignaturePreview`）
- **THEN** 服务按 `file.platform` 选择签名算法，换取带鉴权的限时 `playUrl`

### Requirement: 数据模型

media 服务的数据层 SHALL 基于 MySQL `tj_media` 库的两张表，并通过 goctl 生成的 model 层访问：

| 表 | 用途 | 关键字段 |
|----|------|---------|
| `file` | 文件表（普通文件、图片等） | `id`（bigint PK）、`key`（varchar(255)，文件在云端的唯一标示，如 `aaa.jpg`，业务唯一键）、`filename`（varchar(255)，上传时名称）、`request_id`（varchar(64)，可空）、`status`（tinyint，1-待上传 / 2-已上传未使用 / 3-已使用）、`platform`（tinyint，1-腾讯 / 2-阿里，默认 1）、`create_time` / `update_time`（datetime，自动填充/更新）、`creater` / `updater` / `dep_id`（bigint，默认 0）、`deleted`（tinyint，0=正常, 1=逻辑删除，默认 0，过滤条件） |
| `media` | 媒资表（主要是视频文件） | `id`（bigint PK）、`file_id`（varchar(32)，文件在云端的唯一标示，如 `387702302659783576`，业务外部键）、`filename`（varchar(255)，默认 `''`，列表模糊搜索）、`media_url` / `cover_url`（varchar(255)，默认 `''`）、`duration`（double，秒，默认 0）、`size`（bigint，字节，默认 0）、`request_id`（varchar(32)，默认 `''`）、`status`（tinyint，1-上传中 / 2-已上传，默认 1）、`create_time` / `update_time`、`creater` / `updater` / `dep_id`、`deleted`（逻辑删除） |

存储与模型不变量：

- `file` 与 `media` 之间**无数据库外键关联**：`file` 记录普通文件/图片，`media` 记录视频媒资，两者走不同的上传通道（`FileSave` vs `MediaSave`）
- 跨域引用无数据库外键，靠应用层维护：`course.section.media_id ──→ media.id`（`sql/ddl/tj_course.sql:91`）、`course.catalogue.media_id ──→ media.id`（`sql/ddl/tj_course.sql:123`）、`media.useTimes ←── course.CourseMediaUseInfo`
- `FileVO.path` 与 `MediaVO.useTimes` 在表中**无对应列**，分别由 `key` 拼接对象存储访问前缀、由 course 域回传 `quote_num` 在应用层组装/聚合
- 缓存 key 前缀（goctl 生成）：`file` → `cache:file:id:`，`media` → `cache:media:id:`
- Go 结构体映射：`*_gen.go` 与表字段一一对应（`request_id` 映射可空类型 `sql.NullString`，tinyint 映射 `int64`，`media.duration` 映射 `float64`）
- 模型扩展模式：所有 `*_gen.go` 由 goctl 自动生成**禁止修改**，扩展方法统一放同名自定义 `.go` 文件

> 已知缺口（model 空壳）：`mediamodel.go` 与 `filemodel.go` 两个自定义 model 文件当前均为 goctl 空壳（仅内嵌生成接口，未添加任何扩展方法），可用方法仅为生成的四件套：`Insert(ctx, data)` / `FindOne(ctx, id)`（带缓存）/ `Update(ctx, data)`（全字段更新）/ `Delete(ctx, id)`（物理删除）。待补齐的扩展方法：

| 需求来源 | 需要的方法 | 说明 |
|---------|-----------|------|
| `MediaList` RPC | `FindPage(ctx, offset, limit, name, sortBy, isAsc)` | 分页 + 模糊搜索 + 动态排序 |
| `MediaDelete` / `FileDelete` RPC | `SoftDelete(ctx, id, updater)` | 生成的 `Delete` 是物理删除，与 `deleted` 列语义冲突 |
| 全部查询 | `deleted = 0` 过滤 | 生成的 `FindOne` 不过滤逻辑删除 |
| `MediaSave` 幂等 | `FindOneByFileId(ctx, fileId)` | 按云端 `file_id` 反查，避免重复落库 |
| `FileSave` 幂等 | `FindOneByKey(ctx, key)` | 按云端 `key` 反查 |
| 文件状态流转 | `UpdateStatus(ctx, id, status)` | 1→2→3 单字段更新，避免全字段 `Update` 覆盖 |

#### Scenario: FileVO.path 应用层组装

- **WHEN** 查询文件信息构造 `FileVO`
- **THEN** `path` 不读取数据库列（`file` 表无 `path` 列），由 `key` 拼接对象存储访问前缀（CDN/对象存储域名）得到

#### Scenario: goctl 生成文件禁止修改与扩展位置

- **WHEN** 需要为 `MediaModel` / `FileModel` 新增数据访问方法（如 `FindPage` / `SoftDelete`）
- **THEN** 扩展方法必须写入同名自定义 `.go` 文件（`mediamodel.go` / `filemodel.go`），`*_gen.go` 保持 goctl 生成原样
- **AND** 当前两个自定义文件均为空壳，上述待补齐方法尚不存在

### Requirement: 服务配置与端口

media 服务 SHALL 按以下默认配置部署（API：`apps/media/api/etc/media-api.yaml`；RPC：`apps/media/rpc/etc/media.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| API `Name` / `Host` / `Port` | `media-api` / `0.0.0.0` / `8806` | HTTP 监听 |
| API `Auth.AccessSecret` | `change-me-in-production` | JWT 签名密钥 |
| API `Auth.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| API `MediaRpc.Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `media.rpc` | 自身 RPC 服务发现 |
| RPC `Name` / `ListenOn` | `media.rpc` / `0.0.0.0:8086` | gRPC 监听 |
| RPC `Etcd.Hosts[0]` / `Key` | `127.0.0.1:2379` / `media.rpc` | 服务注册 |
| RPC `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_media?charset=utf8mb4&parseTime=true&loc=Local` | MySQL `tj_media` 库（utf8mb4，时区 Local） |
| RPC `Cache[0]` | `127.0.0.1:6379`，Type `node`，Pass 空 | Redis 单机缓存（`model.NewMediaModel` / `model.NewFileModel` 主键缓存） |

配置结构体：API 层 `Config { rest.RestConf; Auth { AccessSecret, AccessExpire }; MediaRpc zrpc.RpcClientConf }`；RPC 层 `Config { zrpc.RpcServerConf; DataSource string; Cache cache.CacheConf }`。

约束：

- `parseTime=true` 为必需项：`file` / `media` 两表的 `create_time` / `update_time` 均映射为 Go `time.Time`
- JWT：生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造；`media-api` 的 `Auth.AccessSecret` 必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则所有接口鉴权失败
- 认证冲突加注：`media.api` 三组路由（`media` / `signature` / `file`）均声明了 `jwt: Auth`，全部接口需要携带有效 JWT——与聚合文档 `tjxt.openapi.json` 中「认证：否」的标注冲突，以 `.api` 声明为准
- 依赖的外部服务：MySQL `tj_media` 库（DataSource 配置，自建存储）、Redis（Cache 配置，缓存 media/file 主键查询）、`media.rpc`（API 层经 etcd 发现自身 RPC）
- 端口分配：`media-api` 8806（HTTP）、`media.rpc` 8086（gRPC）
- 可观测性注入配置（Telemetry / Prometheus / Log 三段）已注入勿删（configs.md 未记录、以实际 yaml 为准；统一约定见 `specs/infra/observability/spec.md`）

> 已知缺口（对象存储配置，主要阻塞项）：media 服务的核心能力是媒资上传与签名授权，但当前配置中**完全没有对象存储相关配置项**——`media.yaml` 仅 `Name` / `ListenOn` / `Etcd` / `DataSource` / `Cache`，`config.go` 仅 `zrpc.RpcServerConf` / `DataSource` / `Cache` 三个字段；`media-api.yaml` 仅 `Name` / `Host` / `Port` / `Auth` / `MediaRpc`。业务侧实际需要：双平台凭据（腾讯云与阿里云两套 SecretId / SecretKey，对应 `file.platform` 注释）、上传签名参数（Bucket、Region、上传路径前缀、签名有效期）、点播播放参数（点播 AppId、播放域名、防盗链 Key、签名有效期）、VOD 子应用 ID（对应 `media.file_id`）、访问地址前缀（CDN/对象存储域名 + `file.key` 拼接 `FileVO.path`）。建议补齐（当前不存在，仅为缺口说明）：`Storage: { Platform, SecretId, SecretKey, Region, Bucket, UploadPathPrefix, CdnDomain, Vod: { AppId, PlayDomain, PlayKey, SignExpireSec: 3600 } }`。在补齐之前，`SignatureUpload` / `SignaturePreview` / `SignaturePlay` 无法实现真实签名，`FileVO.path` 与 `MediaVO.mediaUrl` / `coverUrl` 缺少地址拼接依据。

#### Scenario: API 层经 etcd 调用自身 RPC

- **WHEN** 客户端请求 `media-api` 的任一 HTTP 接口（8806 端口）
- **THEN** HTTP handler 按 `MediaRpc` 配置（etcd `127.0.0.1:2379`，key `media.rpc`）调用自身 RPC 服务（8086 端口），所有媒资/文件/签名接口最终走自身 RPC

#### Scenario: 生产环境必须替换默认 JWT 密钥

- **WHEN** media 服务以默认 `AccessSecret = change-me-in-production` 上线生产
- **THEN** JWT 可被伪造，生产环境必须替换该默认密钥
- **AND** 替换值必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则所有接口鉴权失败

#### Scenario: 对象存储配置补齐前签名功能受限

- **WHEN** 在未补齐对象存储配置（SecretId / SecretKey / Bucket / Region / Vod 等）的环境下调用签名接口
- **THEN** 签名/上传/播放均走 mock 后端（本地 `http://127.0.0.1:9000`），无法对接真实 COS/OSS，媒资全链路仅为本地可跑的桩

