package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"messenger-api/config"
	"messenger-api/services"
)

func main() {
	config.InitPostgres()
	router := gin.Default()

	router.POST("/api/user/register", services.RegisterHandler)
	router.POST("/api/user/login", services.LoginHandler)

	log.Println("Сервер запущен на :8080")
	router.Run(":8080")
}
