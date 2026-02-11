package observability

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics 保存简化版 Prometheus 指标。
type Metrics struct {
	requestsTotal   atomic.Int64
	requestDuration map[string]*durationBucket
	mu              sync.RWMutex
}

// durationBucket 聚合一条 route+status 维度的耗时统计。
type durationBucket struct {
	count int64
	sum   float64
}

// NewMetrics 创建指标容器。
func NewMetrics() *Metrics {
	return &Metrics{requestDuration: map[string]*durationBucket{}}
}

// ObserveRequest 记录请求数量和耗时。
func (m *Metrics) ObserveRequest(method, path string, status int, d time.Duration) {
	// 总请求数使用原子计数，降低锁竞争。
	m.requestsTotal.Add(1)
	key := fmt.Sprintf("%s|%s|%d", strings.ToUpper(method), path, status)
	m.mu.Lock()
	bucket, ok := m.requestDuration[key]
	if !ok {
		bucket = &durationBucket{}
		m.requestDuration[key] = bucket
	}
	bucket.count++
	bucket.sum += d.Seconds()
	m.mu.Unlock()
}

// Handler 返回 Prometheus 文本格式指标输出。
func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	_, _ = fmt.Fprintf(w, "# HELP app_http_requests_total Total number of HTTP requests.\n")
	_, _ = fmt.Fprintf(w, "# TYPE app_http_requests_total counter\n")
	_, _ = fmt.Fprintf(w, "app_http_requests_total %d\n", m.requestsTotal.Load())

	_, _ = fmt.Fprintf(w, "# HELP app_http_request_duration_seconds_sum Request duration seconds sum by route.\n")
	_, _ = fmt.Fprintf(w, "# TYPE app_http_request_duration_seconds_sum counter\n")
	_, _ = fmt.Fprintf(w, "# HELP app_http_request_duration_seconds_count Request duration count by route.\n")
	_, _ = fmt.Fprintf(w, "# TYPE app_http_request_duration_seconds_count counter\n")

	m.mu.RLock()
	// 先排序再输出，保证抓取结果稳定，便于对比与测试。
	keys := make([]string, 0, len(m.requestDuration))
	for key := range m.requestDuration {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		bucket := m.requestDuration[key]
		parts := strings.Split(key, "|")
		if len(parts) != 3 {
			continue
		}
		method := parts[0]
		path := parts[1]
		status := parts[2]
		_, _ = fmt.Fprintf(w, "app_http_request_duration_seconds_sum{method=\"%s\",path=\"%s\",status=\"%s\"} %.6f\n", method, path, status, bucket.sum)
		_, _ = fmt.Fprintf(w, "app_http_request_duration_seconds_count{method=\"%s\",path=\"%s\",status=\"%s\"} %d\n", method, path, status, bucket.count)
	}
	m.mu.RUnlock()
}
