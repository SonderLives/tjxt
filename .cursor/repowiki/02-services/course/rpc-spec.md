> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/course/rpc/course.proto`

---

# Course RPC Spec

## 服务名

`Course` — 课程管理微服务，覆盖分类(category)、课程(course)、目录(catalogue)、老师(teacher)、题目(subject)五大领域，通过 etcd 服务发现（key: `course.rpc`）。

## RPC 方法总览

### 分类管理 (Category)

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

---

### 课程基础信息 (Course)

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

---

### 课程目录 (Catalogue)

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseCatalogsGet` | `IdRequest { id }` | `CourseAndSectionView` | 查询课程目录+小节 |
| `CourseCatalogueTreeGet` | `CourseCatalogueQueryRequest { id, see, with_practice }` | `CatalogueTreeList` | 递归树形查询目录结构 |
| `CourseCatalogueSave` | `CourseCatalogueSaveRequest { course_id, step, chapters }` | `Empty {}` | 保存课程目录（含章/节/测试） |
| `CourseCatalogueIndexList` | `IdRequest { id }` | `CataSimpleList` | 课程目录 ID 列表 |
| `CourseCatalogueSectionInfo` | `IdRequest { id }` | `CourseSectionInfo` | 目录小节信息 |
| `CourseCatalogueBatchQuery` | `IdsRequest { ids }` | `CataSimpleList` | 批量查询目录 |

---

### 章节/题目 (Subjects)

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseSubjectsGet` | `IdRequest { id }` | `CataSubjectInfoList` | 查询章节下的题目 |
| `CourseSubjectsSave` | `CourseSubjectsSaveRequest { course_id, subjects }` | `Empty {}` | 绑定章节-题目关系 |

---

### 课程媒体 (Media)

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseMediaSave` | `CourseMediaSaveRequest { course_id, medias }` | `Empty {}` | 绑定课程目录-媒资 |

---

### 课程老师 (Teachers)

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseTeachersGet` | `CourseTeachersGetRequest { id, see }` | `TeacherInfoList` | 查询课程老师 |
| `CourseTeachersSave` | `CourseTeachersSaveRequest { id, teachers }` | `Empty {}` | 保存课程老师关系 |

---

### 课程上下架 (Status)

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CourseUpShelf` | `IdRequest { id }` | `Empty {}` | 单个上架 |
| `CourseDownShelf` | `IdRequest { id }` | `Empty {}` | 单个下架 |
| `CourseCheckUpShelf` | `IdRequest { id }` | `Empty {}` | 上架前置校验 |
| `CourseUp` | `IdsRequest { ids }` | `Empty {}` | 批量上架 |
| `CourseDown` | `IdsRequest { ids }` | `Empty {}` | 批量下架 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `course-api` (自身) | HTTP Handler → `courseclient.Course` RPC | 所有课程管理/学员端接口 |
| `learning-api` | `courseclient.Course` RPC | 学习服务查询课程目录、章节、小节 |

---

## 调用典型场景

1. **课程发布流程**：保存基本信息 → 保存目录（章/节）→ 绑定媒资 → 绑定题目 → 绑定老师 → 检查上架 → 上架
2. **学员查看课程**：`CourseFullInfoGet(id, with_catalogue=true)` → 展示课程详情+目录树
3. **学习进行中**：`CourseCatalogsGet(courseId)` → 获取当前课程的章节小节结构
4. **课程搜索索引同步**：`CourseSearchInfoForIndex(id)` → 生成搜索索引数据
5. **批量上下架**：`CourseUp/CourseDown` → 管理后台批量操作