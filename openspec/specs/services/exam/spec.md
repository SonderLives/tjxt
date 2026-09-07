# exam（考试题库）Specification

## Purpose

exam（考试题库）是考试与题库管理微服务，提供题目（题干/选项/答案/解析）的新增、更新、删除、分页检索，以及题目与业务对象（如课程小节）的关联管理。HTTP API 层 `exam-api` 监听 `0.0.0.0:8805`，RPC 服务 `Exam` 监听 `0.0.0.0:8085`（gRPC），经 etcd 以 key `exam.rpc` 注册与发现，API 层经 etcd 调用自身 RPC。数据存储于 MySQL `tj_exam` 库（`question` / `question_detail` / `question_biz` 三张业务表 + Seata `undo_log` 闲置表），主键查询由 Redis 缓存托管。当前仅 `exam-api` 自身消费该 RPC，course / learning 等跨服务消费方尚未接线；JWT 鉴权依赖与签发方 `auth.rpc` 的密钥配置保持一致。

## Requirements

### Requirement: 管理端内容管理接口（笔记 / 互动问题 / 回答评论）

> **状态说明**：下表来自原始 Java 版聚合文档 `docs/tjxt.openapi.json`（api-spec.md，最后同步 2026-08-05）。go-zero 侧 `apps/exam/api/exam.api` 仅实现 7 个题目/业务关联接口；`/admin/notes`、`/admin/questions`、`/admin/replies` 管理端接口在 go-zero 侧**未实现**。源表中「认证：否」为 Java 原始版标注，与 go-zero 侧 `jwt: Auth` 声明不一致，以 `.api` 为准（见「鉴权」Requirement）。源文件中该域无组卷/考试/阅卷分组接口，本规格按实际功能域分组。

exam 服务 SHALL 暴露以下管理端内容管理接口：

| 方法 | 路径 | 说明 | 响应体 |
|------|------|------|--------|
| GET | /admin/notes/page | 管理端分页查询笔记 | R{data: R«PageDTO«NoteAdminVO»»} |
| GET | /admin/notes/{id} | 管理端查询笔记详情 | R{data: R«NoteAdminDetailVO»} |
| PUT | /admin/notes/{id}/hidden/{hidden} | 隐藏指定笔记 | R{data: R} |
| GET | /admin/questions/page | 管理端分页查询互动问题 | R{data: R«PageDTO«QuestionAdminVO»»} |
| GET | /admin/questions/{id} | 管理端根据id查询互动问题 | R{data: R«QuestionAdminVO»} |
| PUT | /admin/questions/{id}/hidden/{hidden} | 隐藏或显示问题 | R{data: R} |
| GET | /admin/replies/page | 分页查询回答或评论 | R{data: R«PageDTO«ReplyVO»»} |
| GET | /admin/replies/{id} | 根据id查询回答或评论 | R{data: R«ReplyVO»} |
| PUT | /admin/replies/{id}/hidden/{hidden} | 隐藏或显示评论 | R{data: R} |

统一约定（引用全局规范）：响应格式为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页为 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码见共享错误码表（`specs/shared/error-codes`）。

#### Scenario: 管理端隐藏指定笔记

- **WHEN** 管理端调用 `PUT /admin/notes/{id}/hidden/{hidden}`
- **THEN** 系统将该笔记置为指定隐藏状态并返回 `R{data: R}`

#### Scenario: 管理端分页查询互动问题

- **WHEN** 管理端调用 `GET /admin/questions/page`
- **THEN** 系统按分页参数返回互动问题分页列表
- **AND** 响应体为 `R{data: R«PageDTO«QuestionAdminVO»»}`

### Requirement: 题库互动问题接口

> **状态说明**：同上，源自 Java 版聚合文档 `docs/tjxt.openapi.json`；go-zero 侧 `.api` 实现其中题目管理相关接口（对应 RPC `SaveQuestion` / `DeleteQuestion` / `GetQuestion` / `ListQuestions`）。

