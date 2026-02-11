package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yang/go-learning-backend/internal/app"
	"github.com/yang/go-learning-backend/internal/cache"
	rediscache "github.com/yang/go-learning-backend/internal/cache/redis"
	"github.com/yang/go-learning-backend/internal/config"
	"github.com/yang/go-learning-backend/internal/handler"
	"github.com/yang/go-learning-backend/internal/service"
	"github.com/yang/go-learning-backend/internal/store/memory"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	todoStore := memory.NewTodoStore()

	todoCache := cache.TodoCache(cache.NewNoopTodoCache())
	if cfg.Redis.Enabled {
		redisClient := rediscache.NewClient(rediscache.Config{
			Addr:        cfg.Redis.Addr,
			Password:    cfg.Redis.Password,
			DB:          cfg.Redis.DB,
			DialTimeout: cfg.Redis.DialTimeout,
			IOTimeout:   cfg.Redis.IOTimeout,
		})
		redisTodoCache := rediscache.NewTodoCache(redisClient)
		if err := redisTodoCache.Ping(context.Background()); err != nil {
			log.Fatalf("redis ping failed: %v", err)
		}
		todoCache = redisTodoCache
		log.Printf("redis cache enabled addr=%s db=%d ttl=%s", cfg.Redis.Addr, cfg.Redis.DB, cfg.Redis.CacheTTL)
	}

	todoService := service.NewTodoService(todoStore, todoCache, cfg.Redis.CacheTTL)
	todoHandler := handler.NewTodoHandlerWithService(todoService)
	router := app.NewRouter(todoHandler)

	addr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server is listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("server exited unexpectedly: %v", err)
	case sig := <-sigCh:
		log.Printf("received signal %s, shutting down...", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		if closeErr := server.Close(); closeErr != nil {
			log.Printf("force close failed: %v", closeErr)
		}
	}
}
