// 详细注释: package week13
package week13

// 详细注释: import (
import (
	// 详细注释: "bytes"
	"bytes"
	// 详细注释: "context"
	"context"
	// 详细注释: "encoding/json"
	"encoding/json"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "strings"
	"strings"
	// 详细注释: "sync"
	"sync"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// 详细注释: var (
var (
	// ErrInvalidItem 表示下单参数非法。
	// 详细注释: ErrInvalidItem = errors.New("invalid item")
	ErrInvalidItem = errors.New("invalid item")

// 详细注释: )
)

// APIError 表示统一错误结构。
// 详细注释: type APIError struct {
type APIError struct {
	// 详细注释: Code    string `json:"code"`
	Code string `json:"code"`
	// 详细注释: Message string `json:"message"`
	Message string `json:"message"`
	// 详细注释: }
}

// WriteSuccess 输出统一成功响应。
// 详细注释: func WriteSuccess(w http.ResponseWriter, status int, requestID string, data any) {
func WriteSuccess(w http.ResponseWriter, status int, requestID string, data any) {
	// 详细注释: writeEnvelope(w, status, requestID, data, nil)
	writeEnvelope(w, status, requestID, data, nil)
	// 详细注释: }
}

// WriteError 输出统一错误响应。
// 详细注释: func WriteError(w http.ResponseWriter, status int, requestID string, code string, message string) {
func WriteError(w http.ResponseWriter, status int, requestID string, code string, message string) {
	// 详细注释: writeEnvelope(w, status, requestID, nil, &APIError{Code: code, Message: message})
	writeEnvelope(w, status, requestID, nil, &APIError{Code: code, Message: message})
	// 详细注释: }
}

// 详细注释: func writeEnvelope(w http.ResponseWriter, status int, requestID string, data any, apiErr *APIError) {
func writeEnvelope(w http.ResponseWriter, status int, requestID string, data any, apiErr *APIError) {
	// 详细注释: w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	// 详细注释: w.WriteHeader(status)
	w.WriteHeader(status)
	// 详细注释: _ = json.NewEncoder(w).Encode(map[string]any{
	_ = json.NewEncoder(w).Encode(map[string]any{
		// 详细注释: "request_id": requestID,
		"request_id": requestID,
		// 详细注释: "data":       data,
		"data": data,
		// 详细注释: "error":      apiErr,
		"error": apiErr,
		// 详细注释: })
	})
	// 详细注释: }
}

// IdempotencyStore 保存幂等请求响应。
// 详细注释: type IdempotencyStore struct {
type IdempotencyStore struct {
	// 详细注释: mu   sync.Mutex
	mu sync.Mutex
	// 详细注释: rows map[string]recordedResponse
	rows map[string]recordedResponse
	// 详细注释: }
}

// 详细注释: type recordedResponse struct {
type recordedResponse struct {
	// 详细注释: status int
	status int
	// 详细注释: body   []byte
	body []byte
	// 详细注释: }
}

// NewIdempotencyStore 创建幂等存储。
// 详细注释: func NewIdempotencyStore() *IdempotencyStore {
func NewIdempotencyStore() *IdempotencyStore {
	// 详细注释: return &IdempotencyStore{rows: make(map[string]recordedResponse)}
	return &IdempotencyStore{rows: make(map[string]recordedResponse)}
	// 详细注释: }
}

// IdempotencyMiddleware 对 POST + Idempotency-Key 启用幂等回放。
// 详细注释: func IdempotencyMiddleware(store *IdempotencyStore, next http.Handler) http.Handler {
func IdempotencyMiddleware(store *IdempotencyStore, next http.Handler) http.Handler {
	// 详细注释: return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: if r.Method != http.MethodPost {
		if r.Method != http.MethodPost {
			// 详细注释: next.ServeHTTP(w, r)
			next.ServeHTTP(w, r)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		// 详细注释: if key == "" {
		if key == "" {
			// 详细注释: next.ServeHTTP(w, r)
			next.ServeHTTP(w, r)
			// 详细注释: return
			return
			// 详细注释: }
		}

		// 详细注释: store.mu.Lock()
		store.mu.Lock()
		// 详细注释: recorded, ok := store.rows[key]
		recorded, ok := store.rows[key]
		// 详细注释: store.mu.Unlock()
		store.mu.Unlock()
		// 详细注释: if ok {
		if ok {
			// 详细注释: w.Header().Set("X-Idempotent-Replay", "1")
			w.Header().Set("X-Idempotent-Replay", "1")
			// 详细注释: w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Type", "application/json")
			// 详细注释: w.WriteHeader(recorded.status)
			w.WriteHeader(recorded.status)
			// 详细注释: _, _ = w.Write(recorded.body)
			_, _ = w.Write(recorded.body)
			// 详细注释: return
			return
			// 详细注释: }
		}

		// 详细注释: rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK, buf: &bytes.Buffer{}}
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK, buf: &bytes.Buffer{}}
		// 详细注释: next.ServeHTTP(rec, r)
		next.ServeHTTP(rec, r)

		// 详细注释: store.mu.Lock()
		store.mu.Lock()
		// 详细注释: store.rows[key] = recordedResponse{status: rec.status, body: rec.buf.Bytes()}
		store.rows[key] = recordedResponse{status: rec.status, body: rec.buf.Bytes()}
		// 详细注释: store.mu.Unlock()
		store.mu.Unlock()
		// 详细注释: })
	})
	// 详细注释: }
}