exam 服务 SHALL 暴露以下题库互动问题接口：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /questions | 新增互动问题 | QuestionFormDTO | R{data: R} |
| GET | /questions/checkName | 校验名称是否有效，存在则无效返回false，不存在返回true | - | R{data: R«boolean»} |
| GET | /questions/list | 查询题目列表 | - | R{data: R«List«QuestionDTO»»} |
| GET | /questions/listOfBiz | 查询业务关联的题目列表 | - | R{data: R«List«QuestionDTO»»} |
| GET | /questions/numOfTeacher | 查询老师出题数量 | - | R{data: R«Map«long,int»»} |
| GET | /questions/page | 分页查询互动问题 | - | R{data: R«PageDTO«QuestionVO»»} |
| GET | /questions/scores | 查询题目分值 | - | R{data: R«Map«long,int»»} |
| DELETE | /questions/{id} | 根据id删除当前用户问题 | - | R{data: R} |
| GET | /questions/{id} | 根据id查询互动问题 | - | R{data: R«QuestionVO»} |
| PUT | /questions/{id} | 修改提问 | QuestionFormDTO | R{data: R} |

#### Scenario: 新增互动问题

- **WHEN** 用户以 `QuestionFormDTO` 调用 `POST /questions`
- **THEN** 系统创建互动问题并返回 `R{data: R}`

#### Scenario: 删除当前用户问题

- **WHEN** 用户调用 `DELETE /questions/{id}`
- **THEN** 系统仅删除该当前用户的问题并返回 `R{data: R}`（权限校验见「题目删除与级联清理」）

### Requirement: 题目业务关联接口

> **状态说明**：同上，源自 Java 版聚合文档 `docs/tjxt.openapi.json`；go-zero 侧 `.api` 实现其中业务关联相关接口（对应 RPC `AddQuestionBiz` / `RemoveQuestionBiz` / `GetQuestionsByBiz`）。

exam 服务 SHALL 暴露以下题目业务关联接口：

| 方法 | 路径 | 说明 | 响应体 |
|------|------|------|--------|
| GET | /question-biz/biz/list | 批量查询与业务有关的题目id | R{data: R«List«QuestionBizDTO»»} |
| GET | /question-biz/biz/{id} | 查询与业务有关的题目id | R{data: R«List«QuestionBizDTO»»} |
| POST | /question-biz/list | 批量保存题目和业务关系（JSON 请求体） | R{data: R} |
| GET | /question-biz/scores | 查询业务下的题目分数和 | R{data: R«Map«long,int»»} |

#### Scenario: 批量保存题目和业务关系

- **WHEN** 客户端以 JSON 请求体调用 `POST /question-biz/list`
- **THEN** 系统批量保存题目与业务对象的关联并返回 `R{data: R}`

#### Scenario: 查询业务下的题目分数和

- **WHEN** 客户端调用 `GET /question-biz/scores`
- **THEN** 系统返回业务下的题目分数和 `R{data: R«Map«long,int»»}`

### Requirement: RPC 服务契约（Exam 服务）

`Exam` 服务 SHALL 作为考试与题库管理微服务经 etcd 服务发现注册（key: `exam.rpc`），监听 `0.0.0.0:8085`（gRPC）。通用消息定义如下：

| Message | 字段 | 用途 |
|---------|------|------|
| `Empty` | (无) | 删除类方法的响应 |
| `IdReq` | `id` (int64) | 单 ID 请求 |
| `IdReply` | `id` (int64) | 单 ID 响应 |
| `PageReq` | `pageNo` (int32), `pageSize` (int32) | 通用分页请求 |

> **预留消息**：`PageReq` 在 proto 中已声明，但 `service Exam` 的 7 个方法**均未引用**（分页参数被内联进 `QuestionListReq` / `QuestionBizListReq`），属于预留消息。
>
> **实现状态**：RPC 层 7/7、API 层 7/7 logic 已全部实现并编译通过（2026-08-06 复核）；以下各规则为依据 `apps/exam/rpc/exam.proto`、`sql/ddl/tj_exam*.sql`、`apps/exam/api/exam.api` 契约推导，建议对照源码最终确认。

消费方 SHALL 为 `exam-api` 自身 API 层（HTTP Handler → `examclient.Exam` RPC，7 个接口全部走自身 RPC）。

> **跨服务消费现状**：经检索全部 13 个服务的 `servicecontext.go`，**当前无任何其他服务 import `examclient`**。尚未接线的逻辑消费方：
>
> - `course` 服务：`apps/course/rpc/course.proto:48` 定义 `CourseSubjectsGet(IdRequest) returns (CataSubjectInfoList)`（章节/小节的题目）；`question_biz.biz_id` 注释明确写「例如小节id」；但 course 侧未 import `examclient`，`exam-api` 侧也未持有 `CourseRpc`，跨域调用未建立。
> - `learning` 服务：学员答题会驱动 `question.answer_times` / `correct_times` 累加；无任何接线证据，纯语义推导。

