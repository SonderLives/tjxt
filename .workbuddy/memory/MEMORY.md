# 天机学堂 tjxt 项目记忆

## 架构
- go-zero v1.10.3，go.work 工作区聚合 **28 个 use 模块**（根 + `pkg` + 13 服务各含 `api`/`rpc` = 26）。每个 api/rpc 各有独立 `go.mod`（2026-08-06 由「每服务单 module」重构为官方标准）。
- module path：`tjxt/apps/<svc>/api`、`tjxt/apps/<svc>/rpc`；data 嵌套为 `tjxt/apps/data/api/data`、`tjxt/apps/data/rpc/data`（配置在 `apps/data/api/data/etc/`、`apps/data/rpc/data/etc/`）。
- 公共库 `tjxt/pkg`：各模块相对 `replace` 引用——普通三级 `../../../pkg`，data 四级 `../../../../pkg`。跨服务依赖在调用方 `require`（伪版本 `v0.0.0-00010101000000-000000000000`）+ `replace` 本地相对路径。
- Model 放 `apps/<svc>/rpc/internal/model/`（custom `*model.go` + 生成 `*model_gen.go` + vars.go）；RPC client 在 `apps/<svc>/rpc/<svc>/<svc>.go`。命名风格 gozero。

## 服务与端口（2026-08-06 重排，唯一无冲突）
API 8801–8813 / RPC 8081–8093 / metrics RPC 9101–9113 / metrics API 9201–9213。

| 服务 | API | RPC | 库 |
|------|-----|-----|-----|
| user | 8801 | 8081 | tj_user |
| auth | 8802 | 8082 | tj_user |
| course | 8803 | 8083 | tj_course |
| learning | 8804 | 8084 | tj_learning |
| exam | 8805 | 8085 | tj_learning |
| media | 8806 | 8086 | tj_media |
| message | 8807 | 8087 | tj_message |
| pay | 8808 | 8088 | tj_pay |
| trade | 8809 | 8089 | tj_trade |
| search | 8810 | 8090 | tj_search |
| data | 8811 | 8091 | 无 DB |
| promotion | 8812 | 8092 | tj_promotion |
| remark | 8813 | 8093 | tj_remark |

改端口时用「旧→新映射单趟 re.sub」（新旧集合重叠，逐条回扫会连锁错改），同步 26 份 yaml + README + `openspec/specs`（architecture/service-topology 与各服务 spec）+ docs。

## go-zero / goctl 规范（用户强约束）
- 生成优先：只手写 `internal/logic`、custom `*model.go`、业务扩展（consts/vars）。禁止手写 handler/types/routes/pb/server/config 骨架，禁止改 `*_gen.go`。
- 开发顺序：DDL(`sql/ddl/tj_<domain>.sql`) → Model → Proto → RPC → .api → API → Logic。Model 已存在则扩展，不覆盖重生成。
- 命令：`cd apps/X/api && goctl api go --api X.api --dir . --style gozero`；`cd apps/X/rpc && goctl rpc protoc X.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style gozero --client=true`；`goctl model mysql ddl -src sql/ddl/tj_X.sql -dir apps/X/rpc/internal/model -cache --style gozero`。
- 覆盖行为：api go 只刷新 types/routes；rpc protoc 只刷新 pb/server/client（`--client=true` 必带）。需强制重生成要先手工删文件。
- `.api` 不支持 `any`（用 `interface{}`）、拒绝顶层数组请求体。统一响应 `result.Write(w,r,data,err)`（`pkg/response.R{Code,Msg,RequestId,Data}`），不用 goctl 生成的 Result 类型。分页 `PageRequest{PageNo,PageSize}`。JWT：`@server(jwt:Auth)`，logic 用 `pkg/auth.GetUserId`/`UserIdFromCtx`。

## 关键坑
- goctl rpc protoc 重跑会在扁平 `internal/logic/` 落 todo 桩并把 `<svc>server.go` 改指扁平包 → 与子包真实实现（如 search 的 `internal/logic/search`）错位，导致「编译通过但返回空结果」。判别：`grep -n "internal/logic" internal/server/*/<svc>server.go`。修复：删扁平桩、保留子包实现。
- custom model 嵌入 `CachedConn`：只能用 `QueryRowsNoCacheCtx`/`QueryRowNoCacheCtx`/`ExecNoCacheCtx`（`QueryRowsCtx` 被遮蔽会编译失败）。
- 批量改 yaml：用 Write 工具写 .py 脚本再执行，勿用 git-bash heredoc（`\n` 会字面化破坏 YAML）。
- go-zero config 不支持 `${ENV}` 占位；`Log.Path` 是相对路径，服务必须从仓库根启动（否则 Loki 采不到）。
- `Makefile` 的 `SERVICES` 只有 11 个，缺 `promotion`/`remark`（无法 `make run-*`）。protoc 在 `/d/Program Files/protoc-35.0-win64/bin`，改 proto 前需 `export PATH`。

