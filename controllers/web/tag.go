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

	data := map[string]interface{}{
		"total": 0,
		"list":  []*models.Tag{},
	}

	data["total"], err = tagService.CountAll()
	if err != nil {
		response.Success(response.Ok, data)
		return
	}

	data["list"], err = tagService.GetAll()
	if err != nil {
		data["total"] = 0
		response.Success(response.Ok, data)
		return
	}
	response.Success(response.Ok, data)
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
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	tagService := services.Tag{
		Id: uint(id),
	}
	var tag *models.Tag
	tag, err = tagService.Get()
	if err != nil {
		response.NotFound(response.TagIsNotExist, nil)
		return
	}
	response.Success(response.Ok, tag)
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
		response.BadRequest(response.TagIsAlreadyExist, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		response.ServerError(response.AddTagFailed, nil)
		return
	}
	response.Created(response.Ok, nil)
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
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	tagService := services.Tag{
		Id:    uint(id),
		Title: title,
	}

	if err = tagService.Exist(); err != nil {
		response.BadRequest(response.TagIsAlreadyExist, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		response.ServerError(response.EditTagFailed, nil)
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
		response.BadRequest(response.InvalidParams, nil)
		return
	}

	tagService := services.Tag{
		Id: uint(id),
	}
	_, err = tagService.Get()
	if err != nil {
		response.NotFound(response.TagIsNotExist, nil)
		return
	}

	articleService := services.Article{
		TagId: []uint{uint(id)},
	}
	total, err := articleService.CountAll()
	if err != nil || total > 0 {
		response.BadRequest(response.DeleteTagFailed, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		response.ServerError(response.DeleteTagFailed, nil)
		return
	}
	response.NoContent()
}
