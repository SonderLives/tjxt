# tjxt — 天机学堂 Go 微服务

基于 **go-zero** 框架的微服务架构，采用 **API + RPC 分层**，通过 **etcd** 服务发现，使用 **MySQL + Redis** 作为存储。
8 个无跨服务调用方的服务（auth/exam/media/message/remark/search/trade/data）已完成 **RPC 折叠**：RPC 进程壳删除，业务实现随迁至各自 `api/<svc>/backend`；trade/search 原挂在 RPC 进程里的 MQ 消费者独立为 `app/mq/trade-mq`、`app/mq/search-mq` 两个后台消费者进程。

## 🏗 架构概览

```
HTTP Client
    │
    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│ API 层 (go-zero REST, 独立 go.mod)                                          │
│                                                                            │
│  auth-api:8802    user-api:8801    course-api:8803   learning-api:8804    │
│  exam-api:8805    media-api:8806    message-api:8807  pay-api:8808        │
│  trade-api:8809   search-api:8810   data-api:8811    promotion-api:8812   │
│  remark-api:8813                                                       │
│                                                                            │
│  折叠服务（auth/exam/media/message/remark/search/trade/data）：             │
│  业务实现在 api/<svc>/backend，同进程内直调，不再出站 RPC                    │
└───────┬───────────────────────────────────────────────┬────────────────────┘
        │ zrpc (etcd 服务发现，仅 5 个保留 RPC 服务)        │ RabbitMQ 消费
        ▼                                               ▼
┌────────────────────────────────────┐  ┌──────────────────────────────────────┐
│ RPC 层 (go-zero zrpc, 独立 go.mod, │  │ MQ 消费者 (app/mq/)                  │
│ 持 model + Redis 缓存)             │  │                                      │
│                                    │  │  trade-mq  支付回写（pay.success →   │
│  user.rpc      course.rpc          │  │            订单已支付 → order.pay）  │
│  learning.rpc  pay.rpc             │  │  search-mq ES 同步（course.up/down   │
│  promotion.rpc                     │  │            → 写/删课程索引）         │
└──────────────┬─────────────────────┘  └──────────────────────────────────────┘
               ▼
     MySQL (sqlx + goctl model --cache)
     Redis (缓存)
     RabbitMQ (事件总线，见下文「事件流」)
```

> ⚠️ **无独立网关**。边缘入口按需自行接入 APISIX / Kong / Nginx。



## 📁 目录结构

