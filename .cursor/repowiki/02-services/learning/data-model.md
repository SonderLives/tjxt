> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_learning.sql`

---

# Learning Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `learninglessonmodel.go` | `learning_lesson` | GrantCourses, RevokeCourses, FindByUserCourse, ListByUser, ListPlansByUser, FindLatestLearnedByUser, CountByCourse, UpdatePlan, RemoveLesson, UpdateLatestLearn, IncrLearnedSections |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `learninglessonmodel.go` 扩展 `LearningLessonModel` 接口）。

`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`。

> learning 是本仓库中唯一在 model 之上再抽一层 `internal/service/learning.go` 的服务，`LearningService` 承载参数校验与 `xerr` 错误映射，详见 `rpc-spec.md`。

---

## 表清单与字段说明

### 1. `learning_lesson` — 学员课程表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | BIGINT | 雪花主键 | PK |
| `user_id` | BIGINT | 学员 ID | `uk_user_course`, `idx_user_status_created` |
| `course_id` | BIGINT | 课程 ID | `uk_user_course`, `idx_course_status` |
| `status` | TINYINT | 0未开始，1学习中，2完成，3失效，默认 0 | `idx_course_status`, `idx_user_status_created` |
| `week_freq` | INT | 每周学习章节数，可空 | - |
| `plan_status` | TINYINT | 0无计划，1计划中，默认 0 | - |
| `learned_sections` | INT | 已学习章节数，默认 0 | - |
| `latest_section_id` | BIGINT | 最近学习小节，可空 | - |
| `latest_learn_time` | DATETIME | 最近学习时间，可空 | - |
| `create_time` | DATETIME | 创建时间，默认 CURRENT_TIMESTAMP | `idx_user_status_created` |
| `expire_time` | DATETIME | 课程失效时间，可空 | - |
| `update_time` | DATETIME | 更新时间，ON UPDATE CURRENT_TIMESTAMP | 自动更新 |

**索引清单**：

| 索引名 | 类型 | 字段 | 用途 |
|--------|------|------|------|
| `PRIMARY` | 主键 | `id` | 主键查询 |
| `uk_user_course` | 唯一 | `user_id`, `course_id` | 保证一人一课一条记录；`GrantCourses` 的 `ON DUPLICATE KEY UPDATE` 幂等依赖此约束；goctl 据此生成 `FindOneByUserIdCourseId` |
| `idx_course_status` | 普通 | `course_id`, `status` | 支撑 `CountByCourse`（统计课程学习人数） |
| `idx_user_status_created` | 普通 | `user_id`, `status`, `create_time` | 支撑 `ListByUser` / `ListPlansByUser`（按用户分页 + `ORDER BY create_time`） |

> 表注释：学员课程表。引擎 InnoDB，字符集 utf8mb4。DDL 使用 `CREATE TABLE IF NOT EXISTS`（区别于 trade 域的 `DROP TABLE IF EXISTS` + `CREATE TABLE`）。

---

## ⚠️ 缺失表（已知缺口）

`sql/ddl/tj_learning.sql` **只定义了 `learning_lesson` 一张表**，但 RPC 契约中存在依赖其他表的接口。以下表 **DDL 未定义、model 未生成**：

| 缺失表 | 依赖它的 RPC 方法 | 说明 |
|--------|------------------|------|
| `learning_plan` | `PlanSave`, `PlanPage` | 学习计划表。当前设计把「每周章节数」直接放在 `learning_lesson.week_freq` + `plan_status` 两个字段上，`UpdatePlan` 也是更新 lesson 表。但 `PlanPageReply` 需要 `week_finished` / `week_points` / `week_total_plan` 三个**本周维度**的聚合值，`learning_lesson` 无法提供「按周切分」的完成量与积分，缺少独立的计划/周统计表。 |
| `learning_record` | `LearningRecordCommit`, `LearningRecordsByCourse` | 学习记录表。proto 头部注释声明「学习记录不再单独建表，通过更新 `learning_lesson` 的 `latest_section_id` / `latest_learn_time` / `learned_sections` 实现」，但 `LearningRecordsByCourse` 的响应 `LearningRecordsReply.records` 是 `repeated LearningRecordDTO { section_id, moment, duration, finished }`，即**逐小节的进度明细**。`learning_lesson` 单行只能保存一个 `latest_section_id`，**无法还原多小节的 moment/duration/finished 列表**。 |

### 缺口影响面

| 受影响的 RPC 方法 | 具体缺失能力 |
|------------------|-------------|
| `LearningRecordsByCourse` | 无法返回 `records` 数组，只能返回 `id` + `latest_section_id` 两个字段 |
| `LearningRecordCommit` | `section_type`（VIDEO/EXAM）、`moment`、`duration`、`commit_time` 四个入参**无处落库**；`UpdateLatestLearn` 虽接收 `moment` / `duration` 形参，但 SQL 中只用到 `sectionID`，两个参数实际被丢弃 |
| `PlanPage` | `week_finished` / `week_points` / `week_total_plan` 无数据源 |
| `LearningLessonVO.sections` | 课程总章节数不在本表，需调 `course.rpc` 补全 |
| `LearningLessonVO.latest_section_name` / `latest_section_index` | 小节名称与序号不在本表，需调 `course.rpc` 补全 |

### 另一处缺口：`learning-api` 有路由但 RPC 无对应方法

| API 路由 | 说明 |
|---------|------|
| `DELETE /lessons/{courseId}` | `api-spec.md` 中列出「删除指定课程信息」，`internal/service/learning.go` 也实现了 `RemoveLesson`、model 实现了 `RemoveLesson`，但 **`learning.proto` 未定义对应的 RPC 方法**，`learning.api` 也未声明该 handler。能力已在下层就绪，上层未打通。 |

---

## 关系图

```
learning_lesson
  ├── user_id    ──→ user 域 user.id（跨库，无外键）
  ├── course_id  ──→ course 域 course.id（跨库，无外键）
  └── latest_section_id ──→ course 域 section.id（跨库，无外键）

  UNIQUE (user_id, course_id)  ← 一人一课一条

