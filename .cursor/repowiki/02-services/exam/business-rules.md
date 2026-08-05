> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/exam/rpc/internal/logic/*.go`, `apps/exam/api/internal/logic/**/*.go`

---

# Exam Business Rules

## ⚠️ 实现状态

本服务的业务 logic **尚未实现**，全部为 goctl 脚手架占位（函数体仅含 `// todo: add your logic here and delete this line` 并返回零值）。

### RPC 层实现状态（`apps/exam/rpc/internal/logic/`）

| Logic 文件 | RPC 方法 | 实现状态 |
|-----------|---------|---------|
| `savequestionlogic.go` | `SaveQuestion` | ❌ 未实现（占位） |
| `deletequestionlogic.go` | `DeleteQuestion` | ❌ 未实现（占位） |
| `getquestionlogic.go` | `GetQuestion` | ❌ 未实现（占位） |
| `listquestionslogic.go` | `ListQuestions` | ❌ 未实现（占位） |
| `addquestionbizlogic.go` | `AddQuestionBiz` | ❌ 未实现（占位） |
| `removequestionbizlogic.go` | `RemoveQuestionBiz` | ❌ 未实现（占位） |
| `getquestionsbybizlogic.go` | `GetQuestionsByBiz` | ❌ 未实现（占位） |

**RPC 层统计：已实现 0 / 总计 7**

### API 层实现状态（`apps/exam/api/internal/logic/`）

| Logic 文件 | Handler 方法 | 实现状态 |
|-----------|-------------|---------|
| `question/savequestionlogic.go` | `SaveQuestion` | ❌ 未实现（占位） |
| `question/deletequestionlogic.go` | `DeleteQuestion` | ❌ 未实现（占位） |
| `question/getquestionlogic.go` | `GetQuestion` | ❌ 未实现（占位） |
| `question/listquestionslogic.go` | `ListQuestions` | ❌ 未实现（占位） |
| `questionbiz/addquestionbizlogic.go` | `AddQuestionBiz` | ❌ 未实现（占位） |
| `questionbiz/removequestionbizlogic.go` | `RemoveQuestionBiz` | ❌ 未实现（占位） |
| `questionbiz/getquestionsbybizlogic.go` | `GetQuestionsByBiz` | ❌ 未实现（占位） |

**API 层统计：已实现 0 / 总计 7**

> **合计：已实现 0 / 总计 14**
>
> 以下各节内容**均为设计意图推导**，依据 `apps/exam/rpc/exam.proto` 的消息定义、`sql/ddl/tj_exam.sql` 与 `sql/ddl/tj_exam_business.sql` 的字段注释与枚举、`apps/exam/api/exam.api` 的路由与校验标签、以及 `docs/tjxt.openapi.json` 的原始 Java 版接口契约。**不代表当前代码行为。**

---

## 1. 题目双表写入 📋 设计意图（待实现）

**核心规则**：一道题目跨 `question`（主体）与 `question_detail`（选项/答案/解析）两表，**共享同一主键**。

| 规则 | 说明 | 依据 |
|------|------|------|
| 共享主键 | `question_detail.id = question.id`，详情表主键非自增 | 两表 DDL 均为 `id bigint NOT NULL` |
| 必须同事务 | 主体与详情写入需原子，避免出现无详情的孤儿题目 | 一对一强关联 |
| 新增/更新分支 | `QuestionSaveReq.id` 在 `.api` 中标注 `optional`，沿用 `id <= 0` 为新增的项目约定 | `exam.api` QuestionSaveReq |
| 统计字段不可传入 | `answer_times` / `correct_times` 不在 `QuestionSaveReq` 中，由系统维护 | proto 字段缺席 |
| 无逻辑删除 | `question` 表**无 `deleted` 列**，只能物理删除 | DDL 字段清单 |

```
流程（SaveQuestion，设计意图）:
  1. 校验 name / answer / analysis 非空（.api 中三者均为必填）
  2. 校验 type ∈ {1..5}、difficulty ∈ {1,2,3}、score > 0
  3. 选择题（type ∈ {1,2,3}）校验 options 为合法 JSON 数组
  4. id <= 0 → 新增：生成雪花 ID
                 → question   插入主体（creater/updater 取 JWT userId）
                 → question_detail 用同一 ID 插入详情
     id >  0 → 更新：FindOne 校验存在 → 两表分别 Update
  5. 步骤 4 需包裹在同一本地事务内
  6. 返回 IdReply{ id }
```

