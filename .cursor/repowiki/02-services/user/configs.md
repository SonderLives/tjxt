> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/user/api/etc/user-api.yaml`, `apps/user/rpc/etc/user.yaml`

---

# User Configs

## API 服务配置 (`apps/user/api/etc/user-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `user-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8808` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `UserRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `UserRpc.Etcd.Key` | `user.rpc` | - | user RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `user.rpc`（HTTP handler → RPC client 调用）
- etcd `127.0.0.1:2379` — 服务发现

> 结构体定义见 `apps/user/api/internal/config/config.go`：内嵌 `rest.RestConf`，`Auth` 为匿名结构体，`UserRpc` 为 `zrpc.RpcClientConf`。

---

## RPC 服务配置 (`apps/user/rpc/etc/user.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `user.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8082` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `user.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_user?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**依赖的外部服务**：
- MySQL `tj_user` 库
- Redis 缓存（节点模式）
- etcd `127.0.0.1:2379` — 服务注册

> 结构体定义见 `apps/user/rpc/internal/config/config.go`：内嵌 `zrpc.RpcServerConf`，外加 `DataSource string` 与 `Cache cache.CacheConf`。RPC 层**不含 Jwt 配置**——令牌签发是 auth 服务的职责。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_user?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_user`
- 字符集: `utf8mb4`
- 时区: `Local`

`parseTime=true` 是必需项：`UserWithDetail.CreateTime` 声明为 `time.Time`，需要驱动直接把 `datetime` 解析为时间类型。

### Redis 缓存

`Cache` 由 `svc.NewServiceContext` 一次性传给两个模型：

| 模型 | 缓存键前缀 |
|------|-----------|
| `UserModel` | `cache:tjUser:user:id:` / `cache:tjUser:user:cellPhone:type:` / `cache:tjUser:user:username:` |
| `UserDetailModel` | goctl 按 `user_detail` 主键生成 |

`UpdateStatus` 走 `ExecNoCacheCtx`，需在业务代码中手动 `DelCacheCtx` 三个键（详见 business-rules.md 第 6 节）。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `user-api` | 8808 | HTTP |
| `user.rpc` | 8082 | gRPC |

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_user) | DataSource 配置 | 自建存储，含 `user` / `user_detail` 两张表 |
| Redis | Cache 配置 | 缓存用户基础信息与资料 |
| etcd | Etcd 配置 | RPC 注册与发现 |

### 作为被依赖方

`user.rpc` 被以下服务在配置中声明为客户端：

| 消费方配置文件 | 配置节 | Key |
|---------------|--------|-----|
| `apps/user/api/etc/user-api.yaml` | `UserRpc` | `user.rpc` |
| `apps/auth/api/etc/auth-api.yaml` | `UserRpc` | `user.rpc` |
