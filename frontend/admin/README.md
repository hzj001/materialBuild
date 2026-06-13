# 建材通管理后台

基于 [pure-admin-thin](https://github.com/pure-admin/pure-admin-thin)（vue-pure-admin 官方精简版）构建，对接建材通 Go 后端 API。

## 功能模块

- **数据概览** — 商家数、订单数、用户数
- **城市管理** — 开通/禁用服务城市
- **商家管理** — 审核入驻、启用/禁用
- **订单监控** — 全平台订单查看

## 本地开发

```bash
# 需要 Node >= 20, pnpm >= 9
pnpm install
pnpm dev
```

访问 http://localhost:4002 ，默认管理员 `13800000000` / `admin123`

需确保后端 API 运行在 `http://localhost:4008`（Vite 会代理 `/api`）。

## Docker 开发

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up frontend-admin backend mysql redis
```

## 生产构建

```bash
pnpm build
# 产物在 dist/，由 Nginx 托管并反代 /api
```

## 定制说明

| 文件 | 说明 |
|---|---|
| `src/api/user.ts` | 登录适配（手机号 + JWT） |
| `src/api/material.ts` | 业务 API |
| `src/router/modules/material.ts` | 业务菜单路由 |
| `src/views/material/` | 业务页面 |
| `.env.development` | 开发配置（已关闭 mock） |
