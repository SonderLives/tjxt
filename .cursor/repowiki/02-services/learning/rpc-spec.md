> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/learning/rpc/learning.proto`

---

# Learning RPC Spec

## 服务名

`Learning` — 学习进度微服务，覆盖我的课表 / 学习计划 / 学习记录，通过 etcd 服务发现（key: `learning.rpc`）。

共 **11 个 RPC 方法**，按业务分为 4 组。

> **设计要点**（proto 头部注释）：学习记录**不再单独建表**，通过更新 `learning_lesson` 的 `latest_section_id` / `latest_learn_time` / `learned_sections` 实现。

## RPC 方法总览

### 课表查询

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `LearningNow` | `Empty {}` | `LearningLessonVO` | 当前正在学习的课程 |
| `LessonPage` | `LessonPageRequest { page_no, page_size, is_asc, sort_by }` | `LessonPageReply { total, pages, list }` | 我的课表分页 |
| `LessonGet` | `LessonRequest { course_id }` | `LearningLessonVO` | 指定课程的学习信息 |
| `LessonValid` | `LessonRequest { course_id }` | `LessonValidReply { lesson_id }` | 学员是否已报名 |
| `LessonCount` | `LessonCountRequest { course_id }` | `LessonCountReply { count }` | 课程的学习人数 |

**响应字段说明（`LearningLessonVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 课表记录 ID（lesson_id） |
| `course_id` | int64 | 课程 ID |
| `course_name` | string | 课程名称（来自 course 服务） |
| `course_cover_url` | string | 课程封面（来自 course 服务） |
| `course_amount` | int64 | 课程数量 |
| `sections` | int32 | 课程总章节数 |
| `learned_sections` | int32 | 已学习章节数 |
| `status` | string | NOT_BEGIN / LEARNING / FINISHED / EXPIRED |
| `plan_status` | string | NO_PLAN / PLAN_RUNNING |
| `week_freq` | int32 | 每周学习章节数 |
| `latest_section_id` | int64 | 最近学习小节 ID |
| `latest_section_name` | string | 最近学习小节名称 |
| `latest_section_index` | int32 | 最近学习小节序号 |
| `create_time` | string | 课程开通时间 |
| `expire_time` | string | 课程失效时间 |
| `latest_learn_time` | string | 最近学习时间 |

> **注意**：`LearningLessonVO.status` / `plan_status` 在 proto 中是**字符串枚举**，而 DDL 中 `learning_lesson.status` / `plan_status` 是 **tinyint**。RPC 层需做数值 ↔ 字符串映射（映射关系见 `data-model.md` 状态枚举章节）。

**`LessonValidReply` 语义**（proto 注释）：返回 `lesson_id`；未报名返回 `0` / `NotFound`，**由调用方判定**。

---

### 学习计划

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `PlanPage` | `LessonPageRequest { page_no, page_size, is_asc, sort_by }` | `PlanPageReply` | 我的学习计划分页 |
| `PlanSave` | `PlanSaveRequest { course_id, freq }` | `Empty {}` | 创建/更新学习计划 |

**响应字段说明（`PlanPageReply`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `total` | int64 | 总记录数 |
| `pages` | int64 | 总页数 |
| `week_finished` | int32 | 本周已完成章节数 |
| `week_points` | int32 | 本周积分 |
| `week_total_plan` | int32 | 本周计划总章节数 |
| `list` | repeated `LearningLessonVO` | 计划中的课程列表 |

**请求字段说明（`PlanSaveRequest`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `course_id` | int64 | 课程 ID |
| `freq` | int64 | 每周学习章节数，落库到 `learning_lesson.week_freq` |

---

### 学习记录

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `LearningRecordCommit` | `LearningRecordCommitRequest { lesson_id, section_id, section_type, moment, duration, commit_time }` | `Empty {}` | 提交学习记录 |
| `LearningRecordsByCourse` | `LessonRequest { course_id }` | `LearningRecordsReply { id, latest_section_id, records }` | 查询某课程的学习记录 |

