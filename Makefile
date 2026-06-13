# 建材通 Docker 快捷命令

.PHONY: dev up down build logs ps clean

# 开发模式（前后分离热更新）
dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

dev-d:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d

# 生产模式
up:
	docker compose up --build -d

down:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

build:
	docker compose build

logs:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

ps:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml ps

# 清除所有容器和数据卷（慎用）
clean:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down -v
