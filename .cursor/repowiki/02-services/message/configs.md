> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/message/api/etc/message-api.yaml`, `apps/message/rpc/etc/message.yaml`

---

# Message Configs

## API 服务配置 (`apps/message/api/etc/message-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `message-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8807` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `MessageRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `MessageRpc.Etcd.Key` | `message.rpc` | - | message RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `message.rpc`（HTTP handler → RPC client 调用，`apps/message/api/internal/svc/servicecontext.go` 注入 `MessageRpc`）
- etcd（服务发现）

> 结构体定义见 `apps/message/api/internal/config/config.go`：`rest.RestConf` + 匿名 `Auth{AccessSecret, AccessExpire}` + `MessageRpc zrpc.RpcClientConf`。**未配置 MySQL / Redis**，API 层不直连存储。

---

## RPC 服务配置 (`apps/message/rpc/etc/message.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `message.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8087` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `message.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_message?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**依赖的外部服务**：
- MySQL `tj_message` 库
- Redis 缓存（节点模式）
- etcd（服务注册）

> 结构体定义见 `apps/message/rpc/internal/config/config.go`：`zrpc.RpcServerConf` + `DataSource string` + `Cache cache.CacheConf`。**未配置 JWT 段**——令牌校验在 API 层的 `rest.WithJwt` 完成，RPC 层不做鉴权。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_message?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_message`
- 字符集: `utf8mb4`
- 时区: `Local`

> 注意：连接串字符集为 `utf8mb4`，而 `sql/ddl/tj_message.sql` 中 7 张表里有 6 张建表字符集为 `utf8`（仅 `notice_task_target` 为 `utf8mb4`），两者不一致。

### JWT 密钥

message 服务**只在 API 层校验 JWT，不签发令牌**。`Auth.AccessSecret` 必须与 auth 服务的 `Jwt.AccessSecret` 保持一致，否则所有接口 401。

| 密钥 | 用途 |
|------|------|
| `Auth.AccessSecret` | 校验 `accessToken`（`routes.go` 中 6 处 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)`） |

`Auth.AccessExpire` 在 message 服务中仅作为配置项被反序列化，代码内无签发逻辑引用。生产环境**必须修改**默认值 `change-me-in-production`。

### 端口分配

| 组件 | 端口 | 说明 |
|------|------|------|
| `message-api` | 8807 | HTTP，对外 |
| `message.rpc` | 8087 | gRPC，集群内 |

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_message) | DataSource 配置 | 自建存储，6 个 Model 共用一个 `sqlx.NewMysql(c.DataSource)` 连接 |
| Redis | Cache 配置 | 6 个 Model 均走 `sqlc.CachedConn`，缓存前缀 `cache:<model>:id:` |
| etcd | Etcd 配置 | RPC 端注册 `message.rpc`，API 端按同名 key 发现 |
| auth 服务 | 无直接连接 | 仅共享 `AccessSecret` 完成离线 JWT 校验，不调用 auth RPC |

### 与其它服务的配置差异

| 差异点 | message | 对照（auth） |
|--------|---------|-------------|
| API 是否依赖其它服务 RPC | 仅依赖自身 `message.rpc` | 额外依赖 `user.rpc` |
| RPC 是否配置 JWT | 否 | 是（`Jwt.AccessSecret` / `RefreshSecret`） |
| RPC 是否签发令牌 | 否 | 是（`SignToken`） |
