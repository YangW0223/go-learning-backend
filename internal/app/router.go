package app

import (
	"net/http"
	"strings"

	"github.com/yang/go-learning-backend/internal/handler"
)

func NewRouter(todoHandler *handler.TodoHandler) http.Handler {
	mux := http.NewServeMux()

	// /ping: Week 01 练习接口，返回固定 {"message":"pong"}。
	// 这个接口适合用来做最小可用连通性检查。
	mux.HandleFunc("GET /ping", handler.Ping)
	mux.HandleFunc("GET /healthz", handler.Health)
	mux.HandleFunc("POST /api/v1/todos", todoHandler.Create)
	mux.HandleFunc("GET /api/v1/todos", todoHandler.List)
	mux.HandleFunc("DELETE /api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
		// 期望路径: /api/v1/todos/{id}
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
		id = strings.Trim(id, "/")
		// 由 handler 做 id 格式校验和错误映射。
		todoHandler.Delete(w, r, id)
	})

	mux.HandleFunc("PATCH /api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
		// 这里只接受 .../{id}/done 这样的路径。
		// 例如 /api/v1/todos/123/done。
		if !strings.HasSuffix(r.URL.Path, "/done") {
			http.NotFound(w, r)
			return
		}

		// 从完整路径中提取 id:
		// 1) 去掉前缀 /api/v1/todos/
		// 2) 去掉后缀 /done
		// 3) 去掉残余斜杠
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
		id = strings.TrimSuffix(id, "/done")
		id = strings.Trim(id, "/")
		if id == "" {
			// 没有提取到有效 id 时返回 404，避免误操作。
			http.NotFound(w, r)
			return
		}

		// id 合法后，交由业务 handler 处理。
		todoHandler.MarkDone(w, r, id)
	})

	return mux
}
