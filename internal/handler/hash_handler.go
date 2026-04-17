package handler

import (
	"net/http"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/domain"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/json"
	"github.com/gin-gonic/gin"
)

type hashHandler struct {
	service service.UrlService
}

func NewHashHandler(service service.UrlService) *hashHandler {
	return &hashHandler{
		service: service,
	}
}

func (h *hashHandler) RedirectToOriginalUrl(c *gin.Context) {
	var dto domain.HashDto
	err := c.ShouldBindUri(&dto)
	if err != nil {
		json.ValidationErrorJsonResponse(c, err)
		return
	}
	url, err := h.service.GetOriginalUrl(c, dto)
	if err != nil {
		json.JsonErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.Redirect(http.StatusFound, url)
}
