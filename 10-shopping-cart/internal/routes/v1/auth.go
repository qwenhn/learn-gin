package v1routes

import (
	"github.com/gin-gonic/gin"

	v1handler "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/handler/v1"
)

type AuthRoutes struct {
	handler *v1handler.AuthHandler
}

func NewAuthRoutes(handler *v1handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		handler: handler,
	}
}

func (ar *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", ar.handler.Login)
		auth.POST("/logout", ar.handler.Logout)
		auth.POST("/refresh", ar.handler.RefreshToken)
		auth.POST("/forgot-password", ar.handler.RequestForgotPassword)
		auth.POST("/reset-password", ar.handler.ResetPassword)
	}
}
