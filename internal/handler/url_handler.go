package handler

import (
	"net/http"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/domain"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/json"
	"github.com/gin-gonic/gin"
)

type urlHandler struct {
	service service.UrlService
}

func NewUrlHandler(service service.UrlService) *urlHandler {
	return &urlHandler{
		service: service,
	}
}

// CreateShortUrl godoc
//
//	@Summary		Create a shortened URL
//	@Description	Accepts a long URL and returns a shortened version
//	@Tags			Url
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		domain.CreateShortUrlDto	true	"Long URL to be shortened"
//	@Success		201		{object}	domain.ShortUrlDto			"Successfully created"
//	@Failure		422		{object}	json.ValidationError		"Validation error"
//	@Failure		400		{object}	json.Error					"Invalid request body"
//	@Failure		500		{object}	json.Error					"Internal server error"
//	@Router			/api/v1/url [post]
func (h *urlHandler) CreateShortUrl(c *gin.Context) {
	var dto domain.CreateShortUrlDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		json.InvalidRequestDataResponse(c, err)
		return
	}
	url, err := h.service.CreateShortUrl(c, dto)
	if err != nil {
		json.ErrorResponse(c, err)
		return
	}

	json.Response(c, http.StatusCreated, domain.ShortUrlDto{
		Url: url,
	})
}
