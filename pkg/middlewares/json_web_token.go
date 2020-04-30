package middlewares

import (
	"blog/models"
	"blog/pkg/response"
	"blog/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JsonWebToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		xsrfToken := c.GetHeader("xsrf-token")
		if xsrfToken == "" {
			response.Unauthorized("token cannot be null", nil)
			c.Abort()
			return
		}

		userInfo, err := utils.ParseJsonWebToken(xsrfToken)
		if err != nil {
			var message string
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				message = "token has expired"
				break
			default:
				message = "invalid token"
				break
			}
			response.Unauthorized(message, nil)
			c.Abort()
			return
		}

		_, err = models.GetUserByUsername(userInfo.Username)
		if err != nil {
			response.Unauthorized("user is not exist", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
