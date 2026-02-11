package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/yang/go-learning-backend/examples/week13"
)

func main() {
	app := week13.NewOrderApp()
	idem := week13.NewIdempotencyStore()
	limiter := week13.NewRateLimiter(2, time.Minute)
	mux := week13.NewGovernedMux(app, idem, limiter)

	call := func(name, key string) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
		req.Header.Set("X-User-ID", "u1")
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		fmt.Printf("%s => status=%d replay=%s body=%s", name, rec.Code, rec.Header().Get("X-Idempotent-Replay"), rec.Body.String())
	}

	call("first", "idem-a")
	call("replay", "idem-a")
	call("third-no-key", "")
}