#### Scenario: API 层经 etcd 调用自身 RPC

- **WHEN** 客户端调用 `exam-api` 的任一 HTTP 接口
- **THEN** HTTP Handler 经 etcd 发现 `exam.rpc` 并通过 `examclient.Exam` 调用对应 RPC 方法

#### Scenario: 跨服务调用未建立

- **WHEN** course 或 learning 服务需要取用 exam 题目数据
- **THEN** 当前无法调用（无服务 import `examclient`，exam 亦无 `CourseRpc` 客户端配置），需先建立跨域 RPC 通道

### Requirement: 题目管理 RPC

`Exam` 服务 SHALL 提供以下题目管理 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveQuestion` | `QuestionSaveReq { id, name, type, cateId1, cateId2, cateId3, difficulty, score, options, answer, analysis }` | `IdReply { id }` | 新增/更新题目，`id` 为空表示新增 |
| `DeleteQuestion` | `IdReq { id }` | `Empty {}` | 删除题目 |
| `GetQuestion` | `IdReq { id }` | `QuestionVO` | 按 ID 查询题目（含详情） |
| `ListQuestions` | `QuestionListReq { pageNo, pageSize, name, type, cateId1, cateId2, difficulty }` | `QuestionListReply { total, list }` | 分页查询题目列表 |

关键字段约束：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 题目 ID，新增时省略 |
| `name` | string | 题干 |
| `type` | int32 | 题目类型：1-单选, 2-多选, 3-不定项选择, 4-判断, 5-主观 |
| `cateId1` | int64 | 1 级课程分类 ID |
| `cateId2` | int64 | 2 级课程分类 ID |
| `cateId3` | int64 | 3 级课程分类 ID |
| `difficulty` | int32 | 难易度：1-简单, 2-中等, 3-困难 |
| `score` | int32 | 分值 |
| `options` | string | 选择题选项，JSON 数组字符串 |
| `answer` | string | 正确答案 |
| `analysis` | string | 答案解析 |

- `QuestionSaveReq` 中**不含** `answerTimes` / `correctTimes` —— 这两个统计字段由系统维护，不允许客户端传入。
- `QuestionListReq` 支持 `cateId1` / `cateId2` 过滤，但**不支持 `cateId3`**（proto 中未定义该过滤字段）。
- `QuestionVO` 是 `question` 与 `question_detail` 两表的**聚合视图**（含只读统计 `answerTimes` / `correctTimes`、`createTime`），需在 logic 层做一对一 JOIN 组装。

#### Scenario: 教师出题

- **WHEN** 管理端提交题干/选项/答案并调用 `SaveQuestion`
- **THEN** 系统同时写 `question` 与 `question_detail` 两表并返回题目 ID

#### Scenario: 题库检索

- **WHEN** 管理端按分类/类型/难度筛选并调用 `ListQuestions`
- **THEN** 系统分页返回 `QuestionListReply{ total, list }`

### Requirement: 题目业务关联 RPC

`Exam` 服务 SHALL 提供以下题目业务关联 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `AddQuestionBiz` | `QuestionBizReq { bizId, questionId }` | `IdReply { id }` | 建立题目与业务对象的关联 |
| `RemoveQuestionBiz` | `QuestionBizReq { bizId, questionId }` | `Empty {}` | 解除题目与业务对象的关联 |
| `GetQuestionsByBiz` | `QuestionBizListReq { bizId, pageNo, pageSize }` | `QuestionListReply { total, list }` | 分页查询某业务下的题目列表 |

- `bizId` 为业务 ID（要关联问题的某业务 id，例如小节 id）；`questionId` 为问题 ID；`pageNo` / `pageSize` 为 int32 分页参数。
- `AddQuestionBiz` 返回 `IdReply { id }`，即 `question_biz` 关联表的自增主键（非题目 ID）；`RemoveQuestionBiz` 按 `(bizId, questionId)` 组合定位，与表上的唯一索引 `biz_id(biz_id, question_id)` 对应。

#### Scenario: 小节挂题

- **WHEN** 课程编辑时为某小节选题并调用 `AddQuestionBiz(bizId=小节id, questionId)`
- **THEN** 系统写入 `question_biz` 关联行并返回关联表主键 id

#### Scenario: 小节撤题

- **WHEN** 调用 `RemoveQuestionBiz(bizId, questionId)`
- **THEN** 系统仅删除关联行，题目本身保留在题库

#### Scenario: 学员练习取题

- **WHEN** 学员进入小节练习并调用 `GetQuestionsByBiz(bizId=小节id)`
- **THEN** 系统分页拉取并返回该小节的题目列表

### Requirement: 题目双表写入

一道题目 SHALL 跨 `question`（主体）与 `question_detail`（选项/答案/解析）两表存储，且共享同一主键：

- 共享主键：`question_detail.id = question.id`，详情表主键非自增。
- 必须同事务：主体与详情写入需原子，避免出现无详情的孤儿题目（一对一强关联）。
- 新增/更新分支：`QuestionSaveReq.id` 在 `.api` 中标注 `optional`，沿用 `id <= 0` 为新增的项目约定。
- 统计字段不可传入：`answer_times` / `correct_times` 不在 `QuestionSaveReq` 中，由系统维护。
- 无逻辑删除：`question` 表无 `deleted` 列，只能物理删除。

#### Scenario: 新增题目（SaveQuestion 流程）

- **WHEN** 客户端提交 `SaveQuestion` 且 `id <= 0`
- **THEN** 系统校验 name / answer / analysis 非空（`.api` 中三者均为必填）、`type ∈ {1..5}`、`difficulty ∈ {1,2,3}`、`score > 0`，且选择题（type ∈ {1,2,3}）的 options 为合法 JSON 数组
- **AND** 系统生成雪花 ID，在同一本地事务内插入 `question` 主体（creater/updater 取 JWT userId）并以同一 ID 插入 `question_detail`，返回 `IdReply{ id }`

#### Scenario: 更新题目

- **WHEN** 客户端提交 `SaveQuestion` 且 `id > 0`
- **THEN** 系统先 `FindOne` 校验题目存在，再对 `question` 与 `question_detail` 两表分别 Update

> **能力缺口**：三个 model 均为 goctl 空壳，未封装事务；`SaveQuestion` 的双表原子写入需先在自定义 model 或 logic 层引入 `sqlx.Transact`。

### Requirement: 题目类型与选项校验

题目类型 SHALL 决定 `options` 与 `answer` 的形态，枚举由 DDL 注释定义：

| type | 含义 | options 要求（推导） | answer 语义（依据 `tj_exam.sql` 注释） |
|------|------|---------------------|------------------------------------|
| 1 | 单选题 | 必填，JSON 数组 | 单个 1~10 的选项序号 |
| 2 | 多选题 | 必填，JSON 数组 | 多个序号，**逗号隔开** |
| 3 | 不定向选择题 | 必填，JSON 数组 | 单个或多个序号，逗号隔开 |
| 4 | 判断题 | 可空 | **1 代表正确，其他代表错误** |
| 5 | 主观题 | 可空 | 参考答案文本 |

补充约束：

- `answer` 列长度为 `varchar(40)`，多选题逗号分隔的序号串需在此长度内；主观题参考答案文本同样受 `varchar(40)` 长度限制。
- `options` 为 MySQL `json` 类型且可空，Go 侧映射为 `sql.NullString`，需手工序列化。

#### Scenario: 选择题保存校验

- **WHEN** 保存 type ∈ {1,2,3} 的选择题
- **THEN** options 必填且为合法 JSON 数组，answer 为 1~10 的选项序号（多个答案用逗号隔开，且不超出 `varchar(40)` 长度）

#### Scenario: 判断题答案语义

- **WHEN** 保存 type=4 的判断题
- **THEN** options 可空，answer 取值 1 代表正确、其他值代表错误

### Requirement: 题目查询与聚合

题目查询 SHALL 返回 `question` + `question_detail` 的聚合视图 `QuestionVO`：

#### Scenario: 查看题目详情（GetQuestion 流程）

- **WHEN** 调用 `GetQuestion` 且校验 id > 0
- **THEN** 系统 `QuestionModel.FindOne(id)` 查询主体，不存在则返回 NotFound
- **AND** 系统以同一 ID 查询 `question_detail`，合并为 `QuestionVO` 并将 createTime 格式化为字符串

#### Scenario: 分页检索题目（ListQuestions 流程）

- **WHEN** 调用 `ListQuestions`（pageNo / pageSize 在 `.api` 中为 optional，需兜底默认值）
- **THEN** 系统按 name 模糊、type / cateId1 / cateId2 / difficulty 精确动态拼条件，分页查 question 主体并查总数
- **AND** 系统批量取详情（FindDetailsByIds）避免 N+1，返回 `QuestionListReply{ total, list }`

补充约束（注意点）：

- 不支持按 `cateId3` 过滤：`QuestionListReq` 中无该字段，虽然表里有 `cate_id3` 列。
- 列表是否返回答案：proto 的 `QuestionVO` 含 `answer` / `analysis`，学员端场景应在上层裁剪，避免泄题。
- 无二级索引：`question` 表除主键外无索引，多条件过滤会全表扫描。

> **能力缺口**：`QuestionModel` 只有 `FindOne(id)`，无 `FindPage` / `FindByIds` / `CountByCondition`，`ListQuestions` 当前**无法实现**。

### Requirement: 题目删除与级联清理

题目删除 SHALL 为物理删除（`question` 无逻辑删除列），并在同一事务内级联清理关联数据：

#### Scenario: 删除题目（DeleteQuestion 流程）

- **WHEN** 调用 `DeleteQuestion` 且校验 id > 0、`FindOne` 确认题目存在
- **THEN** 系统检查 `question_biz` 中是否仍有该题的关联（被小节引用）：有引用时拒绝删除，或先清理关联（策略待定）
- **AND** 系统在同一事务内依次删除 `question_biz`（where question_id = id）、`question_detail`（where id = id）、`question`（where id = id）

补充约束：

- 三表级联：生成的 `Delete` 只删单表，需 logic 层编排。
- 事务保证：三次删除需原子，否则产生孤儿详情/孤儿关联。
- 权限校验：原始接口摘要为「根据 id 删除**当前用户**问题」，应校验 `creater` 与 JWT userId 一致。

#### Scenario: 非本人题目的删除校验

- **WHEN** 用户调用 `DELETE /questions/{id}`（原始摘要「根据 id 删除当前用户问题」）
- **THEN** 系统应校验 `creater` 与 JWT userId 一致，仅允许删除本人题目

> **能力缺口**：`QuestionBizModel` 无 `DeleteByQuestionId` 方法，级联清理无法实现。

### Requirement: 题目与业务关联规则

`question_biz` SHALL 作为多对多关联表使用，`(biz_id, question_id)` 上有唯一索引：

- 同一业务不可重复挂同一题：唯一索引 `biz_id(biz_id, question_id)` 在 DB 层强制。
- 幂等新增：重复 `AddQuestionBiz` 会触发唯一键冲突，应先 `FindOneByBizIdQuestionId` 判重并返回已有 ID。
- 业务 ID 语义：`biz_id` 指向业务对象，DDL 注释「例如小节id」。
- 可空列转换：`biz_id` / `question_id` 在 Go 中是 `sql.NullInt64`，proto 是 `int64`，需显式转换。
- 解除关联不删题：`RemoveQuestionBiz` 只删关联行，题目保留在题库。

#### Scenario: 幂等挂题（AddQuestionBiz 流程）

- **WHEN** 调用 `AddQuestionBiz` 且 bizId > 0、questionId > 0
- **THEN** 系统校验 question 存在（FindOne），再 `FindOneByBizIdQuestionId` 判重；已存在则直接返回其 id（幂等）
- **AND** 不存在则 Insert 关联行（id 由 AUTO_INCREMENT 生成），返回 `IdReply{ id }`（关联表主键，非题目 ID）

#### Scenario: 撤题（RemoveQuestionBiz 流程）

- **WHEN** 调用 `RemoveQuestionBiz` 且 bizId > 0、questionId > 0
- **THEN** 系统 `FindOneByBizIdQuestionId` 定位关联行并 `Delete(关联行.id)`

#### Scenario: 按业务取题（GetQuestionsByBiz 流程）

- **WHEN** 调用 `GetQuestionsByBiz` 且 bizId > 0、pageNo / pageSize 兜底
- **THEN** 系统按 biz_id 分页查 `question_biz` 得到 question_id 列表（命中唯一索引前缀）
- **AND** 系统 `FindByIds` 批量取 question 主体、批量取 question_detail 聚合为 `QuestionVO`，返回 `QuestionListReply{ total, list }`

> **能力缺口**：`QuestionBizModel` 无按 `biz_id` 的列表查询方法（`FindOneByBizIdQuestionId` 需两个入参且只返回单行），`GetQuestionsByBiz` 当前**无法实现**。
>
> **效率缺口**：`RemoveQuestionBiz` 需两次数据库往返（先查后删），补 `DeleteByBizIdQuestionId` 可合并为一次。

### Requirement: 答题统计

`answer_times` / `correct_times` SHALL 作为只读统计字段，由答题行为驱动累加：

- 不接受客户端传入：`QuestionSaveReq` 中无这两个字段。
- 需原子自增：应使用 `UPDATE ... SET answer_times = answer_times + 1`，而非读改写。
- 正确率派生：`correct_times / answer_times`，不落库。

#### Scenario: 并发答题自增

- **WHEN** 多个答题请求并发更新同一题目的统计字段
- **THEN** 系统必须以原子自增方式累加，禁止全字段覆盖式 `Update`

#### Scenario: 正确率派生

- **WHEN** 需要展示题目正确率
- **THEN** 以 `correct_times / answer_times` 派生计算，不落库

> **并发缺口**：生成的 `Update` 是**全字段覆盖**，并发答题时会互相覆盖统计值，必须补自定义 `IncrAnswerTimes` 方法。
>
> **触发方缺口**：proto 中**没有任何上报答题结果的 RPC 方法**，统计字段目前无写入入口。
>
> **已知缺口**：答题统计无写入入口（proto 无上报答题结果的 RPC 方法，`answer_times` / `correct_times` 暂无累加路径）；course ↔ exam RPC 未接线。

### Requirement: 鉴权

全部 7 个 HTTP 接口 SHALL 要求携带有效 JWT（`apps/exam/api/exam.api` 中两个 service 块均声明 `jwt: Auth`）：

| 路由组 | group | 鉴权 |
|--------|-------|------|
| `/questions`, `/questions/:id` | `question` | 需 JWT |
| `/questions/biz`, `/questions/biz/:bizId` | `questionbiz` | 需 JWT |

`creater` / `updater` 应从 JWT 上下文取 userId 写入。

> **冲突说明**：api-spec.md 中「认证」列均标注为「否」，那是从原始 Java 版 `docs/tjxt.openapi.json` 提取的，与 go-zero 侧 `.api` 的 `jwt: Auth` 声明**不一致**，以 `.api` 为准。

#### Scenario: 未携带 JWT 访问

- **WHEN** 客户端未携带有效 JWT 调用任一 `/questions` 或 `/questions/biz` 路由组接口
- **THEN** 请求被 go-zero JWT 中间件拒绝

#### Scenario: 写操作记录创建人

- **WHEN** 已认证用户执行题目新增/更新
- **THEN** `creater` / `updater` 从 JWT 上下文的 userId 写入

### Requirement: 数据模型与存储不变量

exam 域 SHALL 使用 MySQL `tj_exam` 库，DDL 有两份文件：`sql/ddl/tj_exam.sql`（完整库结构，含 Seata 事务基础设施表）与 `sql/ddl/tj_exam_business.sql`（仅业务表，剔除 `undo_log`，供 goctl 生成 model），两者业务表结构完全一致（差异仅在 `undo_log`、注释详略与 Records 注释块）。model 由 `tj_exam_business.sql` 生成（`questiondetailmodel_gen.go` 中 `Answer` 字段注释为简写版，可判定），3 张业务表已 100% 生成 model 并在 `apps/exam/rpc/internal/svc/servicecontext.go:12-14` 完成注入。数据表清单：

| 表 | 职责 | 关键结构 |
|----|------|---------|
| `question` | 题目表（题干/类型/三级分类/难度/分值/统计） | 无 `deleted` 逻辑删除列；除主键外无二级索引；`creater`/`updater` 默认 1；`dep_id` 为唯一可空列 |
| `question_detail` | 题目详情表（options/answer/analysis） | 与 `question` 共享主键一对一（`id` 同值，非自增）；`options` 为 json 可空；`answer` 为 varchar(40) |
| `question_biz` | 问题与业务关联表（多对多） | 唯一索引 `(biz_id, question_id)`；主键自增（AUTO_INCREMENT 起始 147，说明原始 Java 版已有历史数据）；`biz_id`/`question_id` 均可空 |
| `undo_log` | Seata AT 事务回滚日志表（仅 `tj_exam.sql`） | 无主键，唯一索引 `(xid, branch_id)`；`utf8 / utf8_general_ci`、`ROW_FORMAT=COMPACT`（Seata 官方脚本原样保留） |

> **范围说明**：exam 域数据模型实际仅含上述 4 张表；转换任务描述中提及的 `paper` / `paper_question` / `exam` / `answer_sheet` 等组卷/考试表在 data-model.md（源自 `sql/ddl/tj_exam.sql`、`sql/ddl/tj_exam_business.sql`）中**不存在**，本规格不虚构。

数据不变量（SHALL 满足）：

- `question_detail.id = question.id`（共享主键，非自增），写入时必须先确定 `question.id` 再用同一 ID 插入详情表。
- `question_biz` 唯一索引 `(biz_id, question_id)` 防止同一业务重复挂同一题，并支撑按 `biz_id` 的前缀查询；goctl 依据该索引自动生成 `FindOneByBizIdQuestionId(ctx, bizId, questionId)`，是三个 model 中唯一的非主键查询方法。
- `question_biz.biz_id` 跨域引用业务对象 ID（如 course 域小节 id），**无数据库外键**，靠应用层维护一致性。
- `question` 表 DDL 中除主键外未定义任何二级索引，`ListQuestions` 的多条件过滤将全表扫描。
- `question.creater` / `updater` 默认值为 `1`（非其他域的 `0`）；与 auth / media 域不同，exam 域无逻辑删除列、`dep_id` 可空。
- `undo_log` 不生成 model（Seata 框架表，由 Seata 客户端框架直接读写，**预期行为非缺陷**）；项目当前未接入 Seata（`apps/exam/` 下无任何 Seata 依赖与配置），`undo_log` 实际处于**闲置状态**。

Go 结构体映射关键点：

- `Question.DepId` 为 `sql.NullInt64`（唯一可空列）；`QuestionDetail.Options` 为 `sql.NullString`（json 列映射为可空字符串，读写需自行 `json.Marshal` / `Unmarshal`，proto 的 `QuestionVO.options` 为 string 可直接透传）。
- `question_biz.biz_id` / `question_id` 在 Go 中是 `sql.NullInt64`，proto 的 `QuestionBizReq` 是普通 `int64`，logic 层需做 `sql.NullInt64{Int64: x, Valid: true}` 显式转换；`FindOneByBizIdQuestionId` 的入参也是 `sql.NullInt64`。
- `create_time` / `update_time` 映射为 `time.Time`，连接串需 `parseTime=true`。

缓存 key 前缀（goctl 生成）：`cache:question:id:` / `cache:questionDetail:id:` / `cache:questionBiz:id:`。所有 `*_gen.go` 由 goctl 生成**禁止修改**，扩展方法统一放在同名自定义 `.go` 文件中；当前三个自定义 model 文件均为 goctl 空壳，未添加任何扩展方法。

> **待补齐的扩展方法（缺口）**：`FindPage` / `CountByCondition`（`ListQuestions` 多条件分页）；`FindQuestionIdsByBizId` / `FindByIds` / `FindDetailsByIds`（按业务取题与聚合、避免 N+1）；三表级联删除事务（`DeleteQuestion`，生成的 `Delete` 只删单表）；`DeleteByBizIdQuestionId`（`RemoveQuestionBiz` 合并为一次往返）；`IncrAnswerTimes`（答题统计原子自增，全字段 `Update` 会覆盖并发写）；`SaveQuestion` 双表写入的事务封装。

#### Scenario: 详情与主体同 ID 写入

- **WHEN** 新增一道题目
- **THEN** 系统先生成 `question.id`，再以同一 ID 插入 `question_detail`（共享主键，非自增）

#### Scenario: 重复挂题被唯一索引拒绝

- **WHEN** 同一 `biz_id` 下重复插入同一 `question_id` 的关联行
- **THEN** `question_biz` 的唯一索引 `(biz_id, question_id)` 在 DB 层拒绝该插入

### Requirement: 服务配置与可观测性

exam 服务 SHALL 按以下配置运行。

API 服务配置（`apps/exam/api/etc/exam-api.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Name` | `exam-api` | 服务名称 |
| `Host` | `0.0.0.0` | 监听地址 |
| `Port` | `8805` | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| `ExamRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | etcd 地址 |
| `ExamRpc.Etcd.Key` | `exam.rpc` | exam RPC 服务发现 key |

