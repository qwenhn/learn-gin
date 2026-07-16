package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/db"
	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/handler"
	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/repository"
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

	r.GET("/api/v1/users/:uuid", userHandler.GetUserByUuid)
	r.POST("/api/v1/users", userHandler.CreateUser)

	r.Run(":8080")
}
