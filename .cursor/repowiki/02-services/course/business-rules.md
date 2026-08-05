> 版本：v2.0 | 更新：2026-08-05 | 来源：`apps/course/{api,rpc}/internal/logic/**/*logic.go`（已实现，对照 proto / DDL 校对）

---

# Course Business Rules

> 本文档反映 **已落地实现**（76/76 logic 全部实现并通过编译）。与旧版「设计意图」的差异均已校正：状态机仅 1–4（退款态 6 由 trade/pay 外部驱动）、媒资绑定方法为 `CourseMediaSave`、保存流程以草稿表 + step 进度驱动。

---

## 1. 分类管理

**核心规则**：分类为三层树形结构（一级→二级→三级），通过 `parent_id` 自引用。

| 规则 | 说明 |
|------|------|
| 层级计算 | `CategoryAdd` 的 `level` 不再由前端决定，而是**根据 parent_id 自动推导**：无 parent（parent_id=0）为一级，否则取其父级的 level+1 |
| 三级分类 | level=3 的分类下可关联课程，一级/二级分类不直接关联课程 |
| 状态控制 | status=1(正常), 2(禁用)，禁用的分类不影响已有课程 |
| 同级排序 | 通过 `priority` 字段控制同级排序（`buildCategoryTree` 映射为 `Index`） |
| 递归树查询 | `CategoryListAll` 返回递归树结构，`CategoryListOneLevel` 仅返回一级 |
| 递归删除 | `CategoryDelete` 级联删除其下所有子分类（`deleteRecursive`），被课程引用的三级分类不可删 |
| 数量回填 | `CategoryGet` / `CategoryList*` 会回填 `courseNum`（关联课程数）与 `thirdCategoryNum`（下属三级分类数） |

---

## 2. 课程信息保存流程（草稿 + Step 进度）

**核心规则**：课程编辑全程落在 `*_draft` 草稿表，按 Step 逐步推进；仅 `CourseUpShelf` 时才复制到正式表。

### Step 枚举（由代码常量确定）

| Step | 含义 | 推进点 |
|------|------|--------|
| 1 | 基本信息 | `CourseBaseInfoSave`（新增时写 `Step=1`；更新时若 `<1` 补为 1） |
| 2 | 课程目录 | `CourseCatalogueSave`（`advanceStep` 至少置 2，且不回退已有进度） |
| 3 | 课程视频 | `CourseMediaSave`（保存媒资信息，**不直接**推进 step，见缺口） |
| 4 | 课程题目 | `CourseSubjectsSave`（**不直接**推进 step） |
| 5 | 课程老师 | `CourseTeachersSave`（**不直接**推进 step） |

> ⚠️ **已知缺口**：只有 `CourseBaseInfoSave` 与 `CourseCatalogueSave` 会推进 `step`；`CourseMediaSave` / `CourseSubjectsSave` / `CourseTeachersSave` 只落各自数据、不修改 `step`。而 `CourseCheckUpShelf` 要求 `draft.Step >= 5` 才允许上架。因此达到 step=5 目前依赖前端在 `CourseCatalogueSave` 的 `Step` 入参中传入累计进度（或后续补足三者的 step 推进逻辑）。详情见文末「已知缺口」。

### 保存方法映射（对照 API 路由）

```
CourseBaseInfoSave   ← POST /courses/draft/{id}/1   基本信息 + 课程内容(介绍/详情/适用人群)
CourseCatalogueSave  ← POST /courses/draft/{id}/2   章(1)→节(2)→测试(3) 树，递归落库
CourseMediaSave      ← POST /courses/draft/media    绑定媒资到小节(cata_id)
CourseSubjectsSave   ← POST /courses/draft/{id}/4   题目绑定到章节
CourseTeachersSave   ← POST /courses/draft/{id}/5   老师绑定(替换式)
```

### 课程名称校验

`CourseCheckName` 用于编辑时校验课程名称是否重复：传入 `id` 时**排除自身**（draft 表 `FindByNameExceptId`），新增（`id=0`）时全局查重。

---

## 3. 课程目录管理（草稿优先）

**核心规则**：目录编辑先写 `course_catalogue_draft`，上架时整体复制到 `course_catalogue`，支持章/节/测试三种类型，通过 `parent_catalogue_id` 构建树。

