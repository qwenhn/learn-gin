package main

import (
	"github.com/gin-gonic/gin"

	v1Handler "github.com/qwenhn/gin-restful-api/03-route-group/internal/api/v1/handler"
	v2Handler "github.com/qwenhn/gin-restful-api/03-route-group/internal/api/v2/handler"
)

func main() {
	r := gin.Default()

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

	r.Run(":8080")
}
