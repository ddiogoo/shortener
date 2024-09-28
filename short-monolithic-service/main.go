package main

import (
	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/database"
	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := database.Config()
	r, port := routes.Routes(db)
	r.Run(port)
}
