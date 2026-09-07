# learning（学习进度）Specification

## Purpose

learning（学习进度）微服务覆盖我的课表 / 学习计划 / 学习记录（课表查询、报名校验、学习计划、学习进度提交），HTTP 服务 `learning-api` 监听 8804，RPC 服务 `learning.rpc` 监听 8084（gRPC），经 etcd（`127.0.0.1:2379`，key `learning.rpc`）注册与发现；数据库为 MySQL `tj_learning` 库（仅 `learning_lesson` 一张表）。架构为三层 `logic` → `service.LearningService`（本仓库唯一在 model 之上再抽 service 层的服务，承载参数校验与 `xerr` 错误映射）→ `model`；RPC 11 个 + API 9 个 logic 共 20/20 已实现并通过编译。作为交易链路的终点不配置任何下游 RpcClient：只被 learning-api 自身调用，并消费 trade 经 RabbitMQ `order.exchange` 发布的 `order.pay` / `order.refund` 事件开通 / 撤销课程（消费端队列与交换机配置就绪，Consumer 未接线）；课程侧字段（课程名 / 封面 / 章节数等）由 API 层经 `course.rpc` 补全。

## Requirements

### Requirement: HTTP 课表接口

learning-api SHALL 提供以下课表接口（均经 `LearningRpc` 转发）：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /lessons/now | 查询我正在学习的课程，响应 `R{data: R«LearningLessonVO»}` |
| GET | /lessons/page | 分页查询我的课表，响应 `R{data: R«PageDTO«LearningLessonVO»»}` |
| GET | /lessons/{courseId} | 查询指定课程信息，响应 `R{data: R«LearningLessonVO»}` |
| DELETE | /lessons/{courseId} | 删除指定课程信息，响应 `R{data: R}` |
| GET | /lessons/{courseId}/count | 统计课程学习人数，响应 `R{data: R«int»}` |
| GET | /lessons/{courseId}/valid | 校验当前课程是否已经报名，响应 `R{data: R«long»}` |

补充约束：

- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页约定 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码遵循共享错误码规范
- 按 configs.md 记载，learning-api 全部路由声明 `@server (jwt: Auth)`（`Auth.AccessExpire` 7200 秒），`Auth.AccessSecret` 须与 auth.rpc 签发密钥一致（聚合文档 `docs/tjxt.openapi.json` 将上述端点认证列标注为「否」，以配置记载为准）
- `LessonCount` 为公开统计，无需登录；其余接口经 `pkg/auth.UserIdFromCtx(l.ctx)` 从 JWT 取 `user_id`

> 已知缺口：`DELETE /lessons/{courseId}` 在聚合文档中列为「删除指定课程信息」，Service 与 Model 已实现 `RemoveLesson`（置 status=3），但 `learning.proto` 未定义对应 RPC 方法、`learning.api` 也未声明该 handler——能力已在下层就绪，上层未打通（需补 proto/api）。

#### Scenario: 分页查询我的课表

- **WHEN** 学员请求 GET /lessons/page
- **THEN** 返回当前用户的课表分页（`LearningLessonVO` 列表），不展示已失效（status=3）课程
- **AND** 课程侧字段 `course_name` / `course_cover_url` / `course_amount` / `sections` 由 `CourseRpc` 补全

#### Scenario: 报名校验拦截播放

- **WHEN** 播放页请求 GET /lessons/{courseId}/valid
- **THEN** 已报名返回 `lesson_id`；未报名时服务层抛 `xerr.NotFound`，由 API 框架渲染为 404
- **AND** 课程已失效（status=EXPIRED）时返回 Conflict（「课程已失效」）

#### Scenario: 统计课程学习人数

- **WHEN** 请求 GET /lessons/{courseId}/count
- **THEN** 返回该课程的报名人数，统计排除已失效课程（status <> 3），无需登录

### Requirement: HTTP 学习记录与学习计划接口

learning-api SHALL 提供以下学习记录与学习计划接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /learning-records | 提交学习记录，请求体 `LearningRecordFormDTO`，响应 `R{data: R}` |
| GET | /learning-records/course/{courseId} | 查询指定课程的学习记录，响应 `R{data: R«LearningLessonDTO»}` |
| GET | /lessons/plans | 查询我的学习计划，响应 `R{data: R«LearningPlanPageVO»}` |
| POST | /lessons/plans | 创建学习计划，请求体 `LearningPlanDTO`，响应 `R{data: R}` |

补充约束：

- business-rules.md（v1.1）记载 API 层 logic 已实现 9 / 9（课表查询 5 + 学习计划 2 + 学习记录 2），RPC 11 + API 9 共 20/20 logic 全部实现并通过编译（`go build ./...` rc=0）

#### Scenario: 提交学习记录

- **WHEN** 播放器定时携带 `LearningRecordFormDTO` 请求 POST /learning-records
- **THEN** API 层把 section_type 字符串映射为 int32（VIDEO→1，EXAM→2，默认 1），并经 `LearningRpc.LearningRecordCommit` 更新学习进度
- **AND** 首次提交后 `learning_lesson.status` 由 0（未开始）自动跃迁为 1（学习中）

