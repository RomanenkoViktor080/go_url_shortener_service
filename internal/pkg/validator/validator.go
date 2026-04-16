package validator

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
)

var uni *ut.UniversalTranslator

func InitValidatorTranslator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Fatal("binding.Validator.Engine err")
	}

	enLocale := en.New()
	ruLocale := ru.New()
	uni = ut.New(enLocale, enLocale, ruLocale)

	transEN, _ := uni.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(v, transEN)
	transRU, ok := uni.GetTranslator("ru")
	ruTranslations.RegisterDefaultTranslations(v, transRU)
}

func FormatValidationErrors(c *gin.Context, err validator.ValidationErrors) map[string]string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		lang = "en"
	}

	messages := make(map[string]string, len(err))
	trans, _ := uni.GetTranslator(lang)
	for _, e := range err {
		messages[strings.ToLower(e.Field())] = e.Translate(trans)
	}
	return messages
}
