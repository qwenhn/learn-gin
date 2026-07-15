package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/qwenhn/gin-restful-api/08-gorm/internal/db"
	"github.com/qwenhn/gin-restful-api/08-gorm/internal/handler"
	"github.com/qwenhn/gin-restful-api/08-gorm/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := db.InitDB(); err != nil {
		log.Fatal("Unable to connect to the database")
	}

	r := gin.Default()

	userRepo := repository.NewSQLUserRepository(db.DB)
	userHandler := handler.NewUserHandler(userRepo)

	r.GET("/api/v1/users/:id", userHandler.GetUserById)
	r.POST("/api/v1/users", userHandler.CreateUser)

	r.Run(":8080")
}
