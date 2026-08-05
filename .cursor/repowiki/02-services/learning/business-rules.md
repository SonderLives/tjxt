> 版本：v1.1 | 更新：2026-08-06 | 来源：`apps/learning/rpc/internal/logic/*.go`, `apps/learning/api/internal/logic/*.go`, `apps/learning/rpc/internal/service/learning.go`

---

# Learning Business Rules

## ✅ 实现状态（已落地）

> **本服务 20/20 logic 已全部实现并通过编译**（`go build ./...` rc=0）。
>
> 架构为三层：`logic` → `service.LearningService`（已实现，带校验与 `xerr` 映射）→ `model`（已实现，含扩展 SQL）。`logic` 层是薄封装：RPC 层负责从 JWT 取 `user_id`、调用 Service、把 `model.LearningLesson` 转成 `pb` 视图；API 层负责调用 `LearningRpc`、用 `CourseRpc.CourseSimpleInfoList` 回填课程维度字段、再转成对外 `types`。
>
> 本文档的所有规则均为**已落地实现**，标注「⚠️ 缺口」的为真实存在的能力限制（非设计待办）。

### RPC 层（`apps/learning/rpc/internal/logic/`）

| 业务分组 | Logic 方法 | 文件 | 实现状态 |
|---------|-----------|------|---------|
| 课表查询 | `LearningNow` | `learningnowlogic.go` | ✅ 已实现 |
| 课表查询 | `LessonPage` | `lessonpagelogic.go` | ✅ 已实现 |
| 课表查询 | `LessonGet` | `lessongetlogic.go` | ✅ 已实现 |
| 课表查询 | `LessonValid` | `lessonvalidlogic.go` | ✅ 已实现 |
| 课表查询 | `LessonCount` | `lessoncountlogic.go` | ✅ 已实现 |
| 学习计划 | `PlanPage` | `planpagelogic.go` | ✅ 已实现 |
| 学习计划 | `PlanSave` | `plansavelogic.go` | ✅ 已实现 |
| 学习记录 | `LearningRecordCommit` | `learningrecordcommitlogic.go` | ✅ 已实现 |
| 学习记录 | `LearningRecordsByCourse` | `learningrecordsbycourselogic.go` | ✅ 已实现 |
| 内部事件 | `GrantCourses` | `grantcourseslogic.go` | ✅ 已实现 |
| 内部事件 | `RevokeCourses` | `revokecourseslogic.go` | ✅ 已实现 |

**RPC 已实现 11 / 总计 11。**

### API 层（`apps/learning/api/internal/logic/`）

| 业务分组 | Logic 方法 | 文件 | 实现状态 |
|---------|-----------|------|---------|
| 课表查询 | `LearningNow` | `learningnowlogic.go` | ✅ 已实现 |
| 课表查询 | `LessonPage` | `lessonpagelogic.go` | ✅ 已实现 |
| 课表查询 | `LessonGet` | `lessongetlogic.go` | ✅ 已实现 |
| 课表查询 | `LessonValid` | `lessonvalidlogic.go` | ✅ 已实现 |
| 课表查询 | `LessonCount` | `lessoncountlogic.go` | ✅ 已实现 |
| 学习计划 | `PlanPage` | `planpagelogic.go` | ✅ 已实现 |
| 学习计划 | `PlanSave` | `plansavelogic.go` | ✅ 已实现 |
| 学习记录 | `LearningRecordCommit` | `learningrecordcommitlogic.go` | ✅ 已实现 |
| 学习记录 | `LearningRecordsByCourse` | `learningrecordsbycourselogic.go` | ✅ 已实现 |

**API 已实现 9 / 总计 9。**

### 下层实现情况（对照）

| 层 | 文件 | 已实现 | 说明 |
|----|------|-------|------|
| Logic | `internal/logic/*` | 20 / 20 | ✅ 全部实现（本次新增） |
| Service | `internal/service/learning.go` | 11 / 11 | ✅ 全部实现，含参数校验与 `xerr` 错误映射 |
| Model | `internal/model/learninglessonmodel.go` | 11 / 11 | ✅ 全部实现，含 SQL |

---

## 1. 参数校验（Service 层 + Logic 层）

✅ **已实现** — 来源 `internal/service/learning.go` 与 `internal/logic/*`

**Logic 层身份提取**：用户相关接口通过 `pkg/auth.UserIdFromCtx(l.ctx)` 从 JWT 取 `user_id`；内部事件接口 `GrantCourses` / `RevokeCourses` 的 `user_id` 取自请求体（来自 trade 内部 RPC，非 JWT）。`LessonCount` 为公开统计，无需登录。

**Service 层校验**（委托 Model 前先做）：

