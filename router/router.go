package router

import (
	"blog/controllers/app"
	"blog/controllers/web"
	"blog/pkg/middlewares"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CustomValidator())

	router.Use(response.BeforeResponse())

	appRouter := router.Group("/app")
	{
		appRouter.GET("/articles/", app.GetArticles)
		appRouter.GET("/articles/:path", app.GetArticle)
		appRouter.GET("/tags", app.GetTags)
		appRouter.GET("/tags/:tag", app.GetArticles)
	}

	webRouter := router.Group("/web")
	{
		webRouter.POST("/login", web.Login)
		webRouter.DELETE("/logout", web.Logout)

		webRouter.Use(middlewares.JsonWebToken())
		{
			webRouter.GET("/articles", web.GetArticles)
			webRouter.GET("/articles/:id", web.GetArticle)
			webRouter.POST("/articles", web.AddArticle)
			webRouter.PUT("/articles/:id", web.EditArticle)
			webRouter.DELETE("/articles/:id", web.DeleteArticle)

			webRouter.GET("/tags", web.GetTags)
			webRouter.GET("/tags/:id", web.GetTag)
			webRouter.POST("/tags", web.AddTag)
			webRouter.PUT("/tags/:id", web.EditTag)
			webRouter.DELETE("/tags/:id", web.DeleteTag)
		}
	}

	return router
}