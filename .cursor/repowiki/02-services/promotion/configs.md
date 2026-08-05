> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/promotion/api/etc/promotion-api.yaml`, `apps/promotion/rpc/etc/promotion.yaml`

---

# Promotion Configs

## API 服务配置 (`apps/promotion/api/etc/promotion-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `promotion-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8818` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `PromotionRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `PromotionRpc.Etcd.Key` | `promotion.rpc` | - | promotion RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `promotion.rpc`（HTTP handler → RPC client 调用，全部 16 个接口均透传）
- etcd（服务发现）

**配置结构体**（`apps/promotion/api/internal/config/config.go`）：

```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    // PromotionRpc 促销域 RPC 客户端（通过 etcd 服务发现）
    PromotionRpc zrpc.RpcClientConf
}
```

> 全部路由在 `routes.go` 中由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹，即 promotion-api 所有接口都需要携带有效 JWT。

---

## RPC 服务配置 (`apps/promotion/rpc/etc/promotion.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `promotion.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8088` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `promotion.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_promotion?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**依赖的外部服务**：
- MySQL `tj_promotion` 库（`coupon` / `user_coupon` / `coupon_code` 三张表）
- Redis 缓存（节点模式）
- etcd（服务注册）

**配置结构体**（`apps/promotion/rpc/internal/config/config.go`）：

```go
type Config struct {
    zrpc.RpcServerConf
    DataSource string          `json:",optional"`
    Cache      cache.CacheConf `json:",optional"`
}
```

> `DataSource` 与 `Cache` 均标记为 `optional`，缺省时服务仍可启动，但所有 Model 调用都会失败——生产环境必须显式配置。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_promotion?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_promotion`
- 字符集: `utf8mb4`
- 时区: `Local`

> `parseTime=true` 与 `loc=Local` 是必需项：`coupon` / `user_coupon` 表大量使用 `sql.NullTime` 承载发放期与有效期，逻辑层的 `parseTime` / `userCouponExpireTime` 均按 `time.Local` 计算。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `promotion-api` | `8818` | HTTP |
| `promotion.rpc` | `8088` | gRPC |

### 缓存使用方式

`ServiceContext` 中三个 Model 共用同一份 `c.Cache` 配置：

```go
CouponModel:     model.NewCouponModel(conn, c.Cache)
UserCouponModel: model.NewUserCouponModel(conn, c.Cache)
CouponCodeModel: model.NewCouponCodeModel(conn, c.Cache)
```

| 缓存键前缀 | 承载对象 |
|-----------|---------|
| `cache:coupon:id:` | 优惠券主键缓存 |
| `cache:userCoupon:id:` | 用户券主键缓存 |
| `cache:couponCode:id:` | 兑换码主键缓存 |
| `cache:couponCode:code:` | 兑换码按 `code` 的唯一索引缓存 |

自定义扩展方法（分页、批量、条件更新）一律走 `NoCache` 系列 API，并在写操作后手工 `DelCacheCtx` 失效对应键，避免读到脏数据。

### JWT 密钥

生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造。promotion-api 的 `Auth.AccessSecret` 必须与签发方（auth 服务）保持一致，否则所有接口鉴权失败。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_promotion) | DataSource 配置 | 自建存储 |
| Redis | Cache 配置 | 缓存券/用户券/兑换码主键查询 |
| promotion.rpc | RpcClient 配置 | API 层通过 etcd 服务发现调用 |
| etcd | Etcd 配置 | RPC 服务注册与发现 |

> promotion.rpc 不依赖任何其他服务的 RPC——券适用范围 `scopes` 中的课程分类名称由 course 服务维护，本服务仅回传 id（`toCouponDetailVO` 中 `CouponScopeVO` 只填 `Id`、不填 `Name`），刻意避免跨服务强依赖。
