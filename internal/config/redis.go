package config

import (
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:       env.GetString("REDIS_ADDR", "localhost:6379"),
		Username:   env.GetString("REDIS_USERNAME", "user"),
		Password:   env.GetString("REDIS_PASSWORD", "password"),
		MaxRetries: env.GetInt("REDIS_MAX_RETRIES", 3),
		PoolSize:   30,
	})
	return rdb
}
