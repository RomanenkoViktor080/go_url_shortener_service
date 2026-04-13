package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) mount() http.Handler {
	router := gin.Default()
	return router
}
