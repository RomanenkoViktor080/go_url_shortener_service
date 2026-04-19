package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/store"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/apperr"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/cache"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UrlRepository interface {
	CreateShortUrl(ctx context.Context, dto domain.CreateShortUrlDto) (store.Url, error)
	GetUniqueNumbers(ctx context.Context, generateSeries int64) ([]int64, error)
	SaveAllHashes(ctx context.Context, hash []string) (int64, error)
	FindUrl(ctx context.Context, hash string) (string, error)
	DeleteShortUrlBeforeCreatedAt(ctx context.Context, arg store.DeleteShortUrlBeforeCreatedAtParams) error
	CountCaches(ctx context.Context) (int64, error)
}
type urlRep struct {
	store     store.Store
	hashCache cache.HashCache
	urlCache  cache.UrlCache
}

func NewUrlRepository(
	store store.Store,
	hashCache cache.HashCache,
	urlCache cache.UrlCache,
) UrlRepository {
	return &urlRep{
		store:     store,
		hashCache: hashCache,
		urlCache:  urlCache,
	}
}

func (r *urlRep) CreateShortUrl(ctx context.Context, dto domain.CreateShortUrlDto) (store.Url, error) {
	hash, err := r.hashCache.GetHash(ctx)
	if err != nil {
		return store.Url{}, err
	}
	url, err := r.store.CreateShortUrl(ctx, store.CreateShortUrlParams{
		Hash: hash,
		Url:  dto.Url,
	})
	if err != nil {
		return store.Url{}, err
	}
	err = r.urlCache.Set(ctx, url.Hash, url.Url)
	if err != nil {
		slog.Warn(
			"could not cache url",
			"key", url.Hash,
			"value", url.Url,
		)
	}
	return url, err
}

func (r *urlRep) GetUniqueNumbers(ctx context.Context, generateSeries int64) ([]int64, error) {
	return r.store.GetUniqueNumbers(ctx, generateSeries)
}
func (r *urlRep) SaveAllHashes(ctx context.Context, hashes []string) (int64, error) {
	return r.store.SaveAllHashes(ctx, hashes)
}

func (r *urlRep) FindUrl(ctx context.Context, hash string) (string, error) {
	url, err := r.urlCache.Get(ctx, hash)
	if err == nil {
		return url, nil
	}
	if errors.Is(err, redis.Nil) {
		url, err = r.store.FindUrlByHash(ctx, hash)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", &apperr.NotFoundError{Message: "url not found"}
			}
			return "", err
		}
		err = r.urlCache.Set(ctx, hash, url)
		if err != nil {
			slog.Warn("could not cache url",
				"key", hash, "value", url,
			)
		}
	}
	return url, err
}

func (r *urlRep) DeleteShortUrlBeforeCreatedAt(
	ctx context.Context,
	arg store.DeleteShortUrlBeforeCreatedAtParams,
) error {
	return r.store.ExecTx(ctx, func(q store.Querier) error {
		hashes, err := r.store.DeleteShortUrlBeforeCreatedAt(ctx, arg)
		err = r.urlCache.DeleteAll(ctx, hashes)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		_, err = r.store.SaveAllHashes(ctx, hashes)
		return err
	})
}
func (r *urlRep) CountCaches(ctx context.Context) (int64, error) {
	return r.store.CountCaches(ctx)
}
