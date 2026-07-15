package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/demo", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome, Gin Framework"})
	})

	r.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "List Users",
		})
	})

	r.GET("/users/:id", func(ctx *gin.Context) {
		userID := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{
			"data": "User details for ID: " + userID,
		})
	})

	r.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "List Products",
		})
	})

	r.GET("/products/detail/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		price := ctx.Query("price")
		color := ctx.Query("color")

		ctx.JSON(http.StatusOK, gin.H{
			"data":  "Detail Product",
			"name":  name,
			"price": price,
			"color": color,
		})
	})

	r.Run(":8080")
}
