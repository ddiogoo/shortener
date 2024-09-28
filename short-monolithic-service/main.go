package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
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
		ShortCode: "abc123",
	}
}

type ShortenedRequest struct {
	Url string `json:"url"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(
		postgres.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Shortened{})

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/shorten", func(c *gin.Context) {
		var shortened []Shortened
		tx := db.Find(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, shortened)
	})

	r.POST("/shorten", func(c *gin.Context) {
		var req ShortenedRequest
		c.BindJSON(&req)
		shortened := NewShortened(req.Url)
		tx := db.Create(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":           shortened.ID,
			"rowsAffected": tx.RowsAffected,
		})
	})

	r.PUT("/shorten/:id", func(c *gin.Context) {
		var req ShortenedRequest
		c.BindJSON(&req)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}
		shortened := Shortened{
			ID:        uint(id),
			Url:       req.Url,
			ShortCode: "abc1234",
			UpdatedAt: time.Now(),
		}

		tx := db.Save(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":           shortened.ID,
			"rowsAffected": tx.RowsAffected,
		})
	})

	r.DELETE("/shorten/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}
		shortened := Shortened{
			ID: uint(id),
		}
		tx := db.Delete(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":           shortened.ID,
			"rowsAffected": tx.RowsAffected,
		})
	})
	port := func() string {
		if os.Getenv("SERVER_PORT") == "" {
			return ":8000"
		} else {
			return ":" + os.Getenv("SERVER_PORT")
		}
	}()
	r.Run(port)
}