#### Scenario: 创建学习计划

- **WHEN** 学员携带 `LearningPlanDTO` 请求 POST /lessons/plans
- **THEN** 每周学习章节数 `week_freq` 落库 `learning_lesson` 并置 `plan_status=1`（计划中）
- **AND** weekFreq <= 0 时返回 BadRequest（「userID/courseID/weekFreq 非法」）

### Requirement: HTTP 笔记接口

learning-api SHALL 按聚合文档提供以下笔记接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /notes | 新增笔记，请求体 `NoteFormDTO`，响应 `R{data: R}` |
| PUT | /notes/{id} | 更新笔记，请求体 `NoteFormDTO`，响应 `R{data: R}` |
| DELETE | /notes/{id} | 删除我的笔记，响应 `R{data: R}` |
| POST | /notes/gathers/{id} | 采集笔记，响应 `R{data: R}` |
| DELETE | /notes/gathers/{id} | 取消采集笔记，响应 `R{data: R}` |
| GET | /notes/page | 用户端分页查询笔记，响应 `R{data: R«PageDTO«NoteVO»»}` |

> 状态注记：business-rules.md 记载 API 层仅实现课表 / 学习计划 / 学习记录 3 域共 9 个 logic；笔记端点仅见于聚合文档 `docs/tjxt.openapi.json`，无对应 logic 实现记载（以下两个 Requirement 同此状态）。

#### Scenario: 新增与更新笔记

- **WHEN** 学员携带 `NoteFormDTO` 请求 POST /notes（或 PUT /notes/{id}）
- **THEN** 按聚合文档摘要新增 / 更新笔记并返回 `R{data: R}`

#### Scenario: 采集与取消采集笔记

- **WHEN** 学员请求 POST /notes/gathers/{id}（或 DELETE /notes/gathers/{id}）
- **THEN** 按聚合文档摘要采集 / 取消采集该笔记并返回 `R{data: R}`

### Requirement: HTTP 回答与评论接口

learning-api SHALL 按聚合文档提供以下回答与评论接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /replies | 新增回答或评论，请求体 `ReplyDTO`，响应 `R{data: R}` |
| GET | /replies/page | 分页查询回答或评论，响应 `R{data: R«PageDTO«ReplyVO»»}` |

> 状态注记：同「HTTP 笔记接口」，该域端点仅见于聚合文档，无 logic 实现记载。

#### Scenario: 新增回答或评论

- **WHEN** 学员携带 `ReplyDTO` 请求 POST /replies
- **THEN** 按聚合文档摘要新增回答或评论并返回 `R{data: R}`

#### Scenario: 分页查询回答或评论

- **WHEN** 学员请求 GET /replies/page
- **THEN** 按聚合文档摘要分页返回回答或评论列表（`PageDTO«ReplyVO»`）

### Requirement: HTTP 积分与签到接口

learning-api SHALL 按聚合文档提供以下积分与签到接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /points/today | 查询我的今日积分，响应 `R{data: R«List«PointsStatisticsVO»»}` |
| GET | /sign-records | 查询签到记录，响应 `R{data: R«Array«byte»»}` |
| POST | /sign-records | 签到功能接口，响应 `R{data: R«SignResultVO»}` |
| GET | /boards | 分页查询指定赛季的积分排行榜，响应 `R{data: R«PointsBoardVO»}` |
| GET | /boards/seasons/list | 查询赛季信息列表，响应 `R{data: R«List«PointsBoardSeasonVO»»}` |

> 状态注记：同「HTTP 笔记接口」，该域端点仅见于聚合文档，无 logic 实现记载。另 data-model.md 记载 `tj_learning` 库无积分表（`PlanPageReply.week_points` 无数据源），积分 / 签到 / 排行榜域在本服务无存储与配置支撑记载。

#### Scenario: 签到与签到记录

- **WHEN** 学员请求 POST /sign-records
- **THEN** 按聚合文档摘要完成签到并返回 `SignResultVO`
- **AND** GET /sign-records 返回该用户的签到记录（byte 数组形式）

#### Scenario: 查询今日积分与排行榜

- **WHEN** 学员请求 GET /points/today 或 GET /boards
- **THEN** 按聚合文档摘要返回今日积分统计列表（`PointsStatisticsVO`）/ 指定赛季积分排行榜分页（`PointsBoardVO`）
- **AND** GET /boards/seasons/list 返回赛季信息列表

### Requirement: RPC 服务契约

`Learning` 服务（学习进度微服务）SHALL 经 etcd 服务发现（key `learning.rpc`，监听 `0.0.0.0:8084`）对外提供 11 个 RPC 方法，按 4 个业务域分组：