[缺失] learning_plan    ──→ 应关联 learning_lesson.id
[缺失] learning_record  ──→ 应关联 learning_lesson.id + section_id

事件来源：trade → RabbitMQ(order.exchange) → learning → GrantCourses/RevokeCourses
```

---

## 状态枚举

### `learning_lesson.status` — 课程学习状态

DDL 为 tinyint，`learninglessonmodel.go:16-20` 定义了对应常量，proto `LearningLessonVO.status` 为字符串：

| 值 | Go 常量 | proto 字符串 | 含义 |
|----|---------|-------------|------|
| 0 | `LessonStatusNotStart` | `NOT_BEGIN` | 未开始 |
| 1 | `LessonStatusInLearn` | `LEARNING` | 学习中 |
| 2 | `LessonStatusDone` | `FINISHED` | 完成 |
| 3 | `LessonStatusExpired` | `EXPIRED` | 失效 |

### `learning_lesson.plan_status` — 学习计划状态

| 值 | Go 常量 | proto 字符串 | 含义 |
|----|---------|-------------|------|
| 0 | `PlanStatusNone` | `NO_PLAN` | 无计划 |
| 1 | `PlanStatusInPlan` | `PLAN_RUNNING` | 计划中 |

> **映射责任**：DDL/model 用 int64，proto 用 string。数值 ↔ 字符串的转换需在 RPC logic 层完成，当前 logic 为占位实现，映射尚未落地。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
learninglessonmodel_gen.go   ← goctl 生成，只读（Insert/FindOne/FindOneByUserIdCourseId/Update/Delete）
learninglessonmodel.go       ← 手写扩展 11 个业务方法 + 状态常量
```

**扩展方法实现要点**：

| 关注点 | 实现方式 |
|--------|---------|
| 幂等开通 | `GrantCourses` 使用 `INSERT ... ON DUPLICATE KEY UPDATE update_time = NOW()`，依赖 `uk_user_course` 唯一索引 |
| 主键生成 | `idgen.NextID()`（`tjxt/pkg/utils/idgen` 雪花算法） |
| 缓存旁路 | 全部扩展方法使用 `ExecNoCacheCtx` / `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx`，**绕过 goctl 缓存层**直接操作 DB |
| 分页归一 | 私有方法 `listBy`：`pageNo < 1 → 1`，`pageSize < 1 → 10`，`pageSize > 100 → 100` |
| 失效过滤 | `ListByUser` / `ListPlansByUser` / `CountByCourse` / `FindLatestLearnedByUser` 均带 `status <> LessonStatusExpired` |
| 可空字段 | `week_freq` / `latest_section_id` 用 `sql.NullInt64{Int64: v, Valid: v > 0}` 写入 |
| 状态自动跃迁 | `UpdateLatestLearn` 用 `status = IF(status = 0, 1, status)`，首次提交进度时从「未开始」自动变为「学习中」 |
| 计划更新校验 | `UpdatePlan` 检查 `RowsAffected() == 0` 时返回 `sql.ErrNoRows`，由 service 层转 `xerr.NotFound` |