| 规则 | 说明 |
|------|------|
| 类型 | type=1 章, type=2 节, type=3 测试/练习 |
| 层级 | 章的 `parent_catalogue_id=0`，节和测试指向所属章；未传 type 时按层级默认（一级=章、下级=节） |
| 草稿主键 | 草稿目录 id 为**雪花 ID**（非自增），写入时需先 `nextID()` 生成再作为子节点 parent |
| 试看 | `trailer` 标记是否支持试看（0/1） |
| 排序 | `c_index` 优先取请求 `index`，缺省时按同级顺序从 1 编号 |
| 正式表同步 | 上架时 `CourseCatalogueModel` 先 `Delete` 旧数据再按草稿逐行 `Insert`（id 保持一致） |

---

## 4. 课程上下架（状态机）

**核心规则**：课程状态由 `CourseUpShelf` / `CourseDownShelf` / `CourseComplete` 控制，全部基于**正式课程表 `course`**。

### 课程服务内部状态（已实现）

| 状态 | 值 | 说明 | 入口 |
|------|------|------|------|
| 待上架 | 1 | 草稿状态，已保存但未发布 | `CourseBaseInfoSave` 初始值 |
| 已上架 | 2 | 对外可见 | `CourseUpShelf` / `CourseUp` |
| 下架 | 3 | 隐藏但信息保留 | `CourseDownShelf` / `CourseDown` |
| 已完结 | 4 | 课程结束 | `CourseComplete` |

```
  待上架(1) ──CourseUpShelf/CourseUp──▶ 已上架(2)
     ▲                                     │
     │                              CourseDownShelf/CourseDown
     │                                     ▼
     └─────────────────────────────── 下架(3) ──CourseComplete──▶ 已完结(4)
```

> 📌 **关于状态 6（已申请退款）**：DDL 列注释中定义了 `6` 退款态，但 **course 服务自身逻辑不处理该状态**——它由 trade/pay 服务在退款流程中直接写 `course.status`。course 的 `Up/DownShelf` 只涉及 1–4。阅读旧文档的「状态机含 6」需以本段为准。

### 上架前置校验 `CourseCheckUpShelf`

校验 `draft.Step >= 5`（1 基本信息 /2 目录 /3 视频 /4 题目 /5 老师 均已保存）。proto 返回值为 `Empty`，故校验结果**以错误传递**：返回 `nil` = 可上架（API 侧 `existed=true`），返回 `BadRequest` = 不可上架（`existed=false`）。

### 单课上架 `CourseUpShelf`（发布）

把草稿及其子表复制到正式表，状态置为 `已上架(2)`：

1. `course_draft → course`（主键沿用草稿雪花 id；已发布过则 `PublishTimes+1` 并保留原 `creater`）
2. `course_content_draft → course_content`（按 id `Upsert`；草稿缺失写空占位）
3. `course_catalogue_draft → course_catalogue`（先清空旧正式目录，再逐行复制）
4. `course_teacher_draft → course_teacher`（先清空旧正式老师，再逐行复制）

草稿数据**保留**，支持再次编辑后重新发布（重新执行 `CourseUpShelf` 覆盖正式表）。

### 批量操作

`CourseUp` / `CourseDown` 支持批量上/下架，入参为逗号分隔的 `courseIds` 字符串（`IdsRequest`）。

---

## 5. 课程媒体绑定

**核心规则**：媒资信息通过 `CourseMediaSave` 回填到**课程目录草稿小节**（`course_catalogue_draft`）。

| 规则 | 说明 |
|------|------|
| 绑定粒度 | 精确到目录（小节）级别，按 `cata_id` 定位 |
| 回填字段 | `media_id`、`video_name`、`media_duration`、`trailer`(bool→0/1) |
| 媒资存在性 | `cata_id` 必须在课程目录草稿中存在，否则返回 `NotFound` |
| 引用计数 | `CourseMediaUseInfo` 统计各 `media_id` 被目录引用次数（`MediaQuoteList`，来自 `course_catalogue` 正式表 `CountByMediaIds`） |

> ⚠️ 旧文档写作 `CourseMediaBind`，实际 RPC 方法名为 **`CourseMediaSave`**。

