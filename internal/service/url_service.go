package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repository "github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/sqlc"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/cache"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/domain"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
)

var domainName = env.GetString("DOMAIN", "http://localhost:8080")

type UrlService interface {
	CreateShortUrl(context context.Context, dto domain.CreateShortUrlDto) (string, error)
}
type svc struct {
	store repository.Store
	cache cache.HashCache
}

func NewUrlService(
	store repository.Store,
	cache cache.HashCache,
) UrlService {
	return &svc{
		store: store,
		cache: cache,
	}
}
func (svc *svc) CreateShortUrl(ctx context.Context, dto domain.CreateShortUrlDto) (string, error) {
	hash, err := svc.cache.GetHash()
	if err != nil {
		return "", err
	}
	url, err := svc.store.CreateShortUrl(ctx, repository.CreateShortUrlParams{
		Hash: hash,
		Url:  dto.Url,
	})
	if err != nil {
		slog.Warn("error creating short url", "error", err)
		return "", errors.New("error creating short url")
	}

	return buildShortUrl(url.Hash), nil

}

func buildShortUrl(hash string) string {
	return fmt.Sprintf("%s/%s", domainName, hash)
}
