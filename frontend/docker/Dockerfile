# 建材通前端 - 生产镜像
# build arg: APP_DIR = client | merchant | admin
FROM nginx:1.25-alpine

ARG APP_DIR=client

COPY shared/ /usr/share/nginx/html/shared/
COPY ${APP_DIR}/ /usr/share/nginx/html/
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