RPC 服务配置（`apps/exam/rpc/etc/exam.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Name` | `exam.rpc` | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8085` | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | etcd 地址 |
| `Etcd.Key` | `exam.rpc` | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_exam?charset=utf8mb4&parseTime=true&loc=Local` | MySQL 连接串（字符集 utf8mb4、时区 Local，`parseTime=true` 为必需项） |
| `Cache[0].Host` / `Type` / `Pass` | `127.0.0.1:6379` / `node` / (空) | Redis 单机节点缓存 |

端口分配：`exam-api` 8805（HTTP）、`exam.rpc` 8085（gRPC）。

- JWT 密钥：生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造；`exam-api` 的 `Auth.AccessSecret` 必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则所有接口鉴权失败。
- 缓存：`Cache` 传入三个 model 构造函数（`servicecontext.go:21-23`），由 goctl 的 `sqlc.CachedConn` 托管主键缓存；三个 model 共用同一 `sqlx.NewMysql(c.DataSource)` 连接与同一份 `Cache` 配置。
- 可观测性注入配置（**已注入勿删**，configs.md 未记录、以实际 yaml 为准）：`Log`（Mode: file，Encoding: json，Level: info，路径 `logs/exam-api` 与 `logs/exam-rpc`）；`Prometheus`（exam-api `0.0.0.0:9205`、exam.rpc `0.0.0.0:9105`，Path `/metrics`）；`Telemetry`（Name 对应各服务，Endpoint `127.0.0.1:4318`，Sampler `1.0`，Batcher `otlphttp`）。