// RateLimiter 是固定窗口限流器。
// 详细注释: type RateLimiter struct {
type RateLimiter struct {
	// 详细注释: mu      sync.Mutex
	mu sync.Mutex
	// 详细注释: limit   int
	limit int
	// 详细注释: window  time.Duration
	window time.Duration
	// 详细注释: nowFn   func() time.Time
	nowFn func() time.Time
	// 详细注释: counter map[string]windowCounter
	counter map[string]windowCounter
	// 详细注释: }
}

// 详细注释: type windowCounter struct {
type windowCounter struct {
	// 详细注释: count   int
	count int
	// 详细注释: resetAt time.Time
	resetAt time.Time
	// 详细注释: }
}

// NewRateLimiter 创建限流器。
// 详细注释: func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	// 详细注释: return &RateLimiter{limit: limit, window: window, nowFn: time.Now, counter: make(map[string]windowCounter)}
	return &RateLimiter{limit: limit, window: window, nowFn: time.Now, counter: make(map[string]windowCounter)}
	// 详细注释: }
}

// Allow 判断请求是否放行。
// 详细注释: func (l *RateLimiter) Allow(key string) bool {
func (l *RateLimiter) Allow(key string) bool {
	// 详细注释: l.mu.Lock()
	l.mu.Lock()
	// 详细注释: defer l.mu.Unlock()
	defer l.mu.Unlock()
	// 详细注释: now := l.nowFn()
	now := l.nowFn()
	// 详细注释: v, ok := l.counter[key]
	v, ok := l.counter[key]
	// 详细注释: if !ok || now.After(v.resetAt) {
	if !ok || now.After(v.resetAt) {
		// 详细注释: l.counter[key] = windowCounter{count: 1, resetAt: now.Add(l.window)}
		l.counter[key] = windowCounter{count: 1, resetAt: now.Add(l.window)}
		// 详细注释: return true
		return true
		// 详细注释: }
	}
	// 详细注释: if v.count >= l.limit {
	if v.count >= l.limit {
		// 详细注释: return false
		return false
		// 详细注释: }
	}
	// 详细注释: v.count++
	v.count++
	// 详细注释: l.counter[key] = v
	l.counter[key] = v
	// 详细注释: return true
	return true
	// 详细注释: }
}

