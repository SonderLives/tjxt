> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/media/rpc/media.proto`

---

# Media RPC Spec

## 服务名

`Media` — 媒资与文件管理微服务，通过 etcd 服务发现（key: `media.rpc`）。

## RPC 方法总览

### 媒资管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `MediaGet` | `MediaIdRequest { mediaId }` | `MediaVO { id, filename, mediaUrl, coverUrl, duration, size, status, creater, createTime, useTimes }` | 按媒资 ID 查询单条媒资 |
| `MediaList` | `MediaListRequest { pageNo, pageSize, name, sortBy, isAsc }` | `MediaListReply { total, list }` | 分页搜索媒资列表 |
| `MediaSave` | `MediaSaveRequest { id, filename, duration, size, fileId }` | `MediaIdReply { id }` | 视频上传完成后保存媒资信息 |
| `MediaDelete` | `MediaIdRequest { mediaId }` | `Empty {}` | 删除媒资视频 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `mediaId` | int64 | 媒资 ID |
| `id` | int64 | 媒资 ID，新增时省略 |
| `filename` | string | 文件名称 |
| `duration` | double | 视频时长，单位秒 |
| `size` | int64 | 视频大小，单位字节 |
| `fileId` | string | 文件在云端的唯一标示（对应 `media.file_id`） |
| `name` | string | 列表按文件名模糊搜索 |
| `sortBy` | string | 排序字段 |
| `isAsc` | string | 是否升序 |

**响应字段说明（`MediaVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `mediaUrl` | string | 媒体播放地址 |
| `coverUrl` | string | 媒体封面地址 |
| `status` | int32 | 1-上传中，2-已上传 |
| `creater` | string | 创建者 |
| `useTimes` | int32 | 被引用次数（由课程域回填） |

---

### 签名管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SignatureUpload` | `SignatureRequest { mediaId, fileName, mediaType }` | `SignatureVO { token, url, uploadUrl, playUrl }` | 获取上传视频的授权签名 |
| `SignaturePreview` | `SignatureRequest { mediaId, fileName, mediaType }` | `SignatureVO { token, url, uploadUrl, playUrl }` | 管理端获取预览视频的授权签名 |
| `SignaturePlay` | `SignatureRequest { mediaId, fileName, mediaType }` | `SignatureVO { token, url, uploadUrl, playUrl }` | 学员端获取播放视频的授权签名 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `mediaId` | int64 | 媒资 ID，上传场景可为 0 |
| `fileName` | string | 文件名 |
| `mediaType` | string | 媒体类型 |

**响应字段说明（`SignatureVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `token` | string | 云端签名令牌 |
| `url` | string | 资源访问地址 |
| `uploadUrl` | string | 上传地址（上传签名场景） |
| `playUrl` | string | 播放地址（播放/预览签名场景） |

> 三个签名方法共用同一组 `SignatureRequest` / `SignatureVO`，由方法语义区分用途。

---

### 文件管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `FileGet` | `FileIdRequest { id }` | `FileVO { id, key, filename, path, status }` | 按 ID 获取文件信息 |
| `FileSave` | `FileSaveRequest { filename, key, size }` | `FileIdReply { id }` | 上传文件后保存文件记录 |
| `FileDelete` | `FileIdRequest { id }` | `Empty {}` | 删除文件 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 文件 ID |
| `filename` | string | 文件上传时的名称 |
| `key` | string | 文件在云端的唯一标示（如 `aaa.jpg`） |
| `size` | int64 | 文件大小 |

**响应字段说明（`FileVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `path` | string | 文件访问路径（DDL 中无对应列，由 `key` 拼接对象存储前缀得到） |
| `status` | int32 | 1-待上传，2-已上传未使用，3-已使用 |

> proto 中另声明了 `OkVO {}` 空消息，当前 `service Media` 未引用，删除类方法统一返回 `Empty`。

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `media-api` (自身 API 层) | HTTP Handler → `mediaclient.Media` RPC | `apps/media/api/internal/svc/servicecontext.go` 中 import `mediaclient "tjxt/apps/media/rpc/media"`，所有媒资/文件/签名接口最终走自身 RPC |

**尚未接线的逻辑消费方**：

| 潜在消费方 | 耦合证据 | 当前状态 |
|-----------|---------|---------|
| `course` 服务 | `apps/course/rpc/course.proto:52` 定义 `CourseMediaSave(CourseMediaSaveRequest)`，`CourseMediaBind` 含 `media_id` / `video_name` / `media_duration`；`course.proto:36` 定义 `CourseMediaUseInfo(MediaIdsRequest) returns (MediaQuoteList)` 回传 `media_id` → `quote_num`；`sql/ddl/tj_course.sql:91,123` 两处 `media_id` 列 | `apps/course/api/internal/svc/servicecontext.go` **未** import `mediaclient`，`apps/course/api/etc/*.yaml` 亦无 `MediaRpc` 配置项，跨服务调用尚未接线 |

> 说明：course 与 media 之间是**双向逻辑依赖** —— course 侧持有 `media_id` 引用媒资，media 侧 `MediaVO.useTimes` 需要 course 侧的 `CourseMediaUseInfo` 回填引用计数。两侧目前均以本地字段/占位方式表达，未建立实际 RPC 连接。

---

## 调用典型场景

1. **上传视频** → 前端调 `SignatureUpload` 获取上传签名 → 直传对象存储 → 回调 `MediaSave` 落库媒资记录（`status=2` 已上传）
2. **上传普通文件/图片** → 前端调 `FileSave` 落库文件记录 → `status` 由 1（待上传）流转到 2（已上传未使用）
3. **管理端预览** → 管理员在媒资列表点击预览 → 调 `SignaturePreview` 换取带鉴权的 `playUrl`
4. **学员播放** → 学员在小节页播放视频 → 调 `SignaturePlay` 换取限时播放签名
5. **媒资检索** → 管理端调 `MediaList`，按 `name` 模糊搜索、按 `sortBy` / `isAsc` 排序分页
6. **课程绑定媒资** → 课程编辑时选取媒资 → course 侧 `CourseMediaSave` 写入 `media_id`，文件 `status` 应流转为 3（已使用）

---

## 自定义 Model 方法

`apps/media/rpc/internal/model/` 下的两个自定义 model 文件当前**均为 goctl 空壳**，未添加任何扩展方法：

- `mediamodel.go` — `MediaModel interface { mediaModel }`，仅内嵌生成接口
- `filemodel.go` — `FileModel interface { fileModel }`，仅内嵌生成接口

可用方法仅为 `*_gen.go` 生成的四件套：

| 方法 | 签名 | 说明 |
|------|------|------|
| `Insert` | `Insert(ctx, data *Media/*File) (sql.Result, error)` | 插入 |
| `FindOne` | `FindOne(ctx, id int64) (*Media/*File, error)` | 按主键查询（带缓存） |
| `Update` | `Update(ctx, data *Media/*File) error` | 全字段更新 |
| `Delete` | `Delete(ctx, id int64) error` | 物理删除 |

> **缺口**：`MediaList` 分页搜索、按 `file_id` / `key` 反查、`deleted` 软删除过滤等能力，生成方法均不支持，需在自定义 model 中补写（参考 auth 域 `rolemodel.go` 的 `FindPage` / `SoftDelete` 模式）。
