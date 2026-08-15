# 天机学堂 tjxt 项目记忆

## 架构
- go-zero v1.10.3 微服务，**go-zero 官方标准结构**：go.work 工作区聚合 **28 个 use 模块**（根模块 `.` + `pkg` + 13 服务各含 `api`+`rpc` 共 26 个服务模块）。
- 每个可独立部署单元（api 进程 / rpc 进程）各有独立 `go.mod`：**2026-08-06 由「每服务单 module」重构为官方标准**，删除原聚合 `apps/<svc>/go.mod` 与孤儿 `apps/data/api/go.mod`。
- module path 规则：`module tjxt/apps/<svc>/api`、`module tjxt/apps/<svc>/rpc`；data 嵌套为 `module tjxt/apps/data/api/data`、`module tjxt/apps/data/rpc/data`。
- 公共库 `tjxt/pkg`：每个服务模块 `go.mod` 用相对 `replace` 引用——普通服务三级 `replace tjxt/pkg => ../../../pkg`，data 嵌套四级 `replace tjxt/pkg => ../../../../pkg`。
- 跨服务依赖：在调用方 `go.mod` 中 `require`（伪版本 `v0.0.0-00010101000000-000000000000`）+ `replace` 到本地相对路径（如 `trade`：`tjxt/apps/course/rpc => ../../course/rpc`、`tjxt/apps/pay/rpc => ../../pay/rpc`）。
- **IDE 报错根因（已随重构消除）**：原聚合模式下子目录无 go.mod，gopls 需向上找父 module，且 `go.work` 的隐式注入不写回各 go.mod，导致 Cursor 报「tjxt/pkg is not in your go.mod」。现每 api/rpc 独立 go.mod 显式 declare 依赖，重启 gopls 即无红波浪线。
- 命名风格 gozero（无下划线），goctl 默认即此风格。

## 服务清单与端口（2026-08-06 统一重排为连续编号，无冲突）
API 端口 8801–8813 连续唯一；RPC 端口 8081–8093 连续唯一。改端口不影响调用方（RPC 走 etcd Key 发现，如 `course.rpc`）。

| 服务 | API | RPC(ListenOn) | etcd Key | DB |
|------|-----|-----|-----|----|
| user | 8801 | 8081 | user.rpc | tj_user(user/user_detail) |
| auth | 8802 | 8082 | auth.rpc | tj_user |
| course | 8803 | 8083 | course.rpc | tj_course |
| learning | 8804 | 8084 | learning.rpc | tj_learning |
| exam | 8805 | 8085 | exam.rpc | tj_learning |
| media | 8806 | 8086 | media.rpc | tj_media(file/media) |
| message | 8807 | 8087 | message.rpc | tj_message |
| pay | 8808 | 8088 | pay.rpc | tj_pay |
| trade | 8809 | 8089 | trade.rpc | tj_trade |
| search | 8810 | 8090 | search.rpc | tj_search |
| data | 8811 | 8091 | data.rpc | tj_data |
| promotion | 8812 | 8092 | promotion.rpc | tj_promotion |
| remark | 8813 | 8093 | remark.rpc | tj_remark |

- data 为嵌套模块，配置在 `apps/data/api/data/etc/data-api.yaml`、`apps/data/rpc/data/etc/data.yaml`（不是 `apps/data/api/etc/`）。
- 端口在 26 个 `apps/**/etc/*.yaml` + README.md + `.cursor/repowiki/**/*.md` + `docs/**/*.md` 中均已同步；批量改端口时务必用「旧→新映射的单趟 re.sub」，因新旧端口集合重叠，逐条回扫会连锁错改。

## go-zero 开发规范（用户强约束）
- 优先 goctl 生成，只手写 logic/custom model/业务扩展。禁止手写骨架/handler/types/routes/pb。
- 开发顺序：DB表 → Model(goctl) → Proto → RPC生成 → .api → API生成 → Logic。
- Model 已存在则扩展 custom model，不重新生成覆盖。
- goctl .api 语法不支持 `any`，用 `interface{}` 代替。
- 统一响应：handler 用 `result.Write(w,r,data,err)`（pkg/response.R{Code,Msg,RequestId,Data any}），不用生成的 Result 类型。
- 统一分页：PageRequest{PageNo,PageSize}。
- JWT：@server jwt:Auth。

## MySQL（开发环境）
host 127.0.0.1 port 3306 user root pass 0000，库 tj_<domain>。

## 关键文件位置
- DDL：sql/ddl/tj_<domain>.sql（纯 DDL，供 goctl model 用）。
- 迁移：sql/migration/tj_<domain>.sql（含数据）。
- zero-skills：.cursor/skills/zero-skills/ 和 .opencode/skills/zero-skills/（go-zero AI 最佳实践）。
- Model 放 `apps/<svc>/rpc/internal/model/`（custom `*model.go` + 生成 `*model_gen.go` + vars.go）。
- RPC 客户端在 `apps/<svc>/rpc/<svc>/<svc>.go`（goctl 1.10 命名）。