```
tjxt/
├── go.work                  # 工作区，聚合 22 个 use 模块（根 + pkg + 13 api + 5 rpc + 2 mq）
├── go.mod                   # 根模块 (仅声明版本)
├── Makefile                 # 构建/生成/运行/校验一站式命令
├── buf.yaml                 # buf workspace（5 个 rpc/<svc> + 8 个折叠服务 api/<svc> 共 13 个 proto 模块）
├── buf.gen.yaml             # buf generate 模板（折叠服务 pb 类型生成，仅 protoc-gen-go）
├── build.ps1                # PowerShell 全量构建（与 make build 同清单：20 个二进制）
├── docker-compose.yml       # MySQL/Redis/RabbitMQ/etcd/ES/RustFS + 可观测性栈（Jaeger/otel-collector/Prometheus/Loki）一键启动
├── tools/
│   └── goctl-templates/     # goctl 自定义模板（api/rpc 生成统一走该目录，handler 统一 result.Write）
├── pkg/                     # 公共代码库 (module: tjxt/pkg)
│   ├── auth/                # JWT 签发/校验、Claims 定义
│   ├── mq/                  # RabbitMQ 事件总线 (生产者/消费者/事件定义)
│   ├── response/            # 统一响应 R{code,msg,requestId,data} + 分页
│   ├── errorx/              # 业务错误码体系 + HTTP 状态映射
│   ├── idgen/               # 雪花算法 ID 生成器
│   ├── page/                # 分页工具
│   ├── timex/               # 统一时间格式化 Layout/Format/Parse
│   └── go.mod
├── sql/
│   ├── ddl/                 # goctl model 生成用的纯 DDL (按库分文件)
│   └── migration/           # 完整迁移脚本 (含索引/外键/初始数据)
├── docs/                    # 架构/设计/历史文档
├── deployments/             # 观测栈配置 (otel-collector / prometheus / loki)
├── api/
│   └── <svc>/               # 13 个 HTTP API 服务模块（module: tjxt/api/<svc>，auth/course/data/exam/learning/media/message/pay/promotion/remark/search/trade/user）
│       ├── go.mod
│       ├── <svc>.api          # goctl API 定义 (路由/请求/响应)
│       ├── <svc>.go           # main 入口
│       ├── etc/<svc>-api.yaml # API 配置
│       ├── internal/
│       │   ├── config/        # 配置加载
│       │   ├── handler/       # HTTP 处理器 (goctl 生成)
│       │   ├── logic/         # 业务逻辑 (仅手写这里)
│       │   ├── svc/           # 服务上下文 (RPC 客户端等)
│       │   └── types/         # 请求/响应结构体 (goctl 生成)
│       └── backend/           # 仅折叠服务（auth/exam/media/message/remark/search/trade/data）：
│           ├── config/        #   折叠 RPC 的配置段（DataSource/Cache/下游 RPC 客户端等）
│           ├── logic/         #   原 rpc/<svc>/internal/logic 的业务实现（仅手写这里）
│           ├── model/         #   sqlx Model (goctl 生成 + 扩展)
│           ├── pb/            #   生成的 pb 类型（make pb / buf generate 产出）
│           └── svc/           #   服务上下文 (DB/Redis/缓存)
├── rpc/
│   └── <svc>/               # 5 个保留 RPC 服务模块（user/course/learning/pay/promotion，module: tjxt/rpc/<svc>）
│       ├── go.mod
│       ├── <svc>.proto        # gRPC 服务定义
│       ├── <svc>.go           # main 入口
│       ├── etc/<svc>.yaml     # RPC 配置
│       ├── internal/
│       │   ├── config/
│       │   ├── logic/         # RPC 业务逻辑 (仅手写这里)
│       │   ├── model/         # sqlx Model (goctl 生成 + 扩展)
│       │   ├── server/        # gRPC 服务注册
│       │   └── svc/           # 服务上下文 (DB/Redis/缓存)
│       └── pb/                # 生成的 pb.go / _grpc.pb.go
└── app/
    └── mq/                  # MQ 消费者进程（自 rpc 进程独立出来的后台消费者）
        ├── trade-mq/        # 支付回写：消费 pay.success 回写订单已支付，发布 order.pay/order.refund
        │   ├── etc/trade-mq.yaml
        │   └── main.go
        └── search-mq/       # ES 同步：消费 course.up/down 写/删课程索引（复用 api/search/backend 的 svc）
            ├── etc/search-mq.yaml
            └── main.go
```



## 🗂 服务清单


