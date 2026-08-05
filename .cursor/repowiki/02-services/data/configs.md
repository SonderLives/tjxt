> 版本：v1.0 | 更新：2026-08-05 | 来源：apps/data/api/data/etc/data-api.yaml、apps/data/rpc/data/etc/data.yaml、apps/data/api/go.mod、apps/data/api/data/go.mod、apps/data/rpc/data/go.mod

# data（统计大屏）服务配置说明

本文档逐字说明 data 服务各 `etc/*.yaml` 配置项与外部依赖。所有内容均来自真实源码，未做任何臆造。

## 一、服务端口与进程

| 进程 | 配置文件 | 关键配置 | 端口 / 监听 |
| --- | --- | --- | --- |
| API 网关 | `apps/data/api/data/etc/data-api.yaml` | `Name: data-api`、`Port: 8811` | HTTP `0.0.0.0:8811` |
| RPC 服务 | `apps/data/rpc/data/etc/data.yaml` | `Name: data.rpc`、`ListenOn: 0.0.0.0:8091` | gRPC `0.0.0.0:8091` |

- API 与 RPC 通过 etcd 服务发现对接：`DataRpc.Etcd.Key: data.rpc`，etcd 地址 `127.0.0.1:2379`。
- 与 message、search 不同，data 的 RPC 配置**没有** `DataSource`、`Cache` 段（见下方「已知结构性问题」）。

## 二、API 配置逐项说明（data-api.yaml）

源文件 `apps/data/api/data/etc/data-api.yaml` 共 13 行，逐项正确如下：

| 配置项 | 取值 | 说明 |
| --- | --- | --- |
| `Name` | `data-api` | 服务名 |
| `Host` | `0.0.0.0` | 监听地址 |
| `Port` | `8811` | HTTP 端口 |
| `Auth.AccessSecret` | `"change-me-in-production"` | JWT 签名密钥（占位值，生产须替换） |
| `Auth.AccessExpire` | `7200` | JWT 有效期（秒） |
| `DataRpc.Etcd.Hosts` | `127.0.0.1:2379` | etcd 地址 |
| `DataRpc.Etcd.Key` | `data.rpc` | 下游 RPC 服务发现 key |

注意：该文件虽声明了 `Auth.AccessSecret` / `Auth.AccessExpire`，但 `data.api` 中**所有 @server 块均未带 `jwt: Auth`**（见 api-spec.md），即 JWT 配置已声明却未被路由启用——属于「已配置未使用」的悬空项。

## 三、RPC 配置逐项说明（data.yaml）

源文件 `apps/data/rpc/data/etc/data.yaml` 共 7 行，仅三段：

| 配置项 | 取值 | 说明 |
| --- | --- | --- |
| `Name` | `data.rpc` | 服务名 |
| `ListenOn` | `0.0.0.0:8091` | gRPC 监听地址 |
| `Etcd.Hosts` | `127.0.0.1:2379` | etcd 地址 |
| `Etcd.Key` | `data.rpc` | 自身注册 key |

对比 message / search 的 RPC 配置，data 的 `data.yaml` **缺失 `DataSource`（数据库连接）与 `Cache`（Redis）两段**。这与 data 服务「不使用关系型表、不依赖 Model」的设计一致（详见 data-model.md），但意味着目前 RPC 没有可落地的数据源。

## 四、外部依赖

| 依赖 | 用途 | 配置位置 | 状态 |
| --- | --- | --- | --- |
| etcd | API↔RPC 服务发现 | `DataRpc.Etcd` / `Etcd` | 已配置（`127.0.0.1:2379`） |
| MySQL | 关系型存储 | 无（data.yaml 无 DataSource） | 📋 设计意图：统计大屏数据走 Redis，无 MySQL |
| Redis | 大屏数据缓存 | 无（data.yaml 无 Cache） | 📋 设计意图：待实现，Redis Key 见 data-model.md |
| JWT | 接口鉴权 | `Auth.*`（已声明） | ⚠️ 已声明但路由未启用 |

## 五、go.mod 模块结构

data 服务存在**三份** go.mod，模块名与 go.work 引用不一致：

| 路径 | module 名 | 备注 |
| --- | --- | --- |
| `apps/data/api/go.mod` | `api` | ⚠️ 内容残缺，仅 3 行（`module api` / `go 1.26.2` / 空行），无 require，疑似占位残留 |
| `apps/data/api/data/go.mod` | `tjxt/apps/data/api/data` | 实际被 go.work 引用 |
| `apps/data/rpc/data/go.mod` | `tjxt/apps/data/rpc/data` | 实际被 go.work 引用 |

`go.work` 引用的真实模块是后两者（`./apps/data/api/data` 与 `./apps/data/rpc/data`），而非 `apps/data/api` 与 `apps/data/rpc`。

## ⚠️ 六、已知结构性问题

data 服务目录与约定存在多处偏差，逐条如实记录：

1. **目录多套一层**
   - 约定：proto / api / rpc 应位于 `apps/data/api/`、`apps/data/rpc/`。
   - 实际：`data.proto` 在 `apps/data/rpc/data/data.proto`，API 代码在 `apps/data/api/data/...`，比约定多嵌套一层 `data/`。

2. **双（实为三）go.mod**
   - `apps/data/api/go.mod`（module 名 `api`，仅 3 行、无 require，残缺占位）；
   - `apps/data/api/data/go.mod`（module `tjxt/apps/data/api/data`）；
   - `apps/data/rpc/data/go.mod`（module `tjxt/apps/data/rpc/data`）。
   - go.work 仅引用后两者，第一份 `api/go.mod` 为遗留废弃文件。

3. **无 DDL**
   - 全仓 `sql/ddl/` 下存在 `tj_message.sql`、`tj_message_model.sql`、`tj_search.sql`，但**没有** `tj_data.sql`。
   - 与「data 不使用关系型表」一致，但缺乏独立的存储说明文档。

4. **配置缺 DataSource / Cache**
   - `data.yaml` 仅有 Name / ListenOn / Etcd 三段，无数据库连接与 Redis 配置，RPC 暂无可用数据源配置。

5. **JWT 配置悬空**
   - `data-api.yaml` 声明了 `Auth.AccessSecret` / `AccessExpire`，但 `data.api` 路由无 `jwt: Auth`，鉴权未生效（详见 business-rules.md「无鉴权安全缺口」）。
