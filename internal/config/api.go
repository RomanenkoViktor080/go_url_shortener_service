package config

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/handler"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/gin-gonic/gin"
)

func RouterMount(
	urlService service.UrlService,
) http.Handler {
	urlHandler := handler.NewUrlHandler(urlService)
	hashHandler := handler.NewHashHandler(urlService)

	router := gin.Default()

	router.POST("/api/v1/url", urlHandler.CreateShortUrl)
	router.GET("/:hash", hashHandler.RedirectToOriginalUrl)

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
