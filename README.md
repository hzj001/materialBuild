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

## 快速开始：从 GitHub 拉取并部署

> 仓库地址：[https://github.com/hzj001/materialBuild](https://github.com/hzj001/materialBuild)

### 前置要求

- [Git](https://git-scm.com/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) 或 Docker Engine + Compose V2
- 可用端口：`4000` `4001` `4002` `4008` `4306` `4379`

### 方式一：拉取预构建镜像（推荐，最快）

项目已将应用镜像发布到 GitHub Container Registry（GHCR），**无需本地编译**，直接拉取即可运行。

**1. 克隆代码**

```bash
git clone https://github.com/hzj001/materialBuild.git
cd materialBuild
```

**2. 配置环境变量**

```bash
cp .env.example .env
```

`.env` 中可按需修改端口与数据库密码，默认即可本地试用。

**3. 拉取并启动全部服务**

```bash
docker compose -f docker-compose.ghcr.yml pull
docker compose -f docker-compose.ghcr.yml up -d
```

Windows PowerShell 同样适用上述命令。

**4. 验证部署**

| 检查项 | 地址 / 命令 |
|---|---|
| 客户端 | http://localhost:4000 |
| 商家端 | http://localhost:4001 |
| 管理后台 | http://localhost:4002 |
| 后端健康检查 | http://localhost:4008/health |
| 容器状态 | `docker compose -f docker-compose.ghcr.yml ps` |
| 查看日志 | `docker compose -f docker-compose.ghcr.yml logs -f backend` |

默认管理员账号：`13800000000` / `admin123`

**5. 停止 / 重置**

```bash
# 停止服务（保留数据）
docker compose -f docker-compose.ghcr.yml down

# 停止并清空数据库等数据卷（重置为初始状态）
docker compose -f docker-compose.ghcr.yml down -v
```

#### GHCR 预构建镜像列表

| 镜像 | 说明 |
|---|---|
| `ghcr.io/hzj001/materialbuild-backend:latest` | Go API 后端 |
| `ghcr.io/hzj001/materialbuild-client:latest` | 客户端 H5 |
| `ghcr.io/hzj001/materialbuild-merchant:latest` | 商家端 H5 |
| `ghcr.io/hzj001/materialbuild-admin:latest` | vue-pure-admin 管理后台 |

MySQL、Redis 仍使用 Docker Hub 官方镜像（`mysql:8.0`、`redis:7-alpine`）。

> **镜像可见性**：若 `docker pull` 提示未授权，请先在 GitHub 仓库 [Packages](https://github.com/hzj001/materialBuild/pkgs/container/materialbuild-backend) 页面将对应 Package 设为 **Public**，或使用 `gh auth token | docker login ghcr.io -u hzj001 --password-stdin` 登录后再拉取。

> **镜像拉取慢或超时**：在 Docker Desktop → Settings → Docker Engine 添加国内镜像加速：
> ```json
> {
>   "registry-mirrors": [
>     "https://docker.1ms.run",
>     "https://docker.m.daocloud.io"
>   ]
> }
> ```

---

### 方式二：本地构建镜像（适合二次开发）

```bash
git clone https://github.com/hzj001/materialBuild.git
cd materialBuild
cp .env.example .env

# 生产模式：本地 build 全部镜像
docker compose up --build -d

# 开发模式：热更新（后端 Air + 前端挂载）
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Windows 开发模式快捷脚本：

```powershell
.\scripts\docker-dev.ps1 -Detached
```

---

### 方式三：维护者发布（推送代码 + 镜像）

若你是项目维护者，可将代码与容器镜像一并推送到 GitHub：

```powershell
# 1. 登录 GitHub
gh auth login

# 2. 一键推送代码 + 4 个容器镜像到 GHCR
cd materialBuild
.\scripts\publish-github.ps1
```

也可分步执行：

```powershell
git push -u origin main
gh auth token | docker login ghcr.io -u hzj001 --password-stdin
.\scripts\push-images.ps1
```

推送代码后，[GitHub Actions](.github/workflows/docker-publish.yml) 也会自动构建并发布镜像到 GHCR。

---

## Docker 访问地址

| 服务 | 地址 |
|---|---|
| 客户端 | http://localhost:4000 |
| 商家端 | http://localhost:4001 |
| 管理后台 (vue-pure-admin) | http://localhost:4002 |
| 后端 API | http://localhost:4008 |
| MySQL | localhost:4306 |
| Redis | localhost:4379 |

更多 Docker 细节见 [docs/docker.md](docs/docker.md)

## 项目结构

```
materialBuild/
├── docker-compose.yml          # 生产编排（本地 build）
├── docker-compose.dev.yml      # 开发编排（覆盖）
├── docker-compose.ghcr.yml     # 使用 GHCR 预构建镜像（推荐部署）
├── .github/workflows/          # GitHub Actions（自动构建推送镜像）
├── .env.example                # 环境变量模板
├── scripts/
│   ├── docker-dev.ps1          # Windows 开发启动
│   ├── publish-github.ps1      # 一键推送代码 + 镜像
│   └── push-images.ps1         # 推送容器镜像到 GHCR
├── backend/                    # Go API
│   ├── Dockerfile              # 生产镜像
│   ├── Dockerfile.dev          # 开发镜像（Air 热重载）
│   └── config.docker.yaml      # 容器内配置
├── frontend/
│   ├── shared/                 # 三端共享框架
│   ├── client/                 # 客户端 H5
│   ├── merchant/               # 商家端 H5
│   ├── admin/                  # vue-pure-admin 管理后台
│   └── docker/                 # Nginx 配置与 Dockerfile
├── database/migrations/        # MySQL 初始化脚本
└── docs/                       # 架构与 Docker 文档
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
