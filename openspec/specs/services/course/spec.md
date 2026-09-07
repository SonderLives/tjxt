# course（课程管理）Specification

## Purpose

course（课程管理）是课程管理微服务，覆盖分类（category）、课程（course）、目录（catalogue）、老师（teacher）、题目（subject）五大领域，提供课程分类维护、课程信息草稿式编辑（基本信息/目录/媒资/题目/老师五步）、上下架发布与多端课程查询。HTTP API 层 `course-api` 监听 `0.0.0.0:8803`，RPC 服务 `Course` 监听 `0.0.0.0:8083`，经 etcd 以 key `course.rpc` 注册与发现，API 层全部接口经 etcd 调用自身 RPC。数据存储于 MySQL `tj_course` 库（`category` 与 `course`/`course_draft` 等正式表-草稿表成对结构），Redis 单机缓存托管主键查询。除 `course-api` 自身外，`learning-api` 亦消费该 RPC（查询课程目录、章节、小节）；课程上架/下架事件经 RabbitMQ `course.events` 交换机发布，供 search 服务同步 ES 课程索引。

## Requirements

### Requirement: HTTP 课程分类管理接口

> **状态说明**：接口清单源自聚合文档 `docs/tjxt.openapi.json`（api-spec.md，最后同步 2026-08-05），「认证」列标注为「否」；go-zero 侧 `apps/course/api/course.api` 的四个路由组（`/categorys`、`/courses`、`/course`、`/catalogues`）均声明 `jwt: Auth`，以 `.api` 为准。

course 服务 SHALL 暴露以下课程分类管理接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /categorys/add | 新增课程分类 | CategoryAddDTO | R{data: R} |
| GET | /categorys/all | 获取所有的课程分类信息，只包含id,名称，课程分类关系 | - | R{data: R«List«SimpleCategoryVO»»} |
| PUT | /categorys/disableOrEnable | 课程分类停用或启用 | CategoryDisableOrEnableDTO | R{data: R} |
| GET | /categorys/getAllOfOneLevel | 获取所有的课程分类，不分层 | - | R{data: R«List«CategoryVO»»} |
| GET | /categorys/list | 查询课程分类信息 | - | R{data: R«List«CategoryVO»»} |
| PUT | /categorys/update | 更新课程分类 | CategoryUpdateDTO | R{data: R} |
| DELETE | /categorys/{id} | 删除分类信息 | - | R{data: R} |
| GET | /categorys/{id} | 获取课程分类信息 | - | R{data: R«CategoryInfoVO»} |

统一约定（引用全局规范）：响应格式为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页为 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码见共享错误码表（`specs/shared/error-codes`）。以下各 HTTP 接口 Requirement 均沿用该约定，不再重复。

#### Scenario: 新增课程分类

- **WHEN** 客户端以 `CategoryAddDTO` 调用 `POST /categorys/add`
- **THEN** 系统创建分类并按 `parent_id` 自动推导层级（规则见「分类管理业务规则」）

#### Scenario: 停用或启用分类

- **WHEN** 客户端以 `CategoryDisableOrEnableDTO` 调用 `PUT /categorys/disableOrEnable`
- **THEN** 系统切换分类状态（1=正常/2=禁用）并返回 `R{data: R}`，禁用不影响已有课程

### Requirement: HTTP 课程基础信息与搜索接口

course 服务 SHALL 暴露以下课程基础信息与搜索接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /courses/baseInfo/save | 保存课程基本信息 | CourseBaseInfoSaveDTO | R{data: R«CourseSaveVO»} |
| GET | /courses/baseInfo/{id} | 获取课程基础信息 | - | R{data: R«CourseBaseInfoVO»} |
| GET | /courses/checkName | 校验课程名称是否已经存在 | - | R{data: R«NameExistVO»} |
| GET | /courses/page | 管理端课程搜索接口 | - | R{data: R«PageDTO«CoursePageVO»»} |
| GET | /courses/portal | 用户端课程搜索接口 | - | R{data: R«PageDTO«CourseVO»»} |
| GET | /courses/simpleInfo/list | 根根条件列表获取课程信息 | - | R{data: R«List«CourseSimpleInfoDTO»»} |
| GET | /courses/generator | 生成练习id | - | R{data: R«CourseCataIdVO»} |
| DELETE | /courses/delete/{id} | 课程删除 | - | R{data: R} |

- `GET /courses/generator` 摘要原文为「生成练习id」，响应 `CourseCataIdVO`；RPC 侧对应方法 `CourseGenerator` 的说明为生成课程草稿 ID（见「课程信息保存流程」）。
- `GET /courses/simpleInfo/list` 摘要「根根条件列表获取课程信息」为源文档原文。
- 课程删除为逻辑删除（RPC `CourseDelete`）。

#### Scenario: 管理端课程搜索

- **WHEN** 管理端调用 `GET /courses/page`
- **THEN** 系统按分页与过滤条件（keyword/status/free/course_type/三级分类/时间范围）返回 `R{data: R«PageDTO«CoursePageVO»»}`

#### Scenario: 课程名称查重

- **WHEN** 客户端调用 `GET /courses/checkName`
- **THEN** 系统返回 `R{data: R«NameExistVO»}`；编辑时排除自身（草稿表 `FindByNameExceptId`），新增（id=0）时全局查重

#### Scenario: 课程删除

- **WHEN** 客户端调用 `DELETE /courses/delete/{id}`
- **THEN** 系统对课程执行逻辑删除并返回 `R{data: R}`

### Requirement: HTTP 课程目录与章节接口

