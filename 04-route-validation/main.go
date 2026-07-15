package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"

	v1Handler "github.com/qwenhn/gin-restful-api/04-route-validation/internal/api/v1/handler"
	v2Handler "github.com/qwenhn/gin-restful-api/04-route-validation/internal/api/v2/handler"
	"github.com/qwenhn/gin-restful-api/04-route-validation/middleware"
	"github.com/qwenhn/gin-restful-api/04-route-validation/utils"
)

func main() {

	if err := utils.RegisterValidators(); err != nil {
		panic(err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	r := gin.Default()

	go middleware.CleanupClients()

	r.Use(middleware.LoggerMiddleware(), middleware.ApiKeyMiddleware(), middleware.RateLimitingMiddleware())

	v1 := r.Group("/v1")
	{
		users := v1.Group("/users")
		{
			userHandlerV1 := v1Handler.NewUserHandler()
			users.GET("", userHandlerV1.GetAllUsers)
			users.GET("/:id", userHandlerV1.GetUsersById)
			users.POST("", userHandlerV1.CreateUser)
			users.PUT("/:id", userHandlerV1.UpdateUser)
			users.DELETE("/:id", userHandlerV1.DeleteUser)
			users.GET("/uuid/:uid", userHandlerV1.GetByUuid)
		}

		products := v1.Group("/products")
		{
			productHandlerV1 := v1Handler.NewProductHandler()
			products.GET("", productHandlerV1.GetAllProducts)
			products.GET("/:id", productHandlerV1.GetProductsById)
			products.POST("", productHandlerV1.CreateProduct)
			products.PUT("/:id", productHandlerV1.UpdateProduct)
			products.DELETE("/:id", productHandlerV1.DeleteProduct)
			products.GET("/by-slug/:slug", productHandlerV1.GetBySlug)
		}

		categories := v1.Group("/categories")
		{
			categoryHandlerV1 := v1Handler.NewCategoryHandler()
			categories.GET("/:category", categoryHandlerV1.GetByCategory)
		}

		news := v1.Group("/news")
		{
			newHandlerV1 := v1Handler.NewNewHandler()
			news.GET("", newHandlerV1.GetAllNews)
			news.GET("/:slug", middleware.SimpleMiddleware(), newHandlerV1.GetAllNews)
			news.POST("", newHandlerV1.CreateNew)
			news.POST("/upload-file", newHandlerV1.UploadFileNew)
			news.POST("/upload-multiple-file", newHandlerV1.UploadMultipleFileNew)
		}
	}

	v2 := r.Group("/v2")
	{
		users := v2.Group("/users")
		{
			userHandlerV2 := v2Handler.NewUserHandler()
			users.GET("", userHandlerV2.GetAllUsers)
			users.GET("/:id", userHandlerV2.GetUsersById)
			users.POST("", userHandlerV2.CreateUser)
			users.PUT("/:id", userHandlerV2.UpdateUser)
			users.DELETE("/:id", userHandlerV2.DeleteUser)
		}
	}

	r.StaticFS("/images", gin.Dir("./uploads", false)) // do not allow directory listing. If a user visits /images/, they'll receive a 404 instead of a file listing
	r.Run(":8080")
}
