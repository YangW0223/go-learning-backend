package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/config"
	"github.com/yang/go-learning-backend/gin-backend/internal/observability"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/handler"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/middleware"
)

// Handlers 聚合 HTTP handlers。
type Handlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
	Todo *handler.TodoHandler
}

// Build 创建 Gin 路由。
func Build(cfg config.Config, logger *log.Logger, metrics *observability.Metrics, h Handlers, parser middleware.TokenParser) *gin.Engine {
	// 生产环境切换到 release mode，减少调试输出和运行开销。
	if cfg.App.Env == "prod" || cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化空引擎，由我们手动注入中间件链。
	r := gin.New()

	// trusted proxies 影响 ClientIP 解析策略，配置错误时记录告警但不中断启动。
	if err := r.SetTrustedProxies(cfg.HTTP.TrustedProxies); err != nil && logger != nil {
		logger.Printf("level=warn msg=set trusted proxies failed err=%v", err)
	}

	// 全局中间件顺序：
	// request id -> recovery -> cors -> timeout -> rate limit -> access log
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())
	r.Use(middleware.Timeout(cfg.HTTP.RequestTimeout))
	r.Use(middleware.RateLimit(300, time.Minute))
	r.Use(middleware.Logger(logger, metrics))

	// 基础观测与健康检查接口。
	r.GET("/healthz", handler.Healthz)
	r.GET("/readyz", handler.Readyz)
	if metrics != nil {
		r.GET("/metrics", gin.WrapH(http.HandlerFunc(metrics.Handler)))
	}

	// /api/v1 下注册公共接口。
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", h.Auth.Register)
		api.POST("/auth/login", h.Auth.Login)
	}

	// 受保护接口统一走 JWT 中间件。
	authGroup := api.Group("")
	authGroup.Use(middleware.AuthJWT(parser))
	{
		authGroup.GET("/me", h.User.Me)
		authGroup.POST("/todos", h.Todo.Create)
		authGroup.GET("/todos", h.Todo.List)
		authGroup.PATCH("/todos/:id", h.Todo.Update)
		authGroup.DELETE("/todos/:id", h.Todo.Delete)
	}

	return r
}
