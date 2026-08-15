# 可观测性：统一的 OpenTelemetry Collector 中枢

> 版本：v1.0 | 更新：2026-08-15 | 来源：`docker-compose.yml`、`deploy/otel-collector/config.yaml`、`deploy/prometheus/prometheus.yml`、`deploy/loki/loki-config.yaml`、26 个 `apps/**/etc/*.yaml`

## 1. 总览

天机学堂的 **trace / metrics / logs 三类信号统一经一个 `otel-collector` 收口**，再分别落到 Jaeger / Prometheus / Loki。业务代码**零改动**——所有接入都只靠「服务 yaml 注入配置 + collector 配置 + docker-compose」三项完成。

```
┌──────────────────────────────────────────────────────────────────────────┐
│  26 个 Go 服务（宿主机运行，非容器）                                         │
│  Telemetry → 127.0.0.1:4318 (OTLP/HTTP)                                   │
│  Prometheus /metrics (9101-9113 / 9201-9213)                              │
│  Log file → logs/<svc>/*.log (JSON)                                       │
└───────┬───────────────────────┬───────────────────────┬──────────────────┘
        │ trace (OTLP/HTTP 4318)│ metrics (scrape)       │ logs (file mount)
        ▼                       ▼                       ▼
   ┌─────────────────────────────────────────────────────────┐
   │              otel-collector (统一中枢)                    │
   │  receivers: otlp / prometheus / filelog                   │
   │  processors: batch / transform                           │
   │  exporters: otlphttp/jaeger / prometheus / loki           │
   └───┬───────────────────┬───────────────────┬─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
   Jaeger:16686      Prometheus:9090       Loki:3100
   (trace UI)        (metrics UI)         (logs store, 建议接 Grafana)
```

设计原则：**基建（collector/Jaeger/Prometheus/Loki）走 Docker，业务服务在宿主机跑**——这与 etcd / ES / RabbitMQ 一贯的「Go 本机跑、基建 docker 化」约定一致。

---

## 2. 三类信号的管线

### 2.1 链路 Trace（go-zero 原生 OTLP → Collector → Jaeger）

- 各服务 yaml 注入 `Telemetry{Name, Endpoint:127.0.0.1:4318, Sampler:1.0, Batcher:otlphttp}`。
- go-zero 通过 `trace.StartAgent` 自动把 span 以 OTLP/HTTP 推到 `127.0.0.1:4318`（即本机 otel-collector 容器映射端口）。
- collector：`otlp` receiver (4317/4318) → `batch` → `otlphttp/jaeger` exporter → `jaeger:4318`（docker 内网）。
- Jaeger UI：`http://localhost:16686`。

> 注意：Jaeger 容器**不再**把 4317/4318 映射给宿主机（已让位给 collector）；其 OTLP 接收仅在 docker 内网对 collector 开放。

### 2.2 指标 Metrics（scrape 代理：服务 /metrics → Collector → Prometheus）

- 各服务 yaml 注入 `Prometheus{Host:0.0.0.0, Port:<唯一>, Path:/metrics}`，仍各自暴露 `/metrics`，**端口不再被 Prometheus 直接抓**。
- collector：`prometheus` receiver 通过 `host.docker.internal` 抓取 26 个宿主机服务的 `/metrics`（RPC 指标端口 9101–9113，API 指标端口 9201–9213）。
- `prometheus` exporter 在 `:8889` 聚合暴露（开启 `resource_to_telemetry_conversion`）。
- Prometheus 只抓 `otel-collector:8889`（见 `deploy/prometheus/prometheus.yml`，不再直接 target 26 个服务）。

> **为什么用 scrape 代理而非 OTLP metrics 直推？** go-zero v1.10.3 **没有 OTLP metrics exporter**，只有 Prometheus `/metrics`。若要从服务侧直推 OTLP，需要自己写 OTel SDK + 桥接 go-zero 的 client_golang 指标 + `prometheusremotewrite` exporter + 给 Prometheus 开 `--web.enable-remote-write-receiver`——而后端本就是 Prometheus，等于绕一圈又回来。scrape 代理零业务代码、最稳。

### 2.3 日志 Logs（file 采集：go-zero 写 JSON → Collector filelog → Loki）

- 各服务 yaml 注入 `Log{Mode:file, Encoding:json, Path:logs/<svc>, Level:info}`。
- go-zero `file` 模式把 `access.log` / `error.log` / `slow.log` / `stat.log` / `severe.log` 写到 `Path` 下；**相对路径 `logs/<svc>` 相对进程 CWD（仓库根）**，落到 `<repo>/logs/<svc>/`。
- docker-compose 把宿主 `./logs` 只读挂进 collector 容器：`/var/log/tjxt`。
- collector：`filelog` receiver 读 `/var/log/tjxt/**/*.log` → 解析 JSON → `transform/logsvc` 把 service/level 提升为 resource 属性 → `loki` exporter 推到 `loki:3100`。
- Loki 流标签：`service_name`、`level`。查询示例：`{service_name="course-api"}`。

> **为什么用 filelog 而非 OTLP logs 直推？** go-zero **没有 Loki / OTLP logs writer**，只能 console/file 原生日志。filelog 方案同样零业务代码。

---

## 3. 配置文件清单

