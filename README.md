# 建材通 — 三四线城市建材服务平台

类似美团的建材电商 H5 应用，面向三四线城市本地建材供需场景。

## 系统角色

| 端 | 说明 |
|---|---|
| **客户端 H5** | 搜索商品、选品下单、支付、客服咨询、售后、查看商家图文/视频 |
| **商家端 H5** | 商品上架、定价打折、配送设置、收款、客服、售后处理 |
| **管理后台** | vue-pure-admin 框架，城市管理、商家审核、订单监控、数据统计 |

## 技术栈

- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **后端**: Go / Gin / GORM / JWT
- **前端**: 客户端/商家端自研 H5；管理端 [vue-pure-admin](https://github.com/pure-admin/vue-pure-admin)（Element Plus）
- **部署**: Docker Compose 全容器化，前后分离

## Docker 一键启动（推荐）

### 前置要求

- Docker Desktop 或 Docker Engine + Compose V2

### 开发模式（热更新）

```bash
# 复制环境变量
cp .env.example .env

# 启动 materialBuild 全部服务
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Windows PowerShell:

```powershell
.\scripts\docker-dev.ps1 -Detached
```

### 访问地址

| 服务 | 地址 |
|---|---|
| 客户端 | http://localhost:4000 |
| 商家端 | http://localhost:4001 |
| 管理后台 (vue-pure-admin) | http://localhost:4002 |
| 后端 API | http://localhost:4008 |
| MySQL | localhost:4306 |
| Redis | localhost:4379 |

默认管理员：`13800000000` / `admin123`

### 生产模式

```bash
docker compose up --build -d
```

详细说明见 [docs/docker.md](docs/docker.md)

## 项目结构

```
materialBuild/
├── docker-compose.yml        # 生产编排
├── docker-compose.dev.yml    # 开发编排（覆盖）
├── .env.example              # 环境变量模板
├── backend/                  # Go API
│   ├── Dockerfile            # 生产镜像
│   ├── Dockerfile.dev        # 开发镜像（Air 热重载）
│   └── config.docker.yaml    # 容器内配置
├── frontend/
│   ├── shared/               # 三端共享框架
│   ├── client/               # 客户端 H5
│   ├── merchant/             # 商家端 H5
│   ├── admin/                # 管理端 H5
│   └── docker/               # Nginx 配置与 Dockerfile
├── database/migrations/      # MySQL 初始化脚本
└── docs/                     # 架构与 Docker 文档
```

## Docker 服务组（materialBuild）

```
materialBuild-net
├── mb-mysql
├── mb-redis
├── mb-backend          ← Go API，连接 MySQL + Redis
├── mb-frontend-client  ← Nginx 反代 /api → backend
├── mb-frontend-merchant
└── mb-frontend-admin
```

前端通过 Nginx 反向代理访问后端，API 地址统一为 `/api/v1`，无需跨域。

## 本地开发（非 Docker）

```bash
# 数据库
mysql -u root -p < database/migrations/001_init.sql

# 后端（需本地 MySQL + Redis）
cd backend && go run ./cmd/server

# 前端需将 config.js 中 API_BASE 改为 http://localhost:4008/api/v1
```

## API 模块

| 前缀 | 模块 |
|---|---|
| `/api/v1/auth` | 登录注册 |
| `/api/v1/client` | 商品、订单、售后 |
| `/api/v1/merchant` | 商品管理、订单、配送 |
| `/api/v1/admin` | 城市、商家、统计 |
| `/api/v1/common` | 城市列表、分类 |

## 开发路线图

- [x] 项目框架与数据库设计
- [x] 后端 API 骨架
- [x] 三端 H5 页面骨架
- [x] Docker 全容器化部署
- [ ] 文件上传与媒体管理
- [ ] 支付对接
- [ ] WebSocket 在线客服
- [ ] 微信小程序迁移
