package main

import (
	"fmt"
	"log"
	"os"

	"material-build/internal/api"
	"material-build/internal/repository"
	"material-build/pkg/config"
	"material-build/pkg/database"
	redispkg "material-build/pkg/redis"

	"github.com/redis/go-redis/v9"
)

func main() {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	var rdb *redis.Client
	if cfg.Redis.Host != "" {
		client, err := redispkg.Connect(cfg.Redis)
		if err != nil {
			log.Fatalf("Redis 连接失败: %v", err)
		}
		rdb = client
		log.Println("Redis 连接成功")
	}

	repo := repository.New(db)
	router := api.SetupRouter(cfg, repo, rdb)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("建材通 API 启动于 http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
