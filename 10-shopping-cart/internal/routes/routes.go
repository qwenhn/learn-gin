package routes

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/middleware"
	v1routes "github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/routes/v1"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/utils"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/auth"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/pkg/cache"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisCacheService, routes ...Route) {
	httpLogger := utils.NewLoggerWithPath("http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("recovery.log", "warning")
	rateLimiterLogger := utils.NewLoggerWithPath("rate_limiter.log", "warning")

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		middleware.CORSMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
	)

	v1Api := r.Group("/api/v1")

	middleware.InitAuthMiddleware(authService, cacheService)
	protected := v1Api.Group("")
	protected.Use(
		middleware.AuthMiddleware(),
	)

	for _, route := range routes {
		switch route.(type) {
		case *v1routes.AuthRoutes:
			route.Register(v1Api)
		default:
			route.Register(protected)
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