| 业务域 | 方法 | 说明 |
|--------|------|------|
| 课表查询 | LearningNow、LessonPage、LessonGet、LessonValid、LessonCount | 当前正在学课程（Empty → LearningLessonVO）；我的课表分页（LessonPageRequest{page_no, page_size, is_asc, sort_by} → LessonPageReply{total, pages, list}）；指定课程学习信息；报名校验；课程学习人数 |
| 学习计划 | PlanPage、PlanSave | 我的学习计划分页（PlanPageReply）；创建/更新学习计划（PlanSaveRequest{course_id, freq}，freq 落库 `learning_lesson.week_freq`） |
| 学习记录 | LearningRecordCommit、LearningRecordsByCourse | 提交学习记录（LearningRecordCommitRequest{lesson_id, section_id, section_type, moment, duration, commit_time}）；查询某课程的学习记录（LearningRecordsReply{id, latest_section_id, records}） |
| 内部事件（来自 trade） | GrantCourses、RevokeCourses | 开通 / 撤销课程（GrantCoursesRequest{user_id, course_ids}）；既可被 MQ 消费端调用，也可作为内部 RPC 暴露 |

关键字段约定：

- `LearningLessonVO`：id（lesson_id）、course_id、course_name / course_cover_url / course_amount / sections（课程侧，来自 course 服务）、learned_sections、status（NOT_BEGIN / LEARNING / FINISHED / EXPIRED）、plan_status（NO_PLAN / PLAN_RUNNING）、week_freq、latest_section_id / latest_section_name / latest_section_index、create_time / expire_time / latest_learn_time
- `LessonValidReply`：返回 `lesson_id`；未报名返回 `0` / `NotFound`，由调用方判定
- `PlanPageReply`：total、pages、week_finished、week_points、week_total_plan、list（repeated LearningLessonVO）
- 设计要点（proto 头部注释）：学习记录不再单独建表，通过更新 `learning_lesson` 的 `latest_section_id` / `latest_learn_time` / `learned_sections` 实现

消费关系：当前仅 learning-api 自身消费该 RPC（HTTP Handler → `learningclient.Learning`，`servicecontext.go:17` 注入）；learning-api 另注入 `CourseRpc courseclient.Course`（servicecontext.go:25）补全课程侧字段。trade 服务不直连 learning RPC，而是通过 RabbitMQ `order.exchange` 发布 `order.pay` / `order.refund` 事件，由 learning 侧消费后调用 `GrantCourses` / `RevokeCourses`。

> 注（事件命名口径）：本规格按 learning 服务配置记载记录实际拓扑（Exchange `order.exchange`，RoutingKey `order.pay` / `order.refund`，队列 `learning.lesson.pay.queue` / `learning.lesson.refund.queue`）；共享事件契约规格（specs/shared/mq-events/spec.md）按规划口径记载为 `trade.events` / `order.created` / `order.paid`（解锁课程 / 更新学习权限 / 回收权益），两处口径不同，以本服务实际配置记载为准。

#### Scenario: 课程开通链路

- **WHEN** 学员在 trade 支付成功、trade 发布 `order.pay` 事件
- **THEN** learning 消费后调用 `GrantCourses(user_id, course_ids)`，幂等写入 `learning_lesson`
- **AND** 退款成功后 trade 发布 `order.refund`，learning 消费后调用 `RevokeCourses` 将 status 置 3（失效）、plan_status 置 0

#### Scenario: 继续学习

- **WHEN** 前端调用 LearningNow
- **THEN** 返回该用户最近学习的一条课程，跳转到 `latest_section_id` 对应小节

#### Scenario: 提交进度更新课表

- **WHEN** 播放器定时调用 LearningRecordCommit
- **THEN** 更新 `latest_section_id` / `latest_learn_time`，status 由 0（未开始）自动跃迁为 1（学习中）

### Requirement: 参数校验与错误映射

learning SHALL 在 Service 层（`internal/service/learning.go`）委托 Model 前完成参数校验，并做 `xerr` 错误映射；Logic 层负责身份提取：

- Logic 层身份提取：用户相关接口经 `pkg/auth.UserIdFromCtx(l.ctx)` 从 JWT 取 `user_id`；内部事件接口 `GrantCourses` / `RevokeCourses` 的 `user_id` 取自请求体（来自 trade 内部 RPC，非 JWT）；`LessonCount` 为公开统计，无需登录
- Service 层参数校验：

| 规则 | 触发条件 | 返回 |
|------|---------|------|
| 用户与课程 ID 校验 | `userID <= 0` 或 `len(courseIDs) == 0` | `xerr.BadRequestf("userID/courseIDs 非法")` |
| 单课程校验 | `userID <= 0` 或 `courseID <= 0` | `xerr.BadRequestf("userID/courseID 非法")` |
| 用户 ID 校验 | `userID <= 0` | `xerr.BadRequestf("userID 非法")` |
| 课程 ID 校验 | `courseID <= 0` | `xerr.BadRequestf("courseID 非法")` |
| 计划参数校验 | `userID <= 0` 或 `courseID <= 0` 或 `weekFreq <= 0` | `xerr.BadRequestf("userID/courseID/weekFreq 非法")` |
| 记录提交校验 | `lessonID <= 0` 或 `sectionID <= 0` | `xerr.BadRequestf("lessonID/sectionID 非法")` |

- 错误映射：

