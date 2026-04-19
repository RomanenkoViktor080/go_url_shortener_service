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

// RedirectToOriginalUrl godoc
//
//	@Summary		Redirect to original URL
//	@Description	Retrieves the original long URL associated with the provided hash and performs a 302 redirect
//	@Tags			Url
//	@Accept			json
//	@Param			hash	path		string					true	"Short URL hash identifier"
//	@Success		302		{object}	nil						"Redirecting to destination"
//	@Header			302		{string}	Location				"https://example.com"
//	@Failure		422		{object}	json.ValidationError	"Validation error"
//	@Failure		400		{object}	json.Error				"Invalid request body"
//	@Failure		500		{object}	json.Error				"Internal server error"
//	@Router			/{hash} [get]
func (h *hashHandler) RedirectToOriginalUrl(c *gin.Context) {
	var dto domain.URLRedirectDto
	if err := c.ShouldBindUri(&dto); err != nil {
		json.InvalidRequestDataResponse(c, err)
		return
	}
	url, err := h.service.GetOriginalUrl(c, dto)
	if err != nil {
		json.ErrorResponse(c, err)
		return
	}

	c.Redirect(http.StatusFound, url)
}
