> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_exam.sql`, `sql/ddl/tj_exam_business.sql`

---

# Exam Data Model

## DDL 文件说明

exam 域有**两份 DDL 文件**，需区分用途：

| DDL 文件 | 首行声明 | 包含的表 | 用途 |
|---------|---------|---------|------|
| `sql/ddl/tj_exam.sql` | `Pure DDL for goctl model generation (extracted from migration/tj_exam.sql)` | `question`, `question_biz`, `question_detail`, `undo_log` | 完整库结构，含 Seata 事务基础设施表 |
| `sql/ddl/tj_exam_business.sql` | `Pure DDL for goctl model generation (business tables only, undo_log excluded)` | `question`, `question_biz`, `question_detail` | 仅业务表，剔除 `undo_log`，供 goctl 生成 model |

**两份文件的业务表结构完全一致**，差异仅在：

| 差异点 | `tj_exam.sql` | `tj_exam_business.sql` |
|--------|--------------|----------------------|
| `undo_log` 表 | 包含 | 不包含 |
| `question_biz` 表注释 | `问题和业务关联表，例如把小节id和问题id关联，一个小节下可以有多个问题` | `问题和业务关联表`（简写） |
| `question_detail.answer` 列注释 | `选择题正确答案1到10，如果有多个答案，中间使用逗号隔开，如果是判断题，1：代表正确，其他代表错误` | `选择题正确答案`（简写） |
| `-- Records of xxx` 注释块 | 包含 | 不包含 |

**生成来源判定**：`apps/exam/rpc/internal/model/questiondetailmodel_gen.go` 中 `Answer` 字段的行尾注释为 `// 选择题正确答案`（简写版），与 `tj_exam_business.sql` 一致，可判定 **model 是由 `tj_exam_business.sql` 生成的**。

### `tj_exam_business.sql` 的 model 生成情况

| 表 | 是否已生成 model | Model 文件 |
|----|----------------|-----------|
| `question` | ✅ 已生成 | `questionmodel_gen.go` + `questionmodel.go` |
| `question_biz` | ✅ 已生成 | `questionbizmodel_gen.go` + `questionbizmodel.go` |
| `question_detail` | ✅ 已生成 | `questiondetailmodel_gen.go` + `questiondetailmodel.go` |

> **结论**：`tj_exam_business.sql` 的 **3 张业务表已 100% 生成 model，无缺口**。三个 model 也已在 `apps/exam/rpc/internal/svc/servicecontext.go:12-14` 中完成注入。

### `undo_log` 的 model 缺口

| 表 | 是否已生成 model | 说明 |
|----|----------------|------|
| `undo_log` | ❌ **未生成** | 仅存在于 `tj_exam.sql`，`apps/exam/` 全目录检索 `undo_log` / `seata` **零命中** |

> **这是预期行为，非缺陷**：`undo_log` 是 Seata AT 模式的事务回滚日志表（表注释 `AT transaction mode undo table`），由 Seata 客户端框架直接读写，不应生成业务 model。`tj_exam_business.sql` 单独剔除它正是为此。
>
> ⚠️ **真实缺口在别处**：项目当前**未接入 Seata**（`apps/exam/` 下无任何 Seata 依赖与配置），因此 `undo_log` 表虽在 DDL 中定义，实际处于**闲置状态**。若后续 `SaveQuestion` 需要跨 `question` / `question_detail` 双表事务，go-zero 使用本地事务即可；若需跨服务分布式事务（如与 course 域联动挂题），则需先补齐分布式事务框架接入。

---

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `questionmodel.go` | `question` | 无（goctl 空壳，仅内嵌 `questionModel`） |
| `questiondetailmodel.go` | `question_detail` | 无（goctl 空壳，仅内嵌 `questionDetailModel`） |
| `questionbizmodel.go` | `question_biz` | 无（goctl 空壳，仅内嵌 `questionBizModel`） |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `questionmodel.go` 扩展 `QuestionModel` 接口）。

> **当前状态**：三个自定义 model 文件均未添加任何扩展方法。

---

## 表清单与字段说明

### 1. `question` — 题目表

