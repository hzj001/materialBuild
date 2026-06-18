# 建材通后端 - 生产镜像
FROM golang:1.24-alpine AS builder

WORKDIR /build
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/server ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app
COPY --from=builder /build/server .
COPY config.docker.yaml .

RUN mkdir -p /app/uploads

EXPOSE 8080
CMD ["./server"]
