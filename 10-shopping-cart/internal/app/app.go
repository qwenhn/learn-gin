package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/config"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/db"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/db/sqlc"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/routes"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/utils"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/validation"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/auth"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/cache"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/logger"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/mail"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/rabbitmq"
)

type Module interface {
	Routes() routes.Route
}

type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

func NewApplication(cfg *config.Config) (*Application, error) {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Validator init failed")
		return nil, err
	}

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	r := gin.Default()

	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed")
	}

	redisClient := config.NewRedisClient()
	redisCacheService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(redisCacheService)
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory")
		return nil, err
	}

	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to initialize mail service")
		return nil, err
	}

	rabbitmqLogger := utils.NewLoggerWithPath("worker.log", "info")
	rabbitmqService, _ := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"), rabbitmqLogger,
	)

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, redisCacheService, mailService, rabbitmqService),
	}

	routes.RegisterRoutes(r, tokenService, redisCacheService, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}, nil
}

func (a *Application) Run() error {
	server := &http.Server{
		Addr:    a.config.ServerAddr,
		Handler: a.router,
	}

	// Graceful restart or stop
	quit := make(chan os.Signal, 1)

	// SIGINT  - Triggered by Ctrl+C; indicates the user wants to stop the program.
	// SIGTERM - Sent by `kill`, `docker stop`, or Kubernetes; requests a graceful shutdown.
	// SIGHUP  - Commonly used to reload the application's configuration without restarting.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("✅ Server is running at %s", a.config.ServerAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("⛔️ Failed to start server")
		}
	}()

	<-quit // block
	logger.Log.Warn().Msg("⚠️  Shutdown signal received ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("⛔️ Server forced to shutdown")
	}

	logger.Log.Info().Msg("🍺 Server exited gracefully")

	return nil
}

func getModuleRoutes(modules []Module) []routes.Route {
	routes := make([]routes.Route, 0)

	for _, module := range modules {
		routes = append(routes, module.Routes())
	}

	return routes
}
