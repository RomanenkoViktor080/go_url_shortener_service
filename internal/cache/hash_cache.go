package cache

import (
	"context"
	"log/slog"
	"sync/atomic"

	repository "github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/sqlc"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
)

var minLimitPercent = env.GetInt("CACHE_MIN_LIMIT_PERCENT", 20)
var cacheCapacity = env.GetInt("CACHE_CAPACITY", 10000)
var isFilling atomic.Bool
var cache = make(chan string, cacheCapacity)

type HashCache interface {
	GetHash() (string, error)
}

type hashCache struct {
	store repository.Store
	gen   generator.Generator
}

func NewHashCache(
	store repository.Store,
	gen generator.Generator,
) HashCache {
	return &hashCache{
		store: store,
		gen:   gen,
	}
}

func (c *hashCache) GetHash() (string, error) {
	if c.isBelowLimit() && isFilling.CompareAndSwap(false, true) {
		go c.fillCacheBackground()
	}
	return <-cache, nil
}

func (c *hashCache) fillCacheBackground() {
	defer isFilling.Store(false)

	ctx := context.Background()
	fetchSize := cacheCapacity - len(cache)
	hashes := make([]string, 0, fetchSize)
	var err error

	err = c.store.ExecTx(ctx, func(q repository.Querier) error {
		hashes, err = q.GetHashBatch(ctx, int32(fetchSize))
		if err != nil {
			return err
		}

		if len(hashes) < fetchSize {
			_, err = c.gen.GenerateBatchAndSave(ctx, int64(cacheCapacity))
			if err != nil {
				return err
			}
		}
		if len(hashes) < 1 {
			hashes, err = c.gen.GenerateBatch(ctx, 1)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		slog.Error("Error while filling cache", "error", err)
		return
	}
	for _, h := range hashes {
		cache <- h
	}
}

func (c *hashCache) isBelowLimit() bool {
	currentSize := len(cache)
	return float64(currentSize)/float64(cacheCapacity)*100 <= float64(minLimitPercent)
}
