# 天机学堂 Go 微服务

接口契约来自 Apifox 项目 `8658985`，完整 OpenAPI 导出保存在 `specs/tjxt.openapi.json`。当前实现将学习中心作为可运行的参考服务；接口与领域分组可据此按同一模式拆分到认证、课程、交易、支付、优惠券、媒资等服务。

## 架构

- `common`：公共底座，独立 module（错误码 `xerr`、统一响应 `result`、分页 `page`、JWT `auth`、雪花 `idgen`、RabbitMQ 生产者/消费者与订单事件模型），各服务通过 go.mod `replace` 引用。
- `learning`：JWT 鉴权 REST API，覆盖 Apifox 的 `/lessons` 课表接口。
- `learning/rpc`：go-zero zRPC 服务，供课程、交易等内部服务调用。
- `trade`：交易中心 REST API，覆盖购物车、下单/免费课报名、支付（模拟）、退款与订单明细查询（24 个接口，端口 8810）；支付/退款成功后通过 MQ `order.exchange` 发布 `order.pay`/`order.refund` 事件。
- `user`：用户中心 REST API（端口 8808），覆盖认证登录/刷新/退出（`/accounts`）与用户管理（注册、`/users` CRUD、学员/教师/员工分页），共 17 个接口；登录签发 JWT（`userId` 数值 claim + `role`），供 trade 等下游服务消费 `GET /users/{id}`。
- MySQL：学习记录持久化，退款使用软失效，支付开课使用唯一键幂等写入。
- RabbitMQ：监听订单支付与退款事件；队列、交换机均为持久化声明。
- Redis：预留统一缓存配置；etcd：API/RPC 服务注册发现。

## 本地启动

```powershell
docker compose up -d
go run . -f etc/learning-api.yaml
go run .\rpc -f rpc/etc/learning.yaml
go run . -f etc/trade-api.yaml
go run . -f etc/user-api.yaml
```

`trade` 依赖课程服务 `GET /courses/simpleInfo/list` 与用户服务 `GET /users/{id}` 内部接口（`etc/trade-api.yaml` 中 `CourseService`/`UserService`），服务未启动时相关接口返回依赖不可用错误；订单事件发布依赖 RabbitMQ 就绪。

内置演示账户（见 `sql/tj_user.sql`，bcrypt 加密）：`admin / 13500010002`（员工）、`jack / 13500010003`（学员）、`rose / 13500010004`（学员）。种子数据的初始密码为未知明文哈希，无法直接登录；验证学生链路请直接 `POST /students/register` 注册新学员（注册即使用默认密码 `123456`，对应配置 `DefaultPassword`）后登录。员工/管理端功能需先在数据库中为某个 `type=1` 账户写入已知密码的 bcrypt 哈希，或用 `PUT /users/{id}/password/default` 重置（该接口本身需要员工 token）。

生产环境请通过部署平台注入 `Auth.AccessSecret`、数据库 DSN 和 RabbitMQ 密码，勿提交真实密钥。JWT 需包含数值型 `userId` claim。
