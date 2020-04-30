package web

import (
	"blog/pkg/response"
	"blog/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"math"
	"strconv"
	"time"
)

type ArticleForm struct {
	Path        string `json:"path" binding:"required,max=50" comment:"路径"`
	Title       string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
	Content     string `json:"content" binding:"required,max=65535"`
	PublishedAt string `json:"published_at" binding:"required"`
	TagId       []uint `json:"tag_id" binding:"required,dive,numeric,min=1"`
}

func GetArticles(c *gin.Context) {
	fmt.Println("aaa")
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
		response.BadRequest("参数不合法", nil)
		return
	}

	articleService := services.Article{
		Id: uint(id),
	}

	article, err := articleService.GetById()
	if err != nil {
		response.NotFound(err.Error(), "")
		return
	}
	response.Success("success", article)
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
		response.BadRequestWithValidationError(err, nil)
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
		response.BadRequest("标题或路径已存在", nil)
		return
	}

	err = articleService.Add()
	if err != nil {
		response.BadRequest("添加失败", nil)
		return
	}

	response.Created("success", "")
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
		response.BadRequest("参数不合法", "")
		return
	}
	var form ArticleForm

	err = c.ShouldBind(&form)

	if err != nil {
		response.BadRequestWithValidationError(err, err.Error())
		return
	}

	var publishedAt time.Time
	if publishedAt, err = time.ParseInLocation("2006-01-02", form.PublishedAt, time.Local); err != nil {
		response.BadRequest(err.Error(), nil)
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
		response.BadRequest("标题或路径已存在", nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		response.BadRequest("更新失败", "")
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
		response.BadRequest("参数不合法", nil)
		return
	}

	articleService := services.Article{
		Id: uint(id),
	}
	_, err = articleService.GetById()
	if err != nil {
		response.NotFound("1", "")
		return
	}

	err = articleService.Delete()
	if err != nil {
		response.ServerError(err.Error(), "")
		return
	}
	response.NoContent()
}