course 服务 SHALL 暴露以下课程目录与章节接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /catalogues/batchQuery | 根据章节目录批量查询基础信息 | - | R{data: R«List«CataSimpleInfoVO»»} |
| GET | /catalogues/querySectionInfoById/{id} | 获取小节信息 | - | R{data: R«CataSimpleInfoVO»} |
| POST | /courses/catas/save/{id}/{step} | 保存章节 | JSON | R{data: R} |
| GET | /courses/catas/{id} | 获取课程的章节 | - | R{data: R«List«课程目录»»} |
| GET | /courses/catas/index/list/{id} | 根据课程id，查询所有章节的序号 | - | R{data: R«List«CataSimpleInfoVO»»} |
| GET | /courses/{id}/catalogs | 查询课程基本信息、目录、学习进度 | - | R{data: R«CourseAndSectionVO»} |

#### Scenario: 保存章节（按 step 传累计进度）

- **WHEN** 客户端以 JSON 请求体调用 `POST /courses/catas/save/{id}/{step}`
- **THEN** 系统按章(1)/节(2)/测试(3) 树递归保存到目录草稿表，路径中的 `step` 作为累计进度入参（见「课程信息保存流程」的已知缺口）

#### Scenario: 批量查询目录基础信息

- **WHEN** 客户端调用 `GET /catalogues/batchQuery`
- **THEN** 系统按章节目录批量返回基础信息 `R{data: R«List«CataSimpleInfoVO»»}`

### Requirement: HTTP 课程上下架接口

course 服务 SHALL 暴露以下课程上下架接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /courses/upShelf | 课程上架 | CourseIdDTO | R{data: R} |
| POST | /courses/downShelf | 课程下架 | CourseIdDTO | R{data: R} |
| GET | /courses/checkBeforeUpShelf/{id} | 课程上架前校验 | - | R{data: R} |
| POST | /courses/up | 处理指定课程上架失败的问题 | - | R{data: R} |
| POST | /courses/down | 处理指定课程下架失败的问题 | - | R{data: R} |

- `POST /courses/up` / `POST /courses/down` 摘要为源文档原文（处理上架/下架失败问题）；go-zero `.api` 侧入参为 `CourseIdListReq`（批量 ID），RPC 侧为 `CourseUp` / `CourseDown`（`IdsRequest`，逗号分隔 courseIds）。

#### Scenario: 上架前校验

- **WHEN** 客户端调用 `GET /courses/checkBeforeUpShelf/{id}`
- **THEN** 系统（RPC `CourseCheckUpShelf`）校验草稿 `Step >= 5`，校验结果以错误传递：通过时返回 `nil`（API 侧 `existed=true`），不通过时返回 `BadRequest`（`existed=false`）

#### Scenario: 课程上架

- **WHEN** 客户端以 `CourseIdDTO` 调用 `POST /courses/upShelf`
- **THEN** 系统将草稿及其子表复制到正式表并置状态为已上架(2)（流程见「课程上下架状态机」）

### Requirement: HTTP 课程媒资、题目与老师编辑接口

course 服务 SHALL 暴露以下课程媒资、题目与老师编辑接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /courses/media/save/{id} | 课程视频 | JSON | R{data: R} |
| GET | /courses/subjects/get/{id} | 获取小节或练习中的题目（用于编辑） | - | R{data: R«List«CataSimpleSubjectVO»»} |
| POST | /courses/subjects/save/{id} | 保存小节或练习中的题目 | JSON | R{data: R} |
| POST | /courses/teachers/save | 保存老师信息 | 课程老师关系模型 | R{data: R} |
| GET | /courses/teachers/{id} | 查询课程相关的老师信息 | - | R{data: R«List«老师课程信息»»} |

#### Scenario: 保存课程视频（媒资绑定）

- **WHEN** 客户端以 JSON 请求体调用 `POST /courses/media/save/{id}`
- **THEN** 系统将媒资回填到课程目录草稿小节（按 `cata_id` 定位，规则见「课程媒体绑定」）

#### Scenario: 保存小节或练习中的题目

- **WHEN** 客户端以 JSON 请求体调用 `POST /courses/subjects/save/{id}`
- **THEN** 系统保存小节/练习与题目的绑定关系并返回 `R{data: R}`

### Requirement: HTTP 跨服务内部接口

`/course/*` 路由组 SHALL 提供以下跨服务内部查询接口（api-spec.md 权限标签标注 `internal`）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| GET | /course/{id} | 获取课程信息 | - | R{data: R«CourseFullInfoDTO»} |
| GET | /course/{id}/searchInfo | 课程上架时，需要查询课程信息，加入索引库 | - | R{data: R«课程信息»} |
| GET | /course/section/{id} | sectionInfo | - | R{data: R«小节信息，包含课程id和媒资id»} |
| GET | /course/infoByTeacherIds | 通过老师id获取老师负责的课程和出的题目数量 | - | R{data: R«List«老师id和老师对应的课程数，出题数»»} |
| GET | /course/media/useInfo | mediaUserInfo | - | R{data: R«List«媒资被引用情况»»} |
| GET | /course/name | queryCoursesIdByName | - | R{data: R«List«long»»} |

> 权限标签说明：`internal` 为聚合文档标注；go-zero `.api` 侧该组同属 `/course` 前缀并声明 `jwt: Auth`。

#### Scenario: 搜索服务取课程索引数据

- **WHEN** 搜索服务（或上架流程）调用 `GET /course/{id}/searchInfo`
- **THEN** 系统返回用于加入搜索索引库的课程信息（RPC `CourseSearchInfoForIndex`）

#### Scenario: 媒资被引用情况查询

- **WHEN** 媒资服务调用 `GET /course/media/useInfo`
- **THEN** 系统返回各媒资被课程目录引用的情况列表（规则见「课程媒体绑定」）

### Requirement: RPC 服务契约（Course 服务）

`Course` 服务 SHALL 作为课程管理微服务经 etcd 服务发现注册（key: `course.rpc`），监听 `0.0.0.0:8083`（gRPC），覆盖分类（category）、课程（course）、目录（catalogue）、老师（teacher）、题目（subject）五大领域。RPC 层共 76 个 logic，已全部实现并通过编译（business-rules.md v2.0，2026-08-05 复核）。