**请求字段说明（`LearningRecordCommitRequest`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `lesson_id` | int64 | 课表记录 ID |
| `section_id` | int64 | 小节 ID |
| `section_type` | int32 | 0 / 1 / 2，由 api 解析：VIDEO / EXAM 字符串映射 |
| `moment` | int64 | 视频播放进度（秒） |
| `duration` | int64 | 小节总时长（秒） |
| `commit_time` | string | 提交时间 |

**响应字段说明（`LearningRecordDTO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `section_id` | int64 | 小节 ID |
| `moment` | int64 | 播放进度 |
| `duration` | int64 | 总时长 |
| `finished` | bool | 是否已学完 |

> proto 方法注释明确：`LearningRecordCommit` 的作用是「更新 lesson 的 latest_section/learn_time/learned_sections」，不写独立记录表。

---

### 内部事件（来自 trade）

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GrantCourses` | `GrantCoursesRequest { user_id, course_ids }` | `Empty {}` | 通知本服务开通课程 |
| `RevokeCourses` | `GrantCoursesRequest { user_id, course_ids }` | `Empty {}` | 撤销课程 |

**请求字段说明（`GrantCoursesRequest`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `user_id` | int64 | 学员 ID |
| `course_ids` | repeated int64 | 课程 ID 列表 |

> proto 分组注释：`---- 内部：来自 trade 的 mq 事件 / 内部 RPC 调用 ----`。这两个方法既可被 MQ 消费端调用，也可作为内部 RPC 暴露。

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `learning-api` (自身 API 层) | HTTP Handler → `learningclient.Learning` RPC | `apps/learning/api/internal/svc/servicecontext.go:17` 注入 `LearningRpc learningclient.Learning` |

（注：Grep 全部 `apps/*/api/internal/svc/servicecontext.go` 与 `apps/*/rpc/internal/svc/servicecontext.go`，**只有 learning 自身 API 层**引用 `tjxt/apps/learning/rpc/learning`。trade 服务不直连 learning RPC，而是通过 RabbitMQ `order.exchange` 发布 `order.pay` / `order.refund` 事件，由 learning 侧消费后调用 `GrantCourses` / `RevokeCourses`。）

**learning-api 自身依赖的下游**：

| 下游 | 注入位置 | 说明 |
|------|---------|------|
| `learning.rpc` | `apps/learning/api/internal/svc/servicecontext.go:24` | 学习域主逻辑 |
| `course.rpc` | `apps/learning/api/internal/svc/servicecontext.go:25` `CourseRpc courseclient.Course` | 补全 `LearningLessonVO` 中 `course_name` / `course_cover_url` / `sections` 等课程侧字段 |

---

## 调用典型场景

1. **课程开通** → 学员在 trade 支付成功 → trade 发布 `order.pay` 事件 → learning 消费 → `GrantCourses(user_id, course_ids)` 幂等写入 `learning_lesson`
2. **我的课表** → 前端调 `LessonPage` 分页 → learning-api 用 `CourseRpc` 批量补全课程名/封面/章节数 → 返回 `LearningLessonVO` 列表
3. **继续学习** → 前端调 `LearningNow` 取最近学习的一条 → 跳转到 `latest_section_id` 对应小节
4. **提交进度** → 播放器定时调 `LearningRecordCommit` → 更新 `latest_section_id` / `latest_learn_time`，`status` 由 0（未开始）自动跃迁为 1（学习中）
5. **报名校验** → 播放页调 `LessonValid(course_id)` → 返回 `lesson_id`，为 0 或 NotFound 则未报名，拦截播放
6. **学习计划** → 学员调 `PlanSave(course_id, freq)` 设定每周章节数 → `plan_status` 置 1 → `PlanPage` 查看本周完成情况
7. **课程撤销** → 退款成功后 trade 发布 `order.refund` 事件 → learning 消费 → `RevokeCourses` 将 `status` 置 3（失效）、`plan_status` 置 0

---

## 自定义 Model 方法