| 规则 | 触发条件 | 返回 |
|------|---------|------|
| 用户与课程 ID 校验 | `userID <= 0 \|\| len(courseIDs) == 0` | `xerr.BadRequestf("userID/courseIDs 非法")` |
| 单课程校验 | `userID <= 0 \|\| courseID <= 0` | `xerr.BadRequestf("userID/courseID 非法")` |
| 用户 ID 校验 | `userID <= 0` | `xerr.BadRequestf("userID 非法")` |
| 课程 ID 校验 | `courseID <= 0` | `xerr.BadRequestf("courseID 非法")` |
| 计划参数校验 | `userID <= 0 \|\| courseID <= 0 \|\| weekFreq <= 0` | `xerr.BadRequestf("userID/courseID/weekFreq 非法")` |
| 记录提交校验 | `lessonID <= 0 \|\| sectionID <= 0` | `xerr.BadRequestf("lessonID/sectionID 非法")` |

**错误映射**：

| 底层错误 | 映射结果 |
|---------|---------|
| `sql.ErrNoRows`（GetLesson） | `xerr.NotFound("学习记录不存在")` |
| `sql.ErrNoRows`（CreatePlan） | `xerr.NotFound("学习记录不存在")` |
| 其他 DB 错误（GetLesson） | `xerr.Wrapf(err, xerr.CodeInternal, "查询学习记录失败")` |
| 其他 DB 错误（CreatePlan） | `xerr.Wrapf(err, xerr.CodeInternal, "更新学习计划失败")` |

---

## 2. 课程开通与撤销

✅ **已实现（Model + Service + Logic）**

**核心规则**：课程开通幂等，撤销为逻辑失效（不删行）。

| 规则 | 依据 | 说明 |
|------|------|------|
| 开通幂等 | `insertCourseSQL` 的 `ON DUPLICATE KEY UPDATE update_time = NOW()` | 依赖 `uk_user_course` 唯一索引，重复开通只刷新 `update_time` |
| 初始状态 | `insertCourseSQL` 的 VALUES 常量 | `status=0`（未开始）、`week_freq=NULL`、`plan_status=0`、`learned_sections=0`、`latest_section_id=NULL`、`latest_learn_time=NULL`、`expire_time=NULL` |
| 主键生成 | `idgen.NextID()` | 雪花算法，非数据库自增 |
| 逐条插入 | `GrantCourses` 的 `for` 循环 | 每个 `courseID` 单独执行一次 SQL，**非批量 INSERT**，任一条失败立即 return |
| 撤销即失效 | `RevokeCourses` SQL | `status = 3`（失效）、`plan_status = 0`（无计划），不物理删除 |
| 删除亦失效 | `RemoveLesson` SQL | 同样只把 `status` 置 3（该 Service/Model 方法已存在，但 proto/`.api` 未暴露出口，见末节缺口） |

⚠️ **缺口**：两个内部方法的触发源是 trade 通过 RabbitMQ 发布的 `order.pay` / `order.refund` 事件。`learning.yaml` 已配置 `PayQueue` / `RefundQueue`，但**消费端未接线**（详见末节「已知缺口」），目前需手动/其他途径触发开通。

---

## 3. 课表查询与分页

✅ **已实现（Model + Service + Logic）**

**核心规则**：所有列表查询排除已失效课程（`status <> 3`）。

| 规则 | 依据 | 说明 |
|------|------|------|
| 排除失效 | `ListByUser` 的 `status <> LessonStatusExpired` | 课表不展示已退款/已撤销的课程 |
| 计划过滤 | `ListPlansByUser` 追加 `plan_status = PlanStatusInPlan` | 只看已设置计划的课程 |
| 分页归一 | `listBy` 私有方法 | `pageNo < 1 → 1`；`pageSize < 1 → 10`；`pageSize > 100 → 100` |
| 排序 | `listBy` 的 `orderBy` 参数 | 两个列表方法均固定传 `create_time`；`asc=true → ASC`，否则 `DESC` |
| 总数与列表两次查询 | `listBy` 先 `SELECT COUNT(1)` 再 `SELECT ... LIMIT ? OFFSET ?` | 非窗口函数方案 |
| 绕过缓存 | `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx` | 列表查询不走 goctl 缓存 |

**Logic 层职责（本次落地）**：

