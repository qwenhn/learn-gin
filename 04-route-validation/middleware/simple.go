package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Middleware check started")
		ctx.Writer.Write([]byte("Middleware check started")) // Write before the handler response

		ctx.Next()

		log.Println("Middleware check completed")
		ctx.Writer.Write([]byte("Middleware check completed")) // Write after the handler response
	}
}
