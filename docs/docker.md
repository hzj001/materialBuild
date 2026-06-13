# Docker 部署指南

## 架构

```
materialBuild (Docker Compose 项目组)
├── mb-mysql              MySQL 8.0        :4306
├── mb-redis              Redis 7          :4379
├── mb-backend            Go API           :4008
├── mb-frontend-client    Nginx + H5       :4000
├── mb-frontend-merchant  Nginx + H5       :4001
└── mb-frontend-admin     vue-pure-admin   :4002
```

所有容器通过 `materialBuild-net` 桥接网络互通。前端 Nginx 将 `/api/` 反代到 `backend:8080`，实现前后分离。

## 快速启动（开发模式）

```bash
# 1. 复制环境变量
cp .env.example .env

# 2. 启动全部服务（热更新）
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Windows PowerShell
.\scripts\docker-dev.ps1 -Detached
```

| 服务 | 地址 |
|---|---|
| 客户端 | http://localhost:4000 |
| 商家端 | http://localhost:4001 |
| 管理端 | http://localhost:4002 |
| 后端健康检查 | http://localhost:4008/health |
| MySQL | localhost:4306 (user: material) |
| Redis | localhost:4379 |

默认管理员：`13800000000` / `admin123`

## 生产模式

```bash
docker compose up --build -d
```

## 开发特性

| 组件 | 开发方式 |
|---|---|
| 后端 | Air 热重载，源码挂载 `./backend` |
| 前端 | Nginx 静态挂载，改 HTML/JS 刷新即生效 |
| MySQL | 首次启动自动执行 `001_init.sql` |
| Redis | AOF 持久化，数据卷 `materialBuild-redis-data` |

## 常用命令

```bash
# 查看容器状态
docker compose ps

# 查看日志
docker compose logs -f backend

# 停止
docker compose down

# 停止并清除数据卷（重置数据库）
docker compose down -v
```

## 环境变量

见 `.env.example`，主要配置：

- `MYSQL_ROOT_PASSWORD` / `MYSQL_PASSWORD`
- `BACKEND_PORT` / `CLIENT_PORT` / `MERCHANT_PORT` / `ADMIN_PORT`

## 镜像拉取失败？

若 `docker pull` 超时，可在 Docker Desktop → Settings → Docker Engine 配置国内镜像：

```json
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.m.daocloud.io"
  ]
}
```

## 仅本地 Go 开发（不用 Docker 跑后端）

若只在本机 `go run` 后端，需本地启动 MySQL/Redis，并修改 `backend/config.yaml`。
前端 `config.js` 中 `API_BASE` 改为 `http://localhost:4008/api/v1`。
