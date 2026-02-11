// 详细注释: package week14
package week14

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "sync"
	"sync"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// Span 表示一次最小链路追踪记录。
// 详细注释: type Span struct {
type Span struct {
	// 详细注释: TraceID   string
	TraceID string
	// 详细注释: Name      string
	Name string
	// 详细注释: StartTime time.Time
	StartTime time.Time
	// 详细注释: EndTime   time.Time
	EndTime time.Time
	// 详细注释: Error     string
	Error string
	// 详细注释: }
}

// Tracer 记录内存中的 span。
// 详细注释: type Tracer struct {
type Tracer struct {
	// 详细注释: mu     sync.Mutex
	mu sync.Mutex
	// 详细注释: nextID int
	nextID int
	// 详细注释: nowFn  func() time.Time
	nowFn func() time.Time
	// 详细注释: spans  []Span
	spans []Span
	// 详细注释: }
}

// NewTracer 创建 tracer。
// 详细注释: func NewTracer() *Tracer {
func NewTracer() *Tracer {
	// 详细注释: return &Tracer{nextID: 1, nowFn: time.Now}
	return &Tracer{nextID: 1, nowFn: time.Now}
	// 详细注释: }
}

// StartSpan 开始一个 span，并返回结束函数。
// 详细注释: func (t *Tracer) StartSpan(name string) func(err error) {
func (t *Tracer) StartSpan(name string) func(err error) {
	// 详细注释: t.mu.Lock()
	t.mu.Lock()
	// 详细注释: id := fmt.Sprintf("trace-%d", t.nextID)
	id := fmt.Sprintf("trace-%d", t.nextID)
	// 详细注释: t.nextID++
	t.nextID++
	// 详细注释: idx := len(t.spans)
	idx := len(t.spans)
	// 详细注释: t.spans = append(t.spans, Span{TraceID: id, Name: name, StartTime: t.nowFn()})
	t.spans = append(t.spans, Span{TraceID: id, Name: name, StartTime: t.nowFn()})
	// 详细注释: t.mu.Unlock()
	t.mu.Unlock()

	// 详细注释: return func(err error) {
	return func(err error) {
		// 详细注释: t.mu.Lock()
		t.mu.Lock()
		// 详细注释: defer t.mu.Unlock()
		defer t.mu.Unlock()
		// 详细注释: span := t.spans[idx]
		span := t.spans[idx]
		// 详细注释: span.EndTime = t.nowFn()
		span.EndTime = t.nowFn()
		// 详细注释: if err != nil {
		if err != nil {
			// 详细注释: span.Error = err.Error()
			span.Error = err.Error()
			// 详细注释: }
		}
		// 详细注释: t.spans[idx] = span
		t.spans[idx] = span
		// 详细注释: }
	}
	// 详细注释: }
}

// Spans 返回快照。
// 详细注释: func (t *Tracer) Spans() []Span {
func (t *Tracer) Spans() []Span {
	// 详细注释: t.mu.Lock()
	t.mu.Lock()
	// 详细注释: defer t.mu.Unlock()
	defer t.mu.Unlock()
	// 详细注释: out := make([]Span, len(t.spans))
	out := make([]Span, len(t.spans))
	// 详细注释: copy(out, t.spans)
	copy(out, t.spans)
	// 详细注释: return out
	return out
	// 详细注释: }
}

// AuditEvent 表示审计日志事件。
// 详细注释: type AuditEvent struct {
type AuditEvent struct {
	// 详细注释: Actor     string
	Actor string
	// 详细注释: Action    string
	Action string
	// 详细注释: Resource  string
	Resource string
	// 详细注释: Result    string
	Result string
	// 详细注释: Timestamp time.Time
	Timestamp time.Time
	// 详细注释: }
}

// AuditLogger 保存审计日志。
// 详细注释: type AuditLogger struct {
type AuditLogger struct {
	// 详细注释: mu     sync.Mutex
	mu sync.Mutex
	// 详细注释: nowFn  func() time.Time
	nowFn func() time.Time
	// 详细注释: events []AuditEvent
	events []AuditEvent
	// 详细注释: }
}

// NewAuditLogger 创建审计日志器。
// 详细注释: func NewAuditLogger() *AuditLogger {
func NewAuditLogger() *AuditLogger {
	// 详细注释: return &AuditLogger{nowFn: time.Now}
	return &AuditLogger{nowFn: time.Now}
	// 详细注释: }
}

// Log 记录审计事件。
// 详细注释: func (l *AuditLogger) Log(actor, action, resource, result string) {
func (l *AuditLogger) Log(actor, action, resource, result string) {
	// 详细注释: l.mu.Lock()
	l.mu.Lock()
	// 详细注释: defer l.mu.Unlock()
	defer l.mu.Unlock()
	// 详细注释: l.events = append(l.events, AuditEvent{
	l.events = append(l.events, AuditEvent{
		// 详细注释: Actor: actor, Action: action, Resource: resource, Result: result, Timestamp: l.nowFn(),
		Actor: actor, Action: action, Resource: resource, Result: result, Timestamp: l.nowFn(),
		// 详细注释: })
	})
	// 详细注释: }
}

// Events 返回审计日志快照。
// 详细注释: func (l *AuditLogger) Events() []AuditEvent {
func (l *AuditLogger) Events() []AuditEvent {
	// 详细注释: l.mu.Lock()
	l.mu.Lock()
	// 详细注释: defer l.mu.Unlock()
	defer l.mu.Unlock()
	// 详细注释: out := make([]AuditEvent, len(l.events))
	out := make([]AuditEvent, len(l.events))
	// 详细注释: copy(out, l.events)
	copy(out, l.events)
	// 详细注释: return out
	return out
	// 详细注释: }
}

