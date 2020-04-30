package web

import (
	"blog/pkg/response"
	"blog/pkg/utils"
	"blog/services"
	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Username string `json:"username" binding:"required" comment:"用户名"`
	Password string `json:"password" binding:"required" comment:"密码"`
}

func Login(c *gin.Context) {
	var form loginForm
	var err error

	err = c.ShouldBindJSON(&form)
	if err != nil {
		response.BadRequestWithValidationError(err, "")
		return
	}

	userService := services.User{
		Username: form.Username,
	}

	user, err := userService.GetByUsername()
	if err != nil {
		response.BadRequest("用户名或密码错误", nil)
		return
	}

	if !utils.PasswordVerify(user.Password, form.Password) {
		response.BadRequest("用户名或密码错误", nil)
		return
	}

	token, err := utils.GenerateJsonWebToken(&utils.User{
		Username: user.Username,
	})
	if err != nil {
		response.ServerError(err.Error(), nil)
		return
	}

	c.Header("xsrf-token", token)
	response.Created("success", nil)
}

func Logout(c *gin.Context) {

}
