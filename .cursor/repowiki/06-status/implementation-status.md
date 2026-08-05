> 版本：v1.0 | 更新：2026-08-05 | 来源：全仓扫描 `apps/<svc>/{api,rpc}/internal/logic/**/*logic.go`

---

# 实现进度与已知缺口

## 统计口径

递归扫描各服务 `internal/logic/**/*logic.go`，文件内含 goctl 占位注释 `todo: add your logic` 即判定为**未实现**。

---

## 1. 逐服务实现进度

| 服务 | API Logic | RPC Logic | 小计 | 状态 |
|------|-----------|-----------|------|------|
| **auth** | 18/18 | 19/19 | 37/37 | ✅ 完成 |
| **user** | 13/13 | 15/15 | 28/28 | ✅ 完成 |
| **pay** | 16/16 | 17/17 | 33/33 | ✅ 完成 |
| **promotion** | 16/16 | 16/16 | 32/32 | ✅ 完成 |
| **remark** | 2/2 | 2/2 | 4/4 | ✅ 完成 |
| **course** | 38/38 | 38/38 | 76/76 | ✅ 完成 |
| **trade** | 0/36 | 0/37 | 0/73 | ⬜ 未开始 |
| **message** | 0/18 | 0/19 | 0/37 | ⬜ 未开始 |
| **learning** | 0/9 | 0/11 | 0/20 | ⬜ 未开始（Service 层与 Model 扩展已就绪） |
| **media** | 0/10 | 0/10 | 0/20 | ⬜ 未开始 |
| **exam** | 0/7 | 0/7 | 0/14 | ⬜ 未开始 |
| **data** | 0/6 | 0/6 | 0/12 | ⬜ 未开始 |
| **search** | 0/2 | 0/2 | 0/4 | ⬜ 未开始 |
| **合计** | **103/195** | **107/195** | **210/390** | **53.8%** |

**服务维度**：6/13 完成（auth / user / pay / promotion / remark / **course**）。

> 📌 `course` 服务已于 2026-08-05 完成全部 76 个 logic 实现并通过编译（model 层自定义查询 + RPC/API 双层）。其 business-rules.md 已从「设计意图」校正为「已落地实现」（状态机仅 1–4、媒资方法 `CourseMediaSave`、草稿→正式发布流程、step 推进缺口等）。course 仍存在的跨域缺口：未装配 `MediaRpc`/`ExamRpc`/`UserRpc` 客户端，销量/老师详情等字段本地填 0/留空（见 2.5 与 course/business-rules.md 末节）。

---

## 2. 结构性缺口

### 2.1 data 服务目录结构偏离约定

| 问题 | 说明 |
|------|------|
| 多套一层目录 | `apps/data/api/data/...`、`apps/data/rpc/data/...`，其余 12 个服务均为 `apps/<svc>/{api,rpc}/...` |
| 多份 go.mod | 存在三份，其中 `api/go.mod` 仅 3 行、内容残缺 |
| 无 DDL | `sql/ddl/` 下无 `tj_data.sql`，无任何 model |
| 配置缺段 | `data.yaml` 无 `DataSource` / `Cache` 段；JWT 配置已声明但路由未启用 |

**建议**：按其余服务的约定重建目录并统一为单 `go.mod`。

### 2.2 learning 缺表

`sql/ddl/tj_learning.sql` 仅定义 `learning_lesson` 一张表，但 RPC 已声明 `PlanSave`、`CommitLearningRecord` 等接口。

**缺失**：`learning_plan`、`learning_record` 两张表（DDL 未定义、model 未生成）→ 导致 `records` 明细、`week_finished`、`week_points` 无数据源。

### 2.3 media 对象存储配置完全缺失

`config.go` 与 `etc/*.yaml` 中均无 `SecretId` / `SecretKey` / `Bucket` / VOD 等配置项，三个签名类 RPC **无法实现**。这是 media 服务的**首要阻塞项**。

### 2.4 Model 空壳

| 服务 | 情况 |
|------|------|
| trade | 4 个自定义 Model 全为 goctl 空壳，只有基础 CRUD，分页/聚合无 SQL 支撑 |
| media / exam | 5 个自定义 model 均为空壳，缺 `FindPage` / `SoftDelete` / 级联删除 / 统计自增，导致 `MediaList`、`ListQuestions`、`GetQuestionsByBiz` 无法实现 |

### 2.5 跨域 RPC 未接线

