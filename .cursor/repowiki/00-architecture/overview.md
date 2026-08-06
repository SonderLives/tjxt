# tjxt 整体架构概览

> 版本：v1.1 | 更新：2026-08-06

## 系统拓扑

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Client (Web/App)                             │
└──────────────────────────────┬──────────────────────────────────────┘
                               │ HTTPS
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   API Gateway Layer (13 services)                   │
│  user:8801  auth:8802  course:8803  learning:8804  exam:8805         │
│  media:8806  message:8807  pay:8808  trade:8809  search:8810         │
│  data:8811   promotion:8812  remark:8813                            │
└──────────────────────────────┬──────────────────────────────────────┘
                               │ gRPC
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   RPC Service Layer (13 services)                   │
│  user:8081  auth:8082  course:8083  learning:8084  exam:8805         │
│  media:8806  message:8087  pay:8808  trade:8089  search:8090         │
│  data:8091   promotion:8092  remark:8093                            │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
              ┌────────────────┼────────────────┐
              ▼                ▼                ▼
        ┌──────────┐    ┌──────────┐      ┌──────────┐
        │ MySQL    │    │ Redis    │      │ RabbitMQ │
        │ (12 DBs) │    │ (cache)  │      │ (events) │
        └──────────┘    └──────────┘      └──────────┘
```

> 注：`data` 服务不依赖 MySQL（数据为内存/配置驱动），故无独立 DB；其余 12 个服务各独享 `tj_<domain>` 库。

## 模块结构（go-zero 官方标准）

项目采用 **go-zero 官方标准结构**：每个可独立部署的进程（api / rpc）各自拥有独立 `go.mod`，而非「每业务服务一个聚合 module」。整体由 `go.work` 工作区聚合管理。

- **`go.work`** 聚合 **28 个 use 模块**：根模块 `.` + `pkg` 公共库 + 26 个服务模块（13 个服务 × {api, rpc}）。
- **module path 规则**：

  | 模块 | module path |
  |------|-------------|
  | API | `tjxt/apps/<svc>/api` |
  | RPC | `tjxt/apps/<svc>/rpc` |
  | data 嵌套 | `tjxt/apps/data/api/data`、`tjxt/apps/data/rpc/data` |

- **公共库引用**：每个服务模块的 `go.mod` 通过相对 `replace` 引用本地 `pkg`：
  - 普通服务（三级）：`replace tjxt/pkg => ../../../pkg`
  - data 嵌套（四级）：`replace tjxt/pkg => ../../../../pkg`
- **跨服务依赖**：在调用方 `go.mod` 中 `require` + `replace` 到本地相对路径，例如 `trade` 依赖 `course`/`pay`：
  ```go
  require (
      tjxt/apps/course/rpc v0.0.0-00010101000000-000000000000
      tjxt/apps/pay/rpc   v0.0.0-00010101000000-000000000000
  )
  replace (
      tjxt/apps/course/rpc => ../../course/rpc
      tjxt/apps/pay/rpc   => ../../pay/rpc
  )
  ```
- **本地 tjxt 模块版本号**统一用伪版本 `v0.0.0-00010101000000-000000000000`（由 `go.work` 工作区解析到本地目录，不发布到任何 module proxy）。

> 重构记录：原「每服务单 module」结构已于 2026-08-06 拆分为上述官方标准结构。详见 [go-zero 开发规范](../../.cursor/repowiki/01-conventions/go-zero-rules.md) 与 [实现进度](../../.cursor/repowiki/06-status/implementation-status.md) §2.1。

## 服务清单与端口

| 服务 | 领域 | API 端口 | RPC 端口 | 数据库 | 核心职责 |
|------|------|----------|----------|--------|----------|
| user | 用户中心 | 8801 | 8081 | tj_user | 资料/教师学生角色/后台管理 |
| auth | 认证授权 | 8802 | 8082 | tj_user | 登录/注册、JWT、RBAC、菜单权限 |
| course | 课程管理 | 8803 | 8083 | tj_course | 课程/章节/资源 CRUD、发布审核 |
| learning | 学习记录 | 8804 | 8084 | tj_learning | 进度、笔记、收藏、证书 |
| exam | 考试题库 | 8805 | 8085 | tj_exam | 题库/组卷/考试/阅卷 |
| media | 文件媒体 | 8806 | 8086 | tj_media | 上传/签名/转码/存储桶管理 |
| message | 消息通知 | 8807 | 8087 | tj_message | 站内信/短信/模板/定时任务 |
| pay | 支付网关 | 8808 | 8088 | tj_pay | 微信/支付宝/余额支付、对账 |
| trade | 交易订单 | 8809 | 8089 | tj_trade | 订单创建/支付/退款/分账 |
| search | 搜索推荐 | 8810 | 8090 | tj_search | 兴趣标签/ES 索引/推荐 |
| data | 数据看板 | 8811 | 8091 | - | 实时大屏/Top10/今日数据 |
| promotion | 营销优惠 | 8812 | 8092 | tj_promotion | 优惠券/折扣/活动规则 |
| remark | 评论点赞 | 8813 | 8093 | tj_remark | 评论/点赞/收藏 |

## 数据流向原则

1. **单向依赖**：API 层 → RPC 层 → DB/MQ，禁止反向调用
2. **跨服务仅走 RPC Client**：`pkg/` 共享库除外，严禁 `import "tjxt/apps/xxx/internal/..."`
3. **数据库隔离**：每服务独享 `tj_<domain>` 库，跨域查询走对应服务 RPC
4. **异步解耦**：订单支付、消息发送、积分变更等通过 RabbitMQ Event 解耦

## 技术栈锚定

| 组件 | 版本 | 说明 |
|------|------|------|
| go-zero | v1.10.3 | 微服务框架、goctl 代码生成 |
| Go | 1.26.x | 工作区模式 `go.work`（每 api/rpc 独立 module） |
| MySQL | 8.0 | 开发环境 127.0.0.1:3306 root/0000 |
| Redis | 7.x | 缓存、分布式锁、限流 |
| RabbitMQ | 3.12 | 事件总线、延迟队列 |
| etcd | 3.5 | 服务发现、配置中心 |
| Docker Compose | - | 本地一键拉起依赖 |

## 部署拓扑（生产）

- K8s Deployment + Service + HPA
- ConfigMap/Secret 管理配置
- Prometheus + Grafana 监控
- Jaeger 链路追踪
- ELK 日志聚合
