> 版本：v1.2 | 更新：2026-08-06 | 来源：2026-08-06 复核（go build 全模块通过 + 逻辑文件清点）

---

# Media Business Rules

## ⚠️ 实现状态

本服务的业务 logic **已全部实现**：20 个 logic 文件（RPC 10 + API 10）均已落地并编译通过；对象存储后端为 mock（见已知缺口）。下列各方法状态已校正为「已实现」。

### RPC 层实现状态（`apps/media/rpc/internal/logic/`）

| Logic 文件 | RPC 方法 | 实现状态 |
|-----------|---------|---------|
| `mediagetlogic.go` | `MediaGet` | ✅ 已实现 |
| `medialistlogic.go` | `MediaList` | ✅ 已实现 |
| `mediasavelogic.go` | `MediaSave` | ✅ 已实现 |
| `mediadeletelogic.go` | `MediaDelete` | ✅ 已实现 |
| `signatureuploadlogic.go` | `SignatureUpload` | ✅ 已实现 |
| `signaturepreviewlogic.go` | `SignaturePreview` | ✅ 已实现 |
| `signatureplaylogic.go` | `SignaturePlay` | ✅ 已实现 |
| `filegetlogic.go` | `FileGet` | ✅ 已实现 |
| `filesavelogic.go` | `FileSave` | ✅ 已实现 |
| `filedeletelogic.go` | `FileDelete` | ✅ 已实现 |

**RPC 层统计：已实现 10 / 总计 10**

### API 层实现状态（`apps/media/api/internal/logic/`）

| Logic 文件 | Handler 方法 | 实现状态 |
|-----------|-------------|---------|
| `media/mediagetlogic.go` | `MediaGet` | ✅ 已实现 |
| `media/medialistlogic.go` | `MediaList` | ✅ 已实现 |
| `media/mediasavelogic.go` | `MediaSave` | ✅ 已实现 |
| `media/mediadeletelogic.go` | `MediaDelete` | ✅ 已实现 |
| `signature/signatureuploadlogic.go` | `SignatureUpload` | ✅ 已实现 |
| `signature/signaturepreviewlogic.go` | `SignaturePreview` | ✅ 已实现 |
| `signature/signatureplaylogic.go` | `SignaturePlay` | ✅ 已实现 |
| `file/filegetlogic.go` | `FileGet` | ✅ 已实现 |
| `file/filesavelogic.go` | `FileSave` | ✅ 已实现 |
| `file/filedeletelogic.go` | `FileDelete` | ✅ 已实现 |

**API 层统计：已实现 10 / 总计 10**

> **合计：已实现 20 / 总计 20**
>
> 以下各节内容**均为设计意图推导**，依据 `apps/media/rpc/media.proto` 的消息定义、`sql/ddl/tj_media.sql` 的字段注释与状态枚举、`apps/media/api/media.api` 的路由契约、以及 `docs/tjxt.openapi.json` 的原始 Java 版接口摘要。**2026-08-06 复核：logic 已全部实现并编译通过；以下规则为依据 proto/DDL/.api 契约推导，建议对照源码最终确认。**


## 已知缺口

- 对象存储为 mock：`config.go` / `etc/*.yaml` 无 `SecretId/SecretKey/Bucket` 配置，签名/上传/播放均指向本地 `http://127.0.0.1:9000`（`mockBaseURL`）。接入真实 COS/OSS 前媒资全链路仅为本地可跑的桩。
- course ↔ media RPC 未接线：引用计数（`CourseMediaUseInfo`）与媒资绑定（`CourseMediaSave`）跨域未打通。

---

## 1. 文件状态流转 📋 设计意图（契约推导）

**核心规则**：`file.status` 三态流转，由 DDL 注释 `状态：1-待上传 2-已上传,未使用 3-已使用` 定义。

| 规则 | 说明 | 依据 |
|------|------|------|
| 初始态为待上传 | 申请上传时先落库 `status=1` | `file.status` 注释 |
| 上传完成置为未使用 | 对象存储回调后 `status=2` | `file.status` 注释 |
| 被业务引用置为已使用 | 被课程/用户头像等引用后 `status=3` | `file.status` 注释 |
| 平台区分 | `platform` 1-腾讯 / 2-阿里，默认 1 | `file.platform` 注释 |
| 逻辑删除 | 删除时 `deleted=1`，而非物理删除 | `file.deleted` 注释 |

```
状态流转（设计意图）:
  1 待上传 ──上传成功──→ 2 已上传未使用 ──被引用──→ 3 已使用
                                  │                      │
                                  └──── deleted=1 ───────┘
                                       （已使用者应拒绝删除）
```

> ⚠️ 生成的 `FileModel.Delete` 为**物理删除**，与 `deleted` 列的逻辑删除语义冲突，实现时需改用自定义 `SoftDelete`。

## 2. 媒资状态流转 📋 设计意图（契约推导）

**核心规则**：`media.status` 两态流转，由 DDL 注释 `状态：1-上传中，2-已上传` 定义。

| 规则 | 说明 | 依据 |
|------|------|------|
| 默认上传中 | 建记录时 `status=1`（DDL DEFAULT 1） | `media.status` 默认值 |
| 转码/上传完成置为已上传 | 回调后 `status=2`，回填 `media_url` / `cover_url` | `media.status` 注释 |
| `file_id` 关联云端 | `MediaSaveRequest.fileId` 写入 `media.file_id` | proto + DDL |
| 时长与大小由客户端上报 | `duration`(秒) / `size`(字节) 随 `MediaSave` 传入 | `MediaSaveRequest` |
| 逻辑删除 | `MediaDelete` 应置 `deleted=1` | `media.deleted` 注释 |