## goctl 命令（本机已验证）
- API: `cd apps/X/api && goctl api go --api X.api --dir . --style gozero`
- RPC: `cd apps/X/rpc && goctl rpc protoc X.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style gozero`
- Model: `goctl model mysql ddl -src sql/ddl/tj_X.sql -dir apps/X/rpc/internal/model -cache --style gozero`
- goctl api go 不覆盖已有 handler/logic/svc/config，仅刷新 types.go/routes.go。
- goctl rpc protoc 不覆盖已有 logic/svc/config，仅刷新 pb/server，client 默认生成(-c true)。

## 服务实现状态（2026-08-06 实测：go build 全模块通过 + 逻辑文件清点）
- 13/13 服务 API+RPC 契约 100% 落地（API 193 + RPC 201 = 394 接口，394 个 logic 文件 1:1 对应，全代码库 0 处 TODO/panic 占位）。
- 跨服务 RPC 依赖（已接线、编译通过）：trade→{course,pay}、search→course、learning→course。
- 事件驱动：trade 经 RabbitMQ 发布领域事件，Producer 优雅降级（MQ 不可用时为 nil 不阻塞启动）。
- 代码注释标明的「应有但未接线」集成：course→user(教师)、course→learning(课时)、trade→promotion(优惠券)、pay→真实支付网关。
- 待优化（有意桩，非未实现）：media 对象存储为 mock（COS/OSS 未接入）；pay 支付回调 URL 为 demo 占位；trade 优惠券未接入；course 教师/课时未接线；learning section_type 忽略。
- 结构异常（已解决，2026-08-06）：原孤儿 `apps/data/api/go.mod` 已删除，data 现仅 `apps/data/api/data` 与 `apps/data/rpc/data` 两个真实模块，与其余服务同构（仅深一级）。

| 服务 | API | RPC | 编译 | 备注 |
|------|-----|-----|:----:|------|
| auth | 18/18 | 19/19 | ✅ | 已实现 |
| user | 13/13 | 15/15 | ✅ | 已实现 |
| course | 38/38 | 38/38 | ✅ | 已实现；教师(user)/课时(learning)未接线 |
| media | 10/10 | 10/10 | ✅ | 已实现；对象存储为 mock |
| learning | 9/9 | 11/11 | ✅ | 已实现；section_type 占位 |
| exam | 7/7 | 7/7 | ✅ | 已实现 |
| search | 4/4 | 4/4 | ✅ | 已实现；依赖 course |
| trade | 36/36 | 37/37 | ✅ | 已实现；依赖 course,pay；优惠券未接入 |
| pay | 16/16 | 17/17 | ✅ | 已实现；支付网关占位 URL |
| promotion | 16/16 | 16/16 | ✅ | 已实现 |
| message | 18/18 | 19/19 | ✅ | 已实现 |
| remark | 2/2 | 2/2 | ✅ | 已实现 |
| data | 6/6 | 6/6 | ✅ | 已实现（模块 apps/data/api/data、apps/data/rpc/data 编译通过）|
合计 394/394（100%）逻辑文件落地，28 个 use 模块（根+pkg+26 服务模块）独立 `go build ./...` 全部通过。整体代码实现完成度 ≈100%，功能完备度 ≈97%（少量外部集成桩）。完整报告见 docs/completion-analysis.md。

## 批量化实现 stub 服务的可复用工作流（course 已验证）
对「goctl 已生成骨架、logic 全为 todo 占位」的服务，用此流程高效落地：
1. **先补 Model 层**：生成 model 只有 CRUD；业务需 List/Count/Page/FindByX/Upsert/级联 Delete，写 custom `*model.go` 扩展（不覆盖 `_gen.go`）。
2. **致命坑**：自定义 model 嵌入 `CachedConn`，**只能用** `m.QueryRowsNoCacheCtx`/`QueryRowNoCacheCtx`/`ExecNoCacheCtx`；`m.QueryRowsCtx` 被遮蔽会编译失败。
3. **建黄金标准**：挑一个完整域（如 course 的 category，8 RPC+8 API）手写实现，作为质量样板。
4. **写 IMPL_CONTRACT.md**：汇总「可调用方法清单 / 字段类型映射 / 共享 helper / 参考样板文件 / 已知缺口 / 硬性规则」，作为并行子代理的唯一事实源（任务结束即删）。
5. **并行子代理**：按域拆分（如 5 RPC + 4 API 共 9 个 general-purpose 子代理），每个只改非重叠文件集；**严禁子代理跑 go build**（并发写导致冲突），统一由主代理编译修复。
6. **主代理统一编译**：`cd apps/X/rpc && go build ./...` 与 `cd apps/X/api && go build ./...` 必须 rc=0。
7. **抽样核对 + 文档同步**：抽查 2–3 个关键 logic 确为真实实现（非 stub）；同步 repowiki business-rules.md 与 implementation-status.md。
8. 跨域字段（销量/评分/老师详情来自其他服务）当前服务无列则填 0/留空，并在文档「已知缺口」标注未接线的 RPC client。

