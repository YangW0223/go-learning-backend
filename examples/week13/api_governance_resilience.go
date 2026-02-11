package week13

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	// ErrInvalidItem 表示下单参数非法。
	ErrInvalidItem = errors.New("invalid item")
)

// APIError 表示统一错误结构。
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteSuccess 输出统一成功响应。
func WriteSuccess(w http.ResponseWriter, status int, requestID string, data any) {
	writeEnvelope(w, status, requestID, data, nil)
}

// WriteError 输出统一错误响应。
func WriteError(w http.ResponseWriter, status int, requestID string, code string, message string) {
	writeEnvelope(w, status, requestID, nil, &APIError{Code: code, Message: message})
}

func writeEnvelope(w http.ResponseWriter, status int, requestID string, data any, apiErr *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"request_id": requestID,
		"data":       data,
		"error":      apiErr,
	})
}

// IdempotencyStore 保存幂等请求响应。
type IdempotencyStore struct {
	mu   sync.Mutex
	rows map[string]recordedResponse
}

type recordedResponse struct {
	status int
	body   []byte
}

// NewIdempotencyStore 创建幂等存储。
func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{rows: make(map[string]recordedResponse)}
}

// IdempotencyMiddleware 对 POST + Idempotency-Key 启用幂等回放。
func IdempotencyMiddleware(store *IdempotencyStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		store.mu.Lock()
		recorded, ok := store.rows[key]
		store.mu.Unlock()
		if ok {
			w.Header().Set("X-Idempotent-Replay", "1")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(recorded.status)
			_, _ = w.Write(recorded.body)
			return
		}

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK, buf: &bytes.Buffer{}}
		next.ServeHTTP(rec, r)

		store.mu.Lock()
		store.rows[key] = recordedResponse{status: rec.status, body: rec.buf.Bytes()}
		store.mu.Unlock()
	})
}

// RateLimiter 是固定窗口限流器。
type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	nowFn   func() time.Time
	counter map[string]windowCounter
}

type windowCounter struct {
	count   int
	resetAt time.Time
}

// NewRateLimiter 创建限流器。
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{limit: limit, window: window, nowFn: time.Now, counter: make(map[string]windowCounter)}
}

// Allow 判断请求是否放行。
func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.nowFn()
	v, ok := l.counter[key]
	if !ok || now.After(v.resetAt) {
		l.counter[key] = windowCounter{count: 1, resetAt: now.Add(l.window)}
		return true
	}
	if v.count >= l.limit {
		return false
	}
	v.count++
	l.counter[key] = v
	return true
}

// RateLimitMiddleware 对请求做限流。
func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimSpace(r.Header.Get("X-User-ID"))
		if key == "" {
			key = "anonymous"
		}
		if !limiter.Allow(key) {
			WriteError(w, http.StatusTooManyRequests, "req-rate-limit", "RATE_LIMITED", "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Order 表示最小业务对象。
type Order struct {
	ID   string `json:"id"`
	Item string `json:"item"`
}

// OrderApp 负责业务逻辑。
type OrderApp struct {
	mu     sync.Mutex
	nextID int
}

// NewOrderApp 创建业务对象。
func NewOrderApp() *OrderApp {
	return &OrderApp{nextID: 1}
}

// CreateOrder 创建订单。
func (a *OrderApp) CreateOrder(item string) (Order, error) {
	item = strings.TrimSpace(item)
	if len(item) < 2 {
		return Order{}, ErrInvalidItem
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	id := a.nextID
	a.nextID++
	return Order{ID: "ord-" + strconvItoa(id), Item: item}, nil
}

// NewGovernedMux 创建带治理能力的路由。
func NewGovernedMux(app *OrderApp, idem *IdempotencyStore, limiter *RateLimiter) *http.ServeMux {
	base := http.NewServeMux()
	base.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, http.StatusMethodNotAllowed, "req-method", "METHOD_NOT_ALLOWED", "method not allowed")
			return
		}

		var req struct {
			Item string `json:"item"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "req-json", "INVALID_JSON", "invalid json")
			return
		}

		order, err := app.CreateOrder(req.Item)
		if err != nil {
			if errors.Is(err, ErrInvalidItem) {
				WriteError(w, http.StatusBadRequest, "req-item", "INVALID_ITEM", err.Error())
				return
			}
			WriteError(w, http.StatusInternalServerError, "req-internal", "INTERNAL_ERROR", "internal server error")
			return
		}
		WriteSuccess(w, http.StatusCreated, "req-create", order)
	})

	wrapped := IdempotencyMiddleware(idem, RateLimitMiddleware(limiter, base))
	root := http.NewServeMux()
	root.Handle("/", wrapped)
	return root
}

// ShutdownWithTimeout 统一优雅停机入口。
func ShutdownWithTimeout(server *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return server.Shutdown(ctx)
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	buf    *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	_, _ = r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

func strconvItoa(v int) string {
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for v > 0 {
		buf = append([]byte{byte('0' + v%10)}, buf...)
		v /= 10
	}
	return string(buf)
}
