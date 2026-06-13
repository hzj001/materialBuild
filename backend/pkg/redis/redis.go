package redis

import (
	"context"
	"fmt"
	"time"

	"material-build/pkg/config"

	"github.com/redis/go-redis/v9"
)

func Connect(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func Ping(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}