消费方 SHALL 为：

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `course-api`（自身） | HTTP Handler → `courseclient.Course` RPC | 所有课程管理/学员端接口 |
| `learning-api` | `courseclient.Course` RPC | 学习服务查询课程目录、章节、小节 |

RPC 方法共 37 个，按域分组为：分类管理 8 个、课程基础信息 13 个、课程目录 6 个、章节/题目 2 个、课程媒体 1 个、课程老师 2 个、课程上下架 5 个（见以下各 Requirement）。通用消息：`IdRequest{ id }` 单 ID 请求、`IdsRequest{ ids }` 批量 ID 请求、`Empty{}` 写操作响应、`IdResponse{ id }` 生成/保存返回 ID。

#### Scenario: 课程发布全流程

- **WHEN** 管理端发布一门课程
- **THEN** 按顺序执行：保存基本信息 → 保存目录（章/节）→ 绑定媒资 → 绑定题目 → 绑定老师 → 检查上架（`CourseCheckUpShelf`）→ 上架（`CourseUpShelf`）

#### Scenario: learning-api 消费课程 RPC

- **WHEN** learning-api 需要展示学习内容的课程目录、章节、小节
- **THEN** 经 etcd 发现 `course.rpc` 并通过 `courseclient.Course` 调用对应查询方法

#### Scenario: 批量上下架

- **WHEN** 管理后台批量操作课程
- **THEN** 调用 `CourseUp` / `CourseDown`（入参为逗号分隔的 `courseIds` 字符串）完成批量上/下架

### Requirement: RPC 分类管理方法

`Course` 服务 SHALL 提供以下分类管理 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CategoryAdd` | `CategoryAddRequest { name, parent_id, index }` | `Empty {}` | 新增分类 |
| `CategoryUpdate` | `CategoryUpdateRequest { id, name, index }` | `Empty {}` | 更新分类 |
| `CategoryDelete` | `IdRequest { id }` | `Empty {}` | 删除分类 |
| `CategoryDisableOrEnable` | `CategoryStatusRequest { id, status }` | `Empty {}` | 启用/停用 (status: 1=启用, 2=停用) |
| `CategoryGet` | `IdRequest { id }` | `CategoryInfo` | 查询分类详情 |
| `CategoryListAll` | `CategoryListAllRequest { admin, name }` | `CategoryNodeList` | 递归树形查询分类 |
| `CategoryListOneLevel` | `Empty {}` | `CategoryList` | 查一级分类列表 |
| `CategoryListQuery` | `CategoryListQueryRequest { name, status }` | `CategoryList` | 按条件查询分类 |

#### Scenario: 递归树形查询分类

- **WHEN** 调用 `CategoryListAll`（入参 `admin`、`name`）
- **THEN** 系统返回递归树结构 `CategoryNodeList`
- **AND** `CategoryListOneLevel` 仅返回一级分类列表 `CategoryList`

#### Scenario: 分类启用/停用

- **WHEN** 调用 `CategoryDisableOrEnable(id, status)`
- **THEN** 系统按 `status` 切换分类状态（1=启用, 2=停用）

### Requirement: RPC 课程基础信息方法

`Course` 服务 SHALL 提供以下课程基础信息 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseBaseInfoGet` | `IdRequest { id }` | `CourseBaseInfoView` | 课程详情（含统计） |
| `CourseBaseInfoSave` | `CourseBaseInfoSaveRequest { id, name, cover_url, price, free, third_cate_id, ... }` | `IdResponse { id }` | 保存课程基础信息 |
| `CourseCheckName` | `CourseCheckNameRequest { name, id }` | `NameExistReply { existed }` | 校验课程名称是否重复 |
| `CourseDelete` | `IdRequest { id }` | `Empty {}` | 删除课程（逻辑删除） |
| `CoursePageQuery` | `CoursePageQueryRequest { pageNo, pageSize, keyword, status, free, ... }` | `CoursePageQueryReply { total, pages, list }` | 后台分页查询课程 |
| `CoursePortalQuery` | `CoursePortalQueryRequest { pageNo, pageSize, category_id_lv1/2/3, keyword, ... }` | `CoursePageQueryReply` | 门户端分页查询课程 |
| `CourseSimpleInfoList` | `CourseSimpleInfoQueryRequest { ids, third_cate_ids }` | `CourseSimpleInfoListReply` | 批量获取课程简略信息 |
| `CourseFullInfoGet` | `CourseFullInfoGetRequest { id, with_catalogue, with_teachers }` | `CourseFullInfo` | 课程完整信息（含目录/老师） |
| `CourseSearchInfoForIndex` | `IdRequest { id }` | `CourseSearchIndexInfo` | 用于搜索服务的索引数据 |
| `CourseName2Ids` | `CourseNameRequest { name }` | `CourseIdList { ids }` | 按名称查课程 ID |
| `CourseInfoByTeacherIds` | `TeacherIdsRequest { teacher_ids }` | `TeacherCourseCountList` | 按老师 ID 查课程统计 |
| `CourseSectionGet` | `IdRequest { id }` | `CourseSectionInfo { course_id, media_id }` | 查询课程小节信息 |
| `CourseMediaUseInfo` | `MediaIdsRequest { media_ids }` | `MediaQuoteList` | 查媒资被引用的次数 |
| `CourseGenerator` | `Empty {}` | `IdResponse { id }` | 生成课程草稿 ID |

关键字段：`CourseBaseInfoSaveRequest` 含 `id`、`name`、`cover_url`、`price`、`free`、`third_cate_id` 等字段（proto 完整定义见 `apps/course/rpc/course.proto`）。

> 已知缺口：分页/详情中的 `sold` / `enroll_num` / `score` 字段在 course 库无列，填 0，真实值由 trade / learning 服务提供，course 暂未接线对应 RPC（详见「课程查询」）。

#### Scenario: 课程名称查重（RPC）

