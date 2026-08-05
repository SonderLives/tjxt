> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/course/api/etc/course-api.yaml`, `apps/course/rpc/etc/course.yaml`

---

# Course Configs

## API 服务配置 (`apps/course/api/etc/course-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `course-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8812` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `CourseRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `CourseRpc.Etcd.Key` | `course.rpc` | - | course RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `course.rpc`（HTTP handler → RPC client 调用）

---

## RPC 服务配置 (`apps/course/rpc/etc/course.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `course.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8083` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `course.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_course?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**依赖的外部服务**：
- MySQL `tj_course` 库
- Redis 缓存（节点模式）

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_course?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_course`
- 字符集: `utf8mb4`
- 时区: `Local`

### JWT 密钥

JWT 密钥与 auth 服务共用同一值 `Auth.AccessSecret`，确保各服务间可互相校验 Token。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_course) | DataSource 配置 | 自建存储 |
| Redis | Cache 配置 | 缓存分类/课程数据 |