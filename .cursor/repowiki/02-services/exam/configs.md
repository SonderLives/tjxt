> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/exam/api/etc/exam-api.yaml`, `apps/exam/rpc/etc/exam.yaml`

---

# Exam Configs

## API 服务配置 (`apps/exam/api/etc/exam-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `exam-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8805` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `ExamRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `ExamRpc.Etcd.Key` | `exam.rpc` | - | exam RPC 服务发现 key |

**对应配置结构体**（`apps/exam/api/internal/config/config.go`）：

```go
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	ExamRpc zrpc.RpcClientConf
}
```

**依赖的外部服务**：
- 自身 RPC `exam.rpc`（HTTP handler → RPC client 调用）
- etcd（服务发现）

> `Auth` 配置对应 `apps/exam/api/exam.api` 中两个 service 块的 `jwt: Auth` 声明，全部 7 个 HTTP 接口均需鉴权。

---

## RPC 服务配置 (`apps/exam/rpc/etc/exam.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `exam.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8085` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `exam.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_exam?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**对应配置结构体**（`apps/exam/rpc/internal/config/config.go`）：

```go
type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf
}
```

**依赖的外部服务**：
- MySQL `tj_exam` 库
- Redis 缓存（节点模式）
- etcd（服务注册）

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_exam?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_exam`
- 字符集: `utf8mb4`
- 时区: `Local`

`parseTime=true` 为必需项：`question.create_time` / `update_time` 映射为 Go `time.Time`。

> **字符集注意**：`tj_exam` 库中 3 张业务表为 `utf8mb4 / utf8mb4_0900_ai_ci`，但 `undo_log` 表为 `utf8 / utf8_general_ci`（Seata 官方脚本原样保留），连接串的 `charset=utf8mb4` 对业务表适用。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `exam-api` | 8805 | HTTP |
| `exam.rpc` | 8085 | gRPC |

### 缓存

`Cache` 传入三个 model 构造函数（`apps/exam/rpc/internal/svc/servicecontext.go:21-23`），由 goctl 的 `sqlc.CachedConn` 托管主键缓存：

| Model | 缓存 key 前缀 |
|-------|-------------|
| `question` | `cache:question:id:` |
| `questionDetail` | `cache:questionDetail:id:` |
| `questionBiz` | `cache:questionBiz:id:` |

三个 model 共用同一个 `sqlx.NewMysql(c.DataSource)` 连接与同一份 `Cache` 配置。

### JWT 密钥

生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造。

`exam-api` 的 `Auth.AccessSecret` 必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则所有接口鉴权失败。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_exam) | DataSource 配置 | 自建存储，含 question / question_detail / question_biz / undo_log |
| Redis | Cache 配置 | 缓存三张业务表的主键查询 |
| exam.rpc | RpcClient 配置（`ExamRpc`） | API 层通过 etcd 发现自身 RPC |

---

## ⚠️ 配置缺口

### 1. 无分布式事务（Seata）配置

`sql/ddl/tj_exam.sql` 中定义了 `undo_log` 表（表注释 `AT transaction mode undo table`），这是 Seata AT 模式的事务回滚日志表，但：

| 检查项 | 结果 |
|--------|------|
| `apps/exam/rpc/etc/exam.yaml` | 无任何 Seata / TC 地址配置 |
| `apps/exam/rpc/internal/config/config.go` | 仅 `RpcServerConf` / `DataSource` / `Cache` |
| `apps/exam/` 全目录检索 `undo_log` / `seata` | **零命中** |

> **影响**：`undo_log` 表处于闲置状态。`SaveQuestion` / `DeleteQuestion` 的多表操作用 go-zero 本地事务（`sqlx.Transact`）即可，暂不需要 Seata；但若后续需要与 course 域联动（挂题同时更新课程目录）则需补齐分布式事务配置。详见 [data-model.md](./data-model.md)。

### 2. 无跨服务 RPC 客户端配置

| 检查项 | 结果 |
|--------|------|
| `exam-api` 配置 | 仅 `ExamRpc` 一个客户端 |
| `exam.rpc` 配置 | 无任何 RpcClient 配置 |

> **影响**：`question_biz.biz_id` 语义上指向 course 域的小节 ID，但 exam 服务**无法校验该 ID 是否真实存在**（无 `CourseRpc` 配置）。反向亦然：`apps/course/api/internal/svc/servicecontext.go` 未 import `examclient`，course 侧的 `CourseSubjectsGet` 也无法回调 exam 取题。两域间的 RPC 通道尚未建立。

### 3. 无分页默认值配置

`apps/exam/api/exam.api` 中 `QuestionListReq.PageNo` / `PageSize` 与 `QuestionBizListReq.PageNo` / `PageSize` 均标注 `optional`，且 yaml 中**无分页默认值配置项**。

> **影响**：客户端不传时为零值，logic 层必须硬编码兜底（如 `pageNo=1, pageSize=20`），否则分页查询会退化为 `LIMIT 0`。

### 4. 与其他服务的端口/配置对照

| 服务 | API 端口 | RPC 端口 | 备注 |
|------|---------|---------|------|
| `auth` | 8802 | 8082 | 有 `Jwt` 完整配置（含 Refresh） |
| `media` | 8806 | 8086 | 缺对象存储配置 |
| `exam` | **8805** | **8085** | 配置最精简，仅数据库 + 缓存 |

> exam 是三者中配置项最少的服务：**无外部第三方依赖**，所有配置项均已就位，配置层**不构成实现阻塞**（对比 media 的对象存储配置缺口）。当前的实现阻塞项集中在 model 扩展方法层面，详见 [data-model.md](./data-model.md) 的缺口表。