- **WHEN** 调用 `CourseCheckName(name, id)` 且 `id ≠ 0`
- **THEN** 系统在草稿表按 `FindByNameExceptId` 排除自身后查重，返回 `NameExistReply{ existed }`
- **AND** `id=0`（新增）时全局查重

#### Scenario: 批量获取课程简略信息

- **WHEN** 其他服务以 `ids` 或 `third_cate_ids` 调用 `CourseSimpleInfoList`
- **THEN** 系统返回 `CourseSimpleInfoListReply` 供跨服务回填课程信息

### Requirement: RPC 课程目录、媒资、题目与老师方法

`Course` 服务 SHALL 提供以下课程目录、媒资、题目与老师 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseCatalogsGet` | `IdRequest { id }` | `CourseAndSectionView` | 查询课程目录+小节 |
| `CourseCatalogueTreeGet` | `CourseCatalogueQueryRequest { id, see, with_practice }` | `CatalogueTreeList` | 递归树形查询目录结构 |
| `CourseCatalogueSave` | `CourseCatalogueSaveRequest { course_id, step, chapters }` | `Empty {}` | 保存课程目录（含章/节/测试） |
| `CourseCatalogueIndexList` | `IdRequest { id }` | `CataSimpleList` | 课程目录 ID 列表 |
| `CourseCatalogueSectionInfo` | `IdRequest { id }` | `CourseSectionInfo` | 目录小节信息 |
| `CourseCatalogueBatchQuery` | `IdsRequest { ids }` | `CataSimpleList` | 批量查询目录 |
| `CourseSubjectsGet` | `IdRequest { id }` | `CataSubjectInfoList` | 查询章节下的题目 |
| `CourseSubjectsSave` | `CourseSubjectsSaveRequest { course_id, subjects }` | `Empty {}` | 绑定章节-题目关系 |
| `CourseMediaSave` | `CourseMediaSaveRequest { course_id, medias }` | `Empty {}` | 绑定课程目录-媒资 |
| `CourseTeachersGet` | `CourseTeachersGetRequest { id, see }` | `TeacherInfoList` | 查询课程老师 |
| `CourseTeachersSave` | `CourseTeachersSaveRequest { id, teachers }` | `Empty {}` | 保存课程老师关系 |

#### Scenario: 保存课程目录树

- **WHEN** 调用 `CourseCatalogueSave(course_id, step, chapters)`
- **THEN** 系统将章(1)/节(2)/测试(3) 树递归写入目录草稿表（规则见「课程目录管理」）

#### Scenario: 课程目录树双视图

- **WHEN** 分别调用 `CourseCatalogsGet` 与 `CourseCatalogueTreeGet(id, see, with_practice)`
- **THEN** 前者返回课程目录+小节结构 `CourseAndSectionView`，后者递归返回目录树 `CatalogueTreeList`（业务规则标注二者对应正式/草稿视图）

#### Scenario: 绑定媒资到小节

- **WHEN** 调用 `CourseMediaSave(course_id, medias)`
- **THEN** 系统将媒资绑定到课程目录草稿小节（按 `cata_id` 定位）

#### Scenario: 绑定章节题目

- **WHEN** 调用 `CourseSubjectsSave(course_id, subjects)`
- **THEN** 系统保存章节-题目绑定关系
- **AND** `CourseSubjectsGet(id)` 可查询章节下的题目（`CataSubjectInfoList`）

### Requirement: RPC 课程上下架方法

`Course` 服务 SHALL 提供以下课程上下架 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseUpShelf` | `IdRequest { id }` | `Empty {}` | 单个上架 |
| `CourseDownShelf` | `IdRequest { id }` | `Empty {}` | 单个下架 |
| `CourseCheckUpShelf` | `IdRequest { id }` | `Empty {}` | 上架前置校验 |
| `CourseUp` | `IdsRequest { ids }` | `Empty {}` | 批量上架 |
| `CourseDown` | `IdsRequest { ids }` | `Empty {}` | 批量下架 |

#### Scenario: 单课上架（发布）

- **WHEN** 调用 `CourseUpShelf(id)`
- **THEN** 系统将草稿及子表复制到正式表并置已上架(2)（详见「课程上下架状态机」）

#### Scenario: 上架前置校验（RPC）

- **WHEN** 调用 `CourseCheckUpShelf(id)`
- **THEN** 系统校验 `draft.Step >= 5`，结果以错误传递：返回 `nil` = 可上架，返回 `BadRequest` = 不可上架（proto 返回值为 `Empty`）

### Requirement: 分类管理业务规则

分类 SHALL 为三层树形结构（一级→二级→三级），通过 `parent_id` 自引用，并满足：

- 层级计算：`CategoryAdd` 的 `level` 不由前端决定，而是根据 `parent_id` 自动推导——无 parent（parent_id=0）为一级，否则取其父级的 level+1。
- 三级分类：level=3 的分类下可关联课程，一级/二级分类不直接关联课程。
- 状态控制：status=1(正常)、2(禁用)，禁用的分类不影响已有课程。
- 同级排序：通过 `priority` 字段控制同级排序（值越小越前，`buildCategoryTree` 映射为 `Index`）。
- 递归树查询：`CategoryListAll` 返回递归树结构，`CategoryListOneLevel` 仅返回一级。
- 递归删除：`CategoryDelete` 级联删除其下所有子分类（`deleteRecursive`），被课程引用的三级分类不可删。
- 数量回填：`CategoryGet` / `CategoryList*` 会回填 `courseNum`（关联课程数）与 `thirdCategoryNum`（下属三级分类数）。
- 分类存在性：保存课程时 `third_cate_id` 需能落到三级分类（`CourseBaseInfoSave` 入参校验）。

#### Scenario: 新增分类自动推导层级

- **WHEN** 调用 `CategoryAdd` 且 `parent_id=0`
- **THEN** 该分类创建为一级分类
- **AND** 当 `parent_id` 指向已有分类时，level 自动取父级 level+1，前端传入值不生效