| 底层错误 | 映射结果 |
|---------|---------|
| `sql.ErrNoRows`（GetLesson） | `xerr.NotFound("学习记录不存在")` |
| `sql.ErrNoRows`（CreatePlan） | `xerr.NotFound("学习记录不存在")` |
| 其他 DB 错误（GetLesson） | `xerr.Wrapf(err, xerr.CodeInternal, "查询学习记录失败")` |
| 其他 DB 错误（CreatePlan） | `xerr.Wrapf(err, xerr.CodeInternal, "更新学习计划失败")` |

#### Scenario: Service 层参数校验

- **WHEN** Service 收到 `userID <= 0` / `courseID <= 0` / `weekFreq <= 0` / `lessonID <= 0` / `sectionID <= 0` 等非法参数
- **THEN** 返回对应的 `xerr.BadRequestf` 错误，不落库

#### Scenario: 学习记录不存在映射

- **WHEN** GetLesson / CreatePlan 底层返回 `sql.ErrNoRows`
- **THEN** 映射为 `xerr.NotFound("学习记录不存在")`
- **AND** 其他 DB 错误映射为 `xerr.Wrapf(CodeInternal)`（「查询学习记录失败」/「更新学习计划失败」）

### Requirement: 课程开通与撤销规则

learning SHALL 以幂等方式开通课程、以逻辑失效方式撤销课程（不删行），并按以下拓扑消费 trade 的交易事件：

- 开通幂等：`GrantCourses` 使用 `INSERT ... ON DUPLICATE KEY UPDATE update_time = NOW()`，依赖 `uk_user_course` 唯一索引，重复开通只刷新 `update_time`
- 初始状态：status=0（未开始）、week_freq=NULL、plan_status=0、learned_sections=0、latest_section_id=NULL、latest_learn_time=NULL、expire_time=NULL
- 主键生成：`idgen.NextID()` 雪花算法（`tjxt/pkg/utils/idgen`），非数据库自增
- 逐条插入：`GrantCourses` 对每个 courseID 单独执行一次 SQL（非批量 INSERT），任一条失败立即 return
- 撤销即失效：`RevokeCourses` 置 status=3（失效）、plan_status=0（无计划），不物理删除
- 删除亦失效：`RemoveLesson` 同样只把 status 置 3（该 Service/Model 方法已存在，但 proto/.api 未暴露出口）
- MQ 消费拓扑（trade 生产、learning 消费，两端交换机与路由键必须一致）：PayExchange / RefundExchange 均为 `order.exchange`；PayRoutingKey `order.pay` → PayQueue `learning.lesson.pay.queue`；RefundRoutingKey `order.refund` → RefundQueue `learning.lesson.refund.queue`（生产端 trade 不声明队列）

> 已知缺口：`learning.yaml` 已完整配置 PayQueue / RefundQueue / Exchange / RoutingKey 共 6 项 RabbitMQ 参数，`config.go` 也声明了对应结构体，但 `servicecontext.go` 与 `learning.go` 主入口均无 Consumer 初始化 / 消费协程——trade 发布的 `order.pay` / `order.refund` 事件当前无人消费，课程无法自动开通，目前需手动 / 其他途径触发开通（需在 svcCtx 装配 MQ Consumer）；且无 trade 那样的 MQ 初始化失败降级设计（对照 trade `servicecontext.go:42-46`：仅记日志、Producer 置 nil、不阻塞启动）。

#### Scenario: 幂等开通

- **WHEN** 对同一用户同一课程重复调用 GrantCourses
- **THEN** 依赖 `uk_user_course` 唯一索引命中 ON DUPLICATE KEY，只刷新 `update_time`，不重复插入
- **AND** 新开通记录初始 status=0（未开始）、week_freq=NULL、plan_status=0、learned_sections=0

#### Scenario: 撤销与删除均为逻辑失效

- **WHEN** 调用 RevokeCourses（或 RemoveLesson）
- **THEN** 仅将 `learning_lesson.status` 置 3（失效），不物理删除
- **AND** RevokeCourses 同时把 plan_status 置 0（清除学习计划）

#### Scenario: 事件到方法的映射（设计意图）

- **WHEN** `learning.lesson.pay.queue` 收到 `order.pay` 事件
- **THEN** 触发 `GrantCourses(user_id, course_ids)`，幂等写入 `learning_lesson`（status=0）
- **AND** `learning.lesson.refund.queue` 收到 `order.refund` 事件时触发 `RevokeCourses`（status=3、plan_status=0）

#### Scenario: MQ 消费端接线现状

- **WHEN** 审查 learning 消费端接线
- **THEN** learning.yaml 已完整配置 6 项 RabbitMQ 参数、config.go 已声明结构体，但 servicecontext.go 无 Consumer 创建代码、主入口未启动消费协程
- **AND** trade 发布的 `order.pay` / `order.refund` 当前无人消费，课程无法自动开通

### Requirement: 课表查询与分页规则

课表与计划列表查询 SHALL 统一排除已失效课程，并遵循以下分页与装配规则：

