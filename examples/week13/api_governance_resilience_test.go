package week13

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newWeek13Mux(limit int) *http.ServeMux {
	app := NewOrderApp()
	idem := NewIdempotencyStore()
	limiter := NewRateLimiter(limit, time.Minute)
	return NewGovernedMux(app, idem, limiter)
}

// TestCreateOrderSuccess 验证成功创建返回 201。
func TestCreateOrderSuccess(t *testing.T) {
	mux := newWeek13Mux(10)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
	req.Header.Set("X-User-ID", "u1")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201 got %d body=%s", rec.Code, rec.Body.String())
	}
}

// TestIdempotencyReplay 验证重复幂等键返回相同结果。
func TestIdempotencyReplay(t *testing.T) {
	mux := newWeek13Mux(10)
	makeReq := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
		req.Header.Set("X-User-ID", "u1")
		req.Header.Set("Idempotency-Key", "idem-1")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}

	r1 := makeReq()
	r2 := makeReq()
	if r1.Code != http.StatusCreated || r2.Code != http.StatusCreated {
		t.Fatalf("both responses should be 201, got %d %d", r1.Code, r2.Code)
	}
	if r1.Body.String() != r2.Body.String() {
		t.Fatalf("idempotency body should match")
	}
	if r2.Header().Get("X-Idempotent-Replay") != "1" {
		t.Fatalf("second request should be replay")
	}
}

// TestRateLimit429 验证超限返回 429。
func TestRateLimit429(t *testing.T) {
	mux := newWeek13Mux(1)

	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
	req1.Header.Set("X-User-ID", "u1")
	rec1 := httptest.NewRecorder()
	mux.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"pen"}`))
	req2.Header.Set("X-User-ID", "u1")
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("want 429 got %d", rec2.Code)
	}
	if !strings.Contains(rec2.Body.String(), "RATE_LIMITED") {
		t.Fatalf("body should include RATE_LIMITED, got %s", rec2.Body.String())
	}
}
