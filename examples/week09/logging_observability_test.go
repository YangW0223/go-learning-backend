// 详细注释: package week09
package week09

// 详细注释: import (
import (
	// 详细注释: "bytes"
	"bytes"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"
	// 详细注释: "strings"
	"strings"
	// 详细注释: "testing"
	"testing"
	// 详细注释: )
)

// TestWithObservabilityLoggingAndMetrics 验证日志字段和指标计数。
// 详细注释: func TestWithObservabilityLoggingAndMetrics(t *testing.T) {
func TestWithObservabilityLoggingAndMetrics(t *testing.T) {
	// 详细注释: buf := &bytes.Buffer{}
	buf := &bytes.Buffer{}
	// 详细注释: logger := NewJSONLogger(buf)
	logger := NewJSONLogger(buf)
	// 详细注释: metrics := NewMetrics()
	metrics := NewMetrics()
	// 详细注释: h := WithObservability(logger, metrics, NewDemoHandler())
	h := WithObservability(logger, metrics, NewDemoHandler())

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

	// 详细注释: if rec1.Code != http.StatusOK {
	if rec1.Code != http.StatusOK {
		// 详细注释: t.Fatalf("want 200 got %d", rec1.Code)
		t.Fatalf("want 200 got %d", rec1.Code)
		// 详细注释: }
	}
	// 详细注释: if rec2.Code != http.StatusInternalServerError {
	if rec2.Code != http.StatusInternalServerError {
		// 详细注释: t.Fatalf("want 500 got %d", rec2.Code)
		t.Fatalf("want 500 got %d", rec2.Code)
		// 详细注释: }
	}

	// 详细注释: logText := buf.String()
	logText := buf.String()
	// 详细注释: if !strings.Contains(logText, "request_id") {
	if !strings.Contains(logText, "request_id") {
		// 详细注释: t.Fatalf("log should contain request_id: %s", logText)
		t.Fatalf("log should contain request_id: %s", logText)
		// 详细注释: }
	}
	// 详细注释: if !strings.Contains(logText, "\"status\":500") {
	if !strings.Contains(logText, "\"status\":500") {
		// 详细注释: t.Fatalf("log should contain status=500: %s", logText)
		t.Fatalf("log should contain status=500: %s", logText)
		// 详细注释: }
	}

	// 详细注释: if metrics.RequestCount != 2 {
	if metrics.RequestCount != 2 {
		// 详细注释: t.Fatalf("want request count 2 got %d", metrics.RequestCount)
		t.Fatalf("want request count 2 got %d", metrics.RequestCount)
		// 详细注释: }
	}
	// 详细注释: if metrics.ErrorCount != 1 {
	if metrics.ErrorCount != 1 {
		// 详细注释: t.Fatalf("want error count 1 got %d", metrics.ErrorCount)
		t.Fatalf("want error count 1 got %d", metrics.ErrorCount)
		// 详细注释: }
	}
	// 详细注释: }
}
