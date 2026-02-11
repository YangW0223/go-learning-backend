// 详细注释: package week09
package week09

// 详细注释: import (
import (
	// 详细注释: "encoding/json"
	"encoding/json"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "io"
	"io"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "strconv"
	"strconv"
	// 详细注释: "sync"
	"sync"
	// 详细注释: "sync/atomic"
	"sync/atomic"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// JSONLogger 提供最小结构化日志能力。
// 详细注释: type JSONLogger struct {
type JSONLogger struct {
	// 详细注释: mu sync.Mutex
	mu sync.Mutex
	// 详细注释: w  io.Writer
	w io.Writer
	// 详细注释: }
}

// NewJSONLogger 创建日志实例。
// 详细注释: func NewJSONLogger(w io.Writer) *JSONLogger {
func NewJSONLogger(w io.Writer) *JSONLogger {
	// 详细注释: return &JSONLogger{w: w}
	return &JSONLogger{w: w}
	// 详细注释: }
}

// Log 以 JSON 行格式输出日志。
// 详细注释: func (l *JSONLogger) Log(fields map[string]any) {
func (l *JSONLogger) Log(fields map[string]any) {
	// 详细注释: l.mu.Lock()
	l.mu.Lock()
	// 详细注释: defer l.mu.Unlock()
	defer l.mu.Unlock()
	// 详细注释: b, _ := json.Marshal(fields)
	b, _ := json.Marshal(fields)
	// 详细注释: _, _ = l.w.Write(append(b, '\n'))
	_, _ = l.w.Write(append(b, '\n'))
	// 详细注释: }
}

// Metrics 记录请求总数、错误数和耗时分布。
// 详细注释: type Metrics struct {
type Metrics struct {
	// 详细注释: mu            sync.Mutex
	mu sync.Mutex
	// 详细注释: RequestCount  int
	RequestCount int
	// 详细注释: ErrorCount    int
	ErrorCount int
	// 详细注释: LatencyBucket map[string]int
	LatencyBucket map[string]int
	// 详细注释: }
}

// NewMetrics 创建指标实例。
// 详细注释: func NewMetrics() *Metrics {
func NewMetrics() *Metrics {
	// 详细注释: return &Metrics{LatencyBucket: map[string]int{"lt10ms": 0, "lt100ms": 0, "gte100ms": 0}}
	return &Metrics{LatencyBucket: map[string]int{"lt10ms": 0, "lt100ms": 0, "gte100ms": 0}}
	// 详细注释: }
}

// Observe 写入一次请求观测数据。
// 详细注释: func (m *Metrics) Observe(status int, latency time.Duration) {
func (m *Metrics) Observe(status int, latency time.Duration) {
	// 详细注释: m.mu.Lock()
	m.mu.Lock()
	// 详细注释: defer m.mu.Unlock()
	defer m.mu.Unlock()
	// 详细注释: m.RequestCount++
	m.RequestCount++
	// 详细注释: if status >= 500 {
	if status >= 500 {
		// 详细注释: m.ErrorCount++
		m.ErrorCount++
		// 详细注释: }
	}
	// 详细注释: if latency < 10*time.Millisecond {
	if latency < 10*time.Millisecond {
		// 详细注释: m.LatencyBucket["lt10ms"]++
		m.LatencyBucket["lt10ms"]++
		// 详细注释: } else if latency < 100*time.Millisecond {
	} else if latency < 100*time.Millisecond {
		// 详细注释: m.LatencyBucket["lt100ms"]++
		m.LatencyBucket["lt100ms"]++
		// 详细注释: } else {
	} else {
		// 详细注释: m.LatencyBucket["gte100ms"]++
		m.LatencyBucket["gte100ms"]++
		// 详细注释: }
	}
	// 详细注释: }
}

// 详细注释: var requestCounter uint64
var requestCounter uint64

// WithObservability 包装 handler，增加 request_id、日志和指标。
// 详细注释: func WithObservability(logger *JSONLogger, metrics *Metrics, next http.Handler) http.Handler {
func WithObservability(logger *JSONLogger, metrics *Metrics, next http.Handler) http.Handler {
	// 详细注释: return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: requestID := nextRequestID()
		requestID := nextRequestID()
		// 详细注释: start := time.Now()
		start := time.Now()
		// 详细注释: rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		// 详细注释: r.Header.Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)
		// 详细注释: next.ServeHTTP(rec, r)
		next.ServeHTTP(rec, r)

		// 详细注释: latency := time.Since(start)
		latency := time.Since(start)
		// 详细注释: metrics.Observe(rec.status, latency)
		metrics.Observe(rec.status, latency)
		// 详细注释: logger.Log(map[string]any{
		logger.Log(map[string]any{
			// 详细注释: "request_id": requestID,
			"request_id": requestID,
			// 详细注释: "method":     r.Method,
			"method": r.Method,
			// 详细注释: "path":       r.URL.Path,
			"path": r.URL.Path,
			// 详细注释: "status":     rec.status,
			"status": rec.status,
			// 详细注释: "latency_ms": latency.Milliseconds(),
			"latency_ms": latency.Milliseconds(),
			// 详细注释: })
		})
		// 详细注释: })
	})
	// 详细注释: }
}

// NewDemoHandler 返回可制造 500 的示例 handler。
// 详细注释: func NewDemoHandler() http.Handler {
func NewDemoHandler() http.Handler {
	// 详细注释: return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: if r.URL.Query().Get("fail") == "1" {
		if r.URL.Query().Get("fail") == "1" {
			// 详细注释: http.Error(w, "simulated internal error", http.StatusInternalServerError)
			http.Error(w, "simulated internal error", http.StatusInternalServerError)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: _, _ = w.Write([]byte("ok"))
		_, _ = w.Write([]byte("ok"))
		// 详细注释: })
	})
	// 详细注释: }
}

// 详细注释: func nextRequestID() string {
func nextRequestID() string {
	// 详细注释: id := atomic.AddUint64(&requestCounter, 1)
	id := atomic.AddUint64(&requestCounter, 1)
	// 详细注释: return "req-" + strconv.FormatUint(id, 10)
	return "req-" + strconv.FormatUint(id, 10)
	// 详细注释: }
}

// 详细注释: type statusRecorder struct {
type statusRecorder struct {
	// 详细注释: http.ResponseWriter
	http.ResponseWriter
	// 详细注释: status int
	status int
	// 详细注释: }
}

// 详细注释: func (r *statusRecorder) WriteHeader(statusCode int) {
func (r *statusRecorder) WriteHeader(statusCode int) {
	// 详细注释: r.status = statusCode
	r.status = statusCode
	// 详细注释: r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseWriter.WriteHeader(statusCode)
	// 详细注释: }
}

// Snapshot 返回可打印的指标快照。
// 详细注释: func (m *Metrics) Snapshot() string {
func (m *Metrics) Snapshot() string {
	// 详细注释: m.mu.Lock()
	m.mu.Lock()
	// 详细注释: defer m.mu.Unlock()
	defer m.mu.Unlock()
	// 详细注释: return fmt.Sprintf("requests=%d errors=%d buckets=%v", m.RequestCount, m.ErrorCount, m.LatencyBucket)
	return fmt.Sprintf("requests=%d errors=%d buckets=%v", m.RequestCount, m.ErrorCount, m.LatencyBucket)
	// 详细注释: }
}
