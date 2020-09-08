package router

import (
	"blog/controllers/app"
	"blog/controllers/web"
	"blog/pkg/response"
	"reflect"

	"github.com/gin-gonic/gin"
)

type InitRouterInterface interface {
	InitRouter(router *gin.Engine)
}

func Init() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(response.BeforeResponse())

	container := map[string]InitRouterInterface{
		"article":    &app.Article{},
		"tag":        &app.Tag{},
		"login":      &web.Login{},
		"webArticle": &web.Article{},
		"webTag":     &web.Tag{},
	}

	for _, key := range container {
		reflect.ValueOf(key).MethodByName("InitRouter").Call([]reflect.Value{reflect.ValueOf(router)})
	}

	return router
}