| 服务            | API 端口 | RPC 地址         | etcd Key        | 数据库             | 核心职责                  |
| ------------- | ------ | -------------- | --------------- | --------------- | --------------------- |
| **auth**      | 8802   | 无（已折叠进 API）    | —               | tj_auth         | 登录/注册、Token 签发刷新、权限校验 |
| **user**      | 8801   | 127.0.0.1:8081 | `user.rpc`      | tj_user         | 用户档案、教师/学员详情、角色部门     |
| **course**    | 8803   | 127.0.0.1:8083 | `course.rpc`    | tj_course       | 课程/分类/章节/媒体资源 CRUD    |
| **learning**  | 8804   | 127.0.0.1:8084 | `learning.rpc`  | tj_learning     | 学习进度、收藏、笔记、考试记录       |
| **exam**      | 8805   | 无（已折叠进 API）    | —               | tj_exam         | 试卷/题库/考试/评分/记录        |
| **media**     | 8806   | 无（已折叠进 API）    | —               | tj_media        | 文件上传、媒资管理、转码回调        |
| **message**   | 8807   | 无（已折叠进 API）    | —               | tj_message      | 站内信/通知/公告、模板消息        |
| **pay**       | 8808   | 127.0.0.1:8088 | `pay.rpc`       | tj_pay          | 订单支付、退款、回调、账单         |
| **trade**     | 8809   | 无（已折叠进 API）    | —               | tj_trade        | 交易订单、优惠券、分佣结算         |
| **search**    | 8810   | 无（已折叠进 API）    | —               | tj_search       | 课程/用户全文检索、建议词         |
| **data**      | 8811   | 无（已折叠进 API）    | —               | — (Redis)       | 统计大屏、榜单、今日数据聚合        |
| **promotion** | 8812   | 127.0.0.1:8092 | `promotion.rpc` | tj_promotion    | 优惠券、营销活动、秒杀、积分        |
| **remark**    | 8813   | 无（已折叠进 API）    | —               | tj_remark       | 评价/评论、评分、回复           |
| **trade-mq**  | —      | —（MQ 消费者）      | —               | tj_trade        | 支付回写：消费 `pay.success` 回写订单已支付，发布 `order.pay` |
| **search-mq** | —      | —（MQ 消费者）      | —               | tj_search + ES  | ES 同步：消费 `course.up/down` 写/删课程索引（metrics :9110） |


> 📌 各服务端口已在 `api/*/etc/*.yaml`、`rpc/*/etc/*.yaml` 与 `app/mq/*/etc/*.yaml` 中唯一分配，可同时启动（media-api 8806、promotion-api 8812、course.rpc 8083、promotion.rpc 8092；metrics：保留 rpc 9101-9112 中的 5 个、api 9201-9213、search-mq 9110）。
>
> 📌 **RPC 折叠**：auth/exam/media/message/remark/search/trade/data 的 RPC 进程壳已删除，业务实现随迁 `api/<svc>/backend`；这些服务不再注册 etcd RPC 端点，proto 契约保留在 `api/<svc>/<svc>.proto`（pb 类型经 `make pb` 再生成）。



## 📦 前置依赖


| 组件                | 版本要求   | 说明                                                           |
| ----------------- | ------ | ------------------------------------------------------------ |
| **Go**            | ≥ 1.26 | 用到 `go.work` 工作区                                             |
| **MySQL**         | 8.x    | 默认 `root:0000@127.0.0.1:3306`                                |
| **Redis**         | 7.x    | 默认 `127.0.0.1:6379`                                          |
| **etcd**          | 3.5+   | 默认 `127.0.0.1:2379`，zrpc 服务发现                                |
| **Elasticsearch** | 9.x    | search 服务的检索后端，默认 `http://127.0.0.1:9200`（compose 已内置单节点免认证） |
| **RabbitMQ**      | 3.13+  | 事件总线（course→search、pay→trade 已接线）；未配置时服务仅告警并跳过，不阻塞启动         |
| **goctl**         | 1.10.x | API/RPC/model 骨架代码生成工具                                       |
| **buf**           | 1.75.0 | proto lint/breaking 门禁 + 折叠服务 pb 类型生成（`make pb`）             |
| **protoc-gen-go** | 1.36.x | buf generate 的 pb 生成插件（`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`） |


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

# 3. 编译全部 20 个可执行文件 (13 api + 5 rpc + 2 mq) 到 bin/
make build

