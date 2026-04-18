package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/redis/go-redis/v9"
)

type UrlCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	DeleteAll(ctx context.Context, keys []string) error
}

type urlCache struct {
	redis  redis.Client
	prefix string
	ttl    time.Duration
}

func NewUrlCache(
	redis redis.Client,
) UrlCache {
	return &urlCache{
		redis:  redis,
		prefix: env.GetString("CACHE_URL_PREIFIX", "urls:"),
		ttl:    time.Duration(env.GetInt("CACHE_URL_TTL", 86400)) * time.Second,
	}
}

func (urlCache *urlCache) Get(ctx context.Context, hash string) (string, error) {
	value, err := urlCache.redis.Get(ctx, urlCache.getKey(hash)).Result()
	if errors.Is(err, redis.Nil) {
		return value, nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}
func (urlCache *urlCache) Set(ctx context.Context, hash, url string) error {
	err := urlCache.redis.Set(ctx, urlCache.getKey(hash), url, urlCache.ttl).Err()
	if err != nil {
		slog.Error("failed to set data",
			"key", urlCache.getKey(hash),
			"value", url,
			"error", err,
		)
	}
	return err
}
func (urlCache *urlCache) DeleteAll(ctx context.Context, hashes []string) error {
	if len(hashes) == 0 {
		return nil
	}

	keys := make([]string, len(hashes))
	for i, hash := range hashes {
		keys[i] = urlCache.getKey(hash)
	}

	err := urlCache.redis.Del(ctx, keys...).Err()
	if err != nil {
		slog.Error("failed to delete hashes from redis", "error", err)
	}
	return err
}

func (urlCache *urlCache) getKey(key string) string {
	return urlCache.prefix + key
}
