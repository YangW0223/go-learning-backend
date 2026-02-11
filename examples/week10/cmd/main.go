package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/yang/go-learning-backend/examples/week10"
)

func main() {
	svc := week10.NewService()
	mux := week10.NewMux(svc)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"testing system demo"}`))
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	fmt.Printf("POST /api/v1/todos => %d %s", createRec.Code, createRec.Body.String())

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	fmt.Printf("GET /api/v1/todos/1 => %d %s", getRec.Code, getRec.Body.String())
}
