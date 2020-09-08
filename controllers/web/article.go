package web

import (
	"blog/models"
	"blog/pkg/middlewares"
	"blog/pkg/response"
	"blog/services"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ArticleForm struct {
	Path        string `json:"path" binding:"required,max=50" comment:"路径"`
	Title       string `json:"title" binding:"required,max=100" comment:"标题"`
	Description string `json:"description" binding:"required,max=255" comment:"描述"`
	Content     string `json:"content" binding:"required,max=65535" comment:"内容"`
	PublishedAt string `json:"published_at" binding:"required" comment:"发布时间"`
	TagId       []uint `json:"tag_id" binding:"required,dive,numeric,min=1" comment:"标签"`
}

type Article struct{}

func (a *Article) InitRouter(router *gin.Engine) {
	webRouter := router.Group("/web")
	{
		articleRouter := webRouter.Use(middlewares.JsonWebToken())
		{
			articleRouter.GET("/articles", GetArticles)
			articleRouter.GET("/articles/:id", GetArticle)
			articleRouter.POST("/articles", AddArticle)
			articleRouter.PUT("/articles/:id", EditArticle)
			articleRouter.DELETE("/articles/:id", DeleteArticle)
		}
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

	tagId, _ := strconv.Atoi(tag)
	articleService := services.Article{
		TagId:    []uint{uint(tagId)},
		Page:     int(math.Max(1, float64(page))),
		PageSize: int(math.Max(1, float64(pageSize))),
	}

	data := map[string]interface{}{
		"total": 0,
		"list":  []*models.Article{},
	}

	data["total"], err = articleService.CountAll()
	if err != nil {
		response.Success(response.Ok, data)
		return
	}

	data["list"], err = articleService.GetAll()
	if err != nil {
		data["total"] = 0
		response.Success(response.Ok, data)
		return
	}

	response.Success(response.Ok, data)
}

func GetArticle(c *gin.Context) {
	var err error
	paramId := c.Param("id")
	validate := validator.New()
	err = validate.Var(paramId, "required,numeric,min=1")
	if err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	articleService := services.Article{
		Id: uint(id),
	}

	article, err := articleService.GetById()
	if err != nil {
		response.NotFound(response.ArticleIsNotExist, nil)
		return
	}
	response.Success(response.Ok, article)
}

func AddArticle(c *gin.Context) {
	var form ArticleForm
	var err error

	err = c.ShouldBindJSON(&form)
	if err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	var publishedAt time.Time
	if publishedAt, err = time.ParseInLocation("2006-01-02", form.PublishedAt, time.Local); err != nil {
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	articleService := services.Article{
		Path:        form.Path,
		Title:       form.Title,
		Description: form.Description,
		Content:     form.Content,
		PublishedAt: publishedAt,
		TagId:       form.TagId,
	}

	if err = articleService.Exists(); err != nil {
		response.BadRequest(response.PathOrTitleIsAlreadyExist, nil)
		return
	}

	err = articleService.Add()
	if err != nil {
		response.ServerError(response.AddArticleFailed, nil)
		return
	}

	response.Created(response.Ok, nil)
}

func EditArticle(c *gin.Context) {
	paramId := c.Param("id")
	var err error

	validate := validator.New()
	if err = validate.Var(paramId, "required,numeric,min=1"); err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		response.BadRequest(response.InvalidParams, nil)
		return
	}
	var form ArticleForm

	err = c.ShouldBind(&form)

	if err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	var publishedAt time.Time
	if publishedAt, err = time.ParseInLocation("2006-01-02", form.PublishedAt, time.Local); err != nil {
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	articleService := services.Article{
		Id:          uint(id),
		Path:        form.Path,
		Title:       form.Title,
		Description: form.Description,
		Content:     form.Content,
		PublishedAt: publishedAt,
		TagId:       form.TagId,
	}

	if err = articleService.Exists(); err != nil {
		response.BadRequest(response.PathOrTitleIsAlreadyExist, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		response.BadRequest(response.EditArticleFailed, "")
		return
	}
	response.NoContent()
}

func DeleteArticle(c *gin.Context) {
	var err error
	paramId := c.Param("id")
	validate := validator.New()
	err = validate.Var(paramId, "required,numeric,min=1")
	if err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	articleService := services.Article{
		Id: uint(id),
	}
	_, err = articleService.GetById()
	if err != nil {
		response.NotFound(response.ArticleIsNotExist, "")
		return
	}

	err = articleService.Delete()
	if err != nil {
		response.ServerError(response.DeleteArticleFailed, "")
		return
	}
	response.NoContent()
}
