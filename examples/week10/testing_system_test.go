// 详细注释: package week10
package week10

// 详细注释: import (
import (
	// 详细注释: "bytes"
	"bytes"
	// 详细注释: "encoding/json"
	"encoding/json"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"
	// 详细注释: "testing"
	"testing"
	// 详细注释: )
)

// TestServiceCreateTodoTableDriven 验证 service 表驱动测试模式。
// 详细注释: func TestServiceCreateTodoTableDriven(t *testing.T) {
func TestServiceCreateTodoTableDriven(t *testing.T) {
	// 详细注释: svc := NewService()
	svc := NewService()
	// 详细注释: cases := []struct {
	cases := []struct {
		// 详细注释: name    string
		name string
		// 详细注释: title   string
		title string
		// 详细注释: wantErr bool
		wantErr bool
		// 详细注释: }{
	}{
		// 详细注释: {name: "valid", title: "learn testing", wantErr: false},
		{name: "valid", title: "learn testing", wantErr: false},
		// 详细注释: {name: "too short", title: "ab", wantErr: true},
		{name: "too short", title: "ab", wantErr: true},
		// 详细注释: {name: "blank", title: "   ", wantErr: true},
		{name: "blank", title: "   ", wantErr: true},
		// 详细注释: }
	}
	// 详细注释: for _, tc := range cases {
	for _, tc := range cases {
		// 详细注释: _, err := svc.CreateTodo(tc.title)
		_, err := svc.CreateTodo(tc.title)
		// 详细注释: if tc.wantErr && err == nil {
		if tc.wantErr && err == nil {
			// 详细注释: t.Fatalf("%s: expected error", tc.name)
			t.Fatalf("%s: expected error", tc.name)
			// 详细注释: }
		}
		// 详细注释: if !tc.wantErr && err != nil {
		if !tc.wantErr && err != nil {
			// 详细注释: t.Fatalf("%s: unexpected error %v", tc.name, err)
			t.Fatalf("%s: unexpected error %v", tc.name, err)
			// 详细注释: }
		}
		// 详细注释: }
	}
	// 详细注释: }
}

// TestHandlerScenarios 验证 201/400/404 路径。
// 详细注释: func TestHandlerScenarios(t *testing.T) {
func TestHandlerScenarios(t *testing.T) {
	// 详细注释: svc := NewService()
	svc := NewService()
	// 详细注释: mux := NewMux(svc)
	mux := NewMux(svc)

	// 详细注释: createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"write tests"}`))
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"write tests"}`))
	// 详细注释: createRec := httptest.NewRecorder()
	createRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(createRec, createReq)
	mux.ServeHTTP(createRec, createReq)
	// 详细注释: if createRec.Code != http.StatusCreated {
	if createRec.Code != http.StatusCreated {
		// 详细注释: t.Fatalf("create want 201 got %d", createRec.Code)
		t.Fatalf("create want 201 got %d", createRec.Code)
		// 详细注释: }
	}

	// 详细注释: badReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x"}`))
	badReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x"}`))
	// 详细注释: badRec := httptest.NewRecorder()
	badRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(badRec, badReq)
	mux.ServeHTTP(badRec, badReq)
	// 详细注释: if badRec.Code != http.StatusBadRequest {
	if badRec.Code != http.StatusBadRequest {
		// 详细注释: t.Fatalf("bad title want 400 got %d", badRec.Code)
		t.Fatalf("bad title want 400 got %d", badRec.Code)
		// 详细注释: }
	}

	// 详细注释: notFoundReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/999", nil)
	notFoundReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/999", nil)
	// 详细注释: notFoundRec := httptest.NewRecorder()
	notFoundRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(notFoundRec, notFoundReq)
	mux.ServeHTTP(notFoundRec, notFoundReq)
	// 详细注释: if notFoundRec.Code != http.StatusNotFound {
	if notFoundRec.Code != http.StatusNotFound {
		// 详细注释: t.Fatalf("not found want 404 got %d", notFoundRec.Code)
		t.Fatalf("not found want 404 got %d", notFoundRec.Code)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestE2EAPI 验证端到端 API 行为。
// 详细注释: func TestE2EAPI(t *testing.T) {
func TestE2EAPI(t *testing.T) {
	// 详细注释: svc := NewService()
	svc := NewService()
	// 详细注释: mux := NewMux(svc)
	mux := NewMux(svc)

	// 详细注释: createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"e2e todo"}`))
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"e2e todo"}`))
	// 详细注释: createRec := httptest.NewRecorder()
	createRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(createRec, createReq)
	mux.ServeHTTP(createRec, createReq)
	// 详细注释: if createRec.Code != http.StatusCreated {
	if createRec.Code != http.StatusCreated {
		// 详细注释: t.Fatalf("want 201 got %d", createRec.Code)
		t.Fatalf("want 201 got %d", createRec.Code)
		// 详细注释: }
	}

	// 详细注释: var body map[string]map[string]any
	var body map[string]map[string]any
	// 详细注释: if err := json.NewDecoder(createRec.Body).Decode(&body); err != nil {
	if err := json.NewDecoder(createRec.Body).Decode(&body); err != nil {
		// 详细注释: t.Fatalf("decode err: %v", err)
		t.Fatalf("decode err: %v", err)
		// 详细注释: }
	}
	// 详细注释: id, _ := body["data"]["id"].(string)
	id, _ := body["data"]["id"].(string)
	// 详细注释: if id == "" {
	if id == "" {
		// 详细注释: t.Fatalf("expected id in response")
		t.Fatalf("expected id in response")
		// 详细注释: }
	}

	// 详细注释: getReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/"+id, nil)
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/"+id, nil)
	// 详细注释: getRec := httptest.NewRecorder()
	getRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(getRec, getReq)
	mux.ServeHTTP(getRec, getReq)
	// 详细注释: if getRec.Code != http.StatusOK {
	if getRec.Code != http.StatusOK {
		// 详细注释: t.Fatalf("want 200 got %d", getRec.Code)
		t.Fatalf("want 200 got %d", getRec.Code)
		// 详细注释: }
	}
	// 详细注释: }
}
