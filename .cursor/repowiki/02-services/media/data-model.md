> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_media.sql`

---

# Media Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `mediamodel.go` | `media` | 无（goctl 空壳，仅内嵌 `mediaModel`） |
| `filemodel.go` | `file` | 无（goctl 空壳，仅内嵌 `fileModel`） |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `mediamodel.go` 扩展 `MediaModel` 接口）。

> **当前状态**：两个自定义 model 文件均未添加任何扩展方法，可用能力仅为生成的 `Insert` / `FindOne` / `Update` / `Delete` 四件套。

---

## 表清单与字段说明

### 1. `file` — 文件表（普通文件、图片等）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键，文件 id | PK |
| `key` | varchar(255) | 文件在云端的唯一标示，例如：`aaa.jpg` | 业务唯一键 |
| `filename` | varchar(255) | 文件上传时的名称 | - |
| `request_id` | varchar(64) | 请求 id，可空 | - |
| `status` | tinyint | 1-待上传, 2-已上传未使用, 3-已使用 | 状态流转 |
| `platform` | tinyint | 平台：1-腾讯, 2-阿里，默认 1 | - |
| `create_time` | datetime | 创建时间，默认 CURRENT_TIMESTAMP | 自动填充 |
| `update_time` | datetime | 更新时间，ON UPDATE CURRENT_TIMESTAMP | 自动更新 |
| `creater` | bigint | 创建者 ID，默认 0 | - |
| `updater` | bigint | 更新者 ID，默认 0 | - |
| `dep_id` | bigint | 部门 ID，默认 0 | - |
| `deleted` | tinyint | 0=正常, 1=逻辑删除，默认 0 | 过滤条件 |

> **注意**：proto 的 `FileVO` 含 `path` 字段，但表中**无 `path` 列**，需由 `key` 拼接对象存储访问前缀在应用层组装。

---

### 2. `media` — 媒资表（主要是视频文件）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键 | PK |
| `file_id` | varchar(32) | 文件在云端的唯一标示，例如：`387702302659783576` | 业务外部键 |
| `filename` | varchar(255) | 文件名称，默认 `''` | 列表模糊搜索 |
| `media_url` | varchar(255) | 媒体播放地址，默认 `''` | - |
| `cover_url` | varchar(255) | 媒体封面地址，默认 `''` | - |
| `duration` | double | 视频时长，单位秒，默认 0 | - |
| `size` | bigint | 视频大小，单位字节，默认 0 | - |
| `request_id` | varchar(32) | 请求 id，可空，默认 `''` | - |
| `status` | tinyint | 1-上传中, 2-已上传，默认 1 | 状态流转 |
| `create_time` | datetime | 创建时间，默认 CURRENT_TIMESTAMP | 自动填充 |
| `update_time` | datetime | 更新时间，ON UPDATE CURRENT_TIMESTAMP | 自动更新 |
| `creater` | bigint | 创建者 ID，默认 0 | - |
| `updater` | bigint | 更新者 ID，默认 0 | - |
| `dep_id` | bigint | 部门 ID，默认 0 | - |
| `deleted` | tinyint | 0=正常, 1=逻辑删除，默认 0 | 过滤条件 |

> **注意**：proto 的 `MediaVO` 含 `useTimes`（被引用次数），表中**无对应列**，需由 course 域 `CourseMediaUseInfo` 回传 `quote_num` 在应用层聚合。

---

## Go 结构体映射

goctl 生成的结构体与表字段一一对应（`apps/media/rpc/internal/model/*_gen.go`）：

| Go 字段 | Go 类型 | db tag | 备注 |
|---------|---------|--------|------|
| `File.Id` | int64 | `id` | - |
| `File.Key` | string | `key` | - |
| `File.Filename` | string | `filename` | - |
| `File.RequestId` | sql.NullString | `request_id` | 可空列 |
| `File.Status` | int64 | `status` | tinyint → int64 |
| `File.Platform` | int64 | `platform` | tinyint → int64 |
| `File.CreateTime` / `File.UpdateTime` | time.Time | `create_time` / `update_time` | - |
| `File.Creater` / `File.Updater` / `File.DepId` / `File.Deleted` | int64 | 同名 | - |

`Media` 结构体同理，`duration` 映射为 `float64`，`request_id` 映射为可空类型。

---

## 关系图

```
file (独立表)
  key ──→ 对象存储对象名（腾讯/阿里，由 platform 区分）

media (独立表)
  file_id ──→ 云端点播文件标示（非本库外键）

跨域引用（无数据库外键，靠应用层维护）:
  course.section.media_id ──→ media.id     (sql/ddl/tj_course.sql:91)
  course.catalogue.media_id ──→ media.id   (sql/ddl/tj_course.sql:123)
  media.useTimes ←── course.CourseMediaUseInfo(MediaIdsRequest) → MediaQuoteList
```

> `file` 与 `media` 之间**无数据库外键关联**：`file` 记录普通文件/图片，`media` 记录视频媒资，两者走不同的上传通道（`FileSave` vs `MediaSave`）。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
mediamodel_gen.go         ← goctl 生成，只读
mediamodel.go             ← 手写扩展位置（当前为空壳）
```

当前项目自定义 Model 模式：
- `mediamodel.go` — 无扩展方法
- `filemodel.go` — 无扩展方法

**缓存 key 前缀**（由 goctl 生成）：

| Model | 缓存前缀 |
|-------|---------|
| `file` | `cache:file:id:` |
| `media` | `cache:media:id:` |

### 待补齐的扩展方法（缺口）

| 需求来源 | 需要的方法 | 说明 |
|---------|-----------|------|
| `MediaList` RPC | `FindPage(ctx, offset, limit, name, sortBy, isAsc)` | 分页 + 模糊搜索 + 动态排序 |
| `MediaDelete` / `FileDelete` RPC | `SoftDelete(ctx, id, updater)` | 生成的 `Delete` 是物理删除，与 `deleted` 列语义冲突 |
| 全部查询 | `deleted = 0` 过滤 | 生成的 `FindOne` 不过滤逻辑删除 |
| `MediaSave` 幂等 | `FindOneByFileId(ctx, fileId)` | 按云端 `file_id` 反查，避免重复落库 |
| `FileSave` 幂等 | `FindOneByKey(ctx, key)` | 按云端 `key` 反查 |
| 文件状态流转 | `UpdateStatus(ctx, id, status)` | 1→2→3 单字段更新，避免全字段 `Update` 覆盖 |
