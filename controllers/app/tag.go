package app

import (
	"blog/pkg/response"
	"blog/services"
	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	tagService := services.Tag{}
	tags, _ := tagService.GetAll()

	response.Success("success", tags)
}
