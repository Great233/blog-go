package app

import (
	"blog/pkg/response"
	"blog/services"

	"github.com/gin-gonic/gin"
)

func InitTagRouter() {
}

type Tag struct{}

func (t *Tag) InitRouter(router *gin.Engine) {
	appRouter := router.Group("/app")
	{
		appRouter.GET("/tags", GetTags)
		appRouter.GET("/tags/:tag", GetArticles)
	}
}

func GetTags(c *gin.Context) {
	tagService := services.Tag{}
	tags, _ := tagService.GetAll()

	response.Success("success", tags)
}
