package routes

import (
	"net/http"
	"os"
	"strconv"

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
		req := dto.NewShortenedRequest()
		c.BindJSON(&req)
		shortened := model.NewShortenedCreate(req.Url)
		tx := db.Create(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, dto.NewShortenedResponse(shortened.ID, tx.RowsAffected))
	})
	r.PUT("/shorten/:id", func(c *gin.Context) {
		req := dto.NewShortenedRequest()
		c.BindJSON(&req)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}
		shortened := model.NewShortenedUpdate(uint(id), req.Url)
		tx := db.Save(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, dto.NewShortenedResponse(shortened.ID, tx.RowsAffected))
	})
	r.DELETE("/shorten/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}
		shortened := model.NewShortenedDelete(uint(id))
		tx := db.Delete(&shortened)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": tx.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, dto.NewShortenedResponse(shortened.ID, tx.RowsAffected))
	})
	return r, func() string {
		if os.Getenv("PORT") == "" {
			return ":8000"
		} else {
			return ":" + os.Getenv("PORT")
		}
	}()
}
