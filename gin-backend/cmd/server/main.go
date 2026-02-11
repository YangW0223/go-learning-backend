package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yang/go-learning-backend/gin-backend/internal/bootstrap"
)

// main 负责进程级生命周期管理：
// 1) 初始化应用依赖；
// 2) 启动 HTTP 服务；
// 3) 监听退出信号并执行优雅停机；
// 4) 释放数据库等底层资源。
func main() {
	// 先创建根上下文，后续用于应用初始化和子上下文派生。
	ctx := context.Background()

	// bootstrap.New 会完成配置加载、DB/Redis 连接、router 组装等工作。
	app, err := bootstrap.New(ctx)
	if err != nil {
		// 启动阶段失败属于不可恢复错误，直接 panic 让进程退出。
		panic(err)
	}

	// 无论正常退出还是异常退出，都确保关闭资源。
	defer func() {
		if closeErr := app.Close(); closeErr != nil {
			app.Logger.Printf("level=error msg=close app resources failed err=%v", closeErr)
		}
	}()

	// errCh 接收服务 goroutine 上报的运行期错误。
	errCh := make(chan error, 1)
	go func() {
		// 输出启动日志，包含监听地址与环境标识，便于排障。
		app.Logger.Printf("level=info msg=server start addr=%s env=%s", app.Server.Addr, app.Config.App.Env)

		// ListenAndServe 返回 ErrServerClosed 代表正常关闭，不当作异常上报。
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// sigCh 接收系统终止信号，用于触发优雅停机流程。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 两类退出路径：
	// 1) 服务异常退出；
	// 2) 收到系统信号。
	select {
	case err := <-errCh:
		app.Logger.Fatalf("level=error msg=server exited unexpectedly err=%v", err)
	case sig := <-sigCh:
		app.Logger.Printf("level=info msg=received signal signal=%s", sig.String())
	}

	// 派生带超时的 shutdown 上下文，避免停机无限阻塞。
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.Config.HTTP.ShutdownTimeout)
	defer cancel()

	// 优雅停机会停止接收新连接，并等待在途请求完成（直到超时）。
	if err := app.Server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Printf("level=warn msg=graceful shutdown failed err=%v", err)
		// 优雅停机失败后尝试强制关闭，确保进程可退出。
		if closeErr := app.Server.Close(); closeErr != nil {
			app.Logger.Printf("level=error msg=force close failed err=%v", closeErr)
		}
	}

	// 输出最终停机日志。
	app.Logger.Printf("level=info msg=server stopped")
}
