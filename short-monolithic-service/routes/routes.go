package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/database/model"
	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/routes/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) (*gin.Engine, string) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/shorten", func(c *gin.Context) {
		var shortened []model.Shortened
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
		var req dto.ShortenedRequest
		c.BindJSON(&req)
		shortened := model.NewShortened(req.Url)
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
		var req dto.ShortenedRequest
		c.BindJSON(&req)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}
		shortened := model.Shortened{
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
		shortened := model.Shortened{
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
	return r, func() string {
		if os.Getenv("PORT") == "" {
			return ":8000"
		} else {
			return ":" + os.Getenv("PORT")
		}
	}()
}
