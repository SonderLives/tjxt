> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_course.sql`

---

# Course Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 说明 |
|------------------|--------|------|
| `categorymodel.go` | `category` | - |
| `coursemodel.go` | `course` | 正式课程表 |
| `coursedraftmodel.go` | `course_draft` | 课程草稿表 |
| `coursecontentmodel.go` | `course_content` | 课程详情（大文本） |
| `coursecontentdraftmodel.go` | `course_content_draft` | 课程详情草稿 |
| `coursecataloguemodel.go` | `course_catalogue` | 正式目录 |
| `coursecataloguedraftmodel.go` | `course_catalogue_draft` | 目录草稿 |
| `coursecatasubjectdraftmodel.go` | `course_cata_subject_draft` | 目录-题目关系草稿 |
| `courseteachermodel.go` | `course_teacher` | 正式课程老师 |
| `courseteacherdraftmodel.go` | `course_teacher_draft` | 课程老师草稿 |
| `coursesubjectmodel.go` | `course_subject` | 正式课程-题目关系 |
| `subjectmodel.go` | `subject` | 题目 |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。

---

## 表清单与字段说明

### 1. `category` — 课程分类表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 自增主键 | PK |
| `name` | varchar(50) | 分类名称 | 唯一约束 |
| `parent_id` | bigint | 父分类 ID，一级分类=0 | 树形查询 |
| `level` | int | 1/2/3 级分类 | 过滤条件 |
| `priority` | int | 同级排序值，值越小越前 | 排序条件 |
| `status` | tinyint | 1=正常, 2=禁用 | 过滤条件 |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |
| `creater` / `updater` | bigint | 创建/更新人 | - |
| `deleted` | tinyint | 逻辑删除 | 过滤条件 |

---

### 2. `course` — 正式课程表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键 = 课程草稿 ID | PK |
| `name` | varchar(80) | 课程名称 | - |
| `course_type` | tinyint | 1=直播课, 2=录播课 | 过滤条件 |
| `cover_url` | varchar(255) | 封面链接 | - |
| `first_cate_id` | bigint | 一级分类 ID | 外键 → category |
| `second_cate_id` | bigint | 二级分类 ID | - |
| `third_cate_id` | bigint | 三级分类 ID | 外键 → category |
| `free` | tinyint | 0=付费, 1=免费 | 过滤条件 |
| `price` | int | 价格(分) | - |
| `template_type` | tinyint | 1=固定模板, 2=自定义模板 | - |
| `template_url` | varchar(255) | 自定义模板链接 | - |
| `status` | tinyint | 1=待上架, 2=已上架, 3=下架, 4=已完结 | 过滤条件 |
| `purchase_start_time` | datetime | 购买开始时间 | - |
| `purchase_end_time` | datetime | 购买结束时间 | 索引 |
| `step` | tinyint | 信息填写进度 (1~5) | 排序条件 |
| `score` | int | 课程评价得分 (45=4.5星) | 排序条件 |
| `media_duration` | int | 课程总时长 | - |
| `valid_duration` | int | 有效期(月) | - |
| `section_num` | int | 课程总节数 | - |
| `dep_id` | bigint | 部门 ID | 过滤条件 |
| `publish_times` | int | 发布次数 | - |
| `publish_time` | datetime | 最近发布时间 | - |
| `create_time` / `update_time` | datetime | 创建/更新时间 | 自动填充 |
| `creater` / `updater` | bigint | 创建/更新人 | - |
| `deleted` | tinyint | 逻辑删除 | 过滤条件 |

---

### 3. `course_draft` — 课程草稿表

与 `course` 表结构高度相似，差异点：

| 字段 | 正式表(course) | 草稿表(course_draft) |
|------|---------------|---------------------|
| status | 1=待上架, 2=已上架, 3=下架, 4=已完结 | 0=待上架, 1=已上架, 2=下架, 3=已完结 |
| step 默认值 | 无 | 无默认值 |
| can_update | 无 | 有（是否可更新） |
| c_version | 无 | 有（版本号） |
| media_duration | 有固定宽度 | 无固定宽度 |
| publish_time | 有 | 有但无固定宽度 |

---

### 4. `course_content` / `course_content_draft` — 课程详情表

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | bigint | 主键 = course.id |
| `course_introduce` | varchar(512) | 课程介绍 |
| `use_people` | varchar(512) | 适用人群 |
| `course_detail` | varchar(1024) | 课程详情 |
| `dep_id` | bigint | 部门 ID |
| `create_time` / `update_time` | datetime | 创建/更新时间 |
| `creater` / `updater` | bigint | 创建/更新人 |
| `deleted` | tinyint | 逻辑删除 |

---

### 5. `course_catalogue` / `course_catalogue_draft` — 课程目录表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键 | PK |
| `name` | varchar(255) | 目录名称 | - |
| `trailer` | tinyint | 是否试看 | 过滤条件 |
| `course_id` | bigint | 课程 ID | 外键 → course |
| `type` | tinyint | 1=章, 2=节, 3=测试 | 排序条件 |
| `parent_catalogue_id` | bigint | 所属章 ID（章=0） | 树形查询 |
| `media_id` | bigint | 媒资 ID | 外键 → media |
| `video_id` | bigint | 视频 ID | - |
| `video_name` | varchar(255) | 视频名称 | - |
| `living_start_time` / `living_end_time` | datetime | 直播时间（直播课） | - |
| `play_back` | tinyint | 是否支持回放 | - |
| `media_duration` | int | 视频时长(秒) | - |
| `c_index` | int | 排序值 | 排序条件 |
| `can_update` | tinyint | 目录草稿专用，是否可更新 | - |
| `dep_id` | bigint | 部门 ID | 过滤条件 |
| `create_time` / `update_time` | datetime | 创建/更新时间 | - |
| `creater` / `updater` | bigint | 创建/更新人 | - |
| `deleted` | tinyint | 逻辑删除 | - |

---

### 6. `course_cata_subject_draft` — 目录-题目关系草稿

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | bigint | 自增主键 |
| `course_id` | bigint | 课程 ID |
| `cata_id` | bigint | 目录(小节) ID |
| `subject_id` | bigint | 题目 ID |

---

### 7. `course_subject` — 正式课程-题目关系

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | bigint | 自增主键 |
| `course_id` | bigint | 课程 ID |
| `subject_id` | bigint | 题目 ID |

---

### 8. `course_teacher` / `course_teacher_draft` — 课程老师关系表

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | bigint | 主键/自增主键 |
| `course_id` | bigint | 课程 ID |
| `teacher_id` | bigint | 老师 ID（user 域） |
| `is_show` | tinyint | 用户端是否展示 |
| `c_index` | int | 排序序号 |
| `dep_id` | bigint | 部门 ID |
| `create_time` / `update_time` | datetime | 创建/更新时间 |
| `creater` / `updater` | bigint | 创建/更新人 |
| `deleted` | tinyint | 逻辑删除 |

---

### 9. `subject` — 题目表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 自增主键 | PK |
| `name` | varchar(512) | 题干 | - |
| `subject_type` | tinyint | 1=单选, 2=多选, 3=不定向, 4=判断, 5=主观 | 过滤条件 |
| `difficulty` | tinyint | 1=简单, 2=中等, 3=困难 | 过滤条件 |
| `option1`~`option10` | varchar(512) | 选项 1~10 | - |
| `answer` | varchar(40) | 正确答案 | - |
| `analysis` | varchar(1024) | 答案解析 | - |
| `use_times` | int | 引用次数 | 统计 |
| `answer_times` | int | 回答次数 | 统计 |
| `score` | int | 分值 | - |
| `dep_id` | bigint | 部门 ID | 过滤条件 |
| `create_time` / `update_time` | datetime | 创建/更新时间 | - |
| `creater` / `updater` | bigint | 创建/更新人 | - |
| `deleted` | tinyint | 逻辑删除 | - |

---

## 关系图

```
category (1) ── (N) course ── (N) course_teacher ── (N) user(teacher)
  │                      │
  │                      ├── (1) course_content ── (N) media
  │                      ├── (N) course_catalogue ── (N) course_subject ── subject
  │                      └── (N) course_cata_subject_draft
  │
  └── (N) course (二级分类)
        │
        └── (N) course (三级分类)

course_draft ≈ course（同一套结构，用于编辑）
course_catalogue_draft ≈ course_catalogue
course_teacher_draft ≈ course_teacher
course_content_draft ≈ course_content
course_cata_subject_draft ≈ course_subject
```

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成。当前自定义扩展：

| 自定义 Model 文件 | 说明 |
|------------------|------|
| `categorymodel.go` | 标准 CRUD |
| `coursemodel.go` | 标准 CRUD（包装 defaultCourseModel） |
| `coursecataloguemodel.go` | 标准 CRUD |
| `coursesubjectmodel.go` | 标准 CRUD |
| `coursecatasubjectdraftmodel.go` | 标准 CRUD |
| `coursedraftmodel.go` | 标准 CRUD |
| `coursecontentmodel.go` | 标准 CRUD |
| `coursecontentdraftmodel.go` | 标准 CRUD |
| `coursecataloguedraftmodel.go` | 标准 CRUD |
| `courseteachermodel.go` | 标准 CRUD |
| `courseteacherdraftmodel.go` | 标准 CRUD |
| `subjectmodel.go` | 标准 CRUD |

注意：course 域的自定义扩展目前仅有接口包装，无额外业务方法（与 auth 域不同）。