package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yang/go-learning-backend/gin-backend/internal/auth"
	"github.com/yang/go-learning-backend/gin-backend/internal/config"
	"github.com/yang/go-learning-backend/gin-backend/internal/observability"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
	postgresrepo "github.com/yang/go-learning-backend/gin-backend/internal/repository/postgres"
	redisrepo "github.com/yang/go-learning-backend/gin-backend/internal/repository/redis"
	"github.com/yang/go-learning-backend/gin-backend/internal/service"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/handler"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/middleware"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/router"
)

// App 保存运行时组件。
type App struct {
	Config config.Config
	Server *http.Server
	DB     *sql.DB
	Logger *log.Logger
}

// New 初始化应用依赖并返回可运行 App。
func New(ctx context.Context) (*App, error) {
	// 第一步：加载配置（环境变量 -> 结构体）。
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// 第二步：初始化日志与指标组件。
	logger := observability.NewStdLogger()
	metrics := observability.NewMetrics()

	// 第三步：连接 Postgres 并保证表结构存在。
	db, err := postgresrepo.Open(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}
	if err := postgresrepo.EnsureSchema(ctx, db); err != nil {
		// schema 初始化失败时立即关闭连接，避免泄漏。
		_ = db.Close()
		return nil, err
	}

	// 第四步：组装 repository 层。
	userRepo := postgresrepo.NewUserRepository(db)
	todoRepo := postgresrepo.NewTodoRepository(db)

	// 第五步：组装鉴权能力（JWT + AuthService）。
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.Issuer, cfg.Auth.TokenTTL)
	authService := service.NewAuthService(userRepo, jwtManager, cfg.Auth.PasswordMinSize)

	// 第六步：组装缓存层。默认使用 no-op 缓存，避免 Redis 可选能力影响主链路。
	var todoCache repository.TodoCache = redisrepo.NewNoopTodoCache()
	if cfg.Redis.Enabled {
		// 当启用 Redis 时，切换到真实缓存实现，并在启动阶段做连通性校验。
		redisClient := redisrepo.NewClient(redisrepo.Config{
			Addr:             cfg.Redis.Addr,
			Password:         cfg.Redis.Password,
			DB:               cfg.Redis.DB,
			DialTimeout:      cfg.Redis.DialTimeout,
			ReadWriteTimeout: cfg.Redis.ReadWriteTimeout,
		})
		redisCache := redisrepo.NewTodoCache(redisClient)
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		if err := redisCache.Ping(pingCtx); err != nil {
			cancel()
			_ = db.Close()
			return nil, fmt.Errorf("ping redis: %w", err)
		}
		cancel()
		todoCache = redisCache
		logger.Printf("level=info msg=redis cache enabled addr=%s db=%d ttl=%s", cfg.Redis.Addr, cfg.Redis.DB, cfg.Redis.CacheTTL)
	}

	// 第七步：组装业务层。
	todoService := service.NewTodoService(todoRepo, todoCache, cfg.Redis.CacheTTL)

	// 第八步：组装 HTTP handler 层。
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)

	// 第九步：构建 Gin router（包含中间件与路由注册）。
	engine := router.Build(cfg, logger, metrics, router.Handlers{
		Auth: authHandler,
		User: userHandler,
		Todo: todoHandler,
	}, middleware.AuthParserAdapter{Auth: authService})

	// 第十步：构建 HTTP server 并注入超时配置。
	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	return &App{
		Config: cfg,
		Server: server,
		DB:     db,
		Logger: logger,
	}, nil
}

// Close 释放应用资源。
func (a *App) Close() error {
	// 在 app 或 db 为空时直接返回，避免空指针风险。
	if a == nil || a.DB == nil {
		return nil
	}
	// 当前仅需关闭 DB；其他资源（如 Redis）为短连接模式不需要显式 close。
	return a.DB.Close()
}
