package config

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/RomanenkoViktor080/url_shortener_service/docs"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/handler"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouterMount(
	urlService service.UrlService,
) http.Handler {
	initSwagger()

	urlHandler := handler.NewUrlHandler(urlService)
	hashHandler := handler.NewHashHandler(urlService)

	router := gin.Default()

	router.POST("/api/v1/url", urlHandler.CreateShortUrl)
	router.GET("/:hash", hashHandler.RedirectToOriginalUrl)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func (config *config) Run(handler http.Handler) error {
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server started", "port", config.Port)

	return srv.ListenAndServe()
}

func initSwagger() {
	docs.SwaggerInfo.Title = "URL Shortener Service API"
	docs.SwaggerInfo.Description = "A high-performance service for URL shortening"
	docs.SwaggerInfo.Version = "1.0"
	slog.Info(env.GetString("DOMAIN", "localhoasdasdasdasdasdast:8080"))
	docs.SwaggerInfo.Host = env.GetString("DOMAIN", "localhost:8080")

	docs.SwaggerInfo.BasePath = ""
}