> ⚠️ **能力缺口**：三个 model 均为 goctl 空壳，未封装事务；`SaveQuestion` 的双表原子写入需先在自定义 model 或 logic 层引入 `sqlx.Transact`。

## 2. 题目类型与选项校验 📋 设计意图（待实现）

**核心规则**：`type` 决定 `options` 与 `answer` 的形态，由 DDL 注释定义枚举。

| type | 含义 | options 要求（推导） | answer 语义（依据 `tj_exam.sql` 注释） |
|------|------|---------------------|------------------------------------|
| 1 | 单选题 | 必填，JSON 数组 | 单个 1~10 的选项序号 |
| 2 | 多选题 | 必填，JSON 数组 | 多个序号，**逗号隔开** |
| 3 | 不定向选择题 | 必填，JSON 数组 | 单个或多个序号，逗号隔开 |
| 4 | 判断题 | 可空 | **1 代表正确，其他代表错误** |
| 5 | 主观题 | 可空 | 参考答案文本 |

> `answer` 列长度为 `varchar(40)`，多选题逗号分隔的序号串需在此长度内。
>
> `options` 为 MySQL `json` 类型且可空，Go 侧映射为 `sql.NullString`，需手工序列化。

## 3. 题目查询与聚合 📋 设计意图（待实现）

**核心规则**：`QuestionVO` 是 `question` + `question_detail` 的聚合视图。

```
流程（GetQuestion，设计意图）:
  1. 校验 id > 0
  2. QuestionModel.FindOne(id) → 主体，不存在则返回 NotFound
  3. QuestionDetailModel.FindOne(id) → 详情（同一 ID）
  4. 合并为 QuestionVO，createTime 格式化为字符串
```

```
流程（ListQuestions，设计意图）:
  1. pageNo / pageSize 兜底默认值（.api 中均为 optional）
  2. 按 name 模糊、type / cateId1 / cateId2 / difficulty 精确，动态拼条件
  3. 分页查 question 主体 + 查总数
  4. 批量取详情（FindDetailsByIds）避免 N+1
  5. 返回 QuestionListReply{ total, list }
```

| 注意点 | 说明 |
|--------|------|
| 不支持按 `cateId3` 过滤 | `QuestionListReq` 中无该字段，虽然表里有 `cate_id3` 列 |
| 列表是否返回答案 | proto 的 `QuestionVO` 含 `answer` / `analysis`，学员端场景应在上层裁剪，避免泄题 |
| 无二级索引 | `question` 表除主键外无索引，多条件过滤会全表扫描 |

> ⚠️ **能力缺口**：`QuestionModel` 只有 `FindOne(id)`，无 `FindPage` / `FindByIds` / `CountByCondition`，`ListQuestions` 当前**无法实现**。

## 4. 题目删除与级联清理 📋 设计意图（待实现）

**核心规则**：`question` 无逻辑删除列，删除即物理删除，必须级联清理关联数据。

```
流程（DeleteQuestion，设计意图）:
  1. 校验 id > 0
  2. FindOne 校验题目存在
  3. 检查 question_biz 中是否仍有该题的关联（被小节引用）
     → 有引用时：拒绝删除，或先清理关联（策略待定）
  4. 同一事务内依次删除：
     question_biz  (where question_id = id)
     question_detail (where id = id)
     question      (where id = id)
```

| 规则 | 说明 |
|------|------|
| 三表级联 | 生成的 `Delete` 只删单表，需 logic 层编排 |
| 事务保证 | 三次删除需原子，否则产生孤儿详情/孤儿关联 |
| 权限校验 | 原始接口摘要为「根据 id 删除**当前用户**问题」，应校验 `creater` 与 JWT userId 一致 |

> ⚠️ **能力缺口**：`QuestionBizModel` 无 `DeleteByQuestionId` 方法，级联清理无法实现。

## 5. 题目与业务关联 📋 设计意图（待实现）

**核心规则**：`question_biz` 是多对多关联表，`(biz_id, question_id)` 上有唯一索引。

| 规则 | 说明 | 依据 |
|------|------|------|
| 同一业务不可重复挂同一题 | 唯一索引 `biz_id(biz_id, question_id)` 在 DB 层强制 | `question_biz` DDL |
| 幂等新增 | 重复 `AddQuestionBiz` 会触发唯一键冲突，应先 `FindOneByBizIdQuestionId` 判重并返回已有 ID | goctl 生成方法 |
| 业务 ID 语义 | `biz_id` 指向业务对象，DDL 注释「例如小节id」 | `question_biz.biz_id` 注释 |
| 可空列转换 | `biz_id` / `question_id` 在 Go 中是 `sql.NullInt64`，proto 是 `int64`，需显式转换 | `questionbizmodel_gen.go` |
| 解除关联不删题 | `RemoveQuestionBiz` 只删关联行，题目保留在题库 | 方法语义 |

