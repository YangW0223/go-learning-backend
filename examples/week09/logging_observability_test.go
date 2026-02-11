package week09

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestWithObservabilityLoggingAndMetrics 验证日志字段和指标计数。
func TestWithObservabilityLoggingAndMetrics(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewJSONLogger(buf)
	metrics := NewMetrics()
	h := WithObservability(logger, metrics, NewDemoHandler())

	req1 := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/todos?fail=1", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	if rec1.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", rec1.Code)
	}
	if rec2.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 got %d", rec2.Code)
	}

	logText := buf.String()
	if !strings.Contains(logText, "request_id") {
		t.Fatalf("log should contain request_id: %s", logText)
	}
	if !strings.Contains(logText, "\"status\":500") {
		t.Fatalf("log should contain status=500: %s", logText)
	}

	if metrics.RequestCount != 2 {
		t.Fatalf("want request count 2 got %d", metrics.RequestCount)
	}
	if metrics.ErrorCount != 1 {
		t.Fatalf("want error count 1 got %d", metrics.ErrorCount)
	}
}