#### Scenario: 删除被课程引用的三级分类

- **WHEN** 调用 `CategoryDelete` 且目标三级分类已被课程引用
- **THEN** 系统拒绝删除
- **AND** 删除未被引用的分类时级联删除其下所有子分类（`deleteRecursive`）

### Requirement: 课程信息保存流程（草稿 + Step 进度）

课程编辑 SHALL 全程落在 `*_draft` 草稿表，按 Step 逐步推进；仅 `CourseUpShelf` 时才复制到正式表。

Step 枚举（由代码常量确定）：

| Step | 含义 | 推进点 |
|------|------|--------|
| 1 | 基本信息 | `CourseBaseInfoSave`（新增时写 `Step=1`；更新时若 `<1` 补为 1） |
| 2 | 课程目录 | `CourseCatalogueSave`（`advanceStep` 至少置 2，且不回退已有进度） |
| 3 | 课程视频 | `CourseMediaSave`（保存媒资信息，不直接推进 step，见缺口） |
| 4 | 课程题目 | `CourseSubjectsSave`（不直接推进 step） |
| 5 | 课程老师 | `CourseTeachersSave`（不直接推进 step） |

> 已知缺口：只有 `CourseBaseInfoSave` 与 `CourseCatalogueSave` 会推进 `step`；`CourseMediaSave` / `CourseSubjectsSave` / `CourseTeachersSave` 只落各自数据、不修改 `step`。而 `CourseCheckUpShelf` 要求 `draft.Step >= 5` 才允许上架，因此达到 step=5 目前依赖前端在 `CourseCatalogueSave` 的 `Step` 入参中传入累计进度（或后续补足三者的 step 推进逻辑）。

课程名称校验 SHALL 为：`CourseCheckName` 用于编辑时校验课程名称是否重复，传入 `id` 时排除自身（草稿表 `FindByNameExceptId`），新增（`id=0`）时全局查重。

> **路由映射冲突说明**：business-rules.md 的「保存方法映射」将五个保存方法写作 `POST /courses/draft/{id}/1`、`POST /courses/draft/media`、`POST /courses/draft/{id}/4`、`POST /courses/draft/{id}/5` 形态；api-spec.md 与 go-zero 实际路由 `apps/course/api/course.api` 为 `POST /courses/baseInfo/save`、`POST /courses/media/save/{id}`、`POST /courses/subjects/save/{id}`、`POST /courses/teachers/save`、`POST /courses/catas/save/{id}/{step}`，以 api-spec.md / `.api` 为准（本规格 HTTP 接口表已按后者转写）。

#### Scenario: 新增课程初始化 Step

- **WHEN** 调用 `CourseBaseInfoSave` 新增课程
- **THEN** 系统写入课程草稿表并置 `Step=1`
- **AND** 更新已有课程时若 `Step < 1` 则补为 1

#### Scenario: 保存目录推进进度

- **WHEN** 调用 `CourseCatalogueSave`
- **THEN** `advanceStep` 将 step 至少置 2，且不回退已有进度

#### Scenario: 后三步保存不推进进度

- **WHEN** 调用 `CourseMediaSave` / `CourseSubjectsSave` / `CourseTeachersSave`
- **THEN** 系统仅落各自数据、不修改 step（已知缺口：达成 step=5 依赖前端经 `CourseCatalogueSave` 的 Step 入参传累计进度，否则 `CheckUpShelf` 永不通过）

### Requirement: 课程目录管理（草稿优先）

目录编辑 SHALL 先写 `course_catalogue_draft`，上架时整体复制到 `course_catalogue`；支持章/节/测试三种类型，通过 `parent_catalogue_id` 构建树：

- 类型：type=1 章、2 节、3 测试/练习，仅允许 1/2/3，缺省时按层级默认（一级=章、下级=节）。
- 层级：章的 `parent_catalogue_id=0`，节和测试指向所属章。
- 草稿主键：草稿目录 id 为雪花 ID（非自增），写入时需先 `nextID()` 生成再作为子节点 parent。
- 试看：`trailer` 标记是否支持试看（0/1）。
- 排序：`c_index` 优先取请求 `index`，缺省时按同级顺序从 1 编号。
- 正式表同步：上架时 `CourseCatalogueModel` 先 `Delete` 旧数据再按草稿逐行 `Insert`（id 保持一致）。

#### Scenario: 保存章/节/测试树

- **WHEN** 调用 `CourseCatalogueSave` 提交目录树且未传 type
- **THEN** 系统按层级默认类型（一级=章、下级=节）递归落库，草稿目录 id 由 `nextID()` 生成雪花 ID

#### Scenario: 上架同步正式目录

- **WHEN** 课程上架触发正式目录同步
- **THEN** 系统先删除旧正式目录数据，再按草稿逐行 Insert 且 id 保持一致

### Requirement: 课程上下架状态机

课程状态 SHALL 由 `CourseUpShelf` / `CourseDownShelf` / `CourseComplete` 控制，全部基于正式课程表 `course`：

| 状态 | 值 | 说明 | 入口 |
|------|----|------|------|
| 待上架 | 1 | 草稿状态，已保存但未发布 | `CourseBaseInfoSave` 初始值 |
| 已上架 | 2 | 对外可见 | `CourseUpShelf` / `CourseUp` |
| 下架 | 3 | 隐藏但信息保留 | `CourseDownShelf` / `CourseDown` |
| 已完结 | 4 | 课程结束 | `CourseComplete` |

状态流转：待上架(1) ──CourseUpShelf/CourseUp──▶ 已上架(2)；已上架(2) ──CourseDownShelf/CourseDown──▶ 下架(3)；下架(3) ──CourseComplete──▶ 已完结(4)；源图另绘有下架(3) 回到待上架(1) 的回边（未标注触发方法）。

