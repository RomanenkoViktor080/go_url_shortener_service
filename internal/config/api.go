package config

import (
	"net/http"

	repository "github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/sqlc"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/cache"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/handler"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/encoder"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/gin-gonic/gin"
)

func (app *application) mount() http.Handler {
	store := repository.NewStore(app.dbConnection)
	router := gin.Default()
	urlHandler := handler.NewHandler(
		service.NewUrlService(
			store,
			cache.NewHashCache(
				store,
				generator.NewGenerator(
					store,
					encoder.NewBase62Encoder(),
				),
			),
		),
	)
	router.POST("/api/v1/url", urlHandler.CreateShortUrl)

	return router
}
