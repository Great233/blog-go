package web

import (
	"blog/pkg/response"
	"blog/pkg/utils"
	"blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Username string `json:"username" binding:"required" comment:"用户名"`
	Password string `json:"password" binding:"required" comment:"密码"`
}

type Login struct{}

func (a *Login) InitRouter(router *gin.Engine) {
	webRouter := router.Group("/web")
	{
		webRouter.POST("/login", LoginM)
		webRouter.DELETE("/logout", Logout)
	}
}

func LoginM(c *gin.Context) {
	var form loginForm
	var err error

	err = c.ShouldBindJSON(&form)
	if err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	userService := services.User{
		Username: form.Username,
	}

	user, err := userService.GetByUsername()
	if err != nil {
		response.BadRequest(response.UserIsNotExist, nil)
		return
	}

	if !utils.PasswordVerify(user.Password, form.Password) {
		response.BadRequest(response.PasswordVerifyFailed, nil)
		return
	}

	token, err := utils.GenerateJsonWebToken(&utils.User{
		Username: user.Username,
	})
	if err != nil {
		response.ServerError(response.TokenGenerateFailed, nil)
		return
	}

	response.Respond(http.StatusCreated, response.Ok, nil, map[string]string{
		"xsrf-token": token,
	})
}

func Logout(c *gin.Context) {

}
