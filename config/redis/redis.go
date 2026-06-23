package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"domesv2/config"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(cfg *config.Config) {
	if !cfg.Redis.Enabled {
		log.Println("Redis is disabled, skipping initialization")
		return
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		RedisClient = nil
		return
	}

	log.Println("Redis connection established successfully")
}

func GetRedis() *redis.Client {
	return RedisClient
}

func IsRedisEnabled() bool {
	return RedisClient != nil
}

func GetCtx() context.Context {
	return ctx
}
