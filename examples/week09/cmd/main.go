package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/yang/go-learning-backend/examples/week09"
)

func main() {
	buf := &bytes.Buffer{}
	logger := week09.NewJSONLogger(buf)
	metrics := week09.NewMetrics()
	h := week09.WithObservability(logger, metrics, week09.NewDemoHandler())

	req1 := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/todos?fail=1", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	fmt.Println("request1 status:", rec1.Code)
	fmt.Println("request2 status:", rec2.Code)
	fmt.Println("metrics:", metrics.Snapshot())
	fmt.Println("logs:")
	fmt.Print(buf.String())
}