- 排除失效：`ListByUser` 带 `status <> LessonStatusExpired`，课表不展示已退款 / 已撤销的课程
- 计划过滤：`ListPlansByUser` 追加 `plan_status = PlanStatusInPlan`，只看已设置计划的课程
- 分页归一（私有方法 `listBy`）：`pageNo < 1 → 1`；`pageSize < 1 → 10`；`pageSize > 100 → 100`
- 排序：两个列表方法固定按 `create_time`，`asc=true → ASC`，否则 `DESC`
- 总数与列表两次查询：先 `SELECT COUNT(1)` 再 `SELECT ... LIMIT ? OFFSET ?`（非窗口函数方案）
- 绕过缓存：`QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx`，列表查询不走 goctl 缓存
- Logic 层职责：`auth.UserIdFromCtx` 取当前学员 ID；`toLessonVO` 做状态枚举映射（0→NOT_BEGIN / 1→LEARNING / 2→FINISHED / 3→EXPIRED；plan_status 0→NO_PLAN / 1→PLAN_RUNNING）；`calcPages(total, pageSize)` 换算 pages（pageSize<=0 时返回 0）；时间按 `2006-01-02 15:04:05` 格式化，NullTime 无效时返回空串
- 课程侧字段补全：API 层 `enrichLessons` 调 `CourseRpc.CourseSimpleInfoList` 回填 `course_name` / `course_cover_url` / `course_amount` / `sections`；RPC 层不持有课程信息，这些字段留空

> 已知缺口：课程总章节数 `sections` 与最近学习小节名称 / 序号（`latest_section_name` / `latest_section_index`）不在 `learning_lesson` 表，需调 course.rpc 补全；`CourseSectionGet` 仅返回 `media_id`，`latest_section_name` / `latest_section_index` 恒空（需 course 服务补充小节元数据接口）。

#### Scenario: 分页参数归一

- **WHEN** 列表查询 pageNo=0、pageSize=200
- **THEN** 归一为 pageNo=1、pageSize=100（pageNo<1→1，pageSize<1→10，>100→100）

#### Scenario: 失效课程不出现在课表

- **WHEN** 学员查询我的课表（ListByUser / ListPlansByUser）
- **THEN** 查询条件带 status <> 3，已退款 / 已撤销课程不展示
- **AND** ListPlansByUser 仅返回 plan_status=1（计划中）的课程，固定按 create_time 排序

#### Scenario: 课程侧字段补全

- **WHEN** API 层组装 LearningLessonVO
- **THEN** enrichLessons 调 `CourseRpc.CourseSimpleInfoList` 回填 course_name / course_cover_url / course_amount / sections
- **AND** latest_section_name / latest_section_index 无数据源，恒为空串

### Requirement: 「我正在学」规则

「我正在学」SHALL 返回该用户 `latest_learn_time` 最新的一条未失效记录（`FindLatestLearnedByUser`，带 `status <> LessonStatusExpired` 过滤）。

> 已知缺口：`CurrentLesson` / `FindLatestLearnedByUser` 未做 `sql.ErrNoRows → xerr.NotFound` 映射（与 GetLesson 不一致），无记录时 `LearningNow` 直接透传底层错误，`sql.ErrNoRows` 经 go-zero 框架转 500；如需前端友好应补 `errors.Is(err, sql.ErrNoRows)` 判断（可在 LearningNow logic 修复）。

#### Scenario: 返回最近学习的课程

- **WHEN** 调用 CurrentLesson / LearningNow
- **THEN** 返回该用户 latest_learn_time 最新的一条未失效记录

#### Scenario: 无学习记录时的错误行为

- **WHEN** 用户没有任何学习记录
- **THEN** 当前实现透传 `sql.ErrNoRows`，经 go-zero 框架渲染为 500（未映射为 NotFound）

### Requirement: 报名校验规则

报名校验 SHALL 复用 GetLesson 并在其上追加失效判定，返回 lesson_id 供播放页拦截。流程（LessonValid / ValidateLesson）：

1. `GetLesson(userID, courseID)` → 未找到：`xerr.NotFound("学习记录不存在")`；DB 错误：`xerr.Wrapf(CodeInternal, "查询学习记录失败")`
2. `lesson.Status == LessonStatusExpired` → `xerr.Conflict("课程已失效")`
3. 通过后返回 `lesson.Id`

Logic 层把返回的 `lesson_id` 填入 `LessonValidReply.LessonId`；未报名时服务层抛 `xerr.NotFound`，由 API 框架渲染为 404。

#### Scenario: 已报名通过校验

- **WHEN** 调用 LessonValid(course_id) 且该用户已报名且课程未失效
- **THEN** 返回 lesson_id 填入 LessonValidReply.LessonId

#### Scenario: 未报名或课程失效拦截

- **WHEN** 用户未报名（GetLesson 未找到）
- **THEN** 服务层抛 `xerr.NotFound("学习记录不存在")`，由 API 框架渲染为 404
- **AND** 课程已失效（status=EXPIRED）时抛 `xerr.Conflict("课程已失效")`

