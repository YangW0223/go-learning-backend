package week10

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServiceCreateTodoTableDriven 验证 service 表驱动测试模式。
func TestServiceCreateTodoTableDriven(t *testing.T) {
	svc := NewService()
	cases := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{name: "valid", title: "learn testing", wantErr: false},
		{name: "too short", title: "ab", wantErr: true},
		{name: "blank", title: "   ", wantErr: true},
	}
	for _, tc := range cases {
		_, err := svc.CreateTodo(tc.title)
		if tc.wantErr && err == nil {
			t.Fatalf("%s: expected error", tc.name)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("%s: unexpected error %v", tc.name, err)
		}
	}
}

// TestHandlerScenarios 验证 201/400/404 路径。
func TestHandlerScenarios(t *testing.T) {
	svc := NewService()
	mux := NewMux(svc)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"write tests"}`))
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create want 201 got %d", createRec.Code)
	}

	badReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x"}`))
	badRec := httptest.NewRecorder()
	mux.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("bad title want 400 got %d", badRec.Code)
	}

	notFoundReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/999", nil)
	notFoundRec := httptest.NewRecorder()
	mux.ServeHTTP(notFoundRec, notFoundReq)
	if notFoundRec.Code != http.StatusNotFound {
		t.Fatalf("not found want 404 got %d", notFoundRec.Code)
	}
}

// TestE2EAPI 验证端到端 API 行为。
func TestE2EAPI(t *testing.T) {
	svc := NewService()
	mux := NewMux(svc)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"e2e todo"}`))
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("want 201 got %d", createRec.Code)
	}

	var body map[string]map[string]any
	if err := json.NewDecoder(createRec.Body).Decode(&body); err != nil {
		t.Fatalf("decode err: %v", err)
	}
	id, _ := body["data"]["id"].(string)
	if id == "" {
		t.Fatalf("expected id in response")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/"+id, nil)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", getRec.Code)
	}
}