// RateLimitMiddleware 对请求做限流。
// 详细注释: func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
	// 详细注释: return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: key := strings.TrimSpace(r.Header.Get("X-User-ID"))
		key := strings.TrimSpace(r.Header.Get("X-User-ID"))
		// 详细注释: if key == "" {
		if key == "" {
			// 详细注释: key = "anonymous"
			key = "anonymous"
			// 详细注释: }
		}
		// 详细注释: if !limiter.Allow(key) {
		if !limiter.Allow(key) {
			// 详细注释: WriteError(w, http.StatusTooManyRequests, "req-rate-limit", "RATE_LIMITED", "too many requests")
			WriteError(w, http.StatusTooManyRequests, "req-rate-limit", "RATE_LIMITED", "too many requests")
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: next.ServeHTTP(w, r)
		next.ServeHTTP(w, r)
		// 详细注释: })
	})
	// 详细注释: }
}

// Order 表示最小业务对象。
// 详细注释: type Order struct {
type Order struct {
	// 详细注释: ID   string `json:"id"`
	ID string `json:"id"`
	// 详细注释: Item string `json:"item"`
	Item string `json:"item"`
	// 详细注释: }
}

// OrderApp 负责业务逻辑。
// 详细注释: type OrderApp struct {
type OrderApp struct {
	// 详细注释: mu     sync.Mutex
	mu sync.Mutex
	// 详细注释: nextID int
	nextID int
	// 详细注释: }
}

// NewOrderApp 创建业务对象。
// 详细注释: func NewOrderApp() *OrderApp {
func NewOrderApp() *OrderApp {
	// 详细注释: return &OrderApp{nextID: 1}
	return &OrderApp{nextID: 1}
	// 详细注释: }
}

// CreateOrder 创建订单。
// 详细注释: func (a *OrderApp) CreateOrder(item string) (Order, error) {
func (a *OrderApp) CreateOrder(item string) (Order, error) {
	// 详细注释: item = strings.TrimSpace(item)
	item = strings.TrimSpace(item)
	// 详细注释: if len(item) < 2 {
	if len(item) < 2 {
		// 详细注释: return Order{}, ErrInvalidItem
		return Order{}, ErrInvalidItem
		// 详细注释: }
	}
	// 详细注释: a.mu.Lock()
	a.mu.Lock()
	// 详细注释: defer a.mu.Unlock()
	defer a.mu.Unlock()
	// 详细注释: id := a.nextID
	id := a.nextID
	// 详细注释: a.nextID++
	a.nextID++
	// 详细注释: return Order{ID: "ord-" + strconvItoa(id), Item: item}, nil
	return Order{ID: "ord-" + strconvItoa(id), Item: item}, nil
	// 详细注释: }
}

