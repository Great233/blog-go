package services

import (
	"blog/models"
	"errors"
)

type Tag struct {
	Id    uint
	Title string

	Page     int
	PageSize int
}

func (t *Tag) GetAll() ([]*models.Tag, error) {
	offset := (t.Page - 1) * t.PageSize
	if t.Page <= 0 {
		offset = -1
	}
	return models.GetAllTags(offset, t.PageSize, nil)
}

func (t *Tag) CountAll() (int, error) {
	return models.CountAllTags(nil)
}

func (t *Tag) Get() (*models.Tag, error) {

	condition := make(map[string]interface{})

	if t.Id > 0 {
		condition["id"] = t.Id
	}

	if t.Title != "" {
		condition["title"] = t.Title
	}

	return models.GetTag(condition)
}

func (t *Tag) Exist() error {
	condition := make(map[string]interface{})
	if t.Id > 0 {
		condition["id"] = t.Id
	}
	condition["title"] = t.Title
	return models.TagExists(condition)
}

func (t *Tag) Add() error {
	tag := models.Tag{
		Title: t.Title,
	}
	return models.AddTag(tag)
}

func (t *Tag) Edit() error {
	if t.Id <= 0 {
		return errors.New("invalid Tag.Id")
	}
	return models.EditTag(models.Tag{
		Model: models.Model{
			Id: t.Id,
		},
		Title: t.Title,
	})
}

func (t *Tag) Delete() error {
	return models.DeleteTag(models.Tag{
		Model: models.Model{
			Id: t.Id,
		},
	})
}
