package week09

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// JSONLogger 提供最小结构化日志能力。
type JSONLogger struct {
	mu sync.Mutex
	w  io.Writer
}

// NewJSONLogger 创建日志实例。
func NewJSONLogger(w io.Writer) *JSONLogger {
	return &JSONLogger{w: w}
}

// Log 以 JSON 行格式输出日志。
func (l *JSONLogger) Log(fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, _ := json.Marshal(fields)
	_, _ = l.w.Write(append(b, '\n'))
}

// Metrics 记录请求总数、错误数和耗时分布。
type Metrics struct {
	mu            sync.Mutex
	RequestCount  int
	ErrorCount    int
	LatencyBucket map[string]int
}

// NewMetrics 创建指标实例。
func NewMetrics() *Metrics {
	return &Metrics{LatencyBucket: map[string]int{"lt10ms": 0, "lt100ms": 0, "gte100ms": 0}}
}

// Observe 写入一次请求观测数据。
func (m *Metrics) Observe(status int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestCount++
	if status >= 500 {
		m.ErrorCount++
	}
	if latency < 10*time.Millisecond {
		m.LatencyBucket["lt10ms"]++
	} else if latency < 100*time.Millisecond {
		m.LatencyBucket["lt100ms"]++
	} else {
		m.LatencyBucket["gte100ms"]++
	}
}

var requestCounter uint64

// WithObservability 包装 handler，增加 request_id、日志和指标。
func WithObservability(logger *JSONLogger, metrics *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := nextRequestID()
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		r.Header.Set("X-Request-ID", requestID)
		next.ServeHTTP(rec, r)

		latency := time.Since(start)
		metrics.Observe(rec.status, latency)
		logger.Log(map[string]any{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     rec.status,
			"latency_ms": latency.Milliseconds(),
		})
	})
}

// NewDemoHandler 返回可制造 500 的示例 handler。
func NewDemoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fail") == "1" {
			http.Error(w, "simulated internal error", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
}

func nextRequestID() string {
	id := atomic.AddUint64(&requestCounter, 1)
	return "req-" + strconv.FormatUint(id, 10)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Snapshot 返回可打印的指标快照。
func (m *Metrics) Snapshot() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return fmt.Sprintf("requests=%d errors=%d buckets=%v", m.RequestCount, m.ErrorCount, m.LatencyBucket)
}