## goctl 代码生成关键坑（2026-08-15 search/course 重踩）
- **`goctl rpc protoc` 会回填桩并可能错位接线**：重跑 proto 生成时，goctl 在 `internal/logic/`（扁平、包 `logic`）重新生成 todo 桩，并把 `internal/server/<svc>/<svc>server.go` 写成 import 扁平 `internal/logic`（包 `logic`）。但本项目 search 服务的**真实 logic 在 `internal/logic/search/`（包 `searchlogic`）**，`<svc>server.go` 也 import 这个子包。
- **致命后果**：若把真实 logic 留在 `searchlogic` 子包、却让 server 指向扁平 `logic` 桩，搜索会「编译通过但返回空结果」（SearchCourses 返回空 `CourseSearchPageReply`）。判别依据：`grep -n "internal/logic" internal/server/*/<svc>server.go` 看 server 到底 import 哪个包。
- **正确做法**：重跑 goctl 后，确认 server import 的包（本例 `internal/logic/search`）与真实 logic 文件所在包一致；若 goctl 在 `internal/logic/` 落了同名桩文件，**删掉扁平桩，保留 `search/` 子包的真实实现**。不要为了迎合 goctl 把 `searchlogic` 拍平到 `internal/logic/`（下次 goctl 又会错位）。
- 同理 `internal/server/` 下可能出现多余扁平 `<svc>server.go`，main（`<svc>.go`）只 import `internal/server/<svc>` 子包，多余扁平文件直接删。

## search 服务 ES 集成现状（2026-08-15 增强）
- 课程搜索 `SearchCourses` / `GetTopCoursesByCategory` 走 ES（`svc.CourseDoc` + `courseIndexMapping`，索引名 `course`）。
- **全量重建索引**：新增 `search` RPC `ReindexCourses(Empty)→(indexed,total)`，并新增 `course` RPC `CourseSearchIndexInfoList(page_no,page_size)→(total,items)`（复用 `CourseModel.PageQuery(status=上架)`，避免 N+1）。search 侧 `svc.ReindexAll` 用 `esutil.BulkIndexer` 批量 upsert（以 courseId 为文档 ID，幂等）；`initReindex` 在 `NewServiceContext` 启动时后台异步跑一次（best-effort，5min 超时），不必依赖 RabbitMQ 事件即可被搜。
- **中文分词**：默认 analyzer 改为 `cjk`（ES 内置、无需插件、二元分词，优于 standard 单字）；`config.Elasticsearch.SearchAnalyzer` 可选（IK 场景设 `ik_smart`）。改 analyzer 需删旧 `course` 索引重建（`DELETE /course`）。
- 改 proto 后务必 `export PATH="$PATH:/d/Program Files/protoc-35.0-win64/bin"` 再 `goctl rpc protoc`（protoc 在带空格路径，goctl 1.10.1）。
- **增量同步已打通（2026-08-15 补齐 Producer）**：course 服务 `svc.Producer`（`*mq.Producer`，`initProducer` 在 `NewServiceContext` 中 nil 容错创建）在 `publishCourse`（上架）成功末尾、`downCourse`（下架实际状态变更）成功末尾调用 `svcCtx.PublishCourseEvent(ctx,id,up)` → 发布 `event.CourseEvent{CourseID}` 到 `mq.ExchangeCourse`（course.events），routing key `course.up`/`course.down`。search 侧消费者（`svc.MQClient.Start` 在 `search.go` 启动 goroutine）消费后 upsert/delete ES 文档。四条入口（CourseUpShelf/CourseUp 批量/CourseDownShelf/CourseDown 批量）经两个共享 helper 全覆盖。
- **发布语义（best-effort）**：`PublishCourseEvent` 不返回错误——课程主流程（DB 写入）已成功，MQ 发布失败仅 `logx.Errorf` 告警、不回滚课程操作。恢复兜底：search 启动全量 `ReindexAll` + 手动 `ReindexCourses` RPC。MQ 配置缺失/连接不可用时 `Producer=nil`，`PublishCourseEvent` 直接 return（不阻塞课程启动）。
- **关键约束**：`pkg/mq.Producer` 启动时 `amqp091.Dial` 一次性建连，断线不自动重连（与 search 消费者一致）；RabbitMQ 须在 course/search 启动前就绪，否则该侧增量通道空转（待连重启恢复）。course 的 `course.yaml` 已加 `RabbitMQ{Host:127.0.0.1,Port:5672,User/Pass:rabbitmq}`。
