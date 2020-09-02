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
			response.Unauthorized(response.InvalidParams, nil)
			c.Abort()
			return
		}

		userInfo, err := utils.ParseJsonWebToken(xsrfToken)
		if err != nil {
			var message string
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				message = response.TokenHasExpired
				break
			default:
				message = response.TokenAuthFailed
				break
			}
			response.Unauthorized(message, nil)
			c.Abort()
			return
		}

		_, err = models.GetUserByUsername(userInfo.Username)
		if err != nil {
			response.Unauthorized(response.UserIsNotExist, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
