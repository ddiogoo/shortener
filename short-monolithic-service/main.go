package main

import (
	"github.com/ddiogoo/shortener/tree/master/short-monolithic-service/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	r, port := routes.Routes()
	r.Run(port)
}