| 项 | 实现 |
|----|------|
| userId 提取 | `auth.UserIdFromCtx(l.ctx)` 取当前学员 ID |
| 状态枚举映射 | `toLessonVO`：`status` 0→`NOT_BEGIN` / 1→`LEARNING` / 2→`FINISHED` / 3→`EXPIRED`；`plan_status` 0→`NO_PLAN` / 1→`PLAN_RUNNING` |
| pages 计算 | `calcPages(total, pageSize)`，RPC 层 `LessonPageReply.pages` / `PlanPageReply.pages` 由此换算（`pageSize<=0` 时返回 0） |
| 课程侧字段补全 | **API 层** `enrichLessons` 调 `CourseRpc.CourseSimpleInfoList`，回填 `course_name` / `course_cover_url` / `course_amount` / `sections`；RPC 层不持有课程信息，这些字段留空 |
| 时间格式化 | `nullTime` / `timeLayout`（`2006-01-02 15:04:05`）格式化 `create_time` / `expire_time` / `latest_learn_time`；`NullTime` 无效时返回空串 |

---

## 4. 「我正在学」

✅ **已实现（Model + Service + Logic）**

**核心规则**：取该用户 `latest_learn_time` 最新的一条未失效记录（来自 `FindLatestLearnedByUser`）。

⚠️ **注意**：`CurrentLesson` / `FindLatestLearnedByUser` **未做 `sql.ErrNoRows → xerr.NotFound` 映射**（与 `GetLesson` 不一致）。`LearningNow` 直接透传底层错误：无记录时返回 `sql.ErrNoRows` 经 go-zero 框架转 500。如需前端友好，应在此处补 `errors.Is(err, sql.ErrNoRows)` 判断。

---

## 5. 报名校验

✅ **已实现（Service + Logic）**

`ValidateLesson` 复用 `GetLesson`，在其之上追加失效判定：

```
流程（LessonValid / ValidateLesson）:
  1. GetLesson(userID, courseID)
       → 未找到：xerr.NotFound("学习记录不存在")
       → DB 错误：xerr.Wrapf(CodeInternal, "查询学习记录失败")
  2. if lesson.Status == LessonStatusExpired
       → xerr.Conflict("课程已失效")
  3. return lesson.Id
```

Logic 层 `LessonValid` 把返回的 `lesson_id` 填入 `LessonValidReply.LessonId`；未报名时服务层抛 `xerr.NotFound`，由 API 框架渲染为 404。

---

## 6. 学习计划

✅ **已实现（Model + Service + Logic）** ⚠️ 周维度统计部分无数据源

**核心规则**：学习计划不单独建表，直接更新 `learning_lesson.week_freq` + `plan_status`。

| 规则 | 依据 | 说明 |
|------|------|------|
| 计划落在 lesson 表 | `UpdatePlan` 的 UPDATE 语句 | 设置 `week_freq = ?`、`plan_status = 1` |
| weekFreq 必须 > 0 | `CreatePlan` 校验 | `weekFreq <= 0` → `xerr.BadRequestf` |
| 可空写入 | `sql.NullInt64{Int64: weekFreq, Valid: weekFreq > 0}` | `weekFreq <= 0` 时写 NULL |
| 只能给未失效课程设计划 | UPDATE 的 `AND status <> LessonStatusExpired` | - |
| 记录不存在判定 | `RowsAffected() == 0` → `sql.ErrNoRows` | 由 Service 转 `xerr.NotFound("学习记录不存在")` |
| 撤销课程清计划 | `RevokeCourses` 同时置 `plan_status = 0` | - |

**Logic 层（PlanPage）周维度字段**：

| 字段 | 实现现状 |
|------|---------|
| `week_finished` | 无数据源（learning_lesson.learned_sections 是累计值，无法按周切分）→ 填 0 |
| `week_points` | 无数据源（learning 库无积分表）→ 填 0 |
| `week_total_plan` | 由**当页**列表的 `week_freq` 求和得到（`int64` 返回） |

---

## 7. 学习记录提交

✅ **已实现（Model + Service + Logic）** ⚠️ 存在参数丢弃与能力缺口（源自 Model/Service 层，logic 无法补救）

**核心规则**：提交进度只更新 `learning_lesson` 单行，不落记录明细。

| 规则 | 依据 | 说明 |
|------|------|------|
| 状态自动跃迁 | `UpdateLatestLearn` 的 `status = IF(status = 0, 1, status)` | 首次提交进度时从「未开始」自动变为「学习中」，已是 1/2/3 则保持不变 |
| 更新最近小节 | `latest_section_id = ?`（`sql.NullInt64`，`sectionID > 0` 才有效） | - |
| 更新学习时间 | `latest_learn_time = time.Now()` | 用 Go 侧时间，**不用请求传入的 `commit_time`** |
| 章节计数独立 | `IncrLearnedSections` 单独提供 | `learned_sections + 1`，`CommitRecord` **未调用它** |
| 按 lessonID 定位 | `WHERE id = ?` | 不校验该 lesson 是否属于当前 userID |

**Logic 层落地细节**：

| 项 | 实现 |
|----|------|
| 身份 | API 层 `LearningRecordCommit` 调 `auth.UserIdFromCtx` 做鉴权；RPC 层 `CommitRecord` 再从 JWT 取 `userID` 透传给 Service |
| section_type 映射 | API 层把请求字符串 `VIDEO/EXAM` 映射为 `int32`（VIDEO→1，EXAM→2，默认 1）传入 `LearningRecordCommitRequest.SectionType`；Service 层忽略该字段 |

⚠️ **已知问题（真实代码行为，非本层引入）**：

| 问题 | 位置 | 说明 |
|------|------|------|
| `moment` / `duration` 被丢弃 | `learninglessonmodel.go:UpdateLatestLearn` | 签名接收但未在 SQL 中使用 |
| `commitTime` 被丢弃 | `service/learning.go:CommitRecord` | 签名接收但方法体未使用 |
| `userID` 被丢弃 | `service/learning.go:CommitRecord` | 不做 lesson 归属校验，存在越权提交风险 |
| `section_type` 无处落库 | proto `LearningRecordCommitRequest.section_type` | VIDEO / EXAM 无对应字段 |
| `records` 明细无表 | `LearningRecordsByCourse` | `learning_lesson` 单行只能存一个 `latest_section_id`，无法还原多小节明细。`records` 列表当前**恒为空**，仅返回 `id` 与 `latest_section_id` |

---

## 状态说明

### 课程学习状态（`learning_lesson.status`）

| 值 | Go 常量 | proto 字符串 | 含义 | 跃迁触发 |
|----|---------|-------------|------|---------|
| 0 | `LessonStatusNotStart` | `NOT_BEGIN` | 未开始 | `GrantCourses` 开通时的初始值 |
| 1 | `LessonStatusInLearn` | `LEARNING` | 学习中 | `UpdateLatestLearn` 首次提交进度时由 0 自动跃迁 |
| 2 | `LessonStatusDone` | `FINISHED` | 完成 | ⚠️ 无代码路径写入此值 |
| 3 | `LessonStatusExpired` | `EXPIRED` | 失效 | `RevokeCourses` / `RemoveLesson` |

### 学习计划状态（`learning_lesson.plan_status`）

| 值 | Go 常量 | proto 字符串 | 含义 | 跃迁触发 |
|----|---------|-------------|------|---------|
| 0 | `PlanStatusNone` | `NO_PLAN` | 无计划 | 开通时初始值；`RevokeCourses` 时重置 |
| 1 | `PlanStatusInPlan` | `PLAN_RUNNING` | 计划中 | `UpdatePlan` |

### 小节类型（`LearningRecordCommitRequest.section_type`）

| 值 | 含义 | 备注 |
|----|------|------|
| 1 / 2 | VIDEO / EXAM | API 层字符串映射；Service 层忽略 |

---

## 已知缺口汇总（真实存在，非待办）

| 缺口 | 影响 | 是否可在 logic 层修复 |
|------|------|----------------------|
| **MQ 消费端未接线** | `learning.yaml` 已配置 `PayQueue` / `RefundQueue` / Exchange / RoutingKey，但 `servicecontext.go` 与 `learning.go` 均无 Consumer 初始化。trade 发布的 `order.pay` / `order.refund` 事件当前无人消费，课程无法自动开通 | 否（需在 svcCtx 装配 MQ Consumer） |
| **`learning_record` 表缺失** | `LearningRecordsByCourse.records` 恒空；`moment` / `duration` / `section_type` / `commit_time` 无处落库 | 否（需建表 + Model 扩展） |
| **`status=2`(完成) 无写入路径** | Model 层无方法把 `status` 置 2，「学完」状态无法产生 | 否（需 Service/Model 补逻辑） |
| **`RemoveLesson` 无 RPC 出口** | Service 与 Model 已实现 `RemoveLesson`，但 `learning.proto` 与 `learning.api` 均未声明对应方法/路由（`api-spec.md` 却列了 `DELETE /lessons/{courseId}`） | 否（需补 proto/api） |
| **`CommitRecord` 无归属校验** | Service 层接收 `userID` 但不校验 lesson 归属，存在越权提交他人 lesson 进度的风险 | 否（需 Service 层补校验） |
| **`CurrentLesson` 错误映射不一致** | 未把 `sql.ErrNoRows` 转 `xerr.NotFound`，无记录时返回 500 | 是（可在 `LearningNow` logic 补判断） |
| **`week_finished` / `week_points` 无数据源** | 恒为 0 | 否（需周统计表/积分表） |
| **课程小节名/序号无数据源** | `CourseSectionGet` 仅返回 `media_id`，`latest_section_name` / `latest_section_index` 恒空 | 否（需 course 服务补充小节元数据接口） |
