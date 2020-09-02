package services

import (
	"blog/models"
	"time"
)

type Article struct {
	Id          uint
	Path        string
	Title       string
	Description string
	Content     string
	PublishedAt time.Time

	TagId []uint

	Page     int
	PageSize int
}

func (a *Article) GetAll() ([]*models.Article, error) {
	condition := make(map[string]interface{})

	if a.Title != "" {
		condition["title like %?%"] = a.Title
	}

	tagCondition := make(map[string]interface{})
	if len(a.TagId) > 0 && a.TagId[0] > 0 {
		tagCondition["tag_id=?"] = a.TagId[0]
	}

	return models.GetArticles(
		(a.Page-1)*a.PageSize,
		a.PageSize,
		condition,
		tagCondition,
	)
}

func (a *Article) CountAll() (int, error) {
	condition := make(map[string]interface{})
	if a.Title != "" {
		condition["title like %?%"] = a.Title
	}

	tagCondition := make(map[string]interface{})

	if len(a.TagId) > 0 && a.TagId[0] > 0 {
		tagCondition["tag_id=?"] = a.TagId[0]
	}
	return models.CountAll(condition, tagCondition)
}

func (a *Article) GetByPath() (*models.Article, error) {
	return models.GetArticleByPath(a.Path)
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"path":            a.Path,
		"title":           a.Title,
		"description":     a.Description,
		"content":         a.Content,
		"published_year":  a.PublishedAt.Year(),
		"published_month": a.PublishedAt.Month(),
		"published_date":  a.PublishedAt.Day(),
		"tag_id":          a.TagId,
	}

	return models.AddArticle(article)
}

func (a *Article) Edit() error {
	article := map[string]interface{}{
		"path":            a.Path,
		"title":           a.Title,
		"description":     a.Description,
		"content":         a.Content,
		"published_year":  a.PublishedAt.Year(),
		"published_month": a.PublishedAt.Month(),
		"published_date":  a.PublishedAt.Day(),
	}
	return models.EditArticle(a.Id, article, a.TagId)
}

func (a *Article) Exists() error {
	return models.ArticleExists(a.Title, a.Path, a.Id)
}

func (a *Article) GetById() (*models.Article, error) {
	return models.GetArticleById(a.Id)
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.Id)
}
