// 详细注释: package week13
package week13

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
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// 详细注释: func newWeek13Mux(limit int) *http.ServeMux {
func newWeek13Mux(limit int) *http.ServeMux {
	// 详细注释: app := NewOrderApp()
	app := NewOrderApp()
	// 详细注释: idem := NewIdempotencyStore()
	idem := NewIdempotencyStore()
	// 详细注释: limiter := NewRateLimiter(limit, time.Minute)
	limiter := NewRateLimiter(limit, time.Minute)
	// 详细注释: return NewGovernedMux(app, idem, limiter)
	return NewGovernedMux(app, idem, limiter)
	// 详细注释: }
}

// TestCreateOrderSuccess 验证成功创建返回 201。
// 详细注释: func TestCreateOrderSuccess(t *testing.T) {
func TestCreateOrderSuccess(t *testing.T) {
	// 详细注释: mux := newWeek13Mux(10)
	mux := newWeek13Mux(10)
	// 详细注释: req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
	// 详细注释: req.Header.Set("X-User-ID", "u1")
	req.Header.Set("X-User-ID", "u1")
	// 详细注释: rec := httptest.NewRecorder()
	rec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(rec, req)
	mux.ServeHTTP(rec, req)
	// 详细注释: if rec.Code != http.StatusCreated {
	if rec.Code != http.StatusCreated {
		// 详细注释: t.Fatalf("want 201 got %d body=%s", rec.Code, rec.Body.String())
		t.Fatalf("want 201 got %d body=%s", rec.Code, rec.Body.String())
		// 详细注释: }
	}
	// 详细注释: }
}

// TestIdempotencyReplay 验证重复幂等键返回相同结果。
// 详细注释: func TestIdempotencyReplay(t *testing.T) {
func TestIdempotencyReplay(t *testing.T) {
	// 详细注释: mux := newWeek13Mux(10)
	mux := newWeek13Mux(10)
	// 详细注释: makeReq := func() *httptest.ResponseRecorder {
	makeReq := func() *httptest.ResponseRecorder {
		// 详细注释: req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
		// 详细注释: req.Header.Set("X-User-ID", "u1")
		req.Header.Set("X-User-ID", "u1")
		// 详细注释: req.Header.Set("Idempotency-Key", "idem-1")
		req.Header.Set("Idempotency-Key", "idem-1")
		// 详细注释: rec := httptest.NewRecorder()
		rec := httptest.NewRecorder()
		// 详细注释: mux.ServeHTTP(rec, req)
		mux.ServeHTTP(rec, req)
		// 详细注释: return rec
		return rec
		// 详细注释: }
	}

	// 详细注释: r1 := makeReq()
	r1 := makeReq()
	// 详细注释: r2 := makeReq()
	r2 := makeReq()
	// 详细注释: if r1.Code != http.StatusCreated || r2.Code != http.StatusCreated {
	if r1.Code != http.StatusCreated || r2.Code != http.StatusCreated {
		// 详细注释: t.Fatalf("both responses should be 201, got %d %d", r1.Code, r2.Code)
		t.Fatalf("both responses should be 201, got %d %d", r1.Code, r2.Code)
		// 详细注释: }
	}
	// 详细注释: if r1.Body.String() != r2.Body.String() {
	if r1.Body.String() != r2.Body.String() {
		// 详细注释: t.Fatalf("idempotency body should match")
		t.Fatalf("idempotency body should match")
		// 详细注释: }
	}
	// 详细注释: if r2.Header().Get("X-Idempotent-Replay") != "1" {
	if r2.Header().Get("X-Idempotent-Replay") != "1" {
		// 详细注释: t.Fatalf("second request should be replay")
		t.Fatalf("second request should be replay")
		// 详细注释: }
	}
	// 详细注释: }
}

// TestRateLimit429 验证超限返回 429。
// 详细注释: func TestRateLimit429(t *testing.T) {
func TestRateLimit429(t *testing.T) {
	// 详细注释: mux := newWeek13Mux(1)
	mux := newWeek13Mux(1)

	// 详细注释: req1 := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
	// 详细注释: req1.Header.Set("X-User-ID", "u1")
	req1.Header.Set("X-User-ID", "u1")
	// 详细注释: rec1 := httptest.NewRecorder()
	rec1 := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(rec1, req1)
	mux.ServeHTTP(rec1, req1)

	// 详细注释: req2 := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"pen"}`))
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"pen"}`))
	// 详细注释: req2.Header.Set("X-User-ID", "u1")
	req2.Header.Set("X-User-ID", "u1")
	// 详细注释: rec2 := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(rec2, req2)
	mux.ServeHTTP(rec2, req2)

	// 详细注释: if rec2.Code != http.StatusTooManyRequests {
	if rec2.Code != http.StatusTooManyRequests {
		// 详细注释: t.Fatalf("want 429 got %d", rec2.Code)
		t.Fatalf("want 429 got %d", rec2.Code)
		// 详细注释: }
	}
	// 详细注释: if !strings.Contains(rec2.Body.String(), "RATE_LIMITED") {
	if !strings.Contains(rec2.Body.String(), "RATE_LIMITED") {
		// 详细注释: t.Fatalf("body should include RATE_LIMITED, got %s", rec2.Body.String())
		t.Fatalf("body should include RATE_LIMITED, got %s", rec2.Body.String())
		// 详细注释: }
	}
	// 详细注释: }
}
