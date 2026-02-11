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

	// 详细注释: "github.com/yang/go-learning-backend/examples/week09"
	"github.com/yang/go-learning-backend/examples/week09"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: buf := &bytes.Buffer{}
	buf := &bytes.Buffer{}
	// 详细注释: logger := week09.NewJSONLogger(buf)
	logger := week09.NewJSONLogger(buf)
	// 详细注释: metrics := week09.NewMetrics()
	metrics := week09.NewMetrics()
	// 详细注释: h := week09.WithObservability(logger, metrics, week09.NewDemoHandler())
	h := week09.WithObservability(logger, metrics, week09.NewDemoHandler())

	// 详细注释: req1 := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req1 := httptest.NewRequest(http.MethodGet, "/todos", nil)
	// 详细注释: rec1 := httptest.NewRecorder()
	rec1 := httptest.NewRecorder()
	// 详细注释: h.ServeHTTP(rec1, req1)
	h.ServeHTTP(rec1, req1)

	// 详细注释: req2 := httptest.NewRequest(http.MethodGet, "/todos?fail=1", nil)
	req2 := httptest.NewRequest(http.MethodGet, "/todos?fail=1", nil)
	// 详细注释: rec2 := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()
	// 详细注释: h.ServeHTTP(rec2, req2)
	h.ServeHTTP(rec2, req2)

	// 详细注释: fmt.Println("request1 status:", rec1.Code)
	fmt.Println("request1 status:", rec1.Code)
	// 详细注释: fmt.Println("request2 status:", rec2.Code)
	fmt.Println("request2 status:", rec2.Code)
	// 详细注释: fmt.Println("metrics:", metrics.Snapshot())
	fmt.Println("metrics:", metrics.Snapshot())
	// 详细注释: fmt.Println("logs:")
	fmt.Println("logs:")
	// 详细注释: fmt.Print(buf.String())
	fmt.Print(buf.String())
	// 详细注释: }
}