> **关于状态 6（已申请退款）**：DDL 列注释中定义了 `6` 退款态，但 course 服务自身逻辑不处理该状态——它由 trade/pay 服务在退款流程中直接写 `course.status`。course 的 Up/DownShelf 只涉及 1–4。阅读旧文档的「状态机含 6」需以本段为准。
>
> **状态值域说明**（依 data-model.md）：草稿表 `course_draft` 的 status 取值为 0=待上架、1=已上架、2=下架、3=已完结，与正式表 `course` 的 1~4 不同；本状态机 1–4 指正式表。

上架前置校验 SHALL 为 `CourseCheckUpShelf`：校验 `draft.Step >= 5`（1 基本信息 / 2 目录 / 3 视频 / 4 题目 / 5 老师均已保存）。proto 返回值为 `Empty`，故校验结果以错误传递：返回 `nil` = 可上架（API 侧 `existed=true`），返回 `BadRequest` = 不可上架（`existed=false`）。

单课上架 `CourseUpShelf`（发布）SHALL 把草稿及其子表复制到正式表，状态置为已上架(2)：

1. `course_draft → course`（主键沿用草稿雪花 id；已发布过则 `PublishTimes+1` 并保留原 `creater`）
2. `course_content_draft → course_content`（按 id `Upsert`；草稿缺失写空占位）
3. `course_catalogue_draft → course_catalogue`（先清空旧正式目录，再逐行复制）
4. `course_teacher_draft → course_teacher`（先清空旧正式老师，再逐行复制）

草稿数据 SHALL 保留，支持再次编辑后重新发布（重新执行 `CourseUpShelf` 覆盖正式表）。

批量操作 SHALL 为：`CourseUp` / `CourseDown` 支持批量上/下架，入参为逗号分隔的 `courseIds` 字符串（`IdsRequest`）。

#### Scenario: 单课上架复制草稿

- **WHEN** 调用 `CourseUpShelf(id)` 且上架校验通过
- **THEN** 系统依次执行上述 4 步草稿复制并将状态置为已上架(2)
- **AND** 已发布过的课程 `PublishTimes+1` 并保留原 `creater`；课程内容草稿缺失时写空占位

#### Scenario: 重新发布覆盖正式表

- **WHEN** 已上架课程再次编辑草稿并重新执行 `CourseUpShelf`
- **THEN** 草稿数据保留，正式表被最新草稿覆盖

#### Scenario: 批量上下架入参

- **WHEN** 调用 `CourseUp` / `CourseDown`
- **THEN** 入参 `IdsRequest` 携带逗号分隔的 `courseIds` 字符串，系统批量切换状态 2/3

#### Scenario: 退款态由外部服务写入

- **WHEN** trade/pay 退款流程需要标记课程退款态
- **THEN** 由 trade/pay 直接写 `course.status=6`，course 服务的上/下架逻辑不处理该状态

### Requirement: 课程媒体绑定

媒资信息 SHALL 通过 `CourseMediaSave` 回填到课程目录草稿小节（`course_catalogue_draft`）：

- 绑定粒度：精确到目录（小节）级别，按 `cata_id` 定位。
- 回填字段：`media_id`、`video_name`、`media_duration`、`trailer`（bool→0/1）。
- 媒资存在性：`cata_id` 必须在课程目录草稿中存在，否则返回 `NotFound`。
- 引用计数：`CourseMediaUseInfo` 统计各 `media_id` 被目录引用次数（`MediaQuoteList`，来自 `course_catalogue` 正式表 `CountByMediaIds`）。

> 注：旧文档写作 `CourseMediaBind`，实际 RPC 方法名为 **`CourseMediaSave`**。
>
> 已知缺口：保存媒资时仅校验本地目录存在；`servicecontext.go` 未装配 `MediaRpc`，媒资/题目跨服务存在性校验未接线（见 implementation-status 2.5）。

#### Scenario: 媒资回填小节

- **WHEN** 调用 `CourseMediaSave` 且 `cata_id` 存在于课程目录草稿
- **THEN** 系统将 `media_id`、`video_name`、`media_duration`、`trailer`（bool→0/1）回填到对应小节

#### Scenario: cata_id 不存在

- **WHEN** 调用 `CourseMediaSave` 且 `cata_id` 不在课程目录草稿中
- **THEN** 系统返回 `NotFound`

#### Scenario: 媒资引用统计

- **WHEN** 调用 `CourseMediaUseInfo`
- **THEN** 系统基于正式表 `course_catalogue` 的 `CountByMediaIds` 返回各 `media_id` 被目录引用次数（`MediaQuoteList`）

### Requirement: 课程查询

后台管理查询与门户端查询 SHALL 共用同一返回结构 `CoursePageQueryReply`，差异在过滤维度，且均走 `course` 正式表（非草稿）。

后台查询 `CoursePageQuery` 过滤条件：

| 过滤条件 | 说明 |
|---------|------|
| keyword | 课程名称模糊搜索 |
| status | 按状态过滤 |
| free | 免费/付费 |
| course_type | 直播/录播 |
| first_cate_id / second_cate_id / third_cate_id | 三级分类筛选 |
| begin_time / end_time | 时间范围 |

- 分页经 `pkg/utils/page`（`Normalize` + `CalcPages`），结果一次性回填一/二/三级分类名称（来自 `CategoryModel.ListAll`）。
- 门户查询 `CoursePortalQuery` 维度同上，门户通常只查已上架(2)。

其他查询入口：

| 方法 | 用途 |
|------|------|
| `CourseFullInfoGet` | 聚合课程全量信息（基础 + 内容 + 目录 + 老师 + 题目） |
| `CourseSimpleInfoList` | 批量 ID 查简版信息（供其他服务回填） |
| `CourseInfoByTeacherIds` | 按老师 ID 反查其关联课程数 |
| `CourseSearchInfoForIndex` | 供搜索索引构建的结构化信息 |
| `CourseName2Ids` | 名称→ID 列表（批量解析） |
| `CourseCatalogsGet` / `CourseCatalogueTreeGet` | 课程目录树（正式 / 草稿） |
| `CourseSectionGet` / `CourseCatalogueSectionInfo` | 单节/小节详情 |

