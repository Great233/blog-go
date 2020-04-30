package response

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"log"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Context *gin.Context

var Trans ut.Translator

func Init()  {
	cn := zh.New()
	uni := ut.New(cn, cn)
	Trans, _ = uni.GetTranslator("zh")
	err := zhTrans.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), Trans)
	if err != nil {
		log.Fatalf("middlewares.CustomValidator error: %v", err)
	}
}

func BeforeResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		Context = c
	}
}

func Success(message string, data interface{}) {
	Context.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
	return
}

func Created(message string, data interface{}) {
	Context.JSON(201, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
	return
}

func NoContent() {
	Context.JSON(204, Response{
		Code:    0,
		Message: "",
		Data:    "",
	})
	return
}

func BadRequest(message string, data interface{}) {
	Context.JSON(400, Response{
		Code:    40000,
		Message: message,
		Data:    data,
	})
	return
}

func BadRequestWithValidationError(err error, data interface{}) {
	errs := err.(validator.ValidationErrors)
	if len(errs) <= 0 {
		Context.JSON(400, Response{
			Code:    40000,
			Message: "参数不合法",
			Data:    data,
		})
		return
	}
	for _, e := range errs {
		Context.JSON(400, Response{
			Code:    40000,
			Message: e.Translate(Trans),
			Data:    data,
		})
		return
	}
}

func NotFound(message string, data interface{}) {
	Context.JSON(404, Response{
		Code:    40004,
		Message: message,
		Data:    data,
	})
	return
}

func Unauthorized(message string, data interface{}) {
	Context.JSON(401, Response{
		Code:    40001,
		Message: message,
		Data:    data,
	})
}

func ServerError(message string, data interface{}) {
	Context.JSON(500, Response{
		Code:    50000,
		Message: message,
		Data:    data,
	})
}