来源：`sql/ddl/tj_exam.sql` / `sql/ddl/tj_exam_business.sql`（两份一致）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 题目 id | PK |
| `name` | varchar(512) | 题干 | 列表模糊搜索 |
| `type` | tinyint | 题目类型，1：单选题，2：多选题，3：不定向选择题，4：判断题，5：主观题 | 过滤条件 |
| `cate_id1` | bigint | 1 级课程分类 id | 过滤条件 |
| `cate_id2` | bigint | 2 级课程分类 id | 过滤条件 |
| `cate_id3` | bigint | 3 级课程分类 id | - |
| `difficulty` | tinyint | 难易度，1：简单，2：中等，3：困难 | 过滤条件 |
| `answer_times` | int | 回答次数，默认 0 | 统计字段 |
| `correct_times` | int | 回答正确次数，默认 0 | 统计字段 |
| `score` | int | 分值 | - |
| `dep_id` | bigint | 部门 id，**可空** | - |
| `create_time` | datetime | 创建时间，默认 CURRENT_TIMESTAMP | 自动填充 |
| `update_time` | datetime | 更新时间，ON UPDATE CURRENT_TIMESTAMP | 自动更新 |
| `creater` | bigint | 创建人，**默认 1** | - |
| `updater` | bigint | 更新人，**默认 1** | - |

> **注意 1**：`question` 表**没有 `deleted` 逻辑删除列**，与 auth / media 域的约定不同。`DeleteQuestion` 只能物理删除。
>
> **注意 2**：`creater` / `updater` 默认值为 `1`（而非其他域的 `0`），`dep_id` 是本表唯一可空列。
>
> **注意 3**：DDL 中除主键外**未定义任何二级索引**，`ListQuestions` 的多条件过滤将全表扫描。

---

### 2. `question_detail` — 题目详情表

来源：`sql/ddl/tj_exam.sql` / `sql/ddl/tj_exam_business.sql`（列注释详略不同，结构一致）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 题目 id（与 `question.id` **同值**，非自增） | PK |
| `options` | json | 选择题选项，json 数组格式，**可空** | - |
| `answer` | varchar(40) | 选择题正确答案 1 到 10，多个答案用逗号隔开；判断题 1 代表正确，其他代表错误 | - |
| `analysis` | varchar(1024) | 答案解析 | - |

> **关系**：与 `question` 是**共享主键的一对一**关系（`question_detail.id = question.id`），不是自增主键。写入时必须先确定 `question.id` 再用同一 ID 插入详情表。
>
> `answer` 的完整语义只在 `tj_exam.sql` 的列注释中给出，`tj_exam_business.sql` 是简写版。

---

### 3. `question_biz` — 问题与业务关联表

来源：`sql/ddl/tj_exam.sql` / `sql/ddl/tj_exam_business.sql`（表注释详略不同，结构一致）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键，**AUTO_INCREMENT**（起始 147） | PK |
| `biz_id` | bigint | 业务 id，要关联问题的某业务 id，例如小节 id，**可空** | 唯一索引首列 |
| `question_id` | bigint | 问题 id，**可空** | 唯一索引次列 |

**索引**：

| 索引名 | 类型 | 列 | 作用 |
|--------|------|-----|------|
| PRIMARY | 主键 | `id` | 自增主键 |
| `biz_id` | UNIQUE | `(biz_id ASC, question_id ASC)` | 防止同一业务重复挂同一题目；支撑按 `biz_id` 的前缀查询 |

> **关系**：一个业务对象（如小节）可挂多道题，一道题也可被多个业务复用 —— **多对多关联表**。
>
> 唯一索引 `(biz_id, question_id)` 使 goctl 自动生成了 `FindOneByBizIdQuestionId(ctx, bizId, questionId)` 方法，是三个 model 中唯一的非主键查询能力。
>
> `AUTO_INCREMENT = 147` 说明该表在原始 Java 版本中已有历史数据。

---

### 4. `undo_log` — Seata AT 事务回滚日志表

来源：**仅 `sql/ddl/tj_exam.sql`**（`tj_exam_business.sql` 已剔除）

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `branch_id` | bigint | branch transaction id | 唯一索引次列 |
| `xid` | varchar(100) | global transaction id | 唯一索引首列 |
| `context` | varchar(128) | undo_log context, such as serialization | - |
| `rollback_info` | longblob | rollback info | - |
| `log_status` | int | 0: normal status, 1: defense status | - |
| `log_created` | datetime(6) | create datetime | - |
| `log_modified` | datetime(6) | modify datetime | - |

**索引**：

| 索引名 | 类型 | 列 |
|--------|------|-----|
| `ux_undo_log` | UNIQUE | `(xid ASC, branch_id ASC)` |

> **字符集特殊**：本表使用 `utf8 / utf8_general_ci`、`ROW_FORMAT = COMPACT`，与其余业务表的 `utf8mb4 / utf8mb4_0900_ai_ci`、`ROW_FORMAT = Dynamic` 不同 —— 这是 Seata 官方建表脚本的原样保留。
>
> **无主键**：仅有唯一索引，无 PRIMARY KEY，符合 Seata 官方定义。
>
> **不生成 model**：见上文「`undo_log` 的 model 缺口」小节。

---

## Go 结构体映射

