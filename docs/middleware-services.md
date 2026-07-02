# 中间件服务部署文档 (dlz-tjf-zts)

本文档描述一套独立于「建材通」主业务栈的中间件服务集群,通过 Docker Compose 统一编排,项目组名 `dlz-tjf-zts`,包含以下服务:

| 服务 | 用途 | 版本 |
|------|------|------|
| Elasticsearch | 搜索引擎 / 全文检索 | 9.4.2 |
| MongoDB | 文档数据库 | latest (8.x) |
| Kafka | 消息队列 (KRaft 模式) | apache/kafka latest |
| InfluxDB | 时序数据库 | 2.x latest |
| RocketMQ | 消息队列 (namesrv + broker + proxy) | 5.3.2 |
| IoTDB | 时序数据库 (物联网场景) | 1.3.2 standalone |

编排文件:[`docker-compose.dlz-tjf-zts.yml`](../docker-compose.dlz-tjf-zts.yml)
所有容器接入独立桥接网络 `dlz-tjf-zts-net`,数据卷均以 `dlz-tjf-zts-` 前缀命名。

---

## 一、快速开始

```bash
# 启动全部服务
docker compose -f docker-compose.dlz-tjf-zts.yml up -d

# 查看状态
docker compose -f docker-compose.dlz-tjf-zts.yml ps

# 停止全部
docker compose -f docker-compose.dlz-tjf-zts.yml stop

# 停止并删除容器(保留数据卷)
docker compose -f docker-compose.dlz-tjf-zts.yml down

# 彻底删除(含数据卷)
docker compose -f docker-compose.dlz-tjf-zts.yml down -v
```

> 可通过 `.env` 或环境变量覆盖默认账号密码:`MONGO_USER` / `MONGO_PASSWORD` / `INFLUX_USER` / `INFLUX_PASSWORD` / `INFLUX_TOKEN`。

---

## 二、服务清单与连接信息

### 1. Elasticsearch 9.4.2

| 项 | 值 |
|----|----|
| 容器名 | `dlz-tjf-zts-es` |
| 镜像 | `elasticsearch:9.4.2` |
| HTTP 端口 | `9200` |
| 集群通信端口 | `9300` |
| 安全认证 | 已关闭 (`xpack.security.enabled=false`),开发环境免账号 |
| 访问地址 | http://localhost:9200 |
| 模式 | 单节点 (`discovery.type=single-node`),堆内存 512MB |

验证:`curl http://localhost:9200`

### 2. MongoDB (8.x)

| 项 | 值 |
|----|----|
| 容器名 | `dlz-tjf-zts-mongo` |
| 镜像 | `mongo:latest` |
| 端口 | `27017` |
| 管理员用户 | `root` |
| 管理员密码 | `root123456` |
| 认证库 | `admin` |
| 业务库 | `dlz_tjf_zts` |

连接串:

```
mongodb://root:root123456@127.0.0.1:27017/?authSource=admin
```

> 客户端(Navicat / Compass)连接时,主机用 `127.0.0.1`,**验证数据库(Authentication Database)必须填 `admin`**。

### 3. Kafka (KRaft 模式)

| 项 | 值 |
|----|----|
| 容器名 | `dlz-tjf-zts-kafka` |
| 镜像 | `docker.1ms.run/apache/kafka:latest` |
| 宿主机接入端口 (EXTERNAL) | `9092` |
| 容器间接入端口 (PLAINTEXT) | `29092` |
| 认证 | 无 (PLAINTEXT) |

- 宿主机客户端连接:`localhost:9092`
- 容器间(同网络)连接:`dlz-tjf-zts-kafka:29092`

常用命令:

```bash
# 创建 topic
docker exec dlz-tjf-zts-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic <name> --partitions 1 --replication-factor 1

# 列出 topic
docker exec dlz-tjf-zts-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

### 4. InfluxDB 2.x

| 项 | 值 |
|----|----|
| 容器名 | `dlz-tjf-zts-influxdb` |
| 镜像 | `docker.1ms.run/library/influxdb:latest` |
| 端口 | `8086` |
| 用户名 | `admin` |
| 密码 | `admin123456` |
| 组织 (Org) | `dlz-tjf-zts` |
| 默认 Bucket | `default` |
| Admin Token | `dlz-tjf-zts-admin-token` |
| Web UI | http://localhost:8086 |

验证:`curl http://localhost:8086/health`

### 5. RocketMQ 5.3.2

由三个组件构成:

