# 本地依赖与 Docker Compose

> 版本：v1.1 | 更新：2026-08-15

## docker-compose.yml 服务清单

```yaml
services:
  # ---- 基础依赖 ----
  mysql:
    image: mysql:8.0
    ports: ["3306:3306"]
    environment:
      MYSQL_ROOT_PASSWORD: "0000"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./sql/migration:/docker-entrypoint-initdb.d:ro

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
    volumes: ["redis-data:/data"]

  rabbitmq:
    image: rabbitmq:3.13-management
    ports: ["5672:5672", "15672:15672"]
    environment:
      RABBITMQ_DEFAULT_USER: rabbitmq
      RABBITMQ_DEFAULT_PASS: rabbitmq
    volumes: ["rabbitmq-data:/var/lib/rabbitmq"]

  etcd:
    image: bitnami/etcd:3.5
    ports: ["2379:2379"]
    environment:
      ALLOW_NONE_AUTHENTICATION: "yes"

  # ---- 可观测性（三支柱统一经 otel-collector 收口）----
  # Jaeger：链路追踪 UI（span 由 otel-collector 经 OTLP 转发进来）
  jaeger:
    image: jaegertracing/all-in-one:1.57
    environment:
      COLLECTOR_OTLP_ENABLED: "true"
    ports: ["16686:16686"]   # Jaeger UI

  # OpenTelemetry Collector：trace + metrics + logs 中枢
  #   - traces ：go-zero(otlphttp) → 127.0.0.1:4318 → jaeger:4318(内网)
  #   - metrics：prometheus receiver 抓 26 个宿主机 /metrics → 8889 聚合暴露给 Prometheus
  #   - logs   ：filelog 读 ./logs（挂为 /var/log/tjxt）→ loki exporter → loki:3100
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    command: ["--config=/etc/otelcol-contrib/config.yaml"]
    ports: ["4317:4317", "4318:4318", "8889:8889", "13133:13133"]
    extra_hosts: ["host.docker.internal:host-gateway"]   # 抓宿主机上的 Go 服务
    volumes:
      - ./deploy/otel-collector/config.yaml:/etc/otelcol-contrib/config.yaml:ro
      - ./logs:/var/log/tjxt:ro                          # 各 Go 服务写出的 JSON 日志
    depends_on: ["jaeger", "loki"]

  # Prometheus：抓取 otel-collector:8889（聚合后的 metrics）
  prometheus:
    image: prom/prometheus:v2.53.1
    ports: ["9090:9090"]
    extra_hosts: ["host.docker.internal:host-gateway"]
    volumes:
      - ./deploy/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
    depends_on: ["jaeger"]

  # Loki：日志存储（接收 otel-collector 推送）
  loki:
    image: grafana/loki:latest
    command: ["-config.file=/etc/loki/loki-config.yaml"]
    ports: ["3100:3100"]
    volumes:
      - ./deploy/loki/loki-config.yaml:/etc/loki/loki-config.yaml:ro

volumes:
  mysql-data:
  redis-data:
  rabbitmq-data:
```

## 启动与停止

```bash
# 启动全部依赖（含可观测性栈）
make docker-up          # = docker compose up -d

# 查看容器日志
make docker-logs        # = docker compose logs -f

# 停止并删除容器
make docker-down        # = docker compose down
```

## 端口占用表

| 服务 | 容器端口 | 宿主机端口 | 用途 |
|------|----------|------------|------|
| MySQL | 3306 | 3306 | 数据库 |
| Redis | 6379 | 6379 | 缓存 |
| RabbitMQ | 5672 | 5672 | AMQP 协议 |
| RabbitMQ Management | 15672 | 15672 | 管理界面 rabbitmq/rabbitmq |
| etcd | 2379 | 2379 | 客户端端口 |
| Jaeger UI | 16686 | 16686 | 链路追踪界面 |
| otel-collector | 4317 | 4317 | OTLP gRPC（host） |
| otel-collector | 4318 | 4318 | OTLP HTTP（host，go-zero otlphttp 走这里） |
| otel-collector | 8889 | 8889 | 聚合后的 metrics（Prometheus 抓） |
| otel-collector | 13133 | 13133 | health_check |
| Prometheus | 9090 | 9090 | 指标查询界面 |
| Loki | 3100 | 3100 | 日志接收 / 查询 API |

> Go 服务自身的 API/RPC 端口见各 `apps/**/etc/*.yaml`（API 8801–8813、RPC 8081–8093）；
> 各服务 Prometheus `/metrics` 端口 RPC 9101–9113、API 9201–9213（仅宿主机，不直接对 Prometheus 暴露）。

## 健康检查

```bash
# MySQL
mysql -h 127.0.0.1 -P 3306 -u root -p0000 -e "SELECT 1"

# Redis
redis-cli -h 127.0.0.1 -p 6379 ping

# RabbitMQ
curl -u rabbitmq:rabbitmq http://127.0.0.1:15672/api/overview

# etcd
etcdctl --endpoints=127.0.0.1:2379 endpoint health

# 可观测性
curl -s http://127.0.0.1:16686/   >/dev/null && echo "jaeger ok"     # Jaeger UI
curl -s http://127.0.0.1:13133/health_check >/dev/null && echo "otel-collector ok"
curl -s http://127.0.0.1:9090/-/healthy >/dev/null && echo "prometheus ok"
curl -s http://127.0.0.1:3100/ready   >/dev/null && echo "loki ok"   # Loki ready
```

## 可观测性数据流

```
Go 服务 (go-zero)
  trace  → OTLP/HTTP 127.0.0.1:4318 → otel-collector → jaeger:4318        → Jaeger UI :16686
  metric → /metrics (9101-9113/9201-9213) → otel-collector(prometheus recv)
                                                      → prometheus exporter :8889 → Prometheus :9090
  log   → logs/<svc>/*.log (仓库根) → docker 挂 /var/log/tjxt
                                                      → filelog → loki exporter → loki:3100
```

- 各服务 yaml 已注入 `Telemetry`（trace）、`Prometheus`（/metrics）、`Log{Mode:file,...}`（日志文件）。
- 日志路径为**相对路径** `logs/<svc>`，相对进程启动目录（CWD）；服务须从**仓库根**启动，
  否则日志会散落到各服务子目录，collector 的 filelog 采集不到。

## 数据持久化

- 依赖数据（MySQL/Redis/RabbitMQ）存储在 Docker named volumes 中。
- `make docker-down` 不会删除数据；完全重置：`docker compose down -v`。
- Loki 数据在容器内 `/tmp/loki`（单实例开发配置，重启即丢历史日志，符合本地开发预期）。

## 生产环境差异

| 项目 | 本地 | 生产 |
|------|------|------|
| MySQL | 单实例 | 主从/集群 |
| Redis | 单实例 | 哨兵/集群 |
| RabbitMQ | 单节点 | 镜像队列集群 |
| etcd | 单节点 | 3/5 节点集群 |
| 可观测性 | 单 collector + 单 Loki + 单 Prometheus | 各自集群化，collector 锁版本、Loki 接对象存储 |
| 认证 | 无/弱 | TLS + ACL |
| 备份 | 无 | 定期自动备份 |
