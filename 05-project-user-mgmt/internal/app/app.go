package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/config"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/middleware"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/routes"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/validation"
)

type Application struct {
	router *gin.Engine
	config *config.Config
	module []Module
}

type Module interface {
	Routes() routes.Route
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}

	loadEnv()

	r := gin.Default()

	go middleware.CleanupClients()

	modules := []Module{
		NewUserModule(),
	}

	routes.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		router: r,
		config: cfg,
		module: modules,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddr)
}

func loadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("No .env file found")
	}
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
