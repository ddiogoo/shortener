package database

import (
	"os"

	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var models []interface{}

func Config() (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	models = append(models, &model.Shortened{})
	err = autoMigrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}
