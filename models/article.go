package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type Article struct {
	Model

	Path           string     `json:"path"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Content        string     `json:"content"`
	PublishedYear  int        `json:"published_year"`
	PublishedMonth time.Month `json:"published_month"`
	PublishedDate  int        `json:"published_date"`
	Tags           []Tag      `gorm:"many2many:article_tags;"json:"tags"`
}

func GetArticles(offset int,
	limit int,
	conditions map[string]interface{},
	tagConditions map[string]interface{},
) ([]*Article, error) {
	var articles []*Article
	var err error

	var condStrings []string
	var condArgs []interface{}
	for key, value := range conditions {
		condStrings = append(condStrings, key)
		condArgs = append(condArgs, value)
	}

	err = Query().
		Preload("Tags", tagConditions).
		Where(strings.Join(condStrings, " and "), condArgs...).
		Offset(offset).
		Limit(limit).
		Find(&articles).
		Error
	if err != nil {
		return nil, err
	}

	return articles, err
}

func CountAll(conditions map[string]interface{},
	tagCondition map[string]interface{}) (int, error) {
	var count int
	var tag Tag
	var err error

	err = Query().Model(&Tag{}).Where(tagCondition).First(&tag).Error
	if err != nil {
		return 0, nil
	}

	var condStrings []string
	var condArgs []interface{}
	for key, value := range conditions {
		condStrings = append(condStrings, key)
		condArgs = append(condArgs, value)
	}
	condStrings = append(condStrings, "tag_id=?")
	condArgs = append(condArgs, tag.Id)
	err = Query().Model(&Article{}).
		Where(strings.Join(condStrings, " and "), condArgs...).
		Count(&count).Error
	return count, err
}

func GetArticleByPath(path string) (*Article, error) {
	var article Article
	var err error
	err = Query().Preload("Tags").Where("path=?", path).First(&article).Error

	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetArticleById(id uint) (*Article, error) {
	var article Article
	var err error
	err = Query().Preload("Tags").Where("id=?", id).First(&article).Error

	if err != nil {
		return nil, err
	}
	return &article, nil
}

func AddArticle(data map[string]interface{}) error {
	return Query().Transaction(func(tx *gorm.DB) error {
		article := Article{
			Path:           data["path"].(string),
			Title:          data["title"].(string),
			Description:    data["description"].(string),
			Content:        data["content"].(string),
			PublishedYear:  data["published_year"].(int),
			PublishedMonth: data["published_month"].(time.Month),
			PublishedDate:  data["published_date"].(int),
		}

		var err error

		if err = tx.Create(&article).Error; err != nil {
			return err
		}
		err = addArticleTagRelation(tx, article.Id, data["tag_id"].([]uint))
		if err != nil {
			return err
		}
		return nil
	})
}

func addArticleTagRelation(tx *gorm.DB, articleId uint, tags []uint) error {
	tagsNum := len(tags)
	valueStrings := make([]string, 0, tagsNum)
	valueArgs := make([]interface{}, 0, tagsNum*2)

	for _, tagId := range tags {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, articleId)
		valueArgs = append(valueArgs, tagId)
	}

	return tx.Exec(
		fmt.Sprintf(
			"INSERT INTO gb_article_tags (article_id, tag_id) VALUES %s",
			strings.Join(valueStrings, ", "),
		), valueArgs...).Error
}

func ArticleExists(title string, path string, id uint) error {
	var err error
	query := Query()
	if id > 0 {
		query = Query().Where("id!=?", id)
	}

	var count int
	err = query.Model(&Article{}).Where("(title=? or path=?)", title, path).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("title or path is exist")
	}

	return nil
}

func EditArticle(id uint, data interface{}, tags []uint) error {
	return Query().Transaction(func(tx *gorm.DB) error {
		var article Article
		var err error

		if err = tx.Where("id=?", id).Model(&article).Updates(data).Error; err != nil {
			return err
		}

		err = tx.Table("gb_article_tags").
			Where("article_id=?", id).Delete(nil).Error
		if err != nil {
			return err
		}

		err = addArticleTagRelation(tx, id, tags)
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteArticle(id uint) error {
	return Query().Transaction(func(tx *gorm.DB) error {
		var err error
		if err = tx.Where("id=?", id).Delete(&Article{}).Error; err != nil {
			return err
		}

		err = tx.Table("gb_article_tags").
			Where("article_id=?", id).Delete(nil).Error
		if err != nil {
			return err
		}

		return nil
	})
}