| 文件 | 作用 |
|------|------|
| `docker-compose.yml` | 起 jaeger / otel-collector / prometheus / loki 四个可观测容器（外加 mysql/redis/rabbitmq/etcd 基础依赖） |
| `deploy/otel-collector/config.yaml` | collector 全配置：3 个 receiver、2 个 processor、3 条管线（traces/metrics/logs） |
| `deploy/prometheus/prometheus.yml` | Prometheus 仅抓 `otel-collector:8889` + 自身 |
| `deploy/loki/loki-config.yaml` | Loki 单实例最小配置（filesystem + inmemory ring，`auth_enabled:false`） |
| `apps/**/etc/*.yaml`（×26） | 各服务注入 `Telemetry` / `Prometheus` / `Log` 三段 |

### 3.1 服务 yaml 注入约定（以 `apps/course/api/etc/course-api.yaml` 为例）

```yaml
Name: course-api

Log:
  Mode: file
  Encoding: json
  Path: logs/course-api      # 相对仓库根；collector 经 ./logs 挂载采集
  Level: info

Prometheus:
  Host: 0.0.0.0
  Port: 9203                  # 唯一；API 类 9201-9213，RPC 类 9101-9113
  Path: /metrics

Telemetry:
  Name: course-api
  Endpoint: 127.0.0.1:4318    # 本机 collector（OTLP/HTTP）
  Sampler: 1.0
  Batcher: otlphttp

Host: 0.0.0.0
Port: 8803

# …… 其余业务配置（Auth / *Rpc / Cache 等）
```

- `Log.Path` 每个服务**唯一**（如 `user-rpc` / `course-api`），用于从目录区分服务（go-zero JSON 日志不含 service 字段）。
- `Prometheus.Port` 每个服务**唯一**，避免 collector 抓取时端口冲突。

---

## 4. 起停与查询

### 4.1 启动

```bash
make docker-up        # 等价于 docker compose up -d
                       # 现在会一并拉起 jaeger / otel-collector / prometheus / loki
# 确认容器健康：
make docker-logs      # 跟踪依赖容器日志；collector 健康检查 http://localhost:13133
```

> ⚠️ `make docker-up` 的 help 注释仍写「启动 MySQL/Redis/RabbitMQ/etcd」，实际它启动 compose 中**全部**服务（含可观测性四件套）。

然后**从仓库根**启动各 Go 服务（CWD=仓库根，日志才能落到 `<repo>/logs/<svc>` 被采集）：

```bash
make run-user                 # 前台单服务
make run-all / run-all-rpc    # 后台并行启动全部 API / RPC
# 或手动（务必仓库根）：
go run ./apps/course/api -f apps/course/api/etc/course-api.yaml
```

### 4.2 查询界面

| 信号 | 地址 | 说明 |
|------|------|------|
| Trace | http://localhost:16686 | Jaeger UI，按 service=`course-api` 查 span |
| Metrics | http://localhost:9090 | Prometheus UI，`tjxt-services` target 应 UP |
| Logs | http://localhost:3100 | Loki HTTP API；建议接 Grafana 做查询/面板 |

---

## 5. 端口与依赖关系

- `otel-collector` 暴露：`4317`(OTLP gRPC) / `4318`(OTLP HTTP) / `8889`(聚合 metrics) / `13133`(health_check)。
- `otel-collector` 通过 `extra_hosts: host.docker.internal:host-gateway` 访问宿主机上的 26 个 Go 服务；`depends_on: [jaeger, loki]`。
- `prometheus` 同样配 `host.docker.internal`；`depends_on: [jaeger]`（无强依赖，仅排序）。
- `loki` 暴露 `3100`，无外部依赖。

---

## 6. 关键设计决策与踩坑

| 主题 | 结论 / 坑 |
|------|-----------|
| metrics 走 scrape 代理 | go-zero 无 OTLP metrics exporter；scrape 代理零代码、与 Prometheus 后端最契合 |
| logs 走 filelog | go-zero 无 Loki/OTLP logs writer；filelog 零代码 |
| `Log.Path` 必须相对 + 仓库根启动 | 绝对路径无法用 `${ENV}` 参数化（go-zero v1.10.3 config 不支持占位）；相对路径靠 CWD 解析，务必仓库根启动 |
| go-zero JSON 日志字段 | 固定 `@timestamp` / `level` / `content`，**无 service、无 caller** → 服务名只能从文件路径（`regex_parser`）拿 |
| filelog `timestamp.layout` | **必须用 Go 时间布局** `2006-01-02T15:04:05.000Z07:00`，不能用 strftime `%Y-%m-%d` |
| Loki exporter `labels` | 只接受 **resource / scope** 属性；`level` 是 log-record 属性，须用 `transform` 提升为 `resource.attributes["level"]` 才能做标签 |
| collector 镜像版本 | `otel/opentelemetry-collector-contrib:latest` 未锁版本，生产建议 pin 固定 tag |
| 跨进程网络 | 容器访问宿主机服务统一用 `host.docker.internal`（配合 `host-gateway`） |
| 指标端口分配 | RPC 9101–9113、API 9201–9213，连续唯一，已在 26 个 yaml 同步 |

---

## 7. 相关文档

- [Docker 基础设施](../04-infra/docker-compose.md)
- [架构全览](../00-architecture/overview.md)
- [实现进度与已知缺口](../06-status/implementation-status.md)（§2.8 可观测性）
- [本地开发快速上手](../05-development/quickstart.md)（§2 可观测性起停）
