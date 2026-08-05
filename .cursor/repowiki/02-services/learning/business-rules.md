> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/learning/rpc/internal/logic/*.go`, `apps/learning/api/internal/logic/*.go`, `apps/learning/rpc/internal/service/learning.go`

---

# Learning Business Rules

## ⚠️ 实现状态

> **本服务 logic 层业务逻辑尚未实现。** 全部 20 个 logic 文件仍为 goctl 生成的占位实现（函数体只有 `// todo: add your logic here and delete this line` 与零值返回）。
>
> 但与 trade 不同，learning 的**下层已经就绪**：`internal/model/learninglessonmodel.go` 实现了 11 个扩展方法，`internal/service/learning.go` 实现了 11 个带校验的业务方法，且 `LearningService` 已在 `servicecontext.go:24` 完成注入。**缺的只是 logic 层的接线与 DTO 转换。**
>
> 本文档中标注「📋 设计意图（待实现）」的规则，是依据 proto 注释、DDL 表结构、`docs/tjxt.openapi.json` 契约推导出的待实现规则；标注「✅ 已实现」的规则来自 model / service 层的真实代码。

### RPC 层（`apps/learning/rpc/internal/logic/`）

| 业务分组 | Logic 方法 | 文件 | 实现状态 |
|---------|-----------|------|---------|
| 课表查询 | `LearningNow` | `learningnowlogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonPage` | `lessonpagelogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonGet` | `lessongetlogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonValid` | `lessonvalidlogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonCount` | `lessoncountlogic.go` | 未实现-goctl占位 |
| 学习计划 | `PlanPage` | `planpagelogic.go` | 未实现-goctl占位 |
| 学习计划 | `PlanSave` | `plansavelogic.go` | 未实现-goctl占位 |
| 学习记录 | `LearningRecordCommit` | `learningrecordcommitlogic.go` | 未实现-goctl占位 |
| 学习记录 | `LearningRecordsByCourse` | `learningrecordsbycourselogic.go` | 未实现-goctl占位 |
| 内部事件 | `GrantCourses` | `grantcourseslogic.go` | 未实现-goctl占位 |
| 内部事件 | `RevokeCourses` | `revokecourseslogic.go` | 未实现-goctl占位 |

**RPC 已实现 0 / 总计 11。**

### API 层（`apps/learning/api/internal/logic/`）

| 业务分组 | Logic 方法 | 文件 | 实现状态 |
|---------|-----------|------|---------|
| 课表查询 | `LearningNow` | `learningnowlogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonPage` | `lessonpagelogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonGet` | `lessongetlogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonValid` | `lessonvalidlogic.go` | 未实现-goctl占位 |
| 课表查询 | `LessonCount` | `lessoncountlogic.go` | 未实现-goctl占位 |
| 学习计划 | `PlanPage` | `planpagelogic.go` | 未实现-goctl占位 |
| 学习计划 | `PlanSave` | `plansavelogic.go` | 未实现-goctl占位 |
| 学习记录 | `LearningRecordCommit` | `learningrecordcommitlogic.go` | 未实现-goctl占位 |
| 学习记录 | `LearningRecordsByCourse` | `learningrecordsbycourselogic.go` | 未实现-goctl占位 |

**API 已实现 0 / 总计 9。**

### 下层实现情况（对照）

| 层 | 文件 | 已实现 | 说明 |
|----|------|-------|------|
| Service | `internal/service/learning.go` | 11 / 11 | ✅ 全部实现，含参数校验与 `xerr` 错误映射 |
| Model | `internal/model/learninglessonmodel.go` | 11 / 11 | ✅ 全部实现，含 SQL |

### 汇总

| 层 | 已实现 | 总计 | 比例 |
|----|-------|------|------|
| RPC logic | 0 | 11 | 0% |
| API logic | 0 | 9 | 0% |
| **logic 合计** | **0** | **20** | **0%** |
| Service 层 | 11 | 11 | 100% |
| Model 扩展 | 11 | 11 | 100% |

---

## 1. 参数校验（Service 层）

✅ **已实现** — 来源 `internal/service/learning.go`

**核心规则**：所有 Service 方法在委托 Model 之前先做参数校验，非法参数统一返回 `xerr.BadRequestf`。

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

✅ **已实现（Model + Service）** / 📋 **logic 接线待实现**

**核心规则**：课程开通幂等，撤销为逻辑失效（不删行）。

| 规则 | 依据 | 说明 |
|------|------|------|
| 开通幂等 | `insertCourseSQL` 的 `ON DUPLICATE KEY UPDATE update_time = NOW()` | 依赖 `uk_user_course` 唯一索引，重复开通只刷新 `update_time` |
| 初始状态 | `insertCourseSQL` 的 VALUES 常量 | `status=0`（未开始）、`week_freq=NULL`、`plan_status=0`、`learned_sections=0`、`latest_section_id=NULL`、`latest_learn_time=NULL`、`expire_time=NULL` |
| 主键生成 | `idgen.NextID()` | 雪花算法，非数据库自增 |
| 逐条插入 | `GrantCourses` 的 `for` 循环 | 每个 `courseID` 单独执行一次 SQL，**非批量 INSERT**，任一条失败立即 return |
| 撤销即失效 | `RevokeCourses` SQL | `status = 3`（失效）、`plan_status = 0`（无计划），不物理删除 |
| 删除亦失效 | `RemoveLesson` SQL | 同样只把 `status` 置 3 |

```
流程（GrantCourses）— 已实现:
  Service: 校验 userID > 0 且 len(courseIDs) > 0
  Model:   for cid := range courseIDs
             INSERT INTO learning_lesson
               (id, user_id, course_id, status, week_freq, plan_status,
                learned_sections, latest_section_id, latest_learn_time,
                create_time, expire_time, update_time)
             VALUES (idgen.NextID(), userID, cid, 0, NULL, 0, 0, NULL, NULL,
                     NOW(), NULL, NOW())
             ON DUPLICATE KEY UPDATE update_time = NOW()
```

📋 **设计意图（待实现）**：这两个方法的触发源是 trade 通过 RabbitMQ `order.exchange` 发布的 `order.pay` / `order.refund` 事件。当前 `learning.yaml` 已配置 `PayQueue` / `RefundQueue`，但**消费端未接线**（详见本文末「已知缺口」）。

---

## 3. 课表查询与分页

✅ **已实现（Model + Service）** / 📋 **logic 接线待实现**

**核心规则**：所有列表查询排除已失效课程（`status <> 3`）。

| 规则 | 依据 | 说明 |
|------|------|------|
| 排除失效 | `ListByUser` 的 `status <> LessonStatusExpired` | 课表不展示已退款/已撤销的课程 |
| 计划过滤 | `ListPlansByUser` 追加 `plan_status = PlanStatusInPlan` | 只看已设置计划的课程 |
| 分页归一 | `listBy` 私有方法 | `pageNo < 1 → 1`；`pageSize < 1 → 10`；`pageSize > 100 → 100` |
| 排序 | `listBy` 的 `orderBy` 参数 | 两个列表方法均固定传 `create_time`；`asc=true → ASC`，否则 `DESC` |
| 总数与列表两次查询 | `listBy` 先 `SELECT COUNT(1)` 再 `SELECT ... LIMIT ? OFFSET ?` | 非窗口函数方案 |
| 绕过缓存 | `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx` | 列表查询不走 goctl 缓存 |

```
流程（ListByUser）— 已实现:
  1. SELECT COUNT(1) FROM learning_lesson WHERE user_id = ? AND status <> 3
  2. pageNo/pageSize 归一化，offset = (pageNo - 1) * pageSize
  3. SELECT <rows> FROM learning_lesson
       WHERE user_id = ? AND status <> 3
       ORDER BY create_time DESC|ASC
       LIMIT ? OFFSET ?
```

📋 **设计意图（待实现）** — logic 层职责：

| 待实现项 | 说明 |
|---------|------|
| userId 提取 | 从 JWT 上下文取当前学员 ID，Service 层不感知 context 中的身份 |
| 状态枚举映射 | DB 的 `status` int64 → proto 的 `status` string（0→`NOT_BEGIN` / 1→`LEARNING` / 2→`FINISHED` / 3→`EXPIRED`）；`plan_status`（0→`NO_PLAN` / 1→`PLAN_RUNNING`） |
| pages 计算 | `LessonPageReply.pages` 由 `total` 与 `pageSize` 换算，Model 只返回 `total` |
| 课程侧字段补全 | `course_name` / `course_cover_url` / `sections` / `latest_section_name` / `latest_section_index` 不在 `learning_lesson` 表中，须由 learning-api 调 `CourseRpc` 补全 |
| 时间格式化 | DB 的 `DATETIME` → proto 的 `string`（`create_time` / `expire_time` / `latest_learn_time`） |

---

## 4. 「我正在学」

✅ **已实现（Model + Service）** / 📋 **logic 接线待实现**

**核心规则**：取该用户 `latest_learn_time` 最新的一条未失效记录。

| 规则 | 依据 | 说明 |
|------|------|------|
| 必须学过 | `latest_learn_time IS NOT NULL` | 从未提交过进度的课程不会被选中 |
| 排除失效 | `status <> LessonStatusExpired` | - |
| 取最新一条 | `ORDER BY latest_learn_time DESC LIMIT 1` | - |
| 无记录透传 | `FindLatestLearnedByUser` 直接返回 `err` | `CurrentLesson` **未做 `sql.ErrNoRows` → `xerr.NotFound` 映射**（与 `GetLesson` 的处理不一致），logic 层需自行判定 |

---

## 5. 报名校验

✅ **已实现（Service）** / 📋 **logic 接线待实现**

**核心规则**：`ValidateLesson` 复用 `GetLesson`，在其之上追加失效判定。

```
流程（ValidateLesson）— 已实现:
  1. GetLesson(userID, courseID)
       → 未找到：xerr.NotFound("学习记录不存在")
       → DB 错误：xerr.Wrapf(CodeInternal, "查询学习记录失败")
  2. if lesson.Status == LessonStatusExpired
       → xerr.Conflict("课程已失效")
  3. return lesson.Id
```

📋 **设计意图（待实现）**：proto 对 `LessonValid` 的注释是「返回 `lesson_id`；未报名返回 0/NotFound **由调用方判**」。Service 层选择的是抛 `xerr.NotFound` / `xerr.Conflict`，logic 层需决定是透传错误还是降级为 `lesson_id = 0`，二者需与前端契约对齐。

---

## 6. 学习计划

✅ **已实现（Model + Service）** / 📋 **PlanPage 周维度统计无数据源**

**核心规则**：学习计划不单独建表，直接更新 `learning_lesson.week_freq` + `plan_status`。

| 规则 | 依据 | 说明 |
|------|------|------|
| 计划落在 lesson 表 | `UpdatePlan` 的 UPDATE 语句 | 设置 `week_freq = ?`、`plan_status = 1` |
| weekFreq 必须 > 0 | `CreatePlan` 校验 | `weekFreq <= 0` → `xerr.BadRequestf` |
| 可空写入 | `sql.NullInt64{Int64: weekFreq, Valid: weekFreq > 0}` | `weekFreq <= 0` 时写 NULL（但 Service 层已先行拦截） |
| 只能给未失效课程设计划 | UPDATE 的 `AND status <> LessonStatusExpired` | - |
| 记录不存在判定 | `RowsAffected() == 0` → `sql.ErrNoRows` | 由 Service 转 `xerr.NotFound("学习记录不存在")` |
| 撤销课程清计划 | `RevokeCourses` 同时置 `plan_status = 0` | - |

📋 **设计意图（待实现）** — `PlanPage` 的三个周维度字段：

| 字段 | 现状 |
|------|------|
| `week_finished` | 本周已完成章节数，**无数据源**。`learning_lesson.learned_sections` 是累计值，无法按周切分 |
| `week_points` | 本周积分，**无数据源**。learning 库无积分表 |
| `week_total_plan` | 本周计划总章节数，理论上可由计划中课程的 `week_freq` 求和得到 |

---

## 7. 学习记录提交

✅ **已实现（Model + Service）** / ⚠️ **存在参数丢弃与能力缺口**

**核心规则**：提交进度只更新 `learning_lesson` 单行，不落记录明细。

| 规则 | 依据 | 说明 |
|------|------|------|
| 状态自动跃迁 | `UpdateLatestLearn` 的 `status = IF(status = 0, 1, status)` | 首次提交进度时从「未开始」自动变为「学习中」，已是 1/2/3 则保持不变 |
| 更新最近小节 | `latest_section_id = ?`（`sql.NullInt64`，`sectionID > 0` 才有效） | - |
| 更新学习时间 | `latest_learn_time = time.Now()` | 用 Go 侧时间，**不用请求传入的 `commit_time`** |
| 章节计数独立 | `IncrLearnedSections` 单独提供 | `learned_sections + 1`，注释说明「用于学习记录被确认完成时」，`CommitRecord` **未调用它** |
| 按 lessonID 定位 | `WHERE id = ?` | 不校验该 lesson 是否属于当前 userID |

⚠️ **已知问题（真实代码行为，非设计意图）**：

| 问题 | 位置 | 说明 |
|------|------|------|
| `moment` / `duration` 被丢弃 | `learninglessonmodel.go:179-192` | `UpdateLatestLearn` 签名接收 `moment, duration int64`，但 SQL 语句中**未使用这两个参数** |
| `commitTime` 被丢弃 | `service/learning.go:117-122` | `CommitRecord` 签名接收 `commitTime string`，方法体中**未使用** |
| `userID` 被丢弃 | `service/learning.go:117-122` | `CommitRecord` 接收 `userID` 但未做归属校验，存在越权提交他人 lesson 进度的风险 |
| `section_type` 无处落库 | proto `LearningRecordCommitRequest.section_type` | VIDEO / EXAM 类型无对应字段 |

📋 **设计意图（待实现）** — `LearningRecordsByCourse`：

响应 `LearningRecordsReply { id, latest_section_id, records }` 中的 `records` 是逐小节的 `{ section_id, moment, duration, finished }` 列表。`learning_lesson` 单行只能保存一个 `latest_section_id`，**无法还原多小节明细**。该接口需要 `learning_record` 表支撑，当前 DDL 未定义。

---

## 状态说明

### 课程学习状态（`learning_lesson.status`）

| 值 | Go 常量 | proto 字符串 | 含义 | 跃迁触发 |
|----|---------|-------------|------|---------|
| 0 | `LessonStatusNotStart` | `NOT_BEGIN` | 未开始 | `GrantCourses` 开通时的初始值 |
| 1 | `LessonStatusInLearn` | `LEARNING` | 学习中 | `UpdateLatestLearn` 首次提交进度时由 0 自动跃迁 |
| 2 | `LessonStatusDone` | `FINISHED` | 完成 | 📋 待实现：无代码路径写入此值 |
| 3 | `LessonStatusExpired` | `EXPIRED` | 失效 | `RevokeCourses` / `RemoveLesson` |

### 学习计划状态（`learning_lesson.plan_status`）

| 值 | Go 常量 | proto 字符串 | 含义 | 跃迁触发 |
|----|---------|-------------|------|---------|
| 0 | `PlanStatusNone` | `NO_PLAN` | 无计划 | 开通时初始值；`RevokeCourses` 时重置 |
| 1 | `PlanStatusInPlan` | `PLAN_RUNNING` | 计划中 | `UpdatePlan` |

### 小节类型（`LearningRecordCommitRequest.section_type`）

| 值 | 含义 | 备注 |
|----|------|------|
| 0 / 1 / 2 | 由 api 解析：VIDEO / EXAM 字符串映射 | proto 注释原文；`learning.api` 中 `LearningRecordCommitReq.SectionType` 为 string（`VIDEO` / `EXAM`），映射规则待 logic 层确定 |

---

## 已知缺口汇总

| 缺口 | 影响 |
|------|------|
| **logic 层 0/20 接线** | Service + Model 已 100% 就绪，但 20 个 logic 文件全是 goctl 占位，服务实际不可用 |
| **MQ 消费端未接线** | `learning.yaml` 已配置 `PayQueue` / `RefundQueue` / `PayExchange` / `RefundExchange` / 两个 RoutingKey，`config.go` 也声明了 `RabbitMQ` 结构体，但 `servicecontext.go` 与 `learning.go` 中**均无 MQ Consumer 初始化代码**。trade 发布的 `order.pay` / `order.refund` 事件当前无人消费，课程无法自动开通 |
| **缺失 `learning_record` 表** | `LearningRecordsByCourse` 的 `records` 数组无法返回；`LearningRecordCommit` 的 `moment` / `duration` / `section_type` / `commit_time` 无处落库 |
| **缺失 `learning_plan` / 周统计表** | `PlanPageReply` 的 `week_finished` / `week_points` 无数据源 |
| **`status=2`(完成) 无写入路径** | Model 层没有任何方法把 `status` 置为 2，「学完」状态无法产生 |
| **`RemoveLesson` 无 RPC 出口** | `api-spec.md` 列出 `DELETE /lessons/{courseId}`，Service 与 Model 均已实现 `RemoveLesson`，但 `learning.proto` 与 `learning.api` 都未声明对应方法/路由 |
| **`CommitRecord` 无归属校验** | 接收 `userID` 但不校验 lesson 归属，实现 logic 时须补 |
| **`CurrentLesson` 错误映射不一致** | 未像 `GetLesson` 那样把 `sql.ErrNoRows` 转为 `xerr.NotFound` |
