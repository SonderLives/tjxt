> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/search/api/etc/search-api.yaml`, `apps/search/rpc/etc/search.yaml`

---

# Search Configs

## API 服务配置 (`apps/search/api/etc/search-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `search-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8810` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `SearchRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `SearchRpc.Etcd.Key` | `search.rpc` | - | search RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `search.rpc`（HTTP handler → RPC client 调用，`apps/search/api/internal/svc/servicecontext.go` 注入 `SearchRpc`）
- etcd（服务发现）

> 结构体定义见 `apps/search/api/internal/config/config.go`：`rest.RestConf` + 匿名 `Auth{AccessSecret, AccessExpire}` + `SearchRpc zrpc.RpcClientConf`。**未配置 MySQL / Redis**，API 层不直连存储。

---

## RPC 服务配置 (`apps/search/rpc/etc/search.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `search.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8090` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `search.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_search?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**依赖的外部服务**：
- MySQL `tj_search` 库（单表 `interests`）
- Redis 缓存（节点模式）
- etcd（服务注册）

> 结构体定义见 `apps/search/rpc/internal/config/config.go`：`zrpc.RpcServerConf` + `DataSource string` + `Cache cache.CacheConf`。**未配置 JWT 段，也未配置任何 RPC 客户端**——RPC 层不鉴权、不外调。

---

## ⚠️ 缺失的外部依赖配置

| 预期依赖 | 配置现状 | 影响 |
|---------|---------|------|
| Elasticsearch | **无任何配置项**，`config.go` 中无 ES 字段，`go.mod` 中无 ES 客户端依赖（全仓 Grep `elastic` / `olivere` 零命中） | 服务当前不具备任何全文检索能力，实际只是 `interests` 单表的 CRUD 外壳 |
| course RPC 客户端 | 无 | `GET /interests/{id}/courses`（openapi 中已规划的按分类查课程 TOP10）无法实现 |
| MQ / Kafka | 无 | 无法消费 course 域事件做索引同步 |

引入 ES 时需同步改动：`apps/search/rpc/etc/search.yaml`（新增 ES 段）、`apps/search/rpc/internal/config/config.go`（新增字段）、`apps/search/rpc/internal/svc/servicecontext.go`（注入 client）、`apps/search/go.mod`（新增依赖）。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_search?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_search`
- 字符集: `utf8mb4`
- 时区: `Local`

> 连接串字符集 `utf8mb4` 与 `sql/ddl/tj_search.sql` 中 `interests` 表的建表字符集 `utf8mb4` 一致（search 是本仓少数字符集完全对齐的服务）。

### JWT 密钥

search 服务**只在 API 层校验 JWT，不签发令牌**。`Auth.AccessSecret` 必须与 auth 服务的 `Jwt.AccessSecret` 保持一致，否则 `/interests` 两个接口全部 401。

| 密钥 | 用途 |
|------|------|
| `Auth.AccessSecret` | 校验 `accessToken`（`routes.go` 中 1 处 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)`） |

`Auth.AccessExpire` 在 search 服务中仅作为配置项被反序列化，代码内无签发逻辑引用。生产环境**必须修改**默认值 `change-me-in-production`。

### 端口分配

| 组件 | 端口 | 说明 |
|------|------|------|
| `search-api` | 8810 | HTTP，对外 |
| `search.rpc` | 8090 | gRPC，集群内 |

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_search) | DataSource 配置 | 自建存储，仅 `InterestsModel` 一个模型使用 |
| Redis | Cache 配置 | `InterestsModel` 走 `sqlc.CachedConn`，缓存前缀 `cache:interests:id:` |
| etcd | Etcd 配置 | RPC 端注册 `search.rpc`，API 端按同名 key 发现 |
| auth 服务 | 无直接连接 | 仅共享 `AccessSecret` 完成离线 JWT 校验，不调用 auth RPC |
| Elasticsearch | **未接入** | 见上文「⚠️ 缺失的外部依赖配置」 |