// NewGovernedMux 创建带治理能力的路由。
// 详细注释: func NewGovernedMux(app *OrderApp, idem *IdempotencyStore, limiter *RateLimiter) *http.ServeMux {
func NewGovernedMux(app *OrderApp, idem *IdempotencyStore, limiter *RateLimiter) *http.ServeMux {
	// 详细注释: base := http.NewServeMux()
	base := http.NewServeMux()
	// 详细注释: base.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, r *http.Request) {
	base.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: if r.Method != http.MethodPost {
		if r.Method != http.MethodPost {
			// 详细注释: WriteError(w, http.StatusMethodNotAllowed, "req-method", "METHOD_NOT_ALLOWED", "method not allowed")
			WriteError(w, http.StatusMethodNotAllowed, "req-method", "METHOD_NOT_ALLOWED", "method not allowed")
			// 详细注释: return
			return
			// 详细注释: }
		}

		// 详细注释: var req struct {
		var req struct {
			// 详细注释: Item string `json:"item"`
			Item string `json:"item"`
			// 详细注释: }
		}
		// 详细注释: if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// 详细注释: WriteError(w, http.StatusBadRequest, "req-json", "INVALID_JSON", "invalid json")
			WriteError(w, http.StatusBadRequest, "req-json", "INVALID_JSON", "invalid json")
			// 详细注释: return
			return
			// 详细注释: }
		}

		// 详细注释: order, err := app.CreateOrder(req.Item)
		order, err := app.CreateOrder(req.Item)
		// 详细注释: if err != nil {
		if err != nil {
			// 详细注释: if errors.Is(err, ErrInvalidItem) {
			if errors.Is(err, ErrInvalidItem) {
				// 详细注释: WriteError(w, http.StatusBadRequest, "req-item", "INVALID_ITEM", err.Error())
				WriteError(w, http.StatusBadRequest, "req-item", "INVALID_ITEM", err.Error())
				// 详细注释: return
				return
				// 详细注释: }
			}
			// 详细注释: WriteError(w, http.StatusInternalServerError, "req-internal", "INTERNAL_ERROR", "internal server error")
			WriteError(w, http.StatusInternalServerError, "req-internal", "INTERNAL_ERROR", "internal server error")
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: WriteSuccess(w, http.StatusCreated, "req-create", order)
		WriteSuccess(w, http.StatusCreated, "req-create", order)
		// 详细注释: })
	})

	// 详细注释: wrapped := IdempotencyMiddleware(idem, RateLimitMiddleware(limiter, base))
	wrapped := IdempotencyMiddleware(idem, RateLimitMiddleware(limiter, base))
	// 详细注释: root := http.NewServeMux()
	root := http.NewServeMux()
	// 详细注释: root.Handle("/", wrapped)
	root.Handle("/", wrapped)
	// 详细注释: return root
	return root
	// 详细注释: }
}

// ShutdownWithTimeout 统一优雅停机入口。
// 详细注释: func ShutdownWithTimeout(server *http.Server, timeout time.Duration) error {
func ShutdownWithTimeout(server *http.Server, timeout time.Duration) error {
	// 详细注释: ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// 详细注释: defer cancel()
	defer cancel()
	// 详细注释: return server.Shutdown(ctx)
	return server.Shutdown(ctx)
	// 详细注释: }
}

// 详细注释: type responseRecorder struct {
type responseRecorder struct {
	// 详细注释: http.ResponseWriter
	http.ResponseWriter
	// 详细注释: status int
	status int
	// 详细注释: buf    *bytes.Buffer
	buf *bytes.Buffer
	// 详细注释: }
}

// 详细注释: func (r *responseRecorder) WriteHeader(statusCode int) {
func (r *responseRecorder) WriteHeader(statusCode int) {
	// 详细注释: r.status = statusCode
	r.status = statusCode
	// 详细注释: r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseWriter.WriteHeader(statusCode)
	// 详细注释: }
}

// 详细注释: func (r *responseRecorder) Write(b []byte) (int, error) {
func (r *responseRecorder) Write(b []byte) (int, error) {
	// 详细注释: _, _ = r.buf.Write(b)
	_, _ = r.buf.Write(b)
	// 详细注释: return r.ResponseWriter.Write(b)
	return r.ResponseWriter.Write(b)
	// 详细注释: }
}

// 详细注释: func strconvItoa(v int) string {
func strconvItoa(v int) string {
	// 详细注释: if v == 0 {
	if v == 0 {
		// 详细注释: return "0"
		return "0"
		// 详细注释: }
	}
	// 详细注释: buf := make([]byte, 0, 10)
	buf := make([]byte, 0, 10)
	// 详细注释: for v > 0 {
	for v > 0 {
		// 详细注释: buf = append([]byte{byte('0' + v%10)}, buf...)
		buf = append([]byte{byte('0' + v%10)}, buf...)
		// 详细注释: v /= 10
		v /= 10
		// 详细注释: }
	}
	// 详细注释: return string(buf)
	return string(buf)
	// 详细注释: }
}
