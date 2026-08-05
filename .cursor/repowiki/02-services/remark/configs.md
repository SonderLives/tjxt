> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/remark/api/etc/remark-api.yaml`, `apps/remark/rpc/etc/remark.yaml`

---

# Remark Configs

## API 服务配置 (`apps/remark/api/etc/remark-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `remark-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8813` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `RemarkRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `RemarkRpc.Etcd.Key` | `remark.rpc` | - | remark RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `remark.rpc`（HTTP handler → RPC client 调用，两个接口均透传）
- etcd（服务发现）

**配置结构体**（`apps/remark/api/internal/config/config.go`）：

```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    RemarkRpc zrpc.RpcClientConf
}
```

> `routes.go` 中两个路由由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹，remark-api 所有接口都需要携带有效 JWT。

---

## RPC 服务配置 (`apps/remark/rpc/etc/remark.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `remark.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8093` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `remark.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_remark?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**依赖的外部服务**：
- MySQL `tj_remark` 库（仅 `like_record` 一张表）
- Redis 缓存（节点模式）
- etcd（服务注册）

**配置结构体**（`apps/remark/rpc/internal/config/config.go`）：

```go
type Config struct {
    zrpc.RpcServerConf
    DataSource string
    Cache      cache.CacheConf
}
```

> 与 promotion 不同，remark 的 `DataSource` 与 `Cache` **未标记** `optional`，缺失时服务启动即失败——配置缺陷能在启动阶段暴露。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_remark?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_remark`
- 字符集: `utf8mb4`
- 时区: `Local`

> `parseTime=true` 是必需项：`LikeRecord` 结构体中 `CreateTime` / `UpdateTime` 声明为 `time.Time`，未开启时驱动会返回 `[]byte` 导致扫描失败。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `remark-api` | `8813` | HTTP |
| `remark.rpc` | `8093` | gRPC |

### 缓存使用方式

`ServiceContext` 只注入一个 Model：

```go
conn := sqlx.NewMysql(c.DataSource)
LikeRecordModel: model.NewLikeRecordModel(conn, c.Cache)
```

| 缓存键前缀 | 承载对象 | 使用方 |
|-----------|---------|--------|
| `cache:likeRecord:id:` | 点赞记录主键缓存 | `FindOne` / 写操作失效 |
| `cache:likeRecord:userId:bizId:bizType:` | `uk_user_biz` 唯一索引缓存 | `FindOneByUserIdBizIdBizType`（`Like` 逻辑入口） |

自定义扩展方法 `FindLikedBizIds` **刻意不走缓存**（`QueryRowsNoCacheCtx`），保证点赞后列表页立即可见最新状态。

### JWT 密钥

生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造。remark-api 的 `Auth.AccessSecret` 必须与签发方（auth 服务）保持一致，否则所有接口鉴权失败。

`Auth.AccessExpire`（7200 秒）在 remark-api 中仅作配置项声明，实际的令牌过期校验由 go-zero 的 JWT 中间件依据 token 内的 `exp` 完成。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_remark) | DataSource 配置 | 自建存储 |
| Redis | Cache 配置 | 缓存点赞记录主键与唯一索引查询 |
| remark.rpc | RpcClient 配置 | API 层通过 etcd 服务发现调用 |
| etcd | Etcd 配置 | RPC 服务注册与发现 |

> remark.rpc 不依赖任何其他服务的 RPC。`biz_type` + `biz_id` 是多态关联，本服务不感知具体业务表结构，也不做跨服务校验——这是点赞能力可被任意业务复用的前提。
