package models

import (
	"blog/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strings"
	"time"
)

var db *gorm.DB

type Model struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt *int64 `json:"deleted_at"`
}

func Init() {
	var err error

	dbConfig := config.App.Database

	db, err = gorm.Open(
		strings.ToLower(dbConfig.Driver),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
			dbConfig.Charset,
		),
	)

	if err != nil {
		log.Fatalf("%s connect error: %v", dbConfig.Driver, err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return dbConfig.Prefix + defaultTableName
	}

	db.LogMode(config.App.Server.Mode == "debug")

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimestampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimestampForUpdateCallback)
}

func Query() *gorm.DB {
	return db
}

func updateTimestampForCreateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}

	now := time.Now().Unix()

	if createAtField, ok := scope.FieldByName("CreatedAt"); ok {
		if createAtField.IsBlank {
			_ = createAtField.Set(now)
		}
	}

	if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
		if updatedAtField.IsBlank {
			_ = updatedAtField.Set(now)
		}
	}
}

func updateTimestampForUpdateCallback(scope *gorm.Scope) {
	if _, err := scope.Get("gorm:update_column"); !err {
		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			_ = updatedAtField.Set(time.Now().Unix())
		}
	}
}