> 已知缺口：`sold` / `enroll_num` / `score` 在 course 库无列，分页/详情中这些字段填 0，真实值由 trade / learning 服务提供，course 暂未接线对应 RPC。

#### Scenario: 后台分页查询回填分类名称

- **WHEN** 调用 `CoursePageQuery`（keyword/status/free/course_type/三级分类/时间范围任意组合）
- **THEN** 系统分页（`Normalize` + `CalcPages`）返回 `CoursePageQueryReply`
- **AND** 结果一次性回填一/二/三级分类名称（来自 `CategoryModel.ListAll`）

#### Scenario: 门户端仅查已上架

- **WHEN** 调用 `CoursePortalQuery`
- **THEN** 系统从 `course` 正式表查询（门户通常限定已上架(2)），返回 `CoursePageQueryReply`

### Requirement: 课程老师管理

课程老师关系 SHALL 满足：

- 展示控制：`is_show`(0/1) 控制用户端是否展示该老师。
- 排序：`c_index` 按入参顺序从 1 编号。
- 替换式保存：`CourseTeachersSave` 先 `DeleteByCourseId` 再按入参顺序重插，一次请求即全量覆盖。
- 老师详情：`CourseTeachersGet` 仅返回 `teacher_id` / `is_show` / `c_index`。

> 已知缺口：老师姓名/头像来自 user 服务，course 库无对应列，未接线 `UserRpc`，`CourseTeachersGet` 详情字段留空。

#### Scenario: 替换式保存老师

- **WHEN** 调用 `CourseTeachersSave`
- **THEN** 系统先 `DeleteByCourseId` 再按入参顺序重插（全量覆盖），`c_index` 从 1 编号

#### Scenario: 老师详情字段留空

- **WHEN** 调用 `CourseTeachersGet`
- **THEN** 系统仅返回 `teacher_id` / `is_show` / `c_index`，姓名/头像字段留空（user 服务未接线）

### Requirement: 数据模型与表清单

course 域 SHALL 使用 MySQL `tj_course` 库（DDL：`sql/ddl/tj_course.sql`），正式表与草稿表成对组织。所有 `*_gen.go` 由 goctl 生成，禁止修改；12 个自定义 Model 文件（`categorymodel.go`、`coursemodel.go`、`coursedraftmodel.go`、`coursecontentmodel.go`、`coursecontentdraftmodel.go`、`coursecataloguemodel.go`、`coursecataloguedraftmodel.go`、`coursecatasubjectdraftmodel.go`、`courseteachermodel.go`、`courseteacherdraftmodel.go`、`coursesubjectmodel.go`、`subjectmodel.go`）均为标准 CRUD 接口包装，无额外业务方法（与 auth 域不同）。

表清单（表名 / 用途 / 关键字段）：

| 表 | 用途 | 关键字段 |
|----|------|---------|
| `category` | 课程分类表 | `name` varchar(50) 唯一约束；`parent_id`（一级=0，树形查询）；`level` 1/2/3；`priority` 同级排序（值越小越前）；`status` 1=正常 2=禁用；`deleted` 逻辑删除 |
| `course` | 正式课程表 | `id` 主键 = 课程草稿 ID；`name` varchar(80)；`course_type` 1=直播课 2=录播课；`first_cate_id`/`second_cate_id`/`third_cate_id`（外键 → category）；`free` 0=付费 1=免费；`price` 价格(分)；`template_type` 1=固定 2=自定义；`template_url`；`status` 1~4；`step` 1~5；`score`（45=4.5星）；`media_duration`；`valid_duration` 有效期(月)；`section_num`；`dep_id`；`publish_times`/`publish_time`；`purchase_start_time`/`purchase_end_time`（后者有索引） |
| `course_draft` | 课程草稿表 | 结构与 `course` 高度相似；差异：status 0=待上架 1=已上架 2=下架 3=已完结；多 `can_update`（是否可更新）与 `c_version`（版本号）；`media_duration`/`publish_time` 无固定宽度 |
| `course_content` / `course_content_draft` | 课程详情（大文本） | `id` 主键 = course.id；`course_introduce` varchar(512)；`use_people` varchar(512)；`course_detail` varchar(1024)；`dep_id` |
| `course_catalogue` / `course_catalogue_draft` | 课程目录（正式/草稿） | `type` 1=章 2=节 3=测试；`parent_catalogue_id`（章=0，树形查询）；`course_id`（外键 → course）；`media_id`（外键 → media）；`video_id`/`video_name`；`living_start_time`/`living_end_time`（直播课）；`play_back` 回放；`media_duration` 秒；`trailer` 试看；`c_index` 排序；`can_update`（目录草稿专用）；`dep_id` |
| `course_cata_subject_draft` | 目录-题目关系草稿 | `course_id`；`cata_id`（目录/小节 ID）；`subject_id` |
| `course_subject` | 正式课程-题目关系 | `course_id`；`subject_id` |
| `course_teacher` / `course_teacher_draft` | 课程老师关系（正式/草稿） | `course_id`；`teacher_id`（user 域）；`is_show`；`c_index`；`dep_id` |
| `subject` | 题目表 | `name` varchar(512) 题干；`subject_type` 1=单选 2=多选 3=不定向 4=判断 5=主观；`difficulty` 1=简单 2=中等 3=困难；`option1`~`option10` varchar(512)；`answer` varchar(40)；`analysis` varchar(1024)；`use_times`/`answer_times` 统计；`score`；`dep_id` |

