> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/media/api/etc/media-api.yaml`, `apps/media/rpc/etc/media.yaml`

---

# Media Configs

## API 服务配置 (`apps/media/api/etc/media-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `media-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8806` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `MediaRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `MediaRpc.Etcd.Key` | `media.rpc` | - | media RPC 服务发现 key |

**对应配置结构体**（`apps/media/api/internal/config/config.go`）：

```go
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// MediaRpc 媒资 RPC 客户端配置（通过 etcd 服务发现）
	MediaRpc zrpc.RpcClientConf
}
```

**依赖的外部服务**：
- 自身 RPC `media.rpc`（HTTP handler → RPC client 调用）
- etcd（服务发现）

> API 层三组路由（`media` / `signature` / `file`）在 `apps/media/api/media.api` 中均声明了 `jwt: Auth`，全部接口需要携带有效 JWT。

---

## RPC 服务配置 (`apps/media/rpc/etc/media.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `media.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8086` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `media.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_media?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |

**对应配置结构体**（`apps/media/rpc/internal/config/config.go`）：

```go
type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf
}
```

**依赖的外部服务**：
- MySQL `tj_media` 库
- Redis 缓存（节点模式）
- etcd（服务注册）

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_media?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_media`
- 字符集: `utf8mb4`
- 时区: `Local`

`parseTime=true` 为必需项：`file` / `media` 两表的 `create_time` / `update_time` 均映射为 Go `time.Time`。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `media-api` | 8806 | HTTP |
| `media.rpc` | 8086 | gRPC |

### 缓存

`Cache` 传入 `model.NewMediaModel(conn, c.Cache)` / `model.NewFileModel(conn, c.Cache)`，由 goctl 的 `sqlc.CachedConn` 托管主键缓存：

| Model | 缓存 key 前缀 |
|-------|-------------|
| `media` | `cache:media:id:` |
| `file` | `cache:file:id:` |

### JWT 密钥

生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造。

`media-api` 的 `Auth.AccessSecret` 必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则所有接口鉴权失败。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_media) | DataSource 配置 | 自建存储 |
| Redis | Cache 配置 | 缓存 media/file 主键查询 |
| media.rpc | RpcClient 配置（`MediaRpc`） | API 层通过 etcd 发现自身 RPC |

---

## ⚠️ 对象存储配置缺口

`media` 服务的核心能力是**媒资上传与签名授权**，但当前配置中**完全没有对象存储相关配置项**。

### 缺失证据

| 检查项 | 结果 |
|--------|------|
| `apps/media/rpc/etc/media.yaml` | 仅 `Name` / `ListenOn` / `Etcd` / `DataSource` / `Cache`，无存储配置 |
| `apps/media/rpc/internal/config/config.go` | 仅 `zrpc.RpcServerConf` / `DataSource` / `Cache` 三个字段 |
| `apps/media/api/etc/media-api.yaml` | 仅 `Name` / `Host` / `Port` / `Auth` / `MediaRpc` |
| `apps/media/api/internal/config/config.go` | 仅 `rest.RestConf` / `Auth` / `MediaRpc` |

### 业务侧对配置的实际需求

| 需求来源 | 需要的配置 | 说明 |
|---------|-----------|------|
| `sql/ddl/tj_media.sql` — `file.platform` 注释「1-腾讯，2-阿里」 | 双平台凭据 | 至少需要腾讯云与阿里云两套 SecretId / SecretKey |
| `media.proto` — `SignatureUpload` 返回 `token` / `uploadUrl` | 上传签名参数 | Bucket、Region、上传路径前缀、签名有效期 |
| `media.proto` — `SignaturePlay` / `SignaturePreview` 返回 `playUrl` | 点播播放参数 | 点播 AppId、播放域名、防盗链 Key、签名有效期 |
| `media.file_id` — 云端视频唯一标示 | 点播服务配置 | 视频点播（VOD）子应用 ID |
| `FileVO.path`（表中无此列） | 访问地址前缀 | 由配置的 CDN/对象存储域名 + `file.key` 拼接 |

### 建议补齐的配置结构（尚未存在，仅为缺口说明）

```yaml
# 以下配置项当前 media.yaml 中【不存在】，需在实现签名功能前补齐
# Storage:
#   Platform: 1              # 1-腾讯 2-阿里，对应 file.platform
#   SecretId: ""
#   SecretKey: ""
#   Region: ""
#   Bucket: ""
#   UploadPathPrefix: ""     # 上传路径前缀，用于生成 file.key
#   CdnDomain: ""            # 拼接 FileVO.path / MediaVO.mediaUrl
#   Vod:
#     AppId: 0               # 点播子应用，对应 media.file_id 所属应用
#     PlayDomain: ""
#     PlayKey: ""            # 播放防盗链密钥
#     SignExpireSec: 3600    # 播放/预览签名有效期
```

> **结论**：在补齐上述配置结构体与 yaml 条目之前，`SignatureUpload` / `SignaturePreview` / `SignaturePlay` 三个 RPC 方法**无法实现**；`FileVO.path` 与 `MediaVO.mediaUrl` / `coverUrl` 也缺少地址拼接依据。这是 media 服务当前最主要的配置层阻塞项。
