package main

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestExport(t *testing.T) {
	redisPlugin, err := plugins.NewRedisPlugin(context.Background(), &plugins.RedisPluginConfig{
		Redis: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	})
}
