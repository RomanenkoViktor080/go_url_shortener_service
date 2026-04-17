package cache

import (
	"context"
	"log"
	"log/slog"
	"sync/atomic"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/store"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
)

var (
	minLimitPercent = env.GetInt("CACHE_MIN_LIMIT_PERCENT", 20)
	cacheCapacity   = env.GetInt("CACHE_CAPACITY", 10000)
	isFilling       atomic.Bool
)

type HashCache interface {
	GetHash() (string, error)
}

type hashCache struct {
	store store.Store
	gen   generator.Generator
	cache chan string
}

func NewHashCache(
	store store.Store,
	gen generator.Generator,
) HashCache {
	hashCache := &hashCache{
		store: store,
		gen:   gen,
		cache: make(chan string, cacheCapacity),
	}
	hashCache.init()

	return hashCache
}

func (c *hashCache) init() {
	slog.Info("initializing hash cache")
	cxt := context.Background()
	batch, err := c.store.GetHashBatch(cxt, int32(cacheCapacity))
	if err != nil || batch == nil || len(batch) == 0 {
		slog.Info("generating hashes")
		generateBatch, err := c.gen.GenerateBatch(cxt, int64(cacheCapacity))
		if err != nil {
			log.Fatalf("failed to initialize hash cache, error: %v", err)
		}
		batch = generateBatch
	}
	for _, h := range batch {
		c.cache <- h
	}
	slog.Info("hash cache is initialized")
}

func (c *hashCache) GetHash() (string, error) {
	if c.isBelowLimit() && isFilling.CompareAndSwap(false, true) {
		go c.fillCacheBackground()
	}
	return <-c.cache, nil
}

func (c *hashCache) fillCacheBackground() {
	defer isFilling.Store(false)

	ctx := context.Background()
	fetchSize := cacheCapacity - len(c.cache)
	hashes := make([]string, 0, fetchSize)
	var err error

	err = c.store.ExecTx(ctx, func(q store.Querier) error {
		hashes, err = q.GetHashBatch(ctx, int32(fetchSize))
		if err != nil {
			slog.Error("failed to fetch hashes", "error", err)
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
		slog.Error("error while filling cache", "error", err)
		return
	}
	for _, h := range hashes {
		c.cache <- h
	}
}

func (c *hashCache) isBelowLimit() bool {
	currentSize := len(c.cache)
	return float64(currentSize)/float64(cacheCapacity)*100 <= float64(minLimitPercent)
}