goctl 生成的结构体（`apps/exam/rpc/internal/model/*_gen.go`）：

| Go 字段 | Go 类型 | db tag | 备注 |
|---------|---------|--------|------|
| `Question.Id` | int64 | `id` | - |
| `Question.Name` | string | `name` | - |
| `Question.Type` / `Difficulty` | int64 | `type` / `difficulty` | tinyint → int64 |
| `Question.CateId1/2/3` | int64 | `cate_id1/2/3` | - |
| `Question.AnswerTimes` / `CorrectTimes` / `Score` | int64 | 同名 | int → int64 |
| `Question.DepId` | **sql.NullInt64** | `dep_id` | 唯一可空列 |
| `Question.CreateTime` / `UpdateTime` | time.Time | `create_time` / `update_time` | 需 `parseTime=true` |
| `Question.Creater` / `Updater` | int64 | `creater` / `updater` | - |
| `QuestionDetail.Id` | int64 | `id` | 与 `Question.Id` 同值 |
| `QuestionDetail.Options` | **sql.NullString** | `options` | json 列映射为可空字符串 |
| `QuestionDetail.Answer` / `Analysis` | string | 同名 | - |
| `QuestionBiz.Id` | int64 | `id` | 自增 |
| `QuestionBiz.BizId` / `QuestionId` | **sql.NullInt64** | `biz_id` / `question_id` | 均为可空列 |

> ⚠️ **可空类型影响 API 契约**：`question_biz.biz_id` / `question_id` 在 Go 中是 `sql.NullInt64`，而 proto 的 `QuestionBizReq` 是普通 `int64`，logic 层需做 `sql.NullInt64{Int64: x, Valid: true}` 的显式转换。`FindOneByBizIdQuestionId` 的入参也是 `sql.NullInt64`。
>
> ⚠️ **json 列处理**：`options` 是 MySQL `json` 类型但映射为 `sql.NullString`，读写时需自行 `json.Marshal` / `Unmarshal`，proto 的 `QuestionVO.options` 也是 string，可直接透传。

---

## 关系图

```
question (1) ──1:1(共享主键)── question_detail
   id                              id

question (1) ──── (N) question_biz (N) ──── 业务对象
   id                  question_id
                       biz_id ──→ 跨域引用（如 course 小节 id，无外键）

undo_log (独立，Seata 基础设施，与业务表无关联)
```

**跨域引用**：`question_biz.biz_id` 指向业务对象 ID（DDL 注释「例如小节id」），对应 course 域的小节，**无数据库外键**，靠应用层维护一致性。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
questionmodel_gen.go      ← goctl 生成，只读
questionmodel.go          ← 手写扩展位置（当前为空壳）
```

当前项目自定义 Model 模式：
- `questionmodel.go` — 无扩展方法
- `questiondetailmodel.go` — 无扩展方法
- `questionbizmodel.go` — 无扩展方法

**缓存 key 前缀**（由 goctl 生成）：

| Model | 缓存前缀 |
|-------|---------|
| `question` | `cache:question:id:` |
| `questionDetail` | `cache:questionDetail:id:` |
| `questionBiz` | `cache:questionBiz:id:` |

### 待补齐的扩展方法（缺口）

| 需求来源 | 需要的方法 | 说明 |
|---------|-----------|------|
| `ListQuestions` RPC | `FindPage(ctx, offset, limit, name, type, cateId1, cateId2, difficulty)` | 多条件动态分页，生成方法完全不支持 |
| `ListQuestions` RPC | `CountByCondition(ctx, ...)` | 分页总数 |
| `GetQuestionsByBiz` RPC | `FindQuestionIdsByBizId(ctx, bizId, offset, limit)` | 按 `biz_id` 批量取题目 ID（可命中唯一索引前缀） |
| `GetQuestionsByBiz` RPC | `FindByIds(ctx, ids)` | 批量取题目，避免 N+1 查询 |
| `GetQuestion` / `ListQuestions` | `FindDetailsByIds(ctx, ids)` | 批量取详情，聚合成 `QuestionVO` |
| `DeleteQuestion` RPC | 级联删除（事务） | 需同时清理 `question_detail` 与 `question_biz`，生成的 `Delete` 只删单表 |
| `RemoveQuestionBiz` RPC | `DeleteByBizIdQuestionId(ctx, bizId, questionId)` | 现只能先 `FindOneByBizIdQuestionId` 再 `Delete(id)`，两次往返 |
| 答题统计 | `IncrAnswerTimes(ctx, id, correct bool)` | `answer_times` / `correct_times` 原子自增，全字段 `Update` 会覆盖并发写 |
| `SaveQuestion` 双表写入 | 事务封装 | `question` + `question_detail` 需在同一本地事务内写入 |
