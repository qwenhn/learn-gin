package main

import (
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/app"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/config"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}
