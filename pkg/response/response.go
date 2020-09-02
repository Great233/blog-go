package response

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"net/http"
	"reflect"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var context *gin.Context

var trans ut.Translator

func Init()  {
	validate := binding.Validator.Engine().(*validator.Validate)
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		comment := field.Tag.Get("comment")
		if comment == "" {
			comment = field.Name
		}
		return comment
	})

	cn := zh.New()
	uni := ut.New(cn, cn)
	trans, _ = uni.GetTranslator("zh")
	err := zhTrans.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), trans)
	if err != nil {
		log.Fatalf("middlewares.CustomValidator error: %v", err)
	}
}

func BeforeResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		context = c
	}
}

func Respond(statusCode int, message string, data interface{}, headers map[string]string) {
	for key, value := range headers {
		context.Header(key, value)
	}
	context.JSON(statusCode, Response{
		Message: message,
		Data:    data,
	})
}

func Success(message string, data interface{}) {
	Respond(http.StatusOK, message, data, nil)
	return
}

func Created(message string, data interface{}) {
	Respond(http.StatusCreated, message, data, nil)
	return
}

func NoContent() {
	Respond(http.StatusNoContent, "", nil, nil)
	return
}

func BadRequest(message string, data interface{}) {
	Respond(http.StatusBadRequest, message, data, nil)
	return
}

func BadRequestWithValidationError(err error, data interface{}) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok || len(errs) <= 0 {
		Respond(http.StatusBadRequest, InvalidParams, data, nil)
		return
	}
	for _, e := range errs {
		Respond(http.StatusBadRequest, e.Translate(trans), data, nil)
		return
	}
}

func NotFound(message string, data interface{}) {
	Respond(http.StatusNotFound, message, data, nil)
	return
}

func Unauthorized(message string, data interface{}) {
	Respond(http.StatusUnauthorized, message, data, nil)
}

func ServerError(message string, data interface{}) {
	Respond(http.StatusInternalServerError, message, data, nil)
}