> **配置缺口 1（无分布式事务/Seata 配置）**：`tj_exam.sql` 定义了 Seata AT 模式 `undo_log` 表，但 `exam.yaml` 无任何 Seata / TC 地址配置，`config.go` 仅 `RpcServerConf` / `DataSource` / `Cache`，`apps/exam/` 全目录检索 `undo_log` / `seata` **零命中** —— `undo_log` 闲置。`SaveQuestion` / `DeleteQuestion` 的多表操作用 go-zero 本地事务（`sqlx.Transact`）即可；若后续需与 course 域联动（挂题同时更新课程目录）则需补齐分布式事务配置。
>
> **配置缺口 2（无跨服务 RPC 客户端配置）**：`exam-api` 配置仅 `ExamRpc` 一个客户端，`exam.rpc` 配置无任何 RpcClient —— `question_biz.biz_id` 语义上指向 course 域小节 ID，但 exam **无法校验该 ID 是否真实存在**（无 `CourseRpc`）；反向 `apps/course/api/internal/svc/servicecontext.go` 未 import `examclient`，course 侧 `CourseSubjectsGet` 也无法回调 exam 取题。两域间 RPC 通道尚未建立。
>
> **配置缺口 3（无分页默认值配置）**：`QuestionListReq.PageNo` / `PageSize` 与 `QuestionBizListReq.PageNo` / `PageSize` 均标注 `optional`，且 yaml 中无分页默认值配置项 —— 客户端不传时为零值，logic 层必须硬编码兜底（如 `pageNo=1, pageSize=20`），否则分页查询会退化为 `LIMIT 0`。
>
> **配置完备性**：对照 auth（API 8802 / RPC 8082）、media（8806 / 8086）、exam（**8805** / **8085**），exam 是配置项最少的服务：无外部第三方依赖，所有配置项均已就位，配置层**不构成实现阻塞**（对比 media 的对象存储配置缺口）；当前实现阻塞项集中在 model 扩展方法层面。

#### Scenario: 分页参数缺省兜底

- **WHEN** 客户端调用 `ListQuestions` / `GetQuestionsByBiz` 且未传 `pageNo` / `pageSize`
- **THEN** logic 层以硬编码默认值兜底（如 `pageNo=1, pageSize=20`），避免分页查询退化为 `LIMIT 0`

#### Scenario: JWT 密钥不一致导致鉴权失败

- **WHEN** `exam-api` 的 `Auth.AccessSecret` 与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 配置不一致
- **THEN** 所有接口鉴权失败；两者必须配置为一致，且生产环境必须替换默认值 `change-me-in-production`