### Requirement: 学习计划规则

学习计划 SHALL 不单独建表，直接更新 `learning_lesson.week_freq` + `plan_status`：

- 计划落在 lesson 表：`UpdatePlan` 的 UPDATE 语句设置 `week_freq = ?`、`plan_status = 1`
- weekFreq 必须 > 0：`CreatePlan` 校验 `weekFreq <= 0` → `xerr.BadRequestf`
- 可空写入：`sql.NullInt64{Int64: weekFreq, Valid: weekFreq > 0}`，weekFreq <= 0 时写 NULL
- 只能给未失效课程设计划：UPDATE 带 `AND status <> LessonStatusExpired`
- 记录不存在判定：`RowsAffected() == 0` → `sql.ErrNoRows`，由 Service 转 `xerr.NotFound("学习记录不存在")`
- 撤销课程清计划：`RevokeCourses` 同时置 `plan_status = 0`
- PlanPage 周维度字段：`week_total_plan` 由当页列表的 `week_freq` 求和得到（int64 返回）

> 已知缺口：`week_finished` 无数据源（`learned_sections` 是累计值，无法按周切分）恒填 0；`week_points` 无数据源（learning 库无积分表）恒填 0（需独立的计划 / 周统计表与积分表，`learning_plan` 表 DDL 未定义）。

#### Scenario: 创建或更新学习计划

- **WHEN** 调用 PlanSave(course_id, freq)
- **THEN** UPDATE `learning_lesson` 置 week_freq=freq、plan_status=1，仅对 status <> 3（未失效）的记录生效
- **AND** 记录不存在（RowsAffected==0）时返回 `xerr.NotFound("学习记录不存在")`

#### Scenario: 周维度统计现状

- **WHEN** 调用 PlanPage
- **THEN** week_total_plan 为当页列表 week_freq 求和
- **AND** week_finished / week_points 无数据源，恒为 0

### Requirement: 学习记录提交规则

学习进度提交 SHALL 只更新 `learning_lesson` 单行，不落记录明细：

- 状态自动跃迁：`UpdateLatestLearn` 用 `status = IF(status = 0, 1, status)`，首次提交进度时从「未开始」自动变为「学习中」，已是 1/2/3 则保持不变
- 更新最近小节：`latest_section_id = ?`（`sql.NullInt64`，sectionID > 0 才有效）
- 更新学习时间：`latest_learn_time = time.Now()`，用 Go 侧时间，不用请求传入的 `commit_time`
- 章节计数独立：`IncrLearnedSections` 单独提供（`learned_sections + 1`，用于学习记录被确认完成时），`CommitRecord` 未调用它
- 按 lessonID 定位：`WHERE id = ?`，不校验该 lesson 是否属于当前 userID
- section_type 映射：API 层把请求字符串 VIDEO/EXAM 映射为 int32（VIDEO→1，EXAM→2，默认 1）传入 `LearningRecordCommitRequest.SectionType`；Service 层忽略该字段

> 已知问题（真实代码行为，非本层引入）：`moment` / `duration` 在 `UpdateLatestLearn` 签名中接收但 SQL 未使用（实际被丢弃）；`commitTime` 在 Service `CommitRecord` 中接收但方法体未使用；`userID` 接收但不做 lesson 归属校验，存在越权提交他人 lesson 进度的风险；`section_type` 无处落库；`LearningRecordsByCourse` 的 `records` 明细恒为空（`learning_lesson` 单行只能存一个 `latest_section_id`，无法还原多小节明细），仅返回 `id` 与 `latest_section_id`。

#### Scenario: 提交进度状态自动跃迁

- **WHEN** 对 status=0（未开始）的课程提交学习记录
- **THEN** status 自动跃迁为 1（学习中），并更新 latest_section_id 与 latest_learn_time
- **AND** 已是 1/2/3 的记录状态保持不变

#### Scenario: 参数丢弃与归属校验现状

- **WHEN** 提交 LearningRecordCommit
- **THEN** 仅更新 latest_section_id 与 latest_learn_time（服务端 time.Now()），moment / duration / commit_time / section_type 不落库
- **AND** WHERE id = ? 不校验 lesson 归属当前 userID，存在越权提交风险

#### Scenario: 课程学习记录明细恒空

- **WHEN** 调用 LearningRecordsByCourse
- **THEN** records 列表恒为空，仅返回 id 与 latest_section_id
- **AND** 多小节 moment/duration/finished 明细需 `learning_record` 表支撑（当前缺失）

### Requirement: 状态枚举与跃迁

`learning_lesson` 的状态字段 SHALL 以 tinyint 存储、proto 层以字符串枚举对外，跃迁触发如下：

| 值 | Go 常量 | proto 字符串 | 含义 | 跃迁触发 |
|----|---------|-------------|------|---------|
| 0 | `LessonStatusNotStart` | `NOT_BEGIN` | 未开始 | `GrantCourses` 开通时的初始值 |
| 1 | `LessonStatusInLearn` | `LEARNING` | 学习中 | `UpdateLatestLearn` 首次提交进度时由 0 自动跃迁 |
| 2 | `LessonStatusDone` | `FINISHED` | 完成 | 无代码路径写入此值 |
| 3 | `LessonStatusExpired` | `EXPIRED` | 失效 | `RevokeCourses` / `RemoveLesson` |