## 可观测性（2026-08-15 接入，经 otel-collector 收口）
- 26 份 yaml 已注入三段，**勿删勿改**：`Telemetry{Name,Endpoint:127.0.0.1:4318,Sampler:1.0,Batcher:otlphttp}`（v1.10.3 的 Batcher 无 `jaeger`，用 otlphttp/otlpgrpc）、`Prometheus{Host:0.0.0.0,Port 唯一,Path:/metrics}`（Host 为空则不启动）、`Log{Mode:file,Encoding:json,Path:logs/<svc>,Level:info}`。
- go-zero 无 OTLP metrics/logs 导出 → metrics 走「collector scrape 代理」（prometheus receiver 抓 26 个 /metrics，prometheus exporter 在 8889；Prometheus 只抓 `otel-collector:8889`）；logs 走 filelog 采集（JSON 字段固定 `@timestamp`/`level`/`content`，服务名只能从路径正则取）。
- docker-compose 8 容器：MySQL(3306 root/0000)、Redis(6379)、RabbitMQ(5672/15672)、etcd(2379)、Jaeger(16686)、otel-collector(4317/4318/8889/13133)、Prometheus(9090)、Loki(3100)。collector 配置 `deploy/otel-collector/config.yaml`，需 `extra_hosts: host.docker.internal:host-gateway`。
- trace 仅在跨服务 gRPC 调用产生 span；UI：Jaeger 16686 / Prometheus 9090 / Loki 3100。

## search 服务 ES 集成
- `SearchCourses` / `GetTopCoursesByCategory` 走 ES（索引 `course`，analyzer 默认 `cjk`，改 analyzer 需 `DELETE /course` 重建）。
- 全量：`search` RPC `ReindexCourses` + `course` RPC `CourseSearchIndexInfoList`；`initReindex` 在 `NewServiceContext` 启动后台异步跑一次。
- 增量：course 的 `svc.Producer` 在上/下架成功后发 `event.CourseEvent` 到 `mq.ExchangeCourse`（routing key `course.up`/`course.down`），search 消费后 upsert/delete。best-effort：发布失败仅告警不回滚；`Producer=nil` 时直接 return。`pkg/mq.Producer` 启动一次性建连，断线不自动重连，RabbitMQ 须先就绪。

## 实现状态
- 13/13 服务 API 193 + RPC 201 = 394 接口，394 个 logic 文件 1:1 落地，0 处 TODO 占位；28 模块 `go build ./...` 全通过。完成度 ≈100%，功能完备度 ≈97%。
- 有意桩/未接线：media 对象存储 mock；pay 支付回调 URL 占位 + 未接真实网关；trade 优惠券未接入 promotion；course 未接 user(教师)/learning(课时)；learning section_type 忽略。
- 已接线跨服务 RPC：trade→{course,pay}、search→course、learning→course。

## 批量化实现 stub 服务的可复用流程（course 已验证）
1. 先补 custom model（生成物只有 CRUD）；2. 建「黄金标准」域手写实现作样板；3. 写 `IMPL_CONTRACT.md`（方法清单/字段映射/helper/样板/缺口/硬规则）作并行子代理唯一事实源（结束即删）；4. 按域拆分子代理，**严禁子代理跑 go build**；5. 主代理统一编译 `go build ./...`（rpc+api 均 rc=0）；6. 抽样核对 + 同步文档（更新 `openspec/specs/services/<svc>/spec.md` 中对应的 business-rules Requirement）。

## 文档体系（2026-09-08 完成迁移，OpenSpec 为唯一事实源）
- `openspec/`：**共 24 份 spec，`openspec validate --all` 24/24 通过**。结构：`specs/services/<svc>`（13）、`specs/architecture/service-topology`、`specs/conventions/{goctl-codegen,api-contracts,code-style,git-workflow}`、`specs/shared/{error-codes,mq-events,pkg-contracts}`、`specs/infra/{observability,local-env,local-dev-workflow}`；`config.yaml` 存技术栈/端口/契约浓缩；`changes/` 走 spec-driven 变更提案（当前为空）。
- 源 `.cursor/repowiki/`（80 篇）已全量迁移并**删除目录**（80 篇均已在 git 中，可 `git checkout` 恢复）。其中 `06-status/implementation-status.md` 按用户决定不迁移、随目录删除，其 20+ 已知缺口已分散保留在各服务 spec 的「> 已知缺口：」引用块。
- 迁移后约定：新需求走 `openspec new change "<name>"`；文档引用一律指向 `openspec/specs/...`，不再引用任何 wiki 路径。
- openspec CLI 1.12.0 已全局可用（`openspec validate --specs --strict` / `--all`）。spec 格式：`## Purpose` + `## Requirements` + `### Requirement:`（至少一个 `#### Scenario:`，WHEN/THEN/AND 加粗），正文保留 SHALL。

## MySQL（开发）
127.0.0.1:3306 root/0000，库 `tj_<domain>`。DDL `sql/ddl/`，迁移 `sql/migration/`（改表两处都要同步）。
