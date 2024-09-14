package app

import (
	"log"
	"wefdzen/cmd/app/handler"
	"wefdzen/internal/middleware"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	r.Use()
	r.GET("/registration", handler.Registration())
	r.GET("/deleteaccount", handler.DeleteAccount()) //change to DELETE
	r.GET("/lol", middleware.CheckAccessToken(), handler.MainPage())
	r.GET("/login", handler.Login())
	log.Fatal(r.Run(":8080"))
}
