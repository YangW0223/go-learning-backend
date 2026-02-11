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
	// 详细注释: "time"
	"time"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week13"
	"github.com/yang/go-learning-backend/examples/week13"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: app := week13.NewOrderApp()
	app := week13.NewOrderApp()
	// 详细注释: idem := week13.NewIdempotencyStore()
	idem := week13.NewIdempotencyStore()
	// 详细注释: limiter := week13.NewRateLimiter(2, time.Minute)
	limiter := week13.NewRateLimiter(2, time.Minute)
	// 详细注释: mux := week13.NewGovernedMux(app, idem, limiter)
	mux := week13.NewGovernedMux(app, idem, limiter)

	// 详细注释: call := func(name, key string) {
	call := func(name, key string) {
		// 详细注释: req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{"item":"book"}`))
		// 详细注释: req.Header.Set("X-User-ID", "u1")
		req.Header.Set("X-User-ID", "u1")
		// 详细注释: if key != "" {
		if key != "" {
			// 详细注释: req.Header.Set("Idempotency-Key", key)
			req.Header.Set("Idempotency-Key", key)
			// 详细注释: }
		}
		// 详细注释: rec := httptest.NewRecorder()
		rec := httptest.NewRecorder()
		// 详细注释: mux.ServeHTTP(rec, req)
		mux.ServeHTTP(rec, req)
		// 详细注释: fmt.Printf("%s => status=%d replay=%s body=%s", name, rec.Code, rec.Header().Get("X-Idempotent-Replay"), rec.Body.String())
		fmt.Printf("%s => status=%d replay=%s body=%s", name, rec.Code, rec.Header().Get("X-Idempotent-Replay"), rec.Body.String())
		// 详细注释: }
	}

	// 详细注释: call("first", "idem-a")
	call("first", "idem-a")
	// 详细注释: call("replay", "idem-a")
	call("replay", "idem-a")
	// 详细注释: call("third-no-key", "")
	call("third-no-key", "")
	// 详细注释: }
}