// Job 表示异步任务。
// 详细注释: type Job struct {
type Job struct {
	// 详细注释: ID      string
	ID string
	// 详细注释: Payload string
	Payload string
	// 详细注释: }
}

// JobResult 表示任务处理结果。
// 详细注释: type JobResult struct {
type JobResult struct {
	// 详细注释: JobID    string
	JobID string
	// 详细注释: Success  bool
	Success bool
	// 详细注释: Attempts int
	Attempts int
	// 详细注释: FinalErr string
	FinalErr string
	// 详细注释: }
}

// JobHandler 定义任务执行函数。
// 详细注释: type JobHandler func(ctx context.Context, job Job, attempt int) error
type JobHandler func(ctx context.Context, job Job, attempt int) error

// AsyncWorker 处理任务并支持重试。
// 详细注释: type AsyncWorker struct {
type AsyncWorker struct {
	// 详细注释: MaxRetries int
	MaxRetries int
	// 详细注释: Handler    JobHandler
	Handler JobHandler

	// 详细注释: mu          sync.Mutex
	mu sync.Mutex
	// 详细注释: deadLetters []JobResult
	deadLetters []JobResult
	// 详细注释: }
}

// NewAsyncWorker 创建 worker。
// 详细注释: func NewAsyncWorker(maxRetries int, handler JobHandler) *AsyncWorker {
func NewAsyncWorker(maxRetries int, handler JobHandler) *AsyncWorker {
	// 详细注释: return &AsyncWorker{MaxRetries: maxRetries, Handler: handler}
	return &AsyncWorker{MaxRetries: maxRetries, Handler: handler}
	// 详细注释: }
}

// ProcessBatch 处理一批任务，失败任务进入 dead letter。
// 详细注释: func (w *AsyncWorker) ProcessBatch(ctx context.Context, tracer *Tracer, audit *AuditLogger, jobs []Job) []JobResult {
func (w *AsyncWorker) ProcessBatch(ctx context.Context, tracer *Tracer, audit *AuditLogger, jobs []Job) []JobResult {
	// 详细注释: results := make([]JobResult, 0, len(jobs))
	results := make([]JobResult, 0, len(jobs))
	// 详细注释: for _, job := range jobs {
	for _, job := range jobs {
		// 详细注释: end := tracer.StartSpan("job.process." + job.ID)
		end := tracer.StartSpan("job.process." + job.ID)
		// 详细注释: result := w.processOne(ctx, job)
		result := w.processOne(ctx, job)
		// 详细注释: if result.Success {
		if result.Success {
			// 详细注释: audit.Log("worker", "job_process", job.ID, "success")
			audit.Log("worker", "job_process", job.ID, "success")
			// 详细注释: end(nil)
			end(nil)
			// 详细注释: } else {
		} else {
			// 详细注释: audit.Log("worker", "job_process", job.ID, "failed")
			audit.Log("worker", "job_process", job.ID, "failed")
			// 详细注释: end(fmt.Errorf(result.FinalErr))
			end(fmt.Errorf(result.FinalErr))
			// 详细注释: }
		}
		// 详细注释: results = append(results, result)
		results = append(results, result)
		// 详细注释: }
	}
	// 详细注释: return results
	return results
	// 详细注释: }
}

// 详细注释: func (w *AsyncWorker) processOne(ctx context.Context, job Job) JobResult {
func (w *AsyncWorker) processOne(ctx context.Context, job Job) JobResult {
	// 详细注释: maxAttempt := w.MaxRetries + 1
	maxAttempt := w.MaxRetries + 1
	// 详细注释: for attempt := 1; attempt <= maxAttempt; attempt++ {
	for attempt := 1; attempt <= maxAttempt; attempt++ {
		// 详细注释: if err := w.Handler(ctx, job, attempt); err == nil {
		if err := w.Handler(ctx, job, attempt); err == nil {
			// 详细注释: return JobResult{JobID: job.ID, Success: true, Attempts: attempt}
			return JobResult{JobID: job.ID, Success: true, Attempts: attempt}
			// 详细注释: } else if attempt == maxAttempt {
		} else if attempt == maxAttempt {
			// 详细注释: result := JobResult{JobID: job.ID, Success: false, Attempts: attempt, FinalErr: err.Error()}
			result := JobResult{JobID: job.ID, Success: false, Attempts: attempt, FinalErr: err.Error()}
			// 详细注释: w.mu.Lock()
			w.mu.Lock()
			// 详细注释: w.deadLetters = append(w.deadLetters, result)
			w.deadLetters = append(w.deadLetters, result)
			// 详细注释: w.mu.Unlock()
			w.mu.Unlock()
			// 详细注释: return result
			return result
			// 详细注释: }
		}
		// 详细注释: }
	}
	// 详细注释: return JobResult{JobID: job.ID, Success: false, Attempts: maxAttempt, FinalErr: "unexpected"}
	return JobResult{JobID: job.ID, Success: false, Attempts: maxAttempt, FinalErr: "unexpected"}
	// 详细注释: }
}

// DeadLetters 返回死信快照。
// 详细注释: func (w *AsyncWorker) DeadLetters() []JobResult {
func (w *AsyncWorker) DeadLetters() []JobResult {
	// 详细注释: w.mu.Lock()
	w.mu.Lock()
	// 详细注释: defer w.mu.Unlock()
	defer w.mu.Unlock()
	// 详细注释: out := make([]JobResult, len(w.deadLetters))
	out := make([]JobResult, len(w.deadLetters))
	// 详细注释: copy(out, w.deadLetters)
	copy(out, w.deadLetters)
	// 详细注释: return out
	return out
	// 详细注释: }
}
