package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yang/go-learning-backend/internal/app"
	"github.com/yang/go-learning-backend/internal/handler"
	"github.com/yang/go-learning-backend/internal/store/memory"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	todoStore := memory.NewTodoStore()
	todoHandler := handler.NewTodoHandler(todoStore)
	router := app.NewRouter(todoHandler)

	addr := ":" + port
	log.Printf("server is listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
