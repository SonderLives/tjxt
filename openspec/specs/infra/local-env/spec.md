# local-env（本地开发环境）Specification

## Purpose

定义本地开发环境的基础设施组成（docker-compose）：MySQL/Redis/RabbitMQ/etcd 四项基础依赖 + Jaeger/otel-collector/Prometheus/Loki 可观测性栈，以及端口、健康检查与数据持久化约定。本地开发步骤（环境依赖、启动顺序、生成与校验命令）详见 `specs/infra/local-dev-workflow`。

## Requirements

### Requirement: 基础依赖容器

`docker-compose.yml` SHALL 提供以下基础依赖：

| 容器 | 镜像 | 宿主机端口 | 凭证/说明 |
|------|------|------------|-----------|
| mysql | mysql:8.0 | 3306 | root/0000；挂载 `sql/migration` 初始化 |
| redis | redis:7-alpine | 6379 | named volume 持久化 |
| rabbitmq | rabbitmq:3.13-management | 5672 / 15672 | rabbitmq/rabbitmq；15672 为管理界面 |
| etcd | bitnami/etcd:3.5 | 2379 | ALLOW_NONE_AUTHENTICATION=yes |

#### Scenario: 初始化建库

- **WHEN** mysql 容器首次启动（无 named volume 数据）
- **THEN** 执行 `./sql/migration` 下的初始化脚本完成建库建表

### Requirement: 可观测性容器

`docker-compose.yml` SHALL 提供以下可观测性栈（详见 [observability](../observability/spec.md)）：

| 容器 | 镜像 | 宿主机端口 | 说明 |
|------|------|------------|------|
| jaeger | jaegertracing/all-in-one:1.57 | 16686 | 链路 UI；4317/4318 不映射宿主机 |
| otel-collector | otel/opentelemetry-collector-contrib:latest | 4317 / 4318 / 8889 / 13133 | 三信号中枢；挂载 `./logs:/var/log/tjxt:ro`；`extra_hosts: host.docker.internal:host-gateway` |
| prometheus | prom/prometheus:v2.53.1 | 9090 | 只抓 `otel-collector:8889` |
| loki | grafana/loki:latest | 3100 | 接收 collector filelog 推送 |

- `otel-collector` `depends_on: [jaeger, loki]`；`prometheus` `depends_on: [jaeger]`（仅排序）
- 生产环境 collector 镜像建议 pin 固定 tag（当前为 latest）

#### Scenario: 日志目录挂载

- **WHEN** otel-collector 容器启动
- **THEN** 宿主机 `./logs` 被只读挂载为容器内 `/var/log/tjxt`，供 filelog receiver 采集
- **AND** collector 通过 `host.docker.internal`（host-gateway）访问宿主机上运行的 26 个 Go 服务

### Requirement: 启停命令

本地环境 SHALL 通过 Makefile 目标管理：

- `make docker-up` = `docker compose up -d`（含可观测性四件套）
- `make docker-logs` = `docker compose logs -f`
- `make docker-down` = `docker compose down`（不删除数据卷）

#### Scenario: 一键拉起本地环境

- **WHEN** 开发者执行 `make docker-up`
- **THEN** 8 个容器（4 基础依赖 + 4 可观测）全部启动
- **AND** Go 服务随后从仓库根以 `make run-<svc>` / `make run-all` 启动（不在容器内）

### Requirement: 健康检查

各项依赖 SHALL 可用以下方式验证健康：

```bash
mysql -h 127.0.0.1 -P 3306 -u root -p0000 -e "SELECT 1"
redis-cli -h 127.0.0.1 -p 6379 ping
curl -u rabbitmq:rabbitmq http://127.0.0.1:15672/api/overview
etcdctl --endpoints=127.0.0.1:2379 endpoint health
curl -s http://127.0.0.1:13133/health_check   # otel-collector
curl -s http://127.0.0.1:9090/-/healthy        # prometheus
curl -s http://127.0.0.1:3100/ready            # loki
```

#### Scenario: 依赖就绪自检

- **WHEN** 开发者执行上述健康检查命令
- **THEN** 八项依赖（MySQL/Redis/RabbitMQ/etcd/Jaeger/otel-collector/Prometheus/Loki）各自返回成功响应，即可启动 Go 服务

### Requirement: 数据持久化

- MySQL/Redis/RabbitMQ 数据 SHALL 存于 Docker named volumes（`mysql-data`/`redis-data`/`rabbitmq-data`），`make docker-down` 不删除
- Loki 数据在容器内 `/tmp/loki`（单实例开发配置，重启丢历史日志，符合本地开发预期）

#### Scenario: 完全重置

- **WHEN** 开发者需要清空全部本地数据
- **THEN** 执行 `docker compose down -v` 删除容器与数据卷