```
流程（AddQuestionBiz，设计意图）:
  1. 校验 bizId > 0 且 questionId > 0
  2. 校验 question 存在（FindOne）
  3. FindOneByBizIdQuestionId 判重 → 已存在则直接返回其 id（幂等）
  4. Insert 关联行（id 由 AUTO_INCREMENT 生成）
  5. 返回 IdReply{ id }（关联表主键，非题目 ID）

流程（RemoveQuestionBiz，设计意图）:
  1. 校验 bizId > 0 且 questionId > 0
  2. FindOneByBizIdQuestionId 定位关联行
  3. Delete(关联行.id)
```

> ⚠️ **效率缺口**：`RemoveQuestionBiz` 需两次数据库往返（先查后删），补 `DeleteByBizIdQuestionId` 可合并为一次。

```
流程（GetQuestionsByBiz，设计意图）:
  1. 校验 bizId > 0，pageNo / pageSize 兜底
  2. 按 biz_id 分页查 question_biz → question_id 列表（命中唯一索引前缀）
  3. FindByIds 批量取 question 主体
  4. 批量取 question_detail 聚合为 QuestionVO
  5. 返回 QuestionListReply{ total, list }
```

> ⚠️ **能力缺口**：`QuestionBizModel` 无按 `biz_id` 的列表查询方法（`FindOneByBizIdQuestionId` 需两个入参且只返回单行），`GetQuestionsByBiz` 当前**无法实现**。

## 6. 答题统计 📋 设计意图（待实现）

**核心规则**：`answer_times` / `correct_times` 为只读统计字段，由答题行为驱动累加。

| 规则 | 说明 |
|------|------|
| 不接受客户端传入 | `QuestionSaveReq` 中无这两个字段 |
| 需原子自增 | 应使用 `UPDATE ... SET answer_times = answer_times + 1`，而非读改写 |
| 正确率派生 | `correct_times / answer_times`，不落库 |

> ⚠️ **并发缺口**：生成的 `Update` 是**全字段覆盖**，并发答题时会互相覆盖统计值，必须补自定义 `IncrAnswerTimes` 方法。
>
> ⚠️ **触发方缺口**：proto 中**没有任何上报答题结果的 RPC 方法**，统计字段目前无写入入口。

## 7. 鉴权 📋 设计意图（待实现）

`apps/exam/api/exam.api` 中两个 service 块均声明 `jwt: Auth`：

| 路由组 | group | 鉴权 |
|--------|-------|------|
| `/questions`, `/questions/:id` | `question` | 需 JWT |
| `/questions/biz`, `/questions/biz/:bizId` | `questionbiz` | 需 JWT |

**全部 7 个 HTTP 接口都要求携带有效 JWT**，`creater` / `updater` 应从 JWT 上下文取 userId 写入。

> 说明：`.cursor/repowiki/02-services/exam/api-spec.md` 中「认证」列均标注为「否」，那是从原始 Java 版 `docs/tjxt.openapi.json` 提取的，与 go-zero 侧 `.api` 的 `jwt: Auth` 声明**不一致**，以 `.api` 为准。

---

## 状态说明

### `question.type` 题目类型

| 值 | 含义 |
|----|------|
| 1 | 单选题 |
| 2 | 多选题 |
| 3 | 不定向选择题 |
| 4 | 判断题 |
| 5 | 主观题 |

### `question.difficulty` 难易度

| 值 | 含义 |
|----|------|
| 1 | 简单 |
| 2 | 中等 |
| 3 | 困难 |

### `question_detail.answer` 答案语义

| 题型 | 取值规则 |
|------|---------|
| 选择题（1/2/3） | 选项序号 1 到 10，多个答案中间使用逗号隔开 |
| 判断题（4） | 1 代表正确，其他代表错误 |
| 主观题（5） | 参考答案文本（受 `varchar(40)` 长度限制） |

### 与其他域的差异

| 约定 | auth / media 域 | exam 域 |
|------|----------------|---------|
| 逻辑删除列 `deleted` | 有 | **无**（物理删除） |
| `creater` / `updater` 默认值 | 0 | **1** |
| `dep_id` | NOT NULL DEFAULT 0 | **可空** |