| 组件 | 容器名 | 端口(宿主机) | 说明 |
|------|--------|--------------|------|
| NameServer | `dlz-tjf-zts-rmqnamesrv` | `9876` | 注册中心 |
| Broker | `dlz-tjf-zts-rmqbroker` | `10911` / `10909` | 消息存储 |
| Proxy | `dlz-tjf-zts-rmqproxy` | `18080`(gRPC 8080) / `18081`(remoting 8081) | 客户端接入网关 |

- 镜像:`docker.1ms.run/apache/rocketmq:5.3.2`
- NameServer 地址(容器间):`dlz-tjf-zts-rmqnamesrv:9876`
- 5.x gRPC 客户端接入点(宿主机):`localhost:18081`
- 集群名:`DefaultCluster`
- 注意:各组件用 `JAVA_OPT_EXT` 将默认 8G 堆内存下调,broker 以 `user: root` 运行以获得数据卷写权限。

常用命令:

```bash
# 创建 topic
docker exec dlz-tjf-zts-rmqbroker sh mqadmin updateTopic -n dlz-tjf-zts-rmqnamesrv:9876 -t <TopicName> -c DefaultCluster

# 查看 topic 列表
docker exec dlz-tjf-zts-rmqbroker sh mqadmin topicList -n dlz-tjf-zts-rmqnamesrv:9876
```

### 6. Apache IoTDB 1.3.2 (单机 standalone)

| 项 | 值 |
|----|----|
| 容器名 | `dlz-tjf-zts-iotdb` |
| 镜像 | `docker.m.daocloud.io/apache/iotdb:1.3.2-standalone` |
| RPC / 客户端端口 | `6667` |
| 用户名 | `root` |
| 密码 | `root` |
| 数据卷 | `dlz-tjf-zts-iotdb-data`、`dlz-tjf-zts-iotdb-logs` |

- JDBC/Session 连接:`localhost:6667`,账号 `root/root`
- 命令行 CLI(容器内):

```bash
docker exec -it dlz-tjf-zts-iotdb /iotdb/sbin/start-cli.sh -h 127.0.0.1 -p 6667 -u root -pw root
```

示例 SQL:

```sql
CREATE DATABASE root.demo;
INSERT INTO root.demo.d1(timestamp, temp) VALUES(now(), 25.6);
SELECT * FROM root.demo.d1;
SHOW DATABASES;
```

---

## 三、端口一览

| 端口 | 服务 |
|------|------|
| 9200 / 9300 | Elasticsearch |
| 27017 | MongoDB |
| 9092 | Kafka (宿主机接入) |
| 8086 | InfluxDB |
| 9876 | RocketMQ NameServer |
| 10911 / 10909 | RocketMQ Broker |
| 18080 / 18081 | RocketMQ Proxy |
| 6667 | Apache IoTDB |

---

## 四、镜像清单

| 镜像 | 大小(约) | 来源 |
|------|----------|------|
| `elasticsearch:9.4.2` | 2.47GB | Docker Hub (官方镜像加速) |
| `mongo:latest` | 1.3GB | Docker Hub |
| `docker.1ms.run/apache/kafka:latest` | 686MB | 国内源 1ms.run |
| `docker.1ms.run/library/influxdb:latest` | 404MB | 国内源 1ms.run |
| `docker.1ms.run/apache/rocketmq:5.3.2` | ~600MB | 国内源 1ms.run |
| `docker.m.daocloud.io/apache/iotdb:1.3.2-standalone` | ~800MB | 国内源 daocloud |

> 国内拉取建议在 Docker Desktop → Settings → Docker Engine 配置镜像加速:
> ```json
> { "registry-mirrors": ["https://docker.1ms.run", "https://docker.m.daocloud.io"] }
> ```

### 导出/迁移镜像(离线打包)

若需将镜像打包迁移到其他机器:

```bash
# 导出
docker save elasticsearch:9.4.2 mongo:latest docker.1ms.run/apache/kafka:latest docker.1ms.run/library/influxdb:latest docker.1ms.run/apache/rocketmq:5.3.2 docker.m.daocloud.io/apache/iotdb:1.3.2-standalone -o dlz-tjf-zts-images.tar

# 在目标机器导入
docker load -i dlz-tjf-zts-images.tar
```

---

## 五、安全提示

- 以上账号密码 / Token 均为**开发环境默认值**,生产环境务必修改。
- Elasticsearch 当前关闭了安全认证,生产环境应开启 `xpack.security` 并配置 TLS。
- Kafka / RocketMQ 当前为明文无认证,生产环境应配置 SASL/ACL。
- 各服务端口默认绑定 `0.0.0.0`,若无需对外暴露,建议限制为 `127.0.0.1`。
