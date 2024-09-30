package models

import (
	"gorm.io/gorm"
)

type Shortened struct {
	gorm.Model
	Id        uint   `gorm:"primaryKey"`
	Url       string `gorm:"not null"`
	ShortCode string `gorm:"unique"`
	CreatedAt string `gorm:"not null"`
	UpdatedAt string `gorm:"not null"`
}
