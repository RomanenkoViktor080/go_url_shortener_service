package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/domain"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/repository"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
)

var domainName = env.GetString("DOMAIN", "http://localhost:8080")

type UrlService interface {
	CreateShortUrl(context context.Context, dto domain.CreateShortUrlDto) (string, error)
	GetOriginalUrl(context context.Context, dto domain.URLRedirectDto) (string, error)
}
type urlService struct {
	rep repository.UrlRepository
}

func NewUrlService(
	rep repository.UrlRepository,
) UrlService {
	return &urlService{
		rep: rep,
	}
}
func (svc *urlService) CreateShortUrl(ctx context.Context, dto domain.CreateShortUrlDto) (string, error) {
	url, err := svc.rep.CreateShortUrl(ctx, dto)
	if err != nil {
		slog.Warn("error creating short url", "error", err)
		return "", errors.New("error creating short url")
	}

	return buildShortUrl(url.Hash), nil

}

func (svc *urlService) GetOriginalUrl(ctx context.Context, dto domain.URLRedirectDto) (string, error) {
	return svc.rep.FindUrl(ctx, dto.Hash)
}

func buildShortUrl(hash string) string {
	return fmt.Sprintf("%s/%s", domainName, hash)
}
