// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "bytes"
	"bytes"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week10"
	"github.com/yang/go-learning-backend/examples/week10"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: svc := week10.NewService()
	svc := week10.NewService()
	// 详细注释: mux := week10.NewMux(svc)
	mux := week10.NewMux(svc)

	// 详细注释: createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"testing system demo"}`))
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"testing system demo"}`))
	// 详细注释: createRec := httptest.NewRecorder()
	createRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(createRec, createReq)
	mux.ServeHTTP(createRec, createReq)
	// 详细注释: fmt.Printf("POST /api/v1/todos => %d %s", createRec.Code, createRec.Body.String())
	fmt.Printf("POST /api/v1/todos => %d %s", createRec.Code, createRec.Body.String())

	// 详细注释: getReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
	// 详细注释: getRec := httptest.NewRecorder()
	getRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(getRec, getReq)
	mux.ServeHTTP(getRec, getReq)
	// 详细注释: fmt.Printf("GET /api/v1/todos/1 => %d %s", getRec.Code, getRec.Body.String())
	fmt.Printf("GET /api/v1/todos/1 => %d %s", getRec.Code, getRec.Body.String())
	// 详细注释: }
}