`learninglessonmodel.go` 在 `LearningLessonModel` 接口上扩展了 **11 个方法**：

| 方法 | 签名要点 | 说明 |
|------|---------|------|
| `GrantCourses` | `(ctx, userID int64, courseIDs []int64) error` | 为用户开通课程（幂等，`ON DUPLICATE KEY UPDATE update_time = NOW()`） |
| `RevokeCourses` | `(ctx, userID int64, courseIDs []int64) error` | 撤销课程，`status` 置 3（失效）、`plan_status` 置 0 |
| `FindByUserCourse` | `(ctx, userID, courseID int64) (*LearningLesson, error)` | 委托 gen 方法 `FindOneByUserIdCourseId` |
| `ListByUser` | `(ctx, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error)` | 分页查询用户的学习记录（不含已失效） |
| `ListPlansByUser` | `(ctx, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error)` | 只看「已设置计划」的分页 |
| `FindLatestLearnedByUser` | `(ctx, userID int64) (*LearningLesson, error)` | 最近学习的一条（用于「我正在学」） |
| `CountByCourse` | `(ctx, courseID int64) (int64, error)` | 该课程的注册人数 |
| `UpdatePlan` | `(ctx, userID, courseID, weekFreq int64) error` | 设置/更新学习计划（每周章节数） |
| `RemoveLesson` | `(ctx, userID, courseID int64) error` | 删除该用户对这门课的学习记录（`status` 置 3） |
| `UpdateLatestLearn` | `(ctx, lessonID, sectionID, moment, duration int64) error` | 提交学习记录时更新最新学习进度 |
| `IncrLearnedSections` | `(ctx, lessonID int64) error` | `learned_sections + 1`，用于学习记录被确认完成时 |

私有辅助方法 `listBy(ctx, cond, args, pageNo, pageSize, asc, orderBy)` 统一承载分页逻辑：`pageNo < 1` 归一为 1，`pageSize < 1` 归一为 10，`pageSize > 100` 截断为 100。

gen 层（`learninglessonmodel_gen.go`）提供：`Insert`、`FindOne`、`FindOneByUserIdCourseId`、`Update`、`Delete`。

---

## 业务服务层（`internal/service/learning.go`）

learning 是本仓库中**唯一**在 model 之上再抽一层 `service` 的服务。`LearningService` 接口声明 11 个方法，`servicecontext.go:24` 注入：

| Service 方法 | 委托的 Model 方法 | 附加的校验 |
|-------------|------------------|-----------|
| `GrantCourses` | `GrantCourses` | `userID <= 0 \|\| len(courseIDs) == 0` → `xerr.BadRequestf` |
| `RevokeCourses` | `RevokeCourses` | 同上 |
| `RemoveLesson` | `RemoveLesson` | `userID <= 0 \|\| courseID <= 0` → `xerr.BadRequestf` |
| `GetLesson` | `FindByUserCourse` | `sql.ErrNoRows` → `xerr.NotFound("学习记录不存在")` |
| `ListLessons` | `ListByUser` | `userID <= 0` → `xerr.BadRequestf` |
| `CountLessons` | `CountByCourse` | `courseID <= 0` → `xerr.BadRequestf` |
| `CreatePlan` | `UpdatePlan` | 三参数均须 > 0；`sql.ErrNoRows` → `xerr.NotFound` |
| `ValidateLesson` | `GetLesson` | `status == LessonStatusExpired` → `xerr.Conflict("课程已失效")` |
| `CurrentLesson` | `FindLatestLearnedByUser` | `userID <= 0` → `xerr.BadRequestf` |
| `CommitRecord` | `UpdateLatestLearn` | `lessonID <= 0 \|\| sectionID <= 0` → `xerr.BadRequestf` |
| `ListLessonPlans` | `ListPlansByUser` | `userID <= 0` → `xerr.BadRequestf` |

> **接线缺口**：`ServiceContext.LearningService` 已构造完成，但 11 个 RPC logic 文件**尚未调用它**，全部仍是 goctl 占位。详见 `business-rules.md`。
