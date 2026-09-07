# observability（可观测性）Specification

## Purpose

定义 trace / metrics / logs 三类信号统一经 `otel-collector` 收口的可观测性体系：go-zero 原生 OTLP → Jaeger、scrape 代理 → Prometheus、filelog → Loki，业务代码零改动。

## Requirements

### Requirement: 统一收口架构

全部 26 个 Go 服务进程（13 服务 × {api,rpc}，宿主机运行）的三类信号 SHALL 统一经一个 `otel-collector`（`otel/opentelemetry-collector-contrib`）转发：

| 信号 | 服务侧 | collector | 后端 |
|------|--------|-----------|------|
| Trace | go-zero `Telemetry` OTLP/HTTP 推 `127.0.0.1:4318` | `otlp` receiver → `batch` → `otlphttp/jaeger` | Jaeger（UI 16686） |
| Metrics | 各服务暴露 `/metrics`（RPC 9101-9113 / API 9201-9213） | `prometheus` receiver 抓取 → `:8889` 聚合 | Prometheus（UI 9090，只抓 `otel-collector:8889`） |
| Logs | go-zero 写 JSON 日志到 `logs/<svc>/*.log` | `filelog` receiver（挂载 `/var/log/tjxt`）→ `loki` exporter | Loki（3100） |

- 设计约束：go-zero v1.10.3 无 OTLP metrics exporter、无 Loki/OTLP logs writer，因此 metrics 走 scrape 代理、logs 走 filelog，对业务代码零侵入
- 容器访问宿主机服务统一用 `host.docker.internal`（配合 `host-gateway`）

#### Scenario: Prometheus 抓取路径

- **WHEN** Prometheus 抓取指标
- **THEN** 仅抓取 `otel-collector:8889`（聚合暴露，开启 `resource_to_telemetry_conversion`）
- **AND** Prometheus 不直接 target 26 个服务实例

### Requirement: 服务 yaml 注入约定（勿删）

每个服务 `apps/<svc>/{api,rpc}/etc/*.yaml` SHALL 保留以下三段注入配置，禁止删除或改坏字段：

```yaml
Log:
  Mode: file
  Encoding: json
  Path: logs/<svc>      # 相对仓库根，每服务唯一
  Level: info
Prometheus:
  Host: 0.0.0.0
  Port: <唯一>          # RPC 9101-9113 / API 9201-9213
  Path: /metrics
Telemetry:
  Name: <svc>
  Endpoint: 127.0.0.1:4318
  Sampler: 1.0
  Batcher: otlphttp
```

- `Log.Path` 每服务唯一（go-zero JSON 日志无 service 字段，服务名从文件路径提取）
- `Prometheus.Port` 每服务唯一，避免 collector 抓取端口冲突
- 重新生成骨架（`goctl api go` / `goctl rpc protoc`）不会覆盖这些配置（goctl 只刷新 types/routes/pb/server）

#### Scenario: 从仓库根启动服务

- **WHEN** 启动任一 Go 服务
- **THEN** 工作目录必须为仓库根（`Log.Path` 为相对路径），日志落到 `<repo>/logs/<svc>/` 才能被 collector filelog 采集
- **AND** 若从服务子目录启动，日志散落导致采集不到，视为部署错误

### Requirement: 日志采集管线

日志管线 SHALL 按 filelog 方案采集：

1. docker-compose 将宿主 `./logs` 只读挂载进 collector：`/var/log/tjxt`
2. `filelog` receiver 读 `/var/log/tjxt/**/*.log`，解析 JSON（go-zero 字段固定 `@timestamp`/`level`/`content`）
3. `transform/logsvc` 把 service（从路径提取）/level 提升为 resource 属性
4. `loki` exporter 推到 `loki:3100`，流标签 `service_name`、`level`

实施约束：

- filelog `timestamp.layout` 必须用 Go 时间布局 `2006-01-02T15:04:05.000Z07:00`，不能用 strftime
- Loki exporter `labels` 只接受 resource/scope 属性；`level` 须用 transform 提升为 resource 属性才能做标签

#### Scenario: 按服务查询日志

- **WHEN** 开发者查询 course-api 的日志
- **THEN** 通过 Loki 查询 `{service_name="course-api"}` 获得该服务全部日志流

### Requirement: 配置文件清单

可观测性相关配置 SHALL 由以下文件承载：

| 文件 | 作用 |
|------|------|
| `docker-compose.yml` | 起 jaeger / otel-collector / prometheus / loki 四个可观测容器 |
| `deploy/otel-collector/config.yaml` | collector 全配置：3 receiver、2 processor、3 条管线（traces/metrics/logs） |
| `deploy/prometheus/prometheus.yml` | Prometheus 仅抓 `otel-collector:8889` + 自身 |
| `deploy/loki/loki-config.yaml` | Loki 单实例最小配置（filesystem + inmemory ring，`auth_enabled:false`） |
| `apps/**/etc/*.yaml`（×26） | 各服务注入 `Telemetry` / `Prometheus` / `Log` 三段 |

#### Scenario: 重新生成骨架不丢配置

- **WHEN** 开发者对某服务执行 `goctl api go` 或 `goctl rpc protoc` 重新生成骨架
- **THEN** `deploy/` 下的 collector/Prometheus/Loki 配置与服务 yaml 中的三段注入配置均不受影响（goctl 只刷新 types/routes/pb/server）

### Requirement: 端口与起停

可观测性组件的端口分配与启停 SHALL 遵循以下约定：

- `otel-collector` 暴露 `4317`（OTLP gRPC）/ `4318`（OTLP HTTP）/ `8889`（聚合 metrics）/ `13133`（health_check）
- Jaeger 容器不再向宿主机映射 4317/4318（已让位给 collector），其 OTLP 接收仅在 docker 内网对 collector 开放
- 启停：`make docker-up` 一并拉起可观测性四件套；`make docker-logs` 跟踪容器日志

#### Scenario: 本地起停与查询

- **WHEN** 开发者执行 `make docker-up` 后从仓库根启动服务
- **THEN** Trace 可在 http://localhost:16686 按 service 查询 span
- **AND** Metrics 可在 http://localhost:9090 查看（`tjxt-services` target 应 UP）
- **AND** Logs 可经 Loki http://localhost:3100 查询（建议接 Grafana）