关系概览：`category`(1)─(N)`course`（经一级/二级/三级分类自关联）；`course`(1)─(1)`course_content`；`course`(N)─`course_teacher`─(N)`user`(teacher)；`course`(N)─`course_catalogue`─(N)`course_subject`─`subject`；草稿表与正式表五对一一对应（`course_draft ≈ course`、`course_catalogue_draft ≈ course_catalogue`、`course_teacher_draft ≈ course_teacher`、`course_content_draft ≈ course_content`、`course_cata_subject_draft ≈ course_subject`）。

#### Scenario: 草稿编辑、正式发布

- **WHEN** 课程处于编辑态
- **THEN** 数据写入 `course_draft` / `course_content_draft` / `course_catalogue_draft` / `course_teacher_draft` / `course_cata_subject_draft` 草稿表
- **AND** 上架时复制到对应正式表

#### Scenario: 目录树自引用

- **WHEN** 构建课程目录树
- **THEN** 章的 `parent_catalogue_id=0`，节/测试的 `parent_catalogue_id` 指向所属章（type 1/2/3）

### Requirement: 存储不变量（唯一键、软删除、状态机）

course 域数据 SHALL 满足以下不变量：

- 唯一键：`category.name` 有唯一约束（varchar(50)）。
- 主键约定：`course.id` 主键沿用课程草稿雪花 ID；`course_content.id = course.id`（一对一）；草稿目录 `course_catalogue_draft.id` 为雪花 ID（非自增）。
- 软删除：按 data-model.md 字段清单，`category`、`course`、`course_draft`、`course_content`(draft)、`course_catalogue`(draft)、`course_teacher`(draft)、`subject` 均含 `deleted` tinyint 逻辑删除列与 `create_time`/`update_time`（自动填充）、`creater`/`updater` 审计字段；`course_cata_subject_draft` 与 `course_subject` 的字段清单未列 `deleted` 与审计字段。
- 状态机：正式表 `course.status ∈ {1=待上架, 2=已上架, 3=下架, 4=已完结}`；草稿表 `course_draft.status ∈ {0=待上架, 1=已上架, 2=下架, 3=已完结}`；DDL 列注释另定义 6=退款态，由 trade/pay 外部写入，course 服务自身不写。
- 枚举约束：`course.course_type`（1=直播课 2=录播课）、`course.free`（0=付费 1=免费）、`course.template_type`（1=固定 2=自定义）、`category.status`（1=正常 2=禁用）、`category.level`（1/2/3）、`course_catalogue.type`（1=章 2=节 3=测试）、`course.step`（1~5）、`subject.subject_type`（1~5）、`subject.difficulty`（1~3）。
- 跨域语义引用：`course.first/second/third_cate_id → category`、`course_catalogue.media_id → media`（媒资域）、`course_teacher.teacher_id → user`（用户域）均为源文件标注的外键，course 库内无对应表，靠应用层维护一致性。

#### Scenario: 分类名唯一约束

- **WHEN** 新增或更新分类时使用已存在的 `name`
- **THEN** `category.name` 的唯一约束拒绝该写入

#### Scenario: 正式表与草稿表状态值域不同

- **WHEN** 写入课程状态字段
- **THEN** 正式表 `course.status` 取值限定 1~4，草稿表 `course_draft.status` 取值限定 0~3

### Requirement: 服务配置与可观测性

course 服务 SHALL 按以下配置运行。

API 服务配置（`apps/course/api/etc/course-api.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Name` | `course-api` | 服务名称 |
| `Host` | `0.0.0.0` | 监听地址 |
| `Port` | `8803` | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| `CourseRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | etcd 地址 |
| `CourseRpc.Etcd.Key` | `course.rpc` | course RPC 服务发现 key |

RPC 服务配置（`apps/course/rpc/etc/course.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Name` | `course.rpc` | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8083` | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | etcd 地址 |
| `Etcd.Key` | `course.rpc` | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_course?charset=utf8mb4&parseTime=true&loc=Local` | MySQL 连接串（字符集 utf8mb4、时区 Local） |
| `Cache[0].Host` / `Type` / `Pass` | `127.0.0.1:6379` / `node` / (空) | Redis 单机节点缓存 |

- 端口分配：`course-api` 8803（HTTP）、`course.rpc` 8083（gRPC）。
- JWT 密钥：`Auth.AccessSecret` 与 auth 服务共用同一值，确保各服务间可互相校验 Token；生产环境必须修改默认值。
- 依赖外部服务：MySQL `tj_course` 库（自建存储）、Redis 缓存（节点模式，缓存分类/课程数据）。
- 可观测性注入配置（**已注入勿删**，configs.md 未记录、以实际 yaml 为准）：`Log`（Mode: file，Encoding: json，Level: info，路径 `logs/course-api` 与 `logs/course-rpc`）；`Prometheus`（course-api `0.0.0.0:9203`、course.rpc `0.0.0.0:9103`，Path `/metrics`）；`Telemetry`（Name 对应各服务，Endpoint `127.0.0.1:4318`，Sampler `1.0`，Batcher `otlphttp`）。

> 实际 yaml 中另有 RabbitMQ 配置（configs.md 未记录，代码已接线）：用于发布课程上架/下架事件（`course.events` 交换机，`CourseUpShelf`/`CourseDownShelf` 经 `PublishCourseEvent` 发布），供 search 服务消费并增量同步 ES 课程索引；MQ 未配置或不可用时发布失败仅告警，不阻塞课程操作。

#### Scenario: 上下架事件发布到 MQ

- **WHEN** 课程上架或下架成功
- **THEN** 系统经 RabbitMQ 向 `course.events` 交换机发布事件，供 search 服务同步 ES 课程索引
- **AND** MQ 未配置或不可用时发布失败仅告警，不阻塞课程操作

#### Scenario: JWT 密钥不一致导致鉴权失败

- **WHEN** `course-api` 的 `Auth.AccessSecret` 与 auth 服务签发密钥配置不一致
- **THEN** 各接口鉴权失败；两者必须配置为一致，且生产环境必须替换默认值 `change-me-in-production`

