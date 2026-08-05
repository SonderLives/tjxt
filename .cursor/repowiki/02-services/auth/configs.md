> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/auth/api/etc/auth-api.yaml`, `apps/auth/rpc/etc/auth.yaml`

---

# Auth Configs

## API 服务配置 (`apps/auth/api/etc/auth-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `auth-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8802` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `AuthRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `AuthRpc.Etcd.Key` | `auth.rpc` | - | 服务发现 key |
| `UserRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `UserRpc.Etcd.Key` | `user.rpc` | - | user RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `auth.rpc`（HTTP handler → RPC client 调用）
- `user.rpc` — 用户域 RPC，用于登录时调用 `LoginVerify` 验证凭证

---

## RPC 服务配置 (`apps/auth/rpc/etc/auth.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `auth.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8082` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `auth.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_auth?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |
| `Jwt.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Jwt.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `Jwt.RefreshSecret` | `change-me-refresh-in-production` | - | 刷新令牌密钥 |
| `Jwt.RefreshExpire` | `604800` | - | 刷新令牌有效期（秒）= 7 天 |

**依赖的外部服务**：
- MySQL `tj_auth` 库
- Redis 缓存（节点模式）

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_auth?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_auth`
- 字符集: `utf8mb4`
- 时区: `Local`

### JWT 密钥

生产环境**必须修改**默认值，否则 JWT 可被伪造。

| 密钥 | 用途 |
|------|------|
| `AccessSecret` | 签发和校验 `accessToken` |
| `RefreshSecret` | 签发和校验 `refreshToken` |

两个密钥应不同，实现令牌的独立轮换。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_auth) | DataSource 配置 | 自建存储 |
| Redis | Cache 配置 | 缓存角色/菜单/权限数据 |
| user.rpc | RpcClient 配置 | 登录时验证用户凭证 |