package app

import (
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/handler"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/repository"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/routes"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)
	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