| 调用关系 | 状态 |
|---------|------|
| trade → pay | `PayRpc` 已装配且配了 etcd key，但 logic 层**零调用** |
| course → media | `servicecontext.go` 无对应 client |
| course → exam | `servicecontext.go` 无对应 client |
| learning ← MQ | 消费端配置齐全但**未接线**，trade 发出的 `order.pay` / `order.refund` 事件无人消费 |

### 2.6 Makefile 服务清单不全

`SERVICES` / `API_DIRS` / `RPC_DIRS` 三个变量均只列了 11 个服务，**遗漏 `promotion` 与 `remark`**。这两个服务无法通过 `make run-*`、`make api`、`make rpc` 纳入统一流程，需手工执行。

### 2.7 Seata 未接入

`sql/ddl/tj_exam.sql`、`sql/ddl/tj_trade.sql` 中定义了 `undo_log` 表（Seata AT 模式所需），但项目未接入 Seata，该表处于闲置状态，也未生成 model。

---

## 3. 数据一致性 / 正确性隐患

| 服务 | 隐患 | 说明 |
|------|------|------|
| pay | 缓存失效缺口 | `MarkToSuccess` / `MarkToClosed` 走 `ExecNoCacheCtx` 却不调 `DelCacheCtx`，而 `FindOneByBizOrderNo` 是带缓存查询 → **可能读到脏数据** |
| pay | 回调链未落地 | `IncrNotifyTimes` / `SetNotifyStatus` / `IncrNotifyFailedTimes` 无任何调用点；`ApplyRefund` 为 mock 同步成功 |
| pay | 幂等无 DB 兜底 | `biz_refund_order_no`、`channel_code` 均无唯一索引 |
| pay | 时区不一致 | pay 使用 `Asia/Shanghai`，其余服务统一使用 `Local` |
| promotion | 兑换流程无补偿 | 兑换码已核销、库存已扣减后，若用户券落库失败**不回滚**（与领取流程的补偿逻辑不对称） |
| promotion | 状态常量悬空 | `UserCouponStatusExpired` 已定义但无写入点，过期靠 `expire_time` 与 now 比较判定；DDL 注释中的 `refunded` 状态未定义常量，退还实际写回 `unused` |
| promotion | 操作人字段恒 0 | `buildCoupon(in, 0)` 硬编码 operator，`creater` / `updater` 恒为 0 |
| user | 改密校验薄弱 | `UpdateStudentPassword` 不校验 `code` 与原密码 |
| learning | 形参被丢弃 | `UpdateLatestLearn` 的 `moment` / `duration`、`CommitRecord` 的 `commitTime` / `userID` 未使用；`status=2`（完成）无写入路径 |
| learning | 能力无出口 | `RemoveLesson` 下层已实现，但 proto 与 `.api` 均无对应出口 |
| — | proto 字段恒空 | `LoginVerifyResponse.name`、`NotifyPaySuccessRequest.qr_code_url` 从未回填 |

---

## 4. 文档与实现的偏差

| 位置 | 偏差 |
|------|------|
| `02-services/promotion/api-spec.md` | 缺少 `GET /codes/page`（routes.go 中真实存在）；全部接口标注「认证：否」，实际 routes.go 用 `rest.WithJwt` 全局包裹 |
| `02-services/exam/api-spec.md` | 标注「认证：否」与 `.api` 中的 `jwt: Auth` 冲突，**以 `.api` 为准** |
| `02-services/auth/rpc-spec.md` | 「被哪些 API 服务消费」为推测内容；实测 `user.rpc` 仅被 user-api / auth-api 消费，`pay.rpc` 仅被 pay-api / trade-rpc 消费 |
| `apps/promotion/rpc/promotion.proto` | 头注释称「供 promotion-api 与 trade 等内部服务调用」，但全仓 Grep 后 `apps/trade` 并未 import promotion client |
| 各服务 `api-spec.md` | 由 `docs/tjxt.openapi.json` 机械提取，响应体列残留 URL 编码（`%C2%AB` = `«`），可读性差，待重新生成 |

---

## 5. 建议实施顺序（依赖驱动）

```
✅ course（已完成：76/76 logic 实现并编译通过）
  → media（course 依赖素材；须先补对象存储配置）   ← 下一优先级
  → learning（须先补 learning_plan / learning_record 表）
  → trade（须先接线 pay，补 Model 分页/聚合）
  → exam
  → message
  → search
  → data（须先重建目录结构）
```

---

## 相关文档

- [架构全览](../00-architecture/overview.md)
- [服务拓扑](../00-architecture/service-topology.md)
- [快速上手](../05-development/quickstart.md)