---

## 6. 课程查询

**核心规则**：后台管理查询与门户端查询共用同一返回结构 `CoursePageQueryReply`，差异在过滤维度。

### 后台查询 `CoursePageQuery`

| 过滤条件 | 说明 |
|---------|------|
| keyword | 课程名称模糊搜索 |
| status | 按状态过滤 |
| free | 免费/付费 |
| course_type | 直播/录播 |
| first_cate_id / second_cate_id / third_cate_id | 三级分类筛选 |
| begin_time / end_time | 时间范围 |

分页经 `pkg/utils/page`（`Normalize` + `CalcPages`），结果一次性回填一/二/三级分类名称（来自 `CategoryModel.ListAll`）。

### 门户查询 `CoursePortalQuery`

维度同上，门户通常只查 `已上架(2)`。两者均走 `course` 正式表（非草稿）。

### 其他查询入口

| 方法 | 用途 |
|------|------|
| `CourseFullInfoGet` | 聚合课程全量信息（基础 + 内容 + 目录 + 老师 + 题目） |
| `CourseSimpleInfoList` | 批量 ID 查简版信息（供其他服务回填） |
| `CourseInfoByTeacherIds` | 按老师 ID 反查其关联课程数 |
| `CourseSearchInfoForIndex` | 供搜索索引构建的结构化信息 |
| `CourseName2Ids` | 名称→ID 列表（批量解析） |
| `CourseCatalogsGet` / `CourseCatalogueTreeGet` | 课程目录树（正式 / 草稿） |
| `CourseSectionGet` / `CourseCatalogueSectionInfo` | 单节/小节详情 |

---

## 7. 课程老师管理

| 规则 | 说明 |
|------|------|
| 展示控制 | `is_show`(0/1) 控制用户端是否展示该老师 |
| 排序 | `c_index` 按入参顺序从 1 编号 |
| 替换式保存 | `CourseTeachersSave` 先 `DeleteByCourseId` 再按入参顺序重插，一次请求即全量覆盖 |
| 老师详情 | `CourseTeachersGet` 仅返回 `teacher_id` / `is_show` / `c_index`；**老师姓名/头像来自 user 服务**，course 库无对应列，故详情字段留空（已知缺口） |

---

## 关键校验汇总

| 校验项 | 校验方式 | 实现位置 |
|--------|---------|---------|
| 课程名称重复 | `FindByNameExceptId`（编辑排除自身）/ 全局查重（新增） | `CourseCheckName` |
| 上架前置校验 | `draft.Step >= 5` | `CourseCheckUpShelf` |
| 分类存在性 | 保存课程时 `third_cate_id` 需能落到三级分类 | `CourseBaseInfoSave` 入参校验 |
| 目录类型校验 | type 仅允许 1/2/3，缺省按层级默认 | `CourseCatalogueSave` |
| 媒资绑定 | `cata_id` 必须存在于课程目录草稿 | `CourseMediaSave` |
| 分类可被删 | 被课程引用的三级分类拒绝删除 | `CategoryDelete` |

---

## 已知缺口 / 外部依赖（实现边界）

| 缺口 | 影响 | 说明 |
|------|------|------|
| 销量/报名数/评分 | `sold` / `enroll_num` / `score` 在 course 库无列 | 分页/详情中这些字段**填 0**，真实值由 trade / learning 服务提供，course 暂未接线对应 RPC |
| 老师详情 | `CourseTeachersGet` 仅返回 id | 姓名/头像需 user 服务，未接线 `UserRpc`，详情字段留空 |
| 媒资/题目跨服务校验 | 保存媒资/题目时仅校验本地目录存在 | `servicecontext.go` 未装配 `MediaRpc` / `ExamRpc`（见 implementation-status 2.5） |
| Step 推进不完整 | `CourseMediaSave`/`SubjectsSave`/`TeachersSave` 不推进 step | 达到 step=5 依赖前端经 `CourseCatalogueSave.Step` 传累计进度；否则 `CheckUpShelf` 永不通过 |
| 退款态 6 | course 自身不写状态 6 | 由 trade/pay 在服务外写 `course.status`，course 逻辑无对应处理 |
