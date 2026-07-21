package app

import (
	v1handler "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/handler/v1"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/repository"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/routes"
	v1routes "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/routes/v1"
	v1service "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/service/v1"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/auth"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/cache"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/mail"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/rabbitmq"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitMQService rabbitmq.RabbitMQService) *AuthModule {
	userRepo := repository.NewSQLUserRepository(ctx.DB)
	authService := v1service.NewAuthService(userRepo, tokenService, cacheService, mailService, rabbitMQService)
	authHandler := v1handler.NewAuthHandler(authService)
	authRoutes := v1routes.NewAuthRoutes(authHandler)
	return &AuthModule{
		routes: authRoutes,
	}
}

func (am *AuthModule) Routes() routes.Route {
	return am.routes
}
