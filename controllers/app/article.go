package app

import (
	"blog/pkg/response"
	"blog/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"math"
	"strconv"
)

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

	validate := validator.New()
	err = validate.Var(tag, "required,max=20")
	if err != nil {
		response.NotFound("", nil)
		return
	}

	tagService := services.Tag{
		Title: tag,
	}
	tagRecord, err := tagService.Get()
	if err != nil {
		response.NotFound("", nil)
		return
	}

	articleService := services.Article{
		TagId:    []uint{tagRecord.Id},
		Page:     int(math.Max(1, float64(page))),
		PageSize: int(math.Max(1, float64(pageSize))),
	}

	data := make(map[string]interface{})
	data["total"] = 0
	data["list"] = []interface{}{}

	data["total"], err = articleService.CountAll()
	if err != nil {
		log.Fatalf("app.GetTags error: %v", err)
	}

	data["list"], err = articleService.GetAll()

	if err != nil {
		data["total"] = 0
		log.Fatalf("app.GetTags error: %v", err)
	}

	response.Success("success", data)
}

func GetArticle(c *gin.Context) {
	path := c.Param("path")
	var err error

	validate := validator.New()
	err = validate.Var(path, "required,max=50")
	if err != nil {
		response.NotFound("", nil)
		return
	}
	articleService := services.Article{
		Path: path,
	}

	article, err := articleService.GetByPath()
	if err != nil {
		response.ServerError("", "")
	}
	response.Success("success", article)
}
