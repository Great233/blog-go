package models

import (
	"errors"
)

type Tag struct {
	Model
	Title string `json:"title"`
}

func GetAllTags(offset int, limit int, condition map[string]interface{}) ([]*Tag, error) {
	var tags []*Tag

	query := Query()
	if len(condition) > 0 {
		query = query.Where(condition)
	}
	if offset > -1 && limit > -1 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Find(&tags).Error

	if err != nil {
		return nil, err
	}
	return tags, nil
}

func CountAllTags(condition map[string]interface{}) (int, error) {
	query := Query()

	if len(condition) > 0 {
		query = query.Where(condition)
	}
	var total int
	err := query.Count(total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetTag(condition map[string]interface{}) (*Tag, error) {
	var tag Tag
	var err error

	err = Query().Where(condition).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func TagExists(condition map[string]interface{}) error {
	var total int
	query := Query()

	if id, ok := condition["id"]; ok {
		query.Where("id != ", id)
	}

	if title, ok := condition["title"]; ok {
		query.Where("title=?", title)
	} else {
		return errors.New("condition must contains \"title\"")
	}

	err := query.Count(&total).Error
	if err != nil {
		return err
	}
	return nil
}

func AddTag(tag Tag) error {
	return Query().Create(&tag).Error
}

func EditTag(tag Tag) error {
	return Query().Where("id=?", tag.Id).Update(map[string]interface{}{
		"title": tag.Title,
	}).Error
}

func DeleteTag(tag Tag) error {
	return Query().Where("id=?", tag.Id).Delete(tag).Error
}
