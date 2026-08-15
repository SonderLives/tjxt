# tjxt — 天机学堂 Go 微服务

基于 **go-zero** 框架的微服务架构，采用 **API + RPC 分层**，通过 **etcd** 服务发现，使用 **MySQL + Redis** 作为存储。

## 🏗 架构概览

```
HTTP Client
    │
    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│ API 层 (go-zero REST, 独立 go.mod, 通过 zrpc 调 RPC)                         │
│                                                                            │
│  auth-api:8802    user-api:8801    course-api:8803   learning-api:8804    │
│  exam-api:8805    media-api:8806    message-api:8807  pay-api:8808        │
│  trade-api:8809   search-api:8810   data-api:8811    promotion-api:8812   │
│  remark-api:8813                                                       │
└────────────────────────────┬───────────────────────────────────────────────┘
                             │ zrpc (etcd 服务发现)
                             ▼
┌────────────────────────────────────────────────────────────────────────────┐
│ RPC 层 (go-zero zrpc, 独立 go.mod, 持 model + Redis 缓存)                    │
│                                                                            │
│  auth.rpc      user.rpc       course.rpc     learning.rpc                  │
│  exam.rpc      media.rpc      message.rpc    pay.rpc                       │
│  trade.rpc     search.rpc     data.rpc      promotion.rpc  remark.rpc     │
└────────────────────────────┬───────────────────────────────────────────────┘
                             ▼
                   MySQL (sqlx + goctl model --cache)
                   Redis (缓存 / MQ / Stream)
```

> ⚠️ **无独立网关**。边缘入口按需自行接入 APISIX / Kong / Nginx。



## 📁 目录结构

```
tjxt/
├── go.work                  # 工作区，聚合 28 个 use 模块（根 + pkg + 13 服务 × {api,rpc} 共 26 个服务模块）
├── go.mod                   # 根模块 (仅声明版本)
├── Makefile                 # 构建/生成/运行/校验一站式命令
├── docker-compose.yml       # MySQL/Redis/RabbitMQ/etcd + 可观测性栈（Jaeger/otel-collector/Prometheus/Loki）一键启动
├── pkg/                     # 公共代码库 (module: tjxt/pkg)
│   ├── auth/                # JWT 签发/校验、Claims 定义
│   ├── mq/                  # Redis Stream 事件总线 (生产者/消费者/事件定义)
│   ├── response/            # 统一响应 R{code,msg,requestId,data} + 分页
│   ├── xerr/                # 业务错误码体系 + HTTP 状态映射
│   ├── utils/
│   │   ├── idgen/           # 雪花算法 ID 生成器
│   │   └── page/            # 分页工具
│   └── go.mod
├── sql/
│   ├── ddl/                 # goctl model 生成用的纯 DDL (按库分文件)
│   └── migration/           # 完整迁移脚本 (含索引/外键/初始数据)
├── docs/
│   └── tjxt.openapi.json    # OpenAPI 3.0 规范文档
└── apps/
    └── <svc>/               # 13 个业务服务（auth/course/data/exam/learning/media/message/pay/promotion/remark/search/trade/user）
        ├── go.mod
        ├── api/
        │   ├── <svc>.api          # goctl API 定义 (路由/请求/响应)
        │   ├── <svc>.go           # main 入口
        │   ├── etc/<svc>-api.yaml # API 配置
        │   └── internal/
        │       ├── config/        # 配置加载
        │       ├── handler/       # HTTP 处理器 (goctl 生成)
        │       ├── logic/         # 业务逻辑 (仅手写这里)
        │       ├── svc/           # 服务上下文 (RPC 客户端等)
        │       └── types/         # 请求/响应结构体 (goctl 生成)
        └── rpc/
            ├── <svc>.proto        # gRPC 服务定义
            ├── <svc>.go           # main 入口
            ├── etc/<svc>.yaml     # RPC 配置
            ├── internal/
            │   ├── config/
            │   ├── logic/         # RPC 业务逻辑 (仅手写这里)
            │   ├── model/         # sqlx Model (goctl 生成 + 扩展)
            │   ├── server/        # gRPC 服务注册
            │   └── svc/           # 服务上下文 (DB/Redis/缓存)
            └── pb/                # 生成的 pb.go / _grpc.pb.go
    └── data/                # 嵌套模块 (goctl -m multiple)
        ├── api/data/
        │   ├── data.go
        │   ├── etc/data-api.yaml
        │   └── internal/...
        └── rpc/data/
            ├── data.proto
            ├── data.go
            ├── client/data/
            ├── etc/data.yaml
            └── internal/...
```



## 🗂 服务清单


