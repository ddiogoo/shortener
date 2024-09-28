package model

import (
	"time"

	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/database/key"
	"gorm.io/gorm"
)

type Shortened struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Url       string `gorm:"not null"`
	ShortCode string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewShortened(url string) *Shortened {
	return &Shortened{
		Url:       url,
		ShortCode: key.MD5Hash(url),
	}
}
