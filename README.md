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
│  auth-api:8814    user-api:8808    course-api:8812   learning-api:8813    │
│  exam-api:8815    media-api:8811     message-api:8816  pay-api:8811        │
│  trade-api:8810   search-api:8817    data-api:8818                         │
└────────────────────────────┬───────────────────────────────────────────────┘
                             │ zrpc (etcd 服务发现)
                             ▼
┌────────────────────────────────────────────────────────────────────────────┐
│ RPC 层 (go-zero zrpc, 独立 go.mod, 持 model + Redis 缓存)                    │
│                                                                            │
│  auth.rpc      user.rpc       course.rpc     learning.rpc                  │
│  exam.rpc      media.rpc      message.rpc    pay.rpc                       │
│  trade.rpc     search.rpc     data.rpc                                     │
└────────────────────────────┬───────────────────────────────────────────────┘
                             ▼
                   MySQL (sqlx + goctl model --cache)
                   Redis (缓存 / MQ / Stream)
```

> ⚠️ **无独立网关**。边缘入口按需自行接入 APISIX / Kong / Nginx。



## 📁 目录结构

```
tjxt/
├── go.work                  # 工作区，聚合 13 个模块 (12 服务 + pkg)
├── go.mod                   # 根模块 (仅声明版本)
├── Makefile                 # 构建/生成/运行/校验一站式命令
├── docker-compose.yml       # MySQL/Redis/RabbitMQ/etcd 一键启动
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
    └── <svc>/               # 11 个业务服务 + data 嵌套服务
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


| 服务           | API 端口 | RPC 地址         | etcd Key       | 数据库         | 核心职责                  |
| ------------ | ------ | -------------- | -------------- | ----------- | --------------------- |
| **auth**     | 8814   | 127.0.0.1:8083 | `auth.rpc`     | tj_auth     | 登录/注册、Token 签发刷新、权限校验 |
| **user**     | 8808   | 127.0.0.1:8082 | `user.rpc`     | tj_user     | 用户档案、教师/学员详情、角色部门     |
| **course**   | 8812   | 127.0.0.1:8083 | `course.rpc`   | tj_course   | 课程/分类/章节/媒体资源 CRUD    |
| **learning** | 8813   | 127.0.0.1:8085 | `learning.rpc` | tj_learning | 学习进度、收藏、笔记、考试记录       |
| **exam**     | 8815   | 127.0.0.1:8084 | `exam.rpc`     | tj_exam     | 试卷/题库/考试/评分/记录        |
| **media**    | 8811   | 127.0.0.1:8087 | `media.rpc`    | tj_media    | 文件上传、媒资管理、转码回调        |
| **message**  | 8816   | 127.0.0.1:8085 | `message.rpc`  | tj_message  | 站内信/通知/公告、模板消息        |
| **pay**      | 8811   | 127.0.0.1:8081 | `pay.rpc`      | tj_pay      | 订单支付、退款、回调、账单         |
| **trade**    | 8810   | 127.0.0.1:8084 | `trade.rpc`    | tj_trade    | 交易订单、优惠券、分佣结算         |
| **search**   | 8817   | 127.0.0.1:8086 | `search.rpc`   | — (ES)      | 课程/用户全文检索、建议词         |
| **data**     | 8818   | 127.0.0.1:8088 | `data.rpc`     | — (Redis)   | 统计大屏、榜单、今日数据聚合        |


> ⚠️ **端口冲突提示** (骨架阶段先这么跑，后续统一分配)：
>
> - API: `media-api` 与 `pay-api` 撞 `8811`
> - RPC: `exam.rpc`/`trade.rpc` 撞 `8084`，`learning.rpc`/`message.rpc` 撞 `8085`，`auth.rpc`/`course.rpc` 撞 `8083`
> - **同时启动请先修改对应 yaml 端口**



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

# 3. 编译全部 22 个可执行文件 (11 api + 11 rpc) 到 bin/
make build

# 4. 启动单个服务 (新开终端)
cd apps/user/api && go run user.go -f etc/user-api.yaml
cd apps/user/rpc && go run user.go -f etc/user.yaml

# data 是嵌套结构，多一层目录
cd apps/data/api/data && go run data.go -f etc/data-api.yaml
cd apps/data/rpc/data && go run data.go -f etc/data.yaml
```

**一键启动全部服务** (后台独立窗口，日志写入 `.run/`):

```bash
make run-all       # 启动所有 API
make run-all-rpc   # 启动所有 RPC
make stop-all      # 停止全部
make logs          # 实时尾部日志
```



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
make build           # 编译所有服务
make test            # 各模块 go test
make lint            # golangci-lint
make fmt             # gofumpt 格式化
make clean           # 删除 bin/ 与 .run/
make run-api svc=user    # 前台跑单个 API
make run-rpc svc=user    # 前台跑单个 RPC
make run-all / stop-all  # 后台全量启停
make logs            # tail -f .run/*.log
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


| 项目           | 状态                                                                  |
| ------------ | ------------------------------------------------------------------- |
| **骨架**       | ✅ 11 服务 api+rpc 全部 goctl 生成完毕，`go build` 通过                         |
| **Logic 业务** | 🚧 全部是 goctl 默认 `// todo:` 桩，待按服务对齐 Java 版实现                        |
| **中间件**      | 🚧 JWT/统一响应/链路追踪在 `pkg/` 就绪，尚未接入各 svc                               |
| **包名规范**     | ✅ 模块路径统一 `tjxt/apps/<svc>` (data 为 `tjxt/apps/data/{api,rpc}/data`) |
| **事件总线**     | ✅ `pkg/mq/event` 已定常量与发布/订阅抽象，业务侧尚未接入                               |
| **数据库**      | ✅ DDL 就绪，Model 待生成                                                  |




## 🛠 开发规范

1. **只改 logic**: `api/internal/logic/` `rpc/internal/logic/` 以外的生成代码**不要手改**
2. **错误处理**: 业务错误用 `xerr.New(code, msg)` 或 `xerr.Wrap(err, code, msg)`，不要直接返回 `error`
3. **日志**: 使用 go-zero 内置 `logx`，关键路径带上 `trace.TraceIDFromContext(ctx)`
4. **事务**: RPC 层用 `model.Transact(ctx, func(session sqlx.Session) error { ... })`
5. **分页**: 统一用 `pkg/utils/page.ParsePageReq(req)` + `result.Page{List, Total, Pages}`



## 📚 参考资料

- **go-zero 官网**: [https://go-zero.dev](https://go-zero.dev)
- **goctl 手册**: [https://go-zero.dev/docs/tasks](https://go-zero.dev/docs/tasks)

