> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/course/api/internal/logic/*.go`, `apps/course/api/course.api`

---

# Course Business Rules

## 1. 分类管理

**核心规则**：分类为三层树形结构（一级→二级→三级），通过 `parent_id` 自引用。

| 规则 | 说明 |
|------|------|
| 三级分类 | level=3 的分类下可关联课程，一级/二级分类不直接关联课程 |
| 状态控制 | status=1(正常), 2(禁用)，禁用的分类不影响已有课程 |
| 同级排序 | 通过 `index` 字段控制同级排序 |
| 递归树查询 | `CategoryListAll` 返回递归树结构，`CategoryListOneLevel` 仅返回一级 |

---

## 2. 课程信息保存流程

**核心规则**：课程信息按 Step 步骤逐步保存，一个 Step 完成后方可进入下一步。

### Step 流程

```
Step 1: 保存基本信息 (baseInfo)
    ├── 课程名称、封面、价格
    ├── 分类选择（三级分类）
    ├── 售卖方式（付费/免费）
    └── 课程有效期（validDuration）

Step 2: 保存目录结构 (catalogue/save/{id}/2)
    ├── 章 (type=1) → 节 (type=2) → 测试 (type=3)
    └── 递归结构保存

Step 3: 保存课程视频 (catalogue/save/{id}/3)
    └── 绑定媒资到小节 (media_id)

Step 4: 保存题目 (catalogue/save/{id}/4)
    └── 绑定题目到章节

Step 5: 保存课程老师 (catalogue/save/{id}/5 / teachers/save)
    └── 绑定老师到课程

Step 6: 上架
    └── CourseUpShelf → 从 course_draft 复制到 course 表
```

### 课程名称校验

`CourseCheckName` 用于编辑时校验课程名称是否重复，编辑模式下排除自身 ID。

---

## 3. 课程目录管理

**核心规则**：课程目录支持章/节/测试三种类型，通过 `parent_catalogue_id` 构建树形。

| 规则 | 说明 |
|------|------|
| 类型 | type=1 章, type=2 节, type=3 测试 |
| 层级 | 章的 `parent_catalogue_id=0`，节和测试指向所属章 |
| 试看 | trailer 标记是否支持试看 |
| 排序 | `c_index` 控制目录内排序 |

---

## 4. 课程上下架

**核心规则**：课程有 5 种状态，通过 `CourseUpShelf`/`CourseDownShelf` 控制。

```
状态机：

  待上架 (1)  →  CourseUpShelf  →  已上架 (2)
       ↑                        │
       │                CourseDownShelf
       │                        │
  已申请退款 (6)              下架 (3)
       │                        │
       └───── 退款完成 ──────────┘
                            │
                    CourseComplete → 已完结 (4)
```

| 状态 | 值 | 说明 |
|------|------|------|
| 待上架 | 1 | 草稿状态，课程已保存但未发布 |
| 已上架 | 2 | 课程对外可见 |
| 下架 | 3 | 课程隐藏，但信息保留 |
| 已完结 | 4 | 课程结束 |
| 已申请退款 | 6 | 学员已申请退款 |

### 前置校验

`CourseCheckUpShelf` 用于上架前校验课程是否满足条件（需检查 step 进度是否齐全）。

### 批量操作

`CourseUp` / `CourseDown` 支持批量上架/下架，通过逗号分隔的 courseIds 字符串传递。

---

## 5. 课程媒体绑定

**核心规则**：媒资通过 `CourseMediaBind` 绑定到具体章节小节。

| 规则 | 说明 |
|------|------|
| 绑定粒度 | 精确到目录（小节）级别 |
| 视频信息 | 包含 video_name, media_duration, trailer |
| 引用计数 | 通过 `CourseMediaUseInfo` 查询媒资被引用的次数 |

---

## 6. 课程查询

**核心规则**：提供两种查询入口——后台管理查询和门户端查询。

### 后台查询 (`CoursePageQuery`)

| 过滤条件 | 说明 |
|---------|------|
| keyword | 课程名称模糊搜索 |
| status | 按状态过滤 |
| free | 免费/付费 |
| course_type | 直播/录播 |
| first_cate_id / second_cate_id / third_cate_id | 三级分类筛选 |
| begin_time / end_time | 时间范围 |

### 门户查询 (`CoursePortalQuery`)

| 过滤条件 | 说明 |
|---------|------|
| category_id_lv1/lv2/lv3 | 按分类路径过滤 |
| keyword | 名称搜索 |
| free | 免费/付费 |
| type | 课程类型 |
| status | 通常只查已上架 |

---

## 7. 课程老师管理

| 规则 | 说明 |
|------|------|
| 展示控制 | `is_show` 控制用户端是否展示该老师 |
| 排序 | `c_index` 控制老师排序 |
| 批量保存 | 一次请求传入所有老师及展示状态，替换旧关联 |

---

## 关键校验汇总

| 校验项 | 校验方式 |
|--------|---------|
| 课程名称重复 | `CourseCheckName` 模糊匹配 |
| 上架前置校验 | `CourseCheckUpShelf` 检查 step 进度 |
| 分类存在性 | 保存课程时 `third_cate_id` 需存在 |
| 目录类型校验 | type 仅允许 1/2/3 |
| 媒资存在性 | 绑定媒资时需校验媒资存在 |