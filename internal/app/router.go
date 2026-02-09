package app

import (
	"net/http"
	"strings"

	"github.com/yang/go-learning-backend/internal/handler"
)

func NewRouter(todoHandler *handler.TodoHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handler.Health)
	mux.HandleFunc("POST /api/v1/todos", todoHandler.Create)
	mux.HandleFunc("GET /api/v1/todos", todoHandler.List)

	mux.HandleFunc("PATCH /api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/done") {
			http.NotFound(w, r)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
		id = strings.TrimSuffix(id, "/done")
		id = strings.Trim(id, "/")
		if id == "" {
			http.NotFound(w, r)
			return
		}

		todoHandler.MarkDone(w, r, id)
	})

	return mux
}
