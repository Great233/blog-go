package app

import (
	"blog/pkg/response"
	"blog/services"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Article struct{}

func (a *Article) InitRouter(router *gin.Engine) {
	appRouter := router.Group("/app")
	{
		appRouter.GET("/articles", GetArticles)
		appRouter.GET("/articles/:path", GetArticle)
	}
}

func GetArticles(c *gin.Context) {
	pageParam := c.Query("page")
	pageSizeParam := c.Query("size")

	page := 1
	pageSize := 15
	tag := c.Param("tag")
	var err error

	if page, err = strconv.Atoi(pageParam); err != nil {
		page = 1
	}

	if pageSize, err = strconv.Atoi(pageSizeParam); err != nil {
		pageSize = 15
	}

	var tagId []uint
	if tag != "" {
		validate := validator.New()
		err = validate.Var(tag, "max=20")
		if err != nil {
			response.NotFound(response.TagIsNotExist, nil)
			return
		}

		tagService := services.Tag{
			Title: tag,
		}
		tagRecord, err := tagService.Get()
		if err != nil {
			response.NotFound(response.TagIsNotExist, nil)
			return
		}
		tagId = append(tagId, tagRecord.Id)
	}

	articleService := services.Article{
		TagId:    tagId,
		Page:     int(math.Max(1, float64(page))),
		PageSize: int(math.Max(1, float64(pageSize))),
	}

	data := make(map[string]interface{})
	data["total"] = 0
	data["list"] = []interface{}{}

	data["total"], err = articleService.CountAll()
	if err != nil {
		response.Success(response.Ok, data)
		return
	}

	data["list"], err = articleService.GetAll()

	if err != nil {
		data["total"] = 0
		response.ServerError(response.Ok, data)
		return
	}

	response.Success(response.Ok, data)
}

func GetArticle(c *gin.Context) {
	path := c.Param("path")
	var err error

	validate := validator.New()
	err = validate.Var(path, "required,max=50")
	if err != nil {
		response.NotFound(response.ArticleIsNotExist, nil)
		return
	}
	articleService := services.Article{
		Path: path,
	}

	article, err := articleService.GetByPath()
	if err != nil {
		response.NotFound(response.ArticleIsNotExist, nil)
		return
	}
	response.Success(response.Ok, article)
}
