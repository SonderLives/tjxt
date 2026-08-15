> 版本：v1.4 | 更新：2026-08-15 | 来源：全仓 `go build ./...`（逐模块）+ 逻辑文件清点 + 依赖接线扫描（v1.1 的复核）+ 模块拆分重构（go-zero 官方标准结构）

---

# 实现进度与已知缺口

## 统计口径

- 以每个 `go.mod` 模块独立 `go build ./...` 是否通过 + `internal/logic/**/*logic.go` 是否含真实实现（非 goctl `todo: add your logic` 占位）为准。
- 2026-08-06 模块拆分重构：全仓 **0 处** `TODO` / `panic("implement")` / `not implemented` 占位；同日完成「每服务单 module → 每 api/rpc 独立 module」go-zero 官方标准拆分，删除聚合 `apps/<svc>/go.mod` 与孤儿 `apps/data/api/go.mod`，**28 个 use 模块（根 + pkg + 13 服务 × {api,rpc} 共 26 个）独立 `go build ./...` 全部通过**。
- 接口契约总数：API 端点 **193** + RPC 方法 **201** = **394** 个，与 logic 文件 **1:1 对应**，无缺失。

---

## 1. 逐服务实现进度

| 服务 | API Logic | RPC Logic | 小计 | 编译 | 状态 |
|------|-----------|-----------|------|:----:|------|
| **auth** | 18/18 | 19/19 | 37/37 | ✅ | 完成 |
| **user** | 13/13 | 15/15 | 28/28 | ✅ | 完成 |
| **pay** | 16/16 | 17/17 | 33/33 | ✅ | 完成 |
| **promotion** | 16/16 | 16/16 | 32/32 | ✅ | 完成 |
| **remark** | 2/2 | 2/2 | 4/4 | ✅ | 完成 |
| **course** | 38/38 | 38/38 | 76/76 | ✅ | 完成 |
| **trade** | 36/36 | 37/37 | 73/73 | ✅ | 完成 |
| **learning** | 9/9 | 11/11 | 20/20 | ✅ | 完成 |
| **media** | 10/10 | 10/10 | 20/20 | ✅ | 完成（对象存储为 mock） |
| **exam** | 7/7 | 7/7 | 14/14 | ✅ | 完成 |
| **message** | 18/18 | 19/19 | 37/37 | ✅ | 完成 |
| **search** | 4/4 | 4/4 | 8/8 | ✅ | 完成（依赖 course） |
| **data** | 6/6 | 6/6 | 12/12 | ✅* | 完成（*真实模块 `apps/data/api/data`） |
| **合计** | **193/193** | **201/201** | **394/394** | — | **100%** |

**服务维度：13/13 完成。**

代码实现完成度 **≈100%**（所有契约均有可编译的真实逻辑实现）；功能完备度 **≈97%**（少量外部集成桩，见 §2）。

> 📌 **v1.1 → v1.2 关键变化**：`media` / `exam` / `message` / `data` / `search` 五个服务已由「⬜ 未开始」校正为「✅ 完成」。本批次补齐了 RPC/API 薄封装与自定义 Model（media / exam 的自定义 Model 已补齐分页/软删等方法，不再是空壳，见 §2.4 修正）。全仓逻辑文件由 v1.1 的 303/390（77.7%）校正为 **394/394（100%）**。

---

## 2. 结构性缺口与已知问题

### 2.1 data 服务目录结构偏离约定 — 部分解决（模块拆分已修复，嵌套目录仍保留）

| 问题 | 说明 | 状态 |
|------|------|------|
| 多套一层目录 | `apps/data/api/data/...`、`apps/data/rpc/data/...`，其余 12 个服务均为 `apps/<svc>/{api,rpc}/...` | 仍保留（可接受） |
| 多份 go.mod | 原 `apps/data/api/go.mod` 为孤儿（未纳入 go.work，仅 3 行、内容残缺） | ✅ 已删除；现仅 `apps/data/api/data` 与 `apps/data/rpc/data` 两个真实模块，与各服务一致（仅深一级） |
| 无 DDL | `sql/ddl/` 下无 `tj_data.sql`，无任何 model（`data` 服务不依赖 MySQL，数据为内存/配置驱动） | 仍有效 |
| 配置缺段 | `data.yaml` 无 `DataSource` / `Cache` 段；JWT 配置已声明但路由未启用（见 §2.5 data 写接口裸奔） | 仍有效 |

**说明**：2026-08-06 模块拆分重构已统一为 go-zero 官方标准结构——每个 api/rpc 各有独立 `go.mod`，`go.work` 聚合全部 28 个 use 模块。data 的「多一层目录 + 两个真实模块」是历史遗留，现已与其余服务同构（仅路径深一级），不再是孤儿模块问题。

### 2.2 learning 表结构与缺口 — 仍有效（设计演进）

`sql/ddl/tj_learning.sql` 仅定义 `learning_lesson` 一张表。proto 注释已明确：**学习记录不再单独建表**，通过更新 `learning_lesson` 的 `latest_section_id` / `latest_learn_time` / `learned_sections` 实现；原计划里的 `learning_plan` / `learning_record` 表**按设计演进不再需要**。

**仍存在的缺口**（非缺表，而是缺数据源/出口）：
- `LearningRecordsByCourse.records` 明细列表恒空（单行无法还原多小节历史）；
- `PlanPageReply.week_finished` / `week_points` 无数据源（恒为 0）；
- `status=2`（完成）无写入路径；
- `RemoveLesson` 无 proto/`.api` 出口。

详见 `02-services/learning/business-rules.md` 末节「已知缺口汇总」。

### 2.3 media 对象存储配置完全缺失 — 仍有效（首要阻塞项）

`config.go` 与 `etc/*.yaml` 中均无 `SecretId` / `SecretKey` / `Bucket` / VOD 等配置项，三个签名类 RPC **使用本地 mock 实现**（`mockBaseURL = http://127.0.0.1:9000`，签名/上传/播放 URL 均为本地占位）。logic 已全部落地并编译通过，但接入真实 COS/OSS 前媒资全链路仅为本地可跑的桩。

### 2.4 Model 空壳 — v1.2 已修正（media/exam 已补齐）

v1.1 标注 media / exam 自定义 model 为空壳。2026-08-06 复核：
- **media**：`filemodel.go` / `mediamodel.go` 各 5 个自定义方法（分页 / 软删 / 状态更新等）已实现；
- **exam**：`questionmodel.go` / `questionbizmodel.go` / `questiondetailmodel.go` 共 7 个自定义方法（双表事务写入 / 批量查 / 级联删 / 按 biz 列表等）已实现。

**本条在 v1.2 判定为已解决。**

### 2.5 跨域 RPC 未接线 — 仍有效

**已接线（编译通过，zrpc 客户端已实例化）**：
- `trade → {course, pay}`（RPC→RPC）
- `search → course`（RPC→RPC）
- `learning → course`（API→RPC）

**未接线（servicecontext 无对应 client）**：
| 调用关系 | 状态 |
|---------|------|
| course → user（教师信息） | 未接线（course 侧注释标明） |
| course → learning（课时） | 未接线（lesson_id 填 0） |
| course → media / course → exam | 未接线 |
| trade → promotion（优惠券） | 未接线（优惠未接入，实付金额未扣减） |
| pay → 真实支付网关 | 占位（demo 回调 URL，非真实渠道） |

**事件驱动**：`trade` 的 `ServiceContext` 已实例化 `MQProducer`（RabbitMQ），但**全仓 logic 仍未调用 `Publish`**——领域事件未实际发射，故 learning 消费端**无生产者**（v1.1 结论 v1.2 复核仍成立）。Producer 在 MQ 不可用时优雅降级为 nil，不阻塞启动。

### 2.6 Makefile 服务清单不全 — 仍有效（已复核）

`SERVICES` / `API_DIRS` / `RPC_DIRS` 三个变量仍只列了 11 个服务，**遗漏 `promotion` 与 `remark`**（2026-08-06 复核确认）。这两个服务无法通过 `make run-*`、`make api`、`make rpc` 纳入统一流程，需手工执行。

### 2.7 Seata 未接入 — 仍有效

`sql/ddl/tj_exam.sql`、`sql/ddl/tj_trade.sql` 中定义了 `undo_log` 表（Seata AT 模式所需），但项目未接入 Seata，该表处于闲置状态，也未生成 model。

### 2.8 可观测性（trace/metrics/logs 统一经 otel-collector）— 已完成（2026-08-15）

三类信号已统一收口到 `otel-collector`，业务代码零改动：

- **Trace**：26 个服务 yaml 注入 `Telemetry{Endpoint:127.0.0.1:4318, Batcher:otlphttp, Sampler:1.0}`；collector `otlp` receiver → `batch` → `otlphttp/jaeger` → `jaeger:4318`；Jaeger UI 16686。
- **Metrics**：26 个服务 yaml 注入 `Prometheus{Host:0.0.0.0, Port:<唯一>, Path:/metrics}`（RPC 9101–9113 / API 9201–9213）；collector `prometheus` receiver 经 `host.docker.internal` 抓取 → `prometheus` exporter `:8889` 聚合；Prometheus 只抓 `otel-collector:8889`。
- **Logs**：26 个服务 yaml 注入 `Log{Mode:file, Encoding:json, Path:logs/<svc>, Level:info}`；collector `filelog` 读 `/var/log/tjxt/**/*.log`（宿主 `./logs` 挂载）→ `transform/logsvc` 提升 `service.name`/`level` → `loki` exporter → `loki:3100`；流标签 `service_name`/`level`。

配套新增 `deploy/otel-collector/config.yaml`（三管线）、`deploy/prometheus/prometheus.yml`（只抓 collector）、`deploy/loki/loki-config.yaml`（单实例）、`docker-compose.yml` 加 jaeger/otel-collector/prometheus/loki 四服务。详见 [可观测性文档](../04-infra/observability.md)。

> ⚠️ `Log.Path` 为相对路径，依赖进程 CWD=仓库根；务必从仓库根启动服务，否则日志散落到各服务子目录，collector 采集不到。

---

## 3. 数据一致性 / 正确性隐患

> 下列为 v1.1 复核记录，v1.2 未逐条复验代码逻辑，建议代码评审确认；media / exam 已实现后其相关 Model 隐患应已缓解。

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

## 5. 建议实施顺序（集成联调 / 外部对接阶段）

```
✅ 13/13 服务 logic 全部落地并编译通过（代码实现 ≈100%）
✅ 可观测性：trace/metrics/logs 统一经 otel-collector 收口（Jaeger/Prometheus/Loki），业务代码零改动
  → ① media：接入真实对象存储（腾讯云 COS / 阿里 OSS）替换 mock
  → ② pay：接入真实支付网关回调（替换 demo URL）
  → ③ trade → promotion：打通优惠券计算
  → ④ course → user(教师) / learning(课时) / media：补全跨域数据
  → ⑤ trade 发射领域事件（接 learning 消费端，打通 MQ）
  → ⑥ ✅ 已清理 data 孤儿 go.mod（模块拆分重构已完成）；Makefile 补 promotion / remark
  → ⑦ 复核 §3 数据一致性隐患（代码评审）
```

---

## 相关文档

- [架构全览](../00-architecture/overview.md)
- [服务拓扑](../00-architecture/service-topology.md)
- [快速上手](../05-development/quickstart.md)
- [索引](../index.md)
