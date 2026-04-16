package json

import (
	"errors"
	"log/slog"
	"net/http"

	val "github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Error struct {
	Message string `json:"message"`
}
type ValidationError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func JsonErrorResponse(c *gin.Context, code int, msg string, err error) {
	if code >= 500 {
		slog.Error(msg, err)
	}
	JsonResponse(c, code, err)
}
func JsonResponse(c *gin.Context, code int, any any) {
	c.JSON(code, any)
}
func ValidationErrorJsonResponse(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		JsonResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	fields := val.FormatValidationErrors(c, ve)

	c.JSON(http.StatusUnprocessableEntity, ValidationError{
		Message: "validation error",
		Fields:  fields,
	})
}
