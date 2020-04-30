package web

import (
	"blog/models"
	"blog/pkg/response"
	"blog/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

func GetTags(c *gin.Context) {
	pageParam := c.Query("page")
	pageSizeParam := c.Query("size")
	var err error

	page := 1
	pageSize := 15

	if page, err = strconv.Atoi(pageParam); err != nil {
		page = 1
	}

	if pageSize, err = strconv.Atoi(pageSizeParam); err != nil {
		pageSize = 15
	}

	tagService := services.Tag{
		Page:     page,
		PageSize: pageSize,
	}

	var total int
	total, err = tagService.CountAll()
	if err != nil {
		response.ServerError(err.Error(), "")
		return
	}

	data := map[string]interface{}{
		"total": 0,
		"list":  []interface{}{},
	}

	if total <= 0 {
		response.Success("success", data)
		return
	}
	data["total"] = total

	var tags []*models.Tag
	tags, err = tagService.GetAll()
	if err != nil {
		data["total"] = 0
		response.ServerError(err.Error(), data)
		return
	}
	data["list"] = tags
	response.Success("", data)
}

func GetTag(c *gin.Context) {
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

	tagService := services.Tag{
		Id: uint(id),
	}
	var tag *models.Tag
	tag, err = tagService.Get()
	if err != nil {
		response.NotFound("", "")
		return
	}
	response.Success("success", tag)
}

func AddTag(c *gin.Context) {
	title := c.PostForm("title")
	var err error

	validate := validator.New()
	if err = validate.Var(title, "required,max=20"); err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	tagService := services.Tag{
		Title: title,
	}

	if err = tagService.Exist(); err != nil {
		response.BadRequest("标签已存在", "")
		return
	}

	err = tagService.Add()
	if err != nil {
		response.ServerError(err.Error(), "")
		return
	}
	response.Created("success", "")
}

func EditTag(c *gin.Context) {
	paramId := c.Param("id")
	title := c.PostForm("title")
	var err error

	validate := validator.New()
	if err = validate.Var(paramId, "required,numeric,min=1"); err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}
	if err = validate.Var(title, "required,max=20"); err != nil {
		response.BadRequestWithValidationError(err, nil)
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		response.BadRequest("参数不合法", "")
		return
	}

	tagService := services.Tag{
		Id:    uint(id),
		Title: title,
	}

	if err = tagService.Exist(); err != nil {
		response.BadRequest("标签已存在", "")
		return
	}

	err = tagService.Edit()
	if err != nil {
		response.ServerError("", "")
		return
	}
	response.NoContent()
}

func DeleteTag(c *gin.Context) {
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

	tagService := services.Tag{
		Id: uint(id),
	}
	_, err = tagService.Get()
	if err != nil {
		response.NotFound("标签不存在", "")
		return
	}

	articleService := services.Article{
		TagId: []uint{uint(id)},
	}
	total, err := articleService.CountAll()
	if err != nil || total > 0 {
		response.BadRequest("标签下有文章，不能删除", "")
		return
	}

	err = tagService.Delete()
	if err != nil {
		response.ServerError(err.Error(), "")
		return
	}
	response.NoContent()
}
