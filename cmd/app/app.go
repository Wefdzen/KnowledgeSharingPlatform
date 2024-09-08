package app

import (
	"log"
	"wefdzen/cmd/app/handler"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	r.GET("/registration", handler.Registration())
	r.GET("/deleteaccount", handler.DeleteAccount()) //change to DELETE
	r.GET("/login", handler.Login())
	log.Fatal(r.Run(":8080"))
}
