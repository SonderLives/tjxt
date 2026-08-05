> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/exam/rpc/exam.proto`

---

# Exam RPC Spec

## 服务名

`Exam` — 考试与题库管理微服务，通过 etcd 服务发现（key: `exam.rpc`）。

## RPC 方法总览

### 题目管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveQuestion` | `QuestionSaveReq { id, name, type, cateId1, cateId2, cateId3, difficulty, score, options, answer, analysis }` | `IdReply { id }` | 新增/更新题目，`id` 为空表示新增 |
| `DeleteQuestion` | `IdReq { id }` | `Empty {}` | 删除题目 |
| `GetQuestion` | `IdReq { id }` | `QuestionVO` | 按 ID 查询题目（含详情） |
| `ListQuestions` | `QuestionListReq { pageNo, pageSize, name, type, cateId1, cateId2, difficulty }` | `QuestionListReply { total, list }` | 分页查询题目列表 |

**请求字段说明**：

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

> **注意**：`QuestionSaveReq` 中**没有** `cateId3` 之外的层级限制，但**不含** `answerTimes` / `correctTimes` —— 这两个统计字段由系统维护，不允许客户端传入。
>
> **注意**：`QuestionListReq` 支持 `cateId1` / `cateId2` 过滤，但**不支持 `cateId3`**（proto 中未定义该过滤字段）。

**响应字段说明（`QuestionVO`）**：

| 字段 | 类型 | 说明 | 来源表 |
|------|------|------|--------|
| `id` ~ `score` | - | 题目主体字段 | `question` |
| `answerTimes` | int32 | 回答次数（只读统计） | `question` |
| `correctTimes` | int32 | 回答正确次数（只读统计） | `question` |
| `createTime` | string | 创建时间 | `question` |
| `options` | string | 选择题选项 | `question_detail` |
| `answer` | string | 正确答案 | `question_detail` |
| `analysis` | string | 答案解析 | `question_detail` |

> `QuestionVO` 是 `question` 与 `question_detail` 两表的**聚合视图**，需在 logic 层做一对一 JOIN 组装。

---

### 题目业务关联

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `AddQuestionBiz` | `QuestionBizReq { bizId, questionId }` | `IdReply { id }` | 建立题目与业务对象的关联 |
| `RemoveQuestionBiz` | `QuestionBizReq { bizId, questionId }` | `Empty {}` | 解除题目与业务对象的关联 |
| `GetQuestionsByBiz` | `QuestionBizListReq { bizId, pageNo, pageSize }` | `QuestionListReply { total, list }` | 分页查询某业务下的题目列表 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `bizId` | int64 | 业务 ID，要关联问题的某业务 id，例如小节 id |
| `questionId` | int64 | 问题 ID |
| `pageNo` | int32 | 页码 |
| `pageSize` | int32 | 每页条数 |

> `AddQuestionBiz` 返回 `IdReply { id }`，即 `question_biz` 关联表的自增主键；`RemoveQuestionBiz` 按 `(bizId, questionId)` 组合定位，与表上的唯一索引 `biz_id(biz_id, question_id)` 对应。

---

### 通用消息

| Message | 字段 | 用途 |
|---------|------|------|
| `Empty` | (无) | 删除类方法的响应 |
| `IdReq` | `id` (int64) | 单 ID 请求 |
| `IdReply` | `id` (int64) | 单 ID 响应 |
| `PageReq` | `pageNo` (int32), `pageSize` (int32) | 通用分页请求 |

> `PageReq` 在 proto 中已声明，但 `service Exam` 的 7 个方法**均未引用**（分页参数被内联进 `QuestionListReq` / `QuestionBizListReq`），属于预留消息。

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `exam-api` (自身 API 层) | HTTP Handler → `examclient.Exam` RPC | `apps/exam/api/internal/svc/servicecontext.go:5` 中 import `examclient "tjxt/apps/exam/rpc/exam"`，7 个接口全部走自身 RPC |

**跨服务消费情况**：经检索全部 13 个服务的 `servicecontext.go`，**当前无任何其他服务 import `examclient`**。

**尚未接线的逻辑消费方**：

| 潜在消费方 | 耦合证据 | 当前状态 |
|-----------|---------|---------|
| `course` 服务 | `apps/course/rpc/course.proto:48` 定义 `CourseSubjectsGet(IdRequest) returns (CataSubjectInfoList)`（章节/小节的题目）；`question_biz.biz_id` 注释明确写「例如小节id」 | course 侧未 import `examclient`，`exam-api` 侧也未持有 `CourseRpc`，跨域调用未建立 |
| `learning` 服务 | 学员答题会驱动 `question.answer_times` / `correct_times` 累加 | 无任何接线证据，纯语义推导 |

---

## 调用典型场景

1. **教师出题** → 管理端提交题干/选项/答案 → 调 `SaveQuestion` → 同时写 `question` 与 `question_detail` 两表 → 返回题目 ID
2. **题库检索** → 管理端按分类/类型/难度筛选 → 调 `ListQuestions` → 分页返回题目列表
3. **查看题目详情** → 调 `GetQuestion` → 聚合 `question` + `question_detail` 返回含答案解析的完整视图
4. **小节挂题** → 课程编辑时为某小节选题 → 调 `AddQuestionBiz(bizId=小节id, questionId)` → 写入 `question_biz`
5. **小节撤题** → 调 `RemoveQuestionBiz(bizId, questionId)` → 删除关联，题目本身保留在题库
6. **学员练习取题** → 进入小节练习 → 调 `GetQuestionsByBiz(bizId=小节id)` → 分页拉取该小节的题目
7. **删除题目** → 调 `DeleteQuestion` → 需级联清理 `question_detail` 与 `question_biz`

---

## 自定义 Model 方法

`apps/exam/rpc/internal/model/` 下的三个自定义 model 文件当前**均为 goctl 空壳**，未添加任何扩展方法：

- `questionmodel.go` — `QuestionModel interface { questionModel }`，仅内嵌生成接口
- `questiondetailmodel.go` — `QuestionDetailModel interface { questionDetailModel }`，仅内嵌生成接口
- `questionbizmodel.go` — `QuestionBizModel interface { questionBizModel }`，仅内嵌生成接口

可用方法为 `*_gen.go` 生成的内容：

| Model | 生成方法 |
|-------|---------|
| `QuestionModel` | `Insert`, `FindOne(id)`, `Update`, `Delete` |
| `QuestionDetailModel` | `Insert`, `FindOne(id)`, `Update`, `Delete` |
| `QuestionBizModel` | `Insert`, `FindOne(id)`, `FindOneByBizIdQuestionId(bizId, questionId)`, `Update`, `Delete` |

> `FindOneByBizIdQuestionId` 由 goctl 依据 `question_biz` 表上的唯一索引 `biz_id(biz_id, question_id)` 自动生成，是三个 model 中唯一的非主键查询方法，正好可支撑 `RemoveQuestionBiz` 的定位需求。

> **缺口**：`ListQuestions` 的多条件分页、`GetQuestionsByBiz` 的按 `biz_id` 批量查询、`DeleteQuestion` 的级联清理、答题统计字段自增，生成方法均不支持，需在自定义 model 中补写（详见 [data-model.md](./data-model.md)）。
