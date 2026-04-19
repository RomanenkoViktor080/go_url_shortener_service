package json

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/apperr"
	val "github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	internalError       = Error{Message: "internal error"}
	invalidRequestError = Error{Message: "invalid request body"}
)

type Error struct {
	Message string `json:"message"`
}
type ValidationError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func ErrorResponse(c *gin.Context, err error) {
	if e, ok := errors.AsType[*apperr.NotFoundError](err); ok {
		Response(c, http.StatusNotFound, e)
		return
	}
	slog.Error("internal error", "error", err, "trace", debug.Stack())
	Response(c, http.StatusInternalServerError, internalError)
}
func Response(c *gin.Context, code int, any any) {
	c.JSON(code, any)
}
func InvalidRequestDataResponse(c *gin.Context, err error) {
	ve, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		Response(c, http.StatusBadRequest, invalidRequestError)
		return
	}

	fields := val.FormatValidationErrors(c, ve)

	c.JSON(http.StatusUnprocessableEntity, ValidationError{
		Message: "validation error",
		Fields:  fields,
	})
}
