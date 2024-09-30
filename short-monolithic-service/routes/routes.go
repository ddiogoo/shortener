package routes

import (
	"net/http"
	"os"

	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/database"
	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/models"
	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ShortenedResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []models.Shortened `json:"data"`
}

func Routes() (*gin.Engine, string) {
	db, err := database.New[models.Shortened]()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, PingResponse{
			Status:  "ok",
			Message: "welcome to the shortener api",
		})
	})
	r.GET("/all", func(c *gin.Context) {
		results, err := db.RetrieveAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ShortenedResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    []models.Shortened{},
			})
			return
		}
		c.JSON(http.StatusOK, ShortenedResponse{
			Status:  "success",
			Message: "all shortened urls",
			Data:    results,
		})
	})
	r.GET("/:code", func(c *gin.Context) {
		code := c.Param("code")
		result, err := db.RetrieveOne("short_code = ?", code)
		if err != nil {
			c.JSON(http.StatusNotFound, ShortenedResponse{
				Status:  "error",
				Message: "shortened url not found",
				Data:    []models.Shortened{},
			})
			return
		}
		c.JSON(http.StatusOK, ShortenedResponse{
			Status:  "success",
			Message: "shortened url",
			Data:    []models.Shortened{result},
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
