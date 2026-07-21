package main

import (
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/app"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/config"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/utils"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/logger"
)

func main() {
	rootDir := utils.GetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	err := godotenv.Load(filepath.Join(rootDir, ".env"))

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_ENV", "development"),
	})

	if err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env in api process")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application, err := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	// Start server
	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Application run failed")
	}
}