# 4. 启动单个服务（务必从仓库根执行，日志落到 logs/<svc>/）
#    推荐用 Makefile 目标（已改为从仓库根 go run ./api/... ./rpc/... ./app/mq/...，CWD=仓库根）：
make run-user          # 启动 user-api
make run-user-rpc      # 启动 user-rpc
make run-trade-mq      # 启动支付回写消费者 trade-mq
make run-search-mq     # 启动 ES 同步消费者 search-mq
#    或手动（同样要在仓库根执行）：
go run ./api/user -f api/user/etc/user-api.yaml
go run ./rpc/user -f rpc/user/etc/user.yaml
go run ./app/mq/trade-mq -f app/mq/trade-mq/etc/trade-mq.yaml
go run ./api/data -f api/data/etc/data-api.yaml
```

**一键启动全部服务**（后台独立窗口，日志写入 `logs/<svc>/`）：

```bash
make run-all       # 启动所有 API（各开一个独立 PowerShell 窗口）
make run-all-rpc   # 启动 5 个保留 RPC + 2 个 MQ 消费者（各开一个独立 PowerShell 窗口）
# 停止：各窗口直接关闭（Ctrl-C）；当前无 stop-all 目标
# 看服务日志：tail -f logs/course-api/access.log   （替换 <svc> 即可）
# 看基建容器日志：make docker-logs
```



## 🔭 可观测性（trace / metrics / logs 统一经 otel-collector）

各服务 yaml 已注入 `Telemetry`（链路）、`Prometheus`（/metrics）、`Log{Mode:file,...}`（日志文件）。
先 `make docker-up` 起基建，再起各 Go 服务：


| 信号         | 路径                                                                           | 界面/端点                                                      |
| ---------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------- |
| 链路 Trace   | go-zero → `127.0.0.1:4318` → otel-collector → Jaeger                         | Jaeger UI [http://localhost:16686](http://localhost:16686) |
| 指标 Metrics | 各服务 `/metrics`（api 9201–9213、保留 rpc 9101/9103/9104/9108/9112、search-mq 9110） → otel-collector 聚合 `:8889` → Prometheus | Prometheus [http://localhost:9090](http://localhost:9090)  |
| 日志 Logs    | `logs/<svc>/*.log`（仓库根）→ 挂进 collector → Loki                                 | Loki `:3100`（建议接 Grafana 查询）                               |


> ⚠️ 日志路径是**相对路径** `logs/<svc>`，相对进程启动目录（CWD）；服务须从**仓库根**启动，
> 否则日志会散落到各服务子目录，collector 的 filelog 采集不到。



## ⚙️ 代码生成 (核心约束)

**骨架一律** `goctl`/`buf` **生成，不手写。业务逻辑只允许填** `logic/`/`backend/logic/`**。**

生成走**双通道**：5 个保留 RPC 服务用 goctl（`make rpc`），8 个折叠服务只有 pb 类型、用 buf（`make pb`）。

| 场景                     | 命令                                                                                                                                |
| ---------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| 新建 API 骨架              | `goctl api new <svc>`                                                                                                             |
| 修改 `.api` 后重生成         | `goctl api go -api api/<svc>/<svc>.api -dir api/<svc> --home tools/goctl-templates`                                    |
| 新建 RPC 骨架（仅 5 个保留服务）   | `goctl rpc new <svc>`                                                                                                             |
| 修改保留服务 `.proto` 后重生成   | `cd rpc/<svc> && goctl rpc protoc <svc>.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true -m`                                    |
| 修改折叠服务 `.proto` 后重生成   | `cd api/<svc> && buf generate --template ../../buf.gen.yaml`（只产出 pb 类型到 `backend/pb/`）                                      |
| 数据库 Model (带 Redis 缓存) | `goctl model mysql datasource -url="root:0000@tcp(127.0.0.1:3306)/tj_<svc>" -table="<tbl>" -dir <模块>/internal/model -c` |
| Model 需加自定义方法          | 在 `internal/model/<tbl>model.go` 里写扩展 (`_gen.go` 由 goctl 管理)                                                                      |


> **goctl 版本固定（重要）**：仓库生成物出自 **goctl 1.10.1**，API 侧使用仓库自有模板目录 `tools/goctl-templates/`
> （handler 统一调 `pkg/response.Write`，与 goctl 默认模板不同）。日常生成请安装 1.10.1 并带 `--home tools/goctl-templates`：
>
> ```bash
> go install github.com/zeromicro/go-zero/tools/goctl@v1.10.1
> ```
>
> 用其他版本生成会让 routes.go/types.go 的版本头漂移（`make check-gen` 会拦截）。RPC 侧逻辑文件由 goctl 跳过已存在文件保护，server/client/pb 可安全再生成。
>
> **buf 通道（折叠服务 pb 类型）**：模板在仓库根 `buf.gen.yaml`（version v2，仅 `protoc-gen-go` 插件，无 gRPC 桩）。
> 在 `api/<svc>` 目录下执行时 buf 向上找到根 `buf.yaml` workspace，只对当前模块生成；输出路径由 proto 的
> `go_package`（`./backend/pb`）锚定，落到 `api/<svc>/backend/pb/`。折叠服务没有 RPC 进程壳，因此**只出类型不出 server/client**。


**Makefile 快捷方式**:

```bash
make api            # 根据 .api 重生成所有 API 代码
make rpc            # 根据 .proto 重生成 5 个保留 RPC 服务的 pb/server/client/logic（goctl）
make pb             # buf generate 重生成 8 个折叠服务的 pb 类型（api/<svc>/backend/pb）
make model          # 根据 sql/ddl/*.sql 生成 model（输出到 ./model/<db>/）
make generate       # api + rpc + pb 全量重生成
make check-gen      # 重生成 api + pb 并校验与仓库零漂移（api types 与 pb 类型双守卫）
```



## 🧬 Proto 治理（buf）

仓库根 `buf.yaml`（v2 workspace）注册 13 个 buf 模块：5 个保留 RPC 服务（`rpc/<svc>`）+ 8 个折叠服务（`api/<svc>`，proto 在模块根），承担三层职责：

1. **buf lint**：proto 静态检查（命名风格、syntax、package 规范等）。
2. **buf breaking**：对 `master` 的破坏性变更门禁（改字段类型/编号、删字段/消息/RPC 等会直接报错）。
3. **buf generate**：折叠服务的 pb 类型再生成（`make pb`，见上文「代码生成」）。

### buf 的角色边界（重要）

| 事项 | 归属 | 说明 |
| --- | --- | --- |
| proto lint / breaking 门禁 | **buf** | `buf.yaml` 一处配置，本地与 CI 同口径 |
| 折叠服务 pb 类型生成 | **buf** | `buf.gen.yaml` 模板仅含 protoc-gen-go；rpc 进程壳已删，goctl rpc protoc 无法只生成类型，故走 `make pb` |
| 保留服务 pb / zrpc 代码生成 | **goctl** | buf 不接管 zrpc 生成链路：goctl 的 zrpc 骨架（server/logic/client）不是 protoc 插件，无法由 buf 替代；生成入口仍是 `make rpc` |

### 本地命令

```bash
make buf-lint        # buf lint：检查 13 个 proto 模块（规则与豁免见 buf.yaml）
make buf-breaking    # buf breaking：对比 master，禁止破坏性 proto 变更
```

本地需安装 buf（CI 已固定 `v1.75.0`）与 protoc-gen-go：

```bash
scoop install buf    # 或 go install github.com/bufbuild/buf/cmd/buf@v1.75.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
```

> `buf-breaking` 用 `.git#branch=master` 作为基线；master 上还没有 `buf.yaml`，所以命令里带 `--against-config buf.yaml` 复用同一 workspace 配置，否则 against 侧会被当成单一模块（报 "13 images vs 1 images"）。
> 折叠服务的 proto 从 `rpc/<svc>/<svc>.proto` 搬到 `api/<svc>/<svc>.proto`，但都位于各自 buf 模块根、文件名不变，buf 按「模块内相对路径」对比文件身份，故搬家本身不触发 breaking 报错。

### CI 行为

`.github/workflows/ci.yml` 在 go 门禁之前新增三步：安装 buf（`bufbuild/buf-setup-action`，buf 版本固定 `v1.75.0`）→ `buf lint` → `buf breaking --against '.git#branch=master'`。因 breaking 需要对比 master，checkout 已改为 `fetch-depth: 0` 拉全量历史。原有 lint/test/build 步骤语义不变。

### 豁免规则清单（buf.yaml lint.except）

存量 proto 大量字段为驼峰命名、消息被多个 RPC 复用，修复这些规则都需要破坏性变更（改字段名/消息名/服务名/包名或挪目录），会直接改变 goctl 生成的 pb/zrpc 代码与线上 RPC 全限定名（`/package.Service/Method`）。因此对以下 7 条规则显式豁免（原因已注释在 `buf.yaml`）：

| 豁免规则 | 要求 | 豁免原因 |
| --- | --- | --- |
| `PACKAGE_DIRECTORY_MATCH` | 目录与包名层级一致（`rpc/<svc>/<svc>/`，折叠服务为 `api/<svc>/<svc>/`） | 挪动 proto 会破坏 goctl rpc protoc 的执行路径与仓库既有布局 |
| `PACKAGE_VERSION_SUFFIX` | 包名带版本后缀（如 `user.v1`） | 改包名会改变生成代码与 gRPC 全限定名 |
| `SERVICE_SUFFIX` | 服务名以 `Service` 结尾 | 改服务名会改变 RPC 路径，影响 zrpc 客户端与所有调用方 |
| `FIELD_LOWER_SNAKE_CASE` | 字段名 lower_snake_case | 存量字段大量驼峰（`userId`/`pageNo` 等），改字段名改变 JSON wire 名称与生成 Go 结构体字段 |
| `RPC_REQUEST_RESPONSE_UNIQUE` | 请求/响应类型不得跨 RPC 复用 | 拆分复用消息需新建并重命名类型，改变生成代码面 |
| `RPC_REQUEST_STANDARD_NAME` | 请求类型命名 `RpcNameRequest` | 需重命名消息类型，改变生成 Go 类型名 |
| `RPC_RESPONSE_STANDARD_NAME` | 响应类型命名 `RpcNameResponse` | 需重命名消息类型，改变生成 Go 类型名 |

新增 proto 仍受其余规则约束（消息/RPC/枚举命名、syntax 声明、package 规范等）；新增字段请遵守 snake_case，新增消息请按 `RPC 名 + Request/Response` 命名，避免豁免面继续扩大。



## 📋 Makefile 速查

```bash
make help            # 全部命令
make init            # go work sync 并 tidy 所有模块
make verify          # 校验结构：所有服务的 .api/.proto 与 internal/ 目录是否齐全
make build           # 编译所有服务（20 个二进制到 bin/：13 api + 5 rpc + 2 mq）
make test            # 各模块 go test（含 pkg 与 mq 消费者）
make lint            # golangci-lint 逐模块检查
make fmt             # go fmt 所有模块
make ci              # 一条龙门禁：verify + lint + test + build
make buf-lint        # buf lint：检查 13 个 proto 模块（规则与豁免见 buf.yaml）
make buf-breaking    # buf breaking：对比 master，禁止破坏性 proto 变更
make clean           # 删除 bin/ 与各服务下遗留 exe
make check-gen       # 重生成 api + pb 并校验与仓库零漂移（api types 与 pb 类型双守卫）
make run-<svc>       # 前台跑单个 API（如 make run-user）
make run-<svc>-rpc   # 前台跑单个保留 RPC（如 make run-user-rpc）
make run-trade-mq    # 前台跑支付回写消费者
make run-search-mq   # 前台跑 ES 同步消费者
make run-all         # 后台并行启动所有 API（各开独立窗口）
make run-all-rpc     # 后台并行启动 5 个保留 RPC + 2 个 MQ 消费者（各开独立窗口）
make docker-up       # 拉起 docker-compose 全部容器（含可观测性栈）
make docker-down     # 停止并移除容器
make docker-logs     # 跟踪容器日志
```



## 🔄 事件流 (RabbitMQ)

`pkg/mq` 基于 RabbitMQ（amqp091-go）实现：Producer confirm 模式发布；Consumer 泛型注册、手动 ack、处理失败 Nack 重回队列。交换机/路由键常量在 `pkg/mq/mq.go`，消息体定义在 `pkg/mq/event/`。

当前已定义的事件契约与接线状态：


| 交换机              | 路由键            | 消息体                | 发布方                      | 订阅方                                                            | 状态    |
| ---------------- | -------------- | ------------------ | ------------------------ | -------------------------------------------------------------- | ----- |
| `course.events`  | `course.up`    | `CourseEvent`      | course.rpc 课程上架          | app/mq/search-mq（队列 `search.course.up`，写/更新 ES 索引）             | ✅ 已接线 |
| `course.events`  | `course.down`  | `CourseEvent`      | course.rpc 课程下架          | app/mq/search-mq（队列 `search.course.down`，删 ES 索引）              | ✅ 已接线 |
| `pay.exchange`   | `pay.success`  | `PaySuccessEvent`  | pay.rpc 支付成功回调           | app/mq/trade-mq（队列 `trade.order.paid`，条件更新订单为已支付）              | ✅ 已接线 |
| `order.exchange` | `order.pay`    | `OrderPayEvent`    | app/mq/trade-mq 订单置为已支付后发布 | learning.rpc（队列 `learning.lesson.pay.queue`，为用户加课）           | ✅ 已接线 |
| `order.exchange` | `order.refund` | `OrderRefundEvent` | api/trade（backend 退款审批通过） | learning.rpc（队列 `learning.lesson.refund.queue`，撤销课程）          | ✅ 已接线 |


> 支付链路（demo）：`POST /pay-notify/pay-success` 需要 HMAC 签名——`sign = hex(hmac_sha256(PayNotifySecret, "payOrderNo=<payOrderNo>"))`，密钥在 `api/pay/etc/pay-api.yaml` 的 `PayNotifySecret`；本地演示用 `POST /pay-notify/mock-pay`（JWT 鉴权 + 支付单属主校验），二者最终都走 `NotifyPaySuccess` RPC，条件更新 + 幂等。
>
> Java 版的 `mq:course:cate:change`、`mq:media:*`、`mq:exam:record`、`mq:data:today` 等事件 Go 版尚未实现，需要时按上表模式在 `pkg/mq` 扩展常量与消息体。
>
> 事件发布均为 best-effort：发布失败仅告警、不回滚业务操作；消费侧以幂等（条件更新 / ON DUPLICATE KEY）兜底重复投递。索引一致性由 search 启动全量重建 + 手动 `ReindexCourses` RPC 兜底。

**发布示例**（course.rpc 实际用法，见 `rpc/course/internal/svc/servicecontext.go`）:

```go
svcCtx.Producer.Publish(ctx, mq.ExchangeCourse, mq.RoutingKeyCourseUp, event.CourseEvent{CourseID: id})
```

**订阅示例**（search-mq 实际用法，注册逻辑复用 `api/search/backend/svc/servicecontext.go`，进程入口为 `app/mq/search-mq/main.go`）:

```go
client := mq.NewClient(dsn)
mq.Register(client, mq.Binding{
    Queue:      "search.course.up",
    Exchange:   mq.ExchangeCourse,
    RoutingKey: mq.RoutingKeyCourseUp,
}, func(ctx context.Context, msg *event.CourseEvent) error {
    // 写入/更新 ES 索引；返回 error 时消息 Nack 重回队列
    return nil
})

// main 里阻塞启动（app/mq/search-mq/main.go）
ctx.MQClient.Start(context.Background())
```



## 🧩 公共包 (pkg/)


| 包                 | 说明                                                                                   |
| ----------------- | ------------------------------------------------------------------------------------ |
| `pkg/auth`        | JWT `Sign/Parse`、Claims 含 `userId`/`role`                                            |
| `pkg/errorx`      | 业务错误码 `Codexxx` + `Msg/Wrap/HttpStatus`                                              |
| `pkg/response`    | 统一响应 `R{code,msg,requestId,data}` + `Page{list,total,pages}` + `Write(w,r,data,err)` |
| `pkg/mq`          | RabbitMQ 事件总线：confirm 发布、泛型消费者（手动 ack、失败 Nack 重回队列）、交换机/路由键常量                        |
| `pkg/idgen`       | 雪花算法 `NextID()`                                                                      |
| `pkg/page`        | 分页参数归一化                                                                              |
| `pkg/timex`       | 统一时间格式化 Layout/Format/Parse                                                          |
| `pkg/confenv`     | 配置环境变量占位解析（`${VAR}` / `${VAR|default}`）                                              |


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
- 失败: `result.Fail(err)` 自动映射 `errorx` 错误码 → HTTP 状态码
- 所有 handler 统一用 `result.Write(w, r, data, err)`



## 📊 业务实现现状

> 结论（2026-08-06 复核）：13/13 服务、API+RPC 共 394 个接口，全部 logic 已实现并 `go build` 逐模块编译通过；功能完备度≈97%（少量桩为有意预留，见下）。RPC 折叠后：8 个折叠服务的 RPC logic 随迁 `api/<svc>/backend`，进程数 26 → 20（13 api + 5 rpc + 2 mq）。


| 项目           | 状态                                                                                                              |
| ------------ | --------------------------------------------------------------------------------------------------------------- |
| **骨架**       | ✅ 13 api + 5 保留 rpc + 2 mq 全部生成/实现完毕，`go build ./...` 逐模块编译通过                                                   |
| **Logic 业务** | ✅ 394/394 logic 全部实现（API 193 + RPC 201，其中 8 个折叠服务的 RPC logic 位于 `api/<svc>/backend/logic`），全库 0 处 TODO/panic 占位 |
| **中间件**      | ✅ JWT(`@server jwt:Auth`)、统一响应 `result.Write`、errorx 错误码已在各 svc 接入                                                |
| **包名规范**     | ✅ 模块路径统一 `tjxt/api/<svc>` / `tjxt/rpc/<svc>` / `tjxt/app/mq/<svc>`                              |
| **事件总线**     | ✅ 三段链路全部接线：course.rpc→search-mq 上下架同步 ES；pay.rpc→trade-mq 支付成功回写订单（条件更新+幂等）；trade（trade-mq/api 后端）→learning.rpc 支付加课/退款撤课            |
| **数据库**      | ✅ DDL + goctl model（带缓存）已生成，自定义 Model 已扩展                                                                       |
| **跨域 RPC**   | ⚠️ 已接线 trade backend→{course,pay}、search-mq→course、learning→course；其余（如 course→user/learning、trade→promotion、pay→真实网关）尚未接线 |
| **已知缺口**     | ⚠️ media 对象存储为 mock、pay 支付渠道为 demo（mock 网关，签名校验已接）、trade 优惠券未接入、Seata 未接入（undo_log 表闲置）                         |




## 🛠 开发规范

1. **只改 logic**: `api/<svc>/internal/logic/`、`rpc/<svc>/internal/logic/`、`api/<svc>/backend/logic/`（折叠服务）以外的生成代码**不要手改**
2. **错误处理**: 业务错误用 `errorx.New(code, msg)` 或 `errorx.Wrap(err, code, msg)`，不要直接返回 `error`
3. **日志**: 使用 go-zero 内置 `logx`，关键路径带上 `trace.TraceIDFromContext(ctx)`
4. **事务**: 持 model 的一层（保留 rpc 或折叠服务 backend）用 `model.Transact(ctx, func(session sqlx.Session) error { ... })`
5. **分页**: 统一用 `pkg/page.ParsePageReq(req)` + `result.Page{List, Total, Pages}`



## 📚 参考资料

- **go-zero 官网**: [https://go-zero.dev](https://go-zero.dev)
- **goctl 手册**: [https://go-zero.dev/docs/tasks](https://go-zero.dev/docs/tasks)

