package main

import (
	"log"
	"wefdzen/cmd/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.StartServer()
}
