package database

import (
	"os"

	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabase[T any] struct {
	db     *gorm.DB
	models []interface{}
}

func New[T any]() (*PgDatabase[T], error) {
	db, err := gorm.Open(
		postgres.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	database := &PgDatabase[T]{db: db}
	database.addModels()
	err = database.migrate()
	if err != nil {
		return nil, err
	}
	return database, nil
}

func (pg *PgDatabase[T]) addModels() {
	pg.models = append(pg.models, &models.Shortened{})
}

func (pg *PgDatabase[T]) migrate() error {
	for _, model := range pg.models {
		err := pg.db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgDatabase[T]) RetrieveAll() ([]T, error) {
	var result []T
	tx := pg.db.Find(&result)
	if tx.Error != nil {
		return []T{}, tx.Error
	}
	return result, nil
}

func (pg *PgDatabase[T]) RetrieveOne(query interface{}, args ...interface{}) (T, error) {
	var result T
	tx := pg.db.Where(query, args).First(&result)
	if tx.Error != nil {
		return result, tx.Error
	}
	return result, nil
}

func (pg *PgDatabase[T]) Create(data T) error {
	tx := pg.db.Create(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
