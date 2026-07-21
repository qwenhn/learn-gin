package app

import (
	v1handler "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/handler/v1"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/repository"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/routes"
	v1routes "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/routes/v1"
	v1service "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule(ctx *ModuleContext) *UserModule {
	userRepo := repository.NewSQLUserRepository(ctx.DB)
	userService := v1service.NewUserService(userRepo, ctx.Redis)
	userHandler := v1handler.NewUserHandler(userService)
	userRoutes := v1routes.NewUserRoutes(userHandler)
	return &UserModule{
		routes: userRoutes,
	}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
