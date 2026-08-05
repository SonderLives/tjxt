# 本地依赖与 Docker Compose

> 版本：v1.0 | 更新：2026-08-05

## docker-compose.yml 服务清单

```yaml
services:
  mysql:
    image: mysql:8.0
    ports: ["3306:3306"]
    environment:
      MYSQL_ROOT_PASSWORD: "0000"
      MYSQL_DATABASE: "tj_auth"  # 仅示例，实际每服务各自建库
    volumes:
      - mysql_data:/var/lib/mysql
      - ./sql/ddl:/docker-entrypoint-initdb.d  # 可选：自动初始化

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  rabbitmq:
    image: rabbitmq:3.12-management
    ports: ["5672:5672", "15672:15672"]
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

  etcd:
    image: bitnami/etcd:3.5
    ports: ["2379:2379", "2380:2380"]
    environment:
      ETCD_ADVERTISE_CLIENT_URLS: "http://etcd:2379"
      ALLOW_NONE_AUTHENTICATION: "yes"
    volumes:
      - etcd_data:/bitnami/etcd

volumes:
  mysql_data:
  redis_data:
  rabbitmq_data:
  etcd_data:
```

## 启动与停止

```bash
# 启动所有依赖
make docker-up

# 查看日志
make docker-logs

# 停止并删除
make docker-down
```

## 端口占用表

| 服务 | 容器端口 | 宿主机端口 | 用途 |
|------|----------|------------|------|
| MySQL | 3306 | 3306 | 数据库 |
| Redis | 6379 | 6379 | 缓存 |
| RabbitMQ | 5672 | 5672 | AMQP 协议 |
| RabbitMQ Management | 15672 | 15672 | 管理界面 guest/guest |
| etcd | 2379 | 2379 | 客户端端口 |
| etcd | 2380 | 2380 | 集群通信 |

## 健康检查

```bash
# MySQL
mysql -h 127.0.0.1 -P 3306 -u root -p0000 -e "SELECT 1"

# Redis
redis-cli -h 127.0.0.1 -p 6379 ping

# RabbitMQ
curl -u guest:guest http://127.0.0.1:15672/api/overview

# etcd
etcdctl --endpoints=127.0.0.1:2379 endpoint health
```

## 数据持久化

- 所有数据存储在 Docker named volumes 中
- `make docker-down` 不会删除数据
- 完全重置：`docker compose down -v`

## 生产环境差异

| 项目 | 本地 | 生产 |
|------|------|------|
| MySQL | 单实例 | 主从/集群 |
| Redis | 单实例 | 哨兵/集群 |
| RabbitMQ | 单节点 | 镜像队列集群 |
| etcd | 单节点 | 3/5 节点集群 |
| 认证 | 无/弱 | TLS + ACL |
| 备份 | 无 | 定期自动备份 |