package config

import (
	"net/http"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/store"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/cache"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/handler"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/encoder"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/repository"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/gin-gonic/gin"
)

func (app *application) mount() http.Handler {
	store := store.NewStore(app.dbConnection)
	urlRepository := repository.NewUrlRepository(
		store,
		cache.NewHashCache(
			store,
			generator.NewGenerator(
				store,
				encoder.NewBase62Encoder(),
			),
		),
		cache.NewUrlCache(*app.redisClient),
	)
	urlService := service.NewUrlService(
		urlRepository,
	)
	urlHandler := handler.NewUrlHandler(urlService)
	hashHandler := handler.NewHashHandler(urlService)

	router := gin.Default()

	router.POST("/api/v1/url", urlHandler.CreateShortUrl)
	router.GET("/:hash", hashHandler.RedirectToOriginalUrl)

	return router
}