## 3. 媒资保存（新增/更新） 📋 设计意图（契约推导）

**核心规则**：`MediaSaveRequest.id` 为 optional，沿用项目内 `id <= 0` 为新增的约定（参考 auth 域 `SaveRole`）。

```
流程（MediaSave，设计意图）:
  1. 校验 filename / fileId 非空
  2. id <= 0 → 新增（生成雪花 ID，status=1，creater 取 JWT userId）
     id >  0 → 更新（校验记录存在且 deleted=0）
  3. 写入 filename / duration / size / file_id
  4. 返回 MediaIdReply{ id }
```

> ⚠️ 幂等缺口：同一 `fileId` 重复调用 `MediaSave` 会产生重复媒资记录，需补 `FindOneByFileId` 校验（见 [data-model.md](./data-model.md) 缺口表）。

## 4. 媒资列表查询 📋 设计意图（契约推导）

**核心规则**：`MediaListRequest` 支持模糊搜索 + 动态排序 + 分页。

| 字段 | 说明 |
|------|------|
| `pageNo` / `pageSize` | 分页参数，`.api` 中均为 optional，需在 logic 内兜底默认值 |
| `name` | 按 `media.filename` 模糊匹配 |
| `sortBy` | 排序字段名 |
| `isAsc` | 是否升序（proto 中为 **string** 类型，非 bool） |

> ⚠️ 安全提示：`sortBy` 直接来自客户端，拼接 `ORDER BY` 时**必须走白名单**，否则构成 SQL 注入。
>
> ⚠️ 能力缺口：`MediaModel` 无 `FindPage`，分页查询无法实现。

## 5. 签名授权 📋 设计意图（契约推导）

**核心规则**：三个签名方法共用 `SignatureRequest` / `SignatureVO`，按用途返回不同字段组合。

| 方法 | 用途 | 预期返回的关键字段 |
|------|------|------------------|
| `SignatureUpload` | 客户端直传对象存储 | `token`, `uploadUrl` |
| `SignaturePreview` | 管理端预览未发布视频 | `token`, `playUrl` |
| `SignaturePlay` | 学员端播放已发布视频 | `token`, `playUrl` |

| 规则 | 说明 | 依据 |
|------|------|------|
| 签名有时效 | 播放/预览签名应限时，避免地址外泄长期可用 | 通用对象存储实践 |
| 上传签名不需要 mediaId | `SignatureRequest.mediaId` 在上传场景可为 0 | proto 字段可选语义 |
| 平台适配 | 需按 `file.platform`（1-腾讯 / 2-阿里）选择不同签名算法 | `file.platform` 注释 |

> ⚠️ **配置缺口（阻塞项）**：`apps/media/rpc/internal/config/config.go` 仅有 `RpcServerConf` / `DataSource` / `Cache` 三项，**没有任何对象存储的 AppId / SecretId / SecretKey / Bucket / Region / 上传路径配置**，`apps/media/rpc/etc/media.yaml` 中亦无对应条目。签名功能在补齐配置结构前**无法实现**。详见 [configs.md](./configs.md)。

## 6. 文件保存与删除 📋 设计意图（契约推导）

```
流程（FileSave，设计意图）:
  1. 校验 filename / key 非空
  2. 生成 ID，落库 status=1（待上传）或 2（已上传未使用）
  3. platform 取默认 1（腾讯）
  4. 返回 FileIdReply{ id }

流程（FileDelete，设计意图）:
  1. FindOne 校验存在且 deleted=0
  2. status=3（已使用）时应拒绝删除
  3. 置 deleted=1，并清理对象存储中的 key
```

## 7. 跨域引用计数 📋 设计意图（契约推导）

**核心规则**：`MediaVO.useTimes` 表示媒资被课程引用的次数，**不落 media 表**。

| 规则 | 说明 | 依据 |
|------|------|------|
| 引用计数来自 course 域 | 调 `Course.CourseMediaUseInfo(MediaIdsRequest)` 取 `MediaQuoteList{ media_id, quote_num }` | `apps/course/rpc/course.proto:36,350-355` |
| 被引用的媒资不可删除 | `quote_num > 0` 时 `MediaDelete` 应拒绝 | 与 `file.status=3` 语义一致 |
| 绑定动作在 course 侧 | `Course.CourseMediaSave` 写入 `media_id` | `apps/course/rpc/course.proto:52,436-447` |

> ⚠️ **接线缺口**：`apps/course/api/internal/svc/servicecontext.go` 未 import `mediaclient`，`apps/media/*/internal/svc/servicecontext.go` 也未持有 `CourseRpc`，两域间的 RPC 通道尚未建立。

---

## 状态说明

### `file.status` 文件状态

| 值 | 含义 | 可否删除（设计意图） |
|----|------|-------------------|
| 1 | 待上传 | 可删除 |
| 2 | 已上传，未使用 | 可删除 |
| 3 | 已使用 | 应拒绝删除 |

### `media.status` 媒资状态

| 值 | 含义 |
|----|------|
| 1 | 上传中（DDL 默认值） |
| 2 | 已上传 |

### `file.platform` 存储平台

| 值 | 含义 |
|----|------|
| 1 | 腾讯（DDL 默认值） |
| 2 | 阿里 |

### `deleted` 逻辑删除

| 值 | 含义 |
|----|------|
| 0 | 正常（DDL 默认值） |
| 1 | 已逻辑删除 |