| 值 | Go 常量 | proto 字符串 | 含义 | 跃迁触发 |
|----|---------|-------------|------|---------|
| 0 | `PlanStatusNone` | `NO_PLAN` | 无计划 | 开通时初始值；`RevokeCourses` 时重置 |
| 1 | `PlanStatusInPlan` | `PLAN_RUNNING` | 计划中 | `UpdatePlan` |

| 值 | 含义 | 备注 |
|----|------|------|
| 1 / 2 | VIDEO / EXAM | `LearningRecordCommitRequest.section_type`；API 层字符串映射，Service 层忽略 |

> 注（源冲突）：rpc-spec.md 与 data-model.md（v1.0，2026-08-05）记载「`ServiceContext.LearningService` 已构造完成，但 11 个 RPC logic 文件尚未调用它，全部仍是 goctl 占位」「数值 ↔ 字符串的转换需在 RPC logic 层完成，当前 logic 为占位实现，映射尚未落地」；business-rules.md（v1.1，2026-08-06，更新时间更晚）记载 20/20 logic 已实现并通过编译，含 `toLessonVO` 枚举映射。本规格按 business-rules.md 记录已实现行为，两处表述出入如实保留。

#### Scenario: 状态值与字符串映射

- **WHEN** RPC 层组装 LearningLessonVO
- **THEN** status 按数值映射字符串（0→NOT_BEGIN、1→LEARNING、2→FINISHED、3→EXPIRED），plan_status（0→NO_PLAN、1→PLAN_RUNNING）

#### Scenario: FINISHED 无写入路径

- **WHEN** 审查 status=2（完成）的写入路径
- **THEN** Model 层无方法把 status 置 2，「学完」状态无法产生（需 Service / Model 补逻辑）

### Requirement: 数据模型与不变量

learning 数据库 `tj_learning` SHALL 仅包含 `learning_lesson`（学员课程表）一张表，引擎 InnoDB、字符集 utf8mb4，DDL 使用 `CREATE TABLE IF NOT EXISTS`（区别于 trade 域的 `DROP TABLE IF EXISTS` + `CREATE TABLE`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 雪花主键（`idgen.NextID()`，非数据库自增） |
| user_id / course_id | BIGINT | 学员 ID / 课程 ID（跨库引用 user 域 user.id、course 域 course.id，无外键） |
| status | TINYINT | 0未开始，1学习中，2完成，3失效，默认 0 |
| week_freq | INT | 每周学习章节数，可空 |
| plan_status | TINYINT | 0无计划，1计划中，默认 0 |
| learned_sections | INT | 已学习章节数，默认 0 |
| latest_section_id | BIGINT | 最近学习小节，可空（跨库引用 course 域 section.id，无外键） |
| latest_learn_time | DATETIME | 最近学习时间，可空 |
| create_time | DATETIME | 创建时间，默认 CURRENT_TIMESTAMP |
| expire_time | DATETIME | 课程失效时间，可空 |
| update_time | DATETIME | 更新时间，ON UPDATE CURRENT_TIMESTAMP |

索引与不变量：

| 索引名 | 类型 | 字段 | 用途 |
|--------|------|------|------|
| PRIMARY | 主键 | id | 主键查询 |
| uk_user_course | 唯一 | user_id, course_id | 保证一人一课一条记录；`GrantCourses` 的 `ON DUPLICATE KEY UPDATE` 幂等依赖此约束；goctl 据此生成 `FindOneByUserIdCourseId` |
| idx_course_status | 普通 | course_id, status | 支撑 `CountByCourse`（统计课程学习人数） |
| idx_user_status_created | 普通 | user_id, status, create_time | 支撑 `ListByUser` / `ListPlansByUser`（按用户分页 + ORDER BY create_time） |

模型扩展模式：`learninglessonmodel_gen.go` 由 goctl 生成禁止修改（Insert / FindOne / FindOneByUserIdCourseId / Update / Delete），扩展方法统一放在同名自定义 `learninglessonmodel.go`（手写 11 个业务方法：GrantCourses、RevokeCourses、FindByUserCourse、ListByUser、ListPlansByUser、FindLatestLearnedByUser、CountByCourse、UpdatePlan、RemoveLesson、UpdateLatestLearn、IncrLearnedSections，含状态常量）；`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`。缓存策略：11 个自定义扩展方法全部使用 `ExecNoCacheCtx` / `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx` 绕过 goctl 缓存直接操作 DB，缓存实际只在 gen 层基础方法生效。事件来源：trade → RabbitMQ(order.exchange) → learning → GrantCourses / RevokeCourses。

