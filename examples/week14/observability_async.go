package week14

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Span 表示一次最小链路追踪记录。
type Span struct {
	TraceID   string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Error     string
}

// Tracer 记录内存中的 span。
type Tracer struct {
	mu     sync.Mutex
	nextID int
	nowFn  func() time.Time
	spans  []Span
}

// NewTracer 创建 tracer。
func NewTracer() *Tracer {
	return &Tracer{nextID: 1, nowFn: time.Now}
}

// StartSpan 开始一个 span，并返回结束函数。
func (t *Tracer) StartSpan(name string) func(err error) {
	t.mu.Lock()
	id := fmt.Sprintf("trace-%d", t.nextID)
	t.nextID++
	idx := len(t.spans)
	t.spans = append(t.spans, Span{TraceID: id, Name: name, StartTime: t.nowFn()})
	t.mu.Unlock()

	return func(err error) {
		t.mu.Lock()
		defer t.mu.Unlock()
		span := t.spans[idx]
		span.EndTime = t.nowFn()
		if err != nil {
			span.Error = err.Error()
		}
		t.spans[idx] = span
	}
}

// Spans 返回快照。
func (t *Tracer) Spans() []Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Span, len(t.spans))
	copy(out, t.spans)
	return out
}

// AuditEvent 表示审计日志事件。
type AuditEvent struct {
	Actor     string
	Action    string
	Resource  string
	Result    string
	Timestamp time.Time
}

// AuditLogger 保存审计日志。
type AuditLogger struct {
	mu     sync.Mutex
	nowFn  func() time.Time
	events []AuditEvent
}

// NewAuditLogger 创建审计日志器。
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{nowFn: time.Now}
}

// Log 记录审计事件。
func (l *AuditLogger) Log(actor, action, resource, result string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, AuditEvent{
		Actor: actor, Action: action, Resource: resource, Result: result, Timestamp: l.nowFn(),
	})
}

// Events 返回审计日志快照。
func (l *AuditLogger) Events() []AuditEvent {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]AuditEvent, len(l.events))
	copy(out, l.events)
	return out
}

// Job 表示异步任务。
type Job struct {
	ID      string
	Payload string
}

// JobResult 表示任务处理结果。
type JobResult struct {
	JobID    string
	Success  bool
	Attempts int
	FinalErr string
}

// JobHandler 定义任务执行函数。
type JobHandler func(ctx context.Context, job Job, attempt int) error

// AsyncWorker 处理任务并支持重试。
type AsyncWorker struct {
	MaxRetries int
	Handler    JobHandler

	mu          sync.Mutex
	deadLetters []JobResult
}

// NewAsyncWorker 创建 worker。
func NewAsyncWorker(maxRetries int, handler JobHandler) *AsyncWorker {
	return &AsyncWorker{MaxRetries: maxRetries, Handler: handler}
}

// ProcessBatch 处理一批任务，失败任务进入 dead letter。
func (w *AsyncWorker) ProcessBatch(ctx context.Context, tracer *Tracer, audit *AuditLogger, jobs []Job) []JobResult {
	results := make([]JobResult, 0, len(jobs))
	for _, job := range jobs {
		end := tracer.StartSpan("job.process." + job.ID)
		result := w.processOne(ctx, job)
		if result.Success {
			audit.Log("worker", "job_process", job.ID, "success")
			end(nil)
		} else {
			audit.Log("worker", "job_process", job.ID, "failed")
			end(fmt.Errorf(result.FinalErr))
		}
		results = append(results, result)
	}
	return results
}

func (w *AsyncWorker) processOne(ctx context.Context, job Job) JobResult {
	maxAttempt := w.MaxRetries + 1
	for attempt := 1; attempt <= maxAttempt; attempt++ {
		if err := w.Handler(ctx, job, attempt); err == nil {
			return JobResult{JobID: job.ID, Success: true, Attempts: attempt}
		} else if attempt == maxAttempt {
			result := JobResult{JobID: job.ID, Success: false, Attempts: attempt, FinalErr: err.Error()}
			w.mu.Lock()
			w.deadLetters = append(w.deadLetters, result)
			w.mu.Unlock()
			return result
		}
	}
	return JobResult{JobID: job.ID, Success: false, Attempts: maxAttempt, FinalErr: "unexpected"}
}

// DeadLetters 返回死信快照。
func (w *AsyncWorker) DeadLetters() []JobResult {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]JobResult, len(w.deadLetters))
	copy(out, w.deadLetters)
	return out
}
