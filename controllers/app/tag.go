package app

import (
	"blog/pkg/response"
	"blog/services"
	"github.com/gin-gonic/gin"
	"log"
)

func GetTags(c *gin.Context) {
	tagService := services.Tag{}
	tags, err := tagService.GetAll()

	if err != nil {
		log.Fatalf("app.GetTags error: %v", err)
	}

	response.Success("success", tags)
}
