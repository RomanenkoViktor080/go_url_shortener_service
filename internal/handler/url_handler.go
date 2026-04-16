package handler

import (
	"net/http"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/domain"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/json"
	"github.com/gin-gonic/gin"
)

type handler struct {
	service service.UrlService
}

func NewHandler(service service.UrlService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateShortUrl(c *gin.Context) {
	var dto domain.CreateShortUrlDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		json.ValidationErrorJsonResponse(c, err)
		return
	}
	url, err := h.service.CreateShortUrl(c, dto)
	if err != nil {
		json.JsonErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	json.JsonResponse(c, http.StatusCreated, url)
}
