package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
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
