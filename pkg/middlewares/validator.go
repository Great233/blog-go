package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func CustomValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		validate := binding.Validator.Engine().(*validator.Validate)
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			comment := field.Tag.Get("comment")
			if comment == "" {
				comment = field.Name
			}
			return comment
		})
		c.Next()
	}
}