> 已知缺口（缺失表）：`sql/ddl/tj_learning.sql` 只定义 `learning_lesson` 一张表，`learning_plan`（`PlanSave` / `PlanPage` 需要按周切分的完成量与积分）与 `learning_record`（`LearningRecordCommit` / `LearningRecordsByCourse` 的逐小节进度明细）两表 DDL 未定义、model 未生成。缺口影响面：`LearningRecordsByCourse` 无法返回 records 数组（只能返回 id + latest_section_id）；`LearningRecordCommit` 的 section_type / moment / duration / commit_time 四个入参无处落库（`UpdateLatestLearn` 的 moment / duration 形参实际被丢弃）；`PlanPage` 的 week_finished / week_points / week_total_plan 无数据源；`LearningLessonVO.sections` 与 `latest_section_name` / `latest_section_index` 需调 course.rpc 补全。

#### Scenario: 一人一课一条记录

- **WHEN** 同一用户对同一课程开通（GrantCourses）
- **THEN** `uk_user_course` 唯一索引保证 `learning_lesson` 单行，重复开通仅刷新 update_time

#### Scenario: 扩展方法绕过缓存

- **WHEN** 调用 11 个自定义 Model 扩展方法
- **THEN** 全部使用 NoCacheCtx 系列方法直接操作 DB
- **AND** goctl 缓存仅在 gen 层 Insert / FindOne / FindOneByUserIdCourseId / Update / Delete 生效

#### Scenario: 缺失表影响

- **WHEN** 依赖 `learning_plan` / `learning_record` 的能力被调用
- **THEN** 周维度统计与逐小节明细无数据源（PlanPage 周字段、records 数组）
- **AND** 学习记录提交的明细参数（moment / duration / section_type / commit_time）无处落库

### Requirement: 服务配置

learning SHALL 按以下配置部署运行：

| 配置项 | 默认值 |
|--------|--------|
| learning-api Name / Host / Port | learning-api / 0.0.0.0 / 8804（HTTP） |
| learning.rpc Name / ListenOn | learning.rpc / 0.0.0.0:8084（gRPC） |
| etcd | 127.0.0.1:2379，服务注册 / 发现 key `learning.rpc` |
| DataSource | root:0000@tcp(127.0.0.1:3306)/tj_learning?charset=utf8mb4&parseTime=true&loc=Local |
| Redis Cache | 127.0.0.1:6379，node 单机模式，密码为空 |
| RabbitMQ | 127.0.0.1:5672（rabbitmq/rabbitmq）；PayQueue `learning.lesson.pay.queue`、RefundQueue `learning.lesson.refund.queue`；PayExchange / RefundExchange `order.exchange`；PayRoutingKey `order.pay`、RefundRoutingKey `order.refund` |
| LearningRpc（API 层） | etcd 127.0.0.1:2379，key `learning.rpc`（HTTP handler → RPC client 调用） |
| CourseRpc（API 层） | etcd 127.0.0.1:2379，key `course.rpc`（课程目录展示时，需要查询课程名 / 章节数等基础信息） |
| Auth | AccessSecret 默认 `change-me-in-production`，AccessExpire 7200 秒 |

约束规则：

- learning-api 全部路由声明 `@server (jwt: Auth)`，`Auth.AccessSecret` 必须与签发方 auth.rpc 的 `Jwt.AccessSecret` 一致（learning-api 只做 token 校验不签发）；生产环境必须修改默认值，否则 JWT 可被伪造
- learning RPC 不配置任何下游 RpcClient：它是交易链路的终点，只被 learning-api 调用以及被 MQ 事件驱动
- Telemetry / Prometheus / Log 三段可观测性注入配置已注入，勿删（统一约定见 `specs/infra/observability/spec.md`）

> 已知配置缺口：MQ 消费端未接线（6 项 RabbitMQ 参数与 config.go 结构体均已声明，但无 Consumer 创建代码，处于「已声明未使用」状态）；无 MQ 容错策略（对照 trade 的降级设计）；RPC 层无 CourseRpc（课程侧字段补全只在 API 层配置，若需 RPC 层直接返回完整 LearningLessonVO 需在 learning.yaml 与 config.go 补 `CourseRpc zrpc.RpcClientConf`）；无积分服务配置（`PlanPageReply.week_points` 需要积分数据，learning.yaml 与 learning-api.yaml 均无相关服务连接配置）。

#### Scenario: JWT 密钥一致性

- **WHEN** learning-api 校验用户 token
- **THEN** `Auth.AccessSecret` 必须与 auth.rpc 签发密钥一致，否则 token 校验失败
- **AND** 生产环境必须修改默认值 `change-me-in-production`，否则 JWT 可被伪造

#### Scenario: 端口与服务发现

- **WHEN** 部署 learning 服务
- **THEN** learning-api 监听 8804（HTTP）、learning.rpc 监听 8084（gRPC）
- **AND** 经 etcd key `learning.rpc` 注册，调用方通过该 key 服务发现

#### Scenario: 交易链路终点

- **WHEN** 审查 learning RPC 的下游依赖配置
- **THEN** learning.yaml 不配置任何下游 RpcClient
- **AND** learning.rpc 只被 learning-api 调用，以及被 MQ 事件驱动（消费端当前未接线）

