package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Shortened struct {
	gorm.Model
	Url       string `gorm:"not null"`
	ShortCode string `gorm:"not null"`
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

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
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
		var shortened Shortened
		tx := db.First(&shortened, c.Param("id"))
		if tx.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		shortened.Url = req.Url
		shortened.UpdatedAt = time.Now()
		tx = db.Save(&shortened)
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