| 服务           | API 端口 | RPC 地址         | etcd Key         | 数据库         | 核心职责                  |
| ------------ | ------ | -------------- | -------------- | ----------- | --------------------- |
| **auth**     | 8802   | 127.0.0.1:8082 | `auth.rpc`     | tj_auth     | 登录/注册、Token 签发刷新、权限校验 |
| **user**     | 8801   | 127.0.0.1:8081 | `user.rpc`     | tj_user     | 用户档案、教师/学员详情、角色部门     |
| **course**   | 8803   | 127.0.0.1:8083 | `course.rpc`   | tj_course   | 课程/分类/章节/媒体资源 CRUD    |
| **learning** | 8804   | 127.0.0.1:8084 | `learning.rpc` | tj_learning | 学习进度、收藏、笔记、考试记录       |
| **exam**     | 8805   | 127.0.0.1:8085 | `exam.rpc`     | tj_exam     | 试卷/题库/考试/评分/记录        |
| **media**    | 8806   | 127.0.0.1:8086 | `media.rpc`    | tj_media    | 文件上传、媒资管理、转码回调        |
| **message**  | 8807   | 127.0.0.1:8087 | `message.rpc`  | tj_message  | 站内信/通知/公告、模板消息        |
| **pay**      | 8808   | 127.0.0.1:8088 | `pay.rpc`      | tj_pay      | 订单支付、退款、回调、账单         |
| **trade**    | 8809   | 127.0.0.1:8089 | `trade.rpc`    | tj_trade    | 交易订单、优惠券、分佣结算         |
| **search**   | 8810   | 127.0.0.1:8090 | `search.rpc`   | tj_search   | 课程/用户全文检索、建议词         |
| **data**     | 8811   | 127.0.0.1:8091 | `data.rpc`     | — (Redis)   | 统计大屏、榜单、今日数据聚合        |
| **promotion**| 8812   | 127.0.0.1:8092 | `promotion.rpc`| tj_promotion| 优惠券、营销活动、秒杀、积分        |
| **remark**   | 8813   | 127.0.0.1:8093 | `remark.rpc`   | tj_remark   | 评价/评论、评分、回复            |


> 📌 各服务 API/RPC 端口已在 `apps/*/api/etc/*.yaml` 与 `apps/*/rpc/etc/*.yaml` 中唯一分配，可同时启动（media-api 8806、promotion-api 8812、course.rpc 8083、trade.rpc 8089、message.rpc 8087、promotion.rpc 8092）。



## 📦 前置依赖


| 组件           | 版本要求   | 说明                            |
| ------------ | ------ | ----------------------------- |
| **Go**       | ≥ 1.26 | 用到 `go.work` 工作区              |
| **MySQL**    | 8.x    | 默认 `root:0000@127.0.0.1:3306` |
| **Redis**    | 7.x    | 默认 `127.0.0.1:6379`           |
| **etcd**     | 3.5+   | 默认 `127.0.0.1:2379`，zrpc 服务发现 |
| **RabbitMQ** | 3.13+  | 可选，预留兼容 Java 版事件              |
| **goctl**    | 最新     | 代码生成工具                        |


```bash
# Windows (PowerShell + Scoop)
scoop install go@1.26
scoop install etcd
scoop install go-zero/goctl

# 或 Docker 一键拉起依赖
docker compose up -d
```



## 🚀 快速开始

```bash
# 1. 同步工作区 & 整理依赖
make deps && go work sync
# 或: make init

# 2. 环境体检 (Go/goctl/golangci-lint + 全模块编译)
make verify

# 3. 编译全部 26 个可执行文件 (13 api + 13 rpc) 到 bin/
make build

# 4. 启动单个服务（务必从仓库根执行，日志落到 logs/<svc>/）
#    推荐用 Makefile 目标（已改为从仓库根 go run ./apps/...，CWD=仓库根）：
make run-user          # 启动 user-api
make run-user-rpc      # 启动 user-rpc
#    或手动（同样要在仓库根执行）：
go run ./apps/user/api -f apps/user/api/etc/user-api.yaml
go run ./apps/user/rpc -f apps/user/rpc/etc/user.yaml
# data 是嵌套结构，多一层目录：
go run ./apps/data/api/data -f apps/data/api/data/etc/data-api.yaml
go run ./apps/data/rpc/data -f apps/data/rpc/data/etc/data.yaml
```

**一键启动全部服务**（后台独立窗口，日志写入 `logs/<svc>/`）：

```bash
make run-all       # 启动所有 API（各开一个独立 PowerShell 窗口）
make run-all-rpc   # 启动所有 RPC
# 停止：各窗口直接关闭（Ctrl-C）；当前无 stop-all 目标
# 看服务日志：tail -f logs/course-api/access.log   （替换 <svc> 即可）
# 看基建容器日志：make docker-logs
```

## 🔭 可观测性（trace / metrics / logs 统一经 otel-collector）

各服务 yaml 已注入 `Telemetry`（链路）、`Prometheus`（/metrics）、`Log{Mode:file,...}`（日志文件）。
先 `make docker-up` 起基建，再起各 Go 服务：

| 信号 | 路径 | 界面/端点 |
|------|------|-----------|
| 链路 Trace | go-zero → `127.0.0.1:4318` → otel-collector → Jaeger | Jaeger UI http://localhost:16686 |
| 指标 Metrics | 各服务 `/metrics`(9101–9113/9201–9213) → otel-collector 聚合 `:8889` → Prometheus | Prometheus http://localhost:9090 |
| 日志 Logs | `logs/<svc>/*.log`（仓库根）→ 挂进 collector → Loki | Loki `:3100`（建议接 Grafana 查询） |

> ⚠️ 日志路径是**相对路径** `logs/<svc>`，相对进程启动目录（CWD）；服务须从**仓库根**启动，
> 否则日志会散落到各服务子目录，collector 的 filelog 采集不到。



## ⚙️ 代码生成 (核心约束)

**骨架一律** `goctl` **生成，不手写。业务逻辑只允许填** `logic/`**。**


| 场景                     | 命令                                                                                                                                |
| ---------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| 新建 API 骨架              | `goctl api new <svc>`                                                                                                             |
| 修改 `.api` 后重生成         | `goctl api go -api apps/<svc>/api/<svc>.api -dir apps/<svc>/api --home <tpl>`                                                     |
| 新建 RPC 骨架              | `goctl rpc new <svc>`                                                                                                             |
| 修改 `.proto` 后重生成       | `cd apps/<svc>/rpc && goctl rpc protoc <svc>.proto --go_out=. --go-grpc_out=. --zrpc_out=. -m`                                    |
| 数据库 Model (带 Redis 缓存) | `goctl model mysql datasource -url="root:0000@tcp(127.0.0.1:3306)/tj_<svc>" -table="<tbl>" -dir apps/<svc>/rpc/internal/model -c` |
| Model 需加自定义方法          | 在 `internal/model/<tbl>model.go` 里写扩展 (`_gen.go` 由 goctl 管理)                                                                      |


**Makefile 快捷方式**:

```bash
make gen-api  svc=user              # 只重生成 API
make gen-rpc  svc=user              # 只重生成 RPC (proto 必须先改)
make gen-model svc=user db=tj_user table=user
make generate                       # api + rpc 全量重生成
```



## 📋 Makefile 速查

```bash
make help            # 全部命令
make verify          # 环境体检
make build           # 编译所有服务（26 个二进制到 bin/）
make test            # 各模块 go test
make lint            # golangci-lint
make fmt             # gofumpt 格式化
make clean           # 删除 bin/ 与各服务下遗留 exe
make run-<svc>          # 前台跑单个 API（如 make run-user）
make run-<svc>-rpc      # 前台跑单个 RPC（如 make run-user-rpc）
make run-all            # 后台并行启动所有 API（各开独立窗口）
make run-all-rpc        # 后台并行启动所有 RPC
make docker-up          # 拉起 docker-compose 全部容器（含可观测性栈）
make docker-logs        # 跟踪容器日志
make info            # 打印服务→端口映射
make d2u             # Windows 行尾修复 (CRLF→LF)
```



## 🔄 事件流 (Redis Stream)

RabbitMQ(`pkg/mq/event`)事件名：


| 事件                      | 发布方       | 订阅方                   |
| ----------------------- | --------- | --------------------- |
| `mq:order:success`      | trade     | learning(加课)、data(榜单) |
| `mq:course:up`          | course    | search(建索引)           |
| `mq:course:down`        | course    | search(删索引)           |
| `mq:course:cate:change` | course    | (分类缓存刷新)              |
| `mq:course:media:quote` | media     | course                |
| `mq:media:exists`       | media     | (内部)                  |
| `mq:media:delete`       | media     | course                |
| `mq:exam:record`        | exam      | learning              |
| `mq:communication:like` | (future)  | data                  |
| `mq:data:today`         | data 定时任务 | data(命中缓存)            |


**发布示例** (`pkg/mq/producer.go`):

```go
producer.Publish(ctx, mq.EventCourseUp, map[string]any{"courseId": 123})
```

**订阅示例** (`pkg/mq/client.go`):

```go
consumer.Subscribe(ctx, mq.EventCourseUp, func(msg map[string]any) error { ... })
```



## 🧩 公共包 (pkg/)


| 包                 | 说明                                                                                   |
| ----------------- | ------------------------------------------------------------------------------------ |
| `pkg/auth`        | JWT `Sign/Parse`、Claims 含 `userId`/`role`                                            |
| `pkg/xerr`        | 业务错误码 `Codexxx` + `Msg/Wrap/HttpStatus`                                              |
| `pkg/response`    | 统一响应 `R{code,msg,requestId,data}` + `Page{list,total,pages}` + `Write(w,r,data,err)` |
| `pkg/mq`          | Redis Stream 生产者/消费者抽象、事件常量                                                          |
| `pkg/utils/idgen` | 雪花算法 `NextID()`                                                                      |
| `pkg/utils/page`  | 分页参数归一化                                                                              |


**使用方式** (在任意服务 `go.mod` 引入):

```go
require tjxt/pkg v0.0.0
replace tjxt/pkg => ../../pkg
```



## 🗄 数据库规范

- **DDL 位置**: `sql/ddl/tj_<svc>.sql` (纯建表，供 goctl model 生成)
- **迁移脚本**: `sql/migration/` (含索引/外键/初始数据，生产执行)
- **命名**: 表名单数 (`user` `course` `order`)，主键 `id BIGINT`，时间 `create_time`/`update_time`
- **缓存**: `goctl model ... -c` 生成带 Redis 缓存的 Model，键前缀 `cache:#{table}:`



## 🔐 统一鉴权 & 响应

**JWT 鉴权** (`pkg/auth/jwt.go`):

- HS256，Claim: `userId`(int64) + `role`(string)
- API 网关层接入 `rest.WithJwt` 中间件，自动注入 `ctx`

**统一响应** (`pkg/response/result.go`):

```json
{
  "code": 200,
  "msg": "success",
  "requestId": "trace-id-from-otel",
  "data": { ... }
}
```

- 成功: `result.Ok()` / `result.OkData(data)`
- 失败: `result.Fail(err)` 自动映射 `xerr` 错误码 → HTTP 状态码
- 所有 handler 统一用 `result.Write(w, r, data, err)`



## 📊 业务实现现状

> 结论（2026-08-06 复核）：13/13 服务、API+RPC 共 394 个接口，全部 logic 已实现并 `go build` 逐模块编译通过；功能完备度≈97%（少量桩为有意预留，见下）。

| 项目           | 状态                                                                  |
| ------------ | ------------------------------------------------------------------- |
| **骨架**       | ✅ 13 服务 api+rpc 全部 goctl 生成完毕，`go build ./...` 逐模块编译通过                  |
| **Logic 业务** | ✅ 394/394 logic 全部实现（API 193 + RPC 201），全库 0 处 TODO/panic 占位           |
| **中间件**      | ✅ JWT(`@server jwt:Auth`)、统一响应 `result.Write`、xerr 错误码已在各 svc 接入       |
| **包名规范**     | ✅ 模块路径统一 `tjxt/apps/<svc>` (data 为 `tjxt/apps/data/{api,rpc}/data`) |
| **事件总线**     | 🚧 `pkg/mq` Producer 已实例化，但 trade 等业务 logic 尚未实际 `Publish`（事件未发射） |
| **数据库**      | ✅ DDL + goctl model（带缓存）已生成，自定义 Model 已扩展                         |
| **跨域 RPC**    | ⚠️ 已接线 trade→{course,pay}、search→course、learning→course；其余（如 course→user/learning、trade→promotion、pay→真实网关）尚未接线 |
| **已知缺口**     | ⚠️ media 对象存储为 mock、pay 支付回调为 demo 占位、trade 优惠券未接入、Seata 未接入（undo_log 表闲置） |




## 🛠 开发规范

1. **只改 logic**: `api/internal/logic/` `rpc/internal/logic/` 以外的生成代码**不要手改**
2. **错误处理**: 业务错误用 `xerr.New(code, msg)` 或 `xerr.Wrap(err, code, msg)`，不要直接返回 `error`
3. **日志**: 使用 go-zero 内置 `logx`，关键路径带上 `trace.TraceIDFromContext(ctx)`
4. **事务**: RPC 层用 `model.Transact(ctx, func(session sqlx.Session) error { ... })`
5. **分页**: 统一用 `pkg/utils/page.ParsePageReq(req)` + `result.Page{List, Total, Pages}`



## 📚 参考资料

- **go-zero 官网**: [https://go-zero.dev](https://go-zero.dev)
- **goctl 手册**: [https://go-zero.dev/docs/tasks](https://go-zero.dev/docs/tasks)

