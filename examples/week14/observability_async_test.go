// 详细注释: package week14
package week14

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "testing"
	"testing"
	// 详细注释: )
)

// TestProcessBatchRetriesAndDeadLetter 验证重试和死信行为。
// 详细注释: func TestProcessBatchRetriesAndDeadLetter(t *testing.T) {
func TestProcessBatchRetriesAndDeadLetter(t *testing.T) {
	// 详细注释: attemptByJob := map[string]int{}
	attemptByJob := map[string]int{}
	// 详细注释: handler := func(_ context.Context, job Job, _ int) error {
	handler := func(_ context.Context, job Job, _ int) error {
		// 详细注释: attemptByJob[job.ID]++
		attemptByJob[job.ID]++
		// 详细注释: if job.ID == "job-success-after-retry" && attemptByJob[job.ID] < 2 {
		if job.ID == "job-success-after-retry" && attemptByJob[job.ID] < 2 {
			// 详细注释: return errors.New("temporary failure")
			return errors.New("temporary failure")
			// 详细注释: }
		}
		// 详细注释: if job.ID == "job-always-fail" {
		if job.ID == "job-always-fail" {
			// 详细注释: return errors.New("permanent failure")
			return errors.New("permanent failure")
			// 详细注释: }
		}
		// 详细注释: return nil
		return nil
		// 详细注释: }
	}

	// 详细注释: worker := NewAsyncWorker(2, handler)
	worker := NewAsyncWorker(2, handler)
	// 详细注释: tracer := NewTracer()
	tracer := NewTracer()
	// 详细注释: audit := NewAuditLogger()
	audit := NewAuditLogger()

	// 详细注释: results := worker.ProcessBatch(context.Background(), tracer, audit, []Job{
	results := worker.ProcessBatch(context.Background(), tracer, audit, []Job{
		// 详细注释: {ID: "job-success-after-retry", Payload: "todo-1"},
		{ID: "job-success-after-retry", Payload: "todo-1"},
		// 详细注释: {ID: "job-always-fail", Payload: "todo-2"},
		{ID: "job-always-fail", Payload: "todo-2"},
		// 详细注释: })
	})

	// 详细注释: if len(results) != 2 {
	if len(results) != 2 {
		// 详细注释: t.Fatalf("want 2 results got %d", len(results))
		t.Fatalf("want 2 results got %d", len(results))
		// 详细注释: }
	}
	// 详细注释: if !results[0].Success || results[0].Attempts != 2 {
	if !results[0].Success || results[0].Attempts != 2 {
		// 详细注释: t.Fatalf("unexpected first result: %+v", results[0])
		t.Fatalf("unexpected first result: %+v", results[0])
		// 详细注释: }
	}
	// 详细注释: if results[1].Success || results[1].Attempts != 3 {
	if results[1].Success || results[1].Attempts != 3 {
		// 详细注释: t.Fatalf("unexpected second result: %+v", results[1])
		t.Fatalf("unexpected second result: %+v", results[1])
		// 详细注释: }
	}

	// 详细注释: dead := worker.DeadLetters()
	dead := worker.DeadLetters()
	// 详细注释: if len(dead) != 1 || dead[0].JobID != "job-always-fail" {
	if len(dead) != 1 || dead[0].JobID != "job-always-fail" {
		// 详细注释: t.Fatalf("unexpected dead letters: %+v", dead)
		t.Fatalf("unexpected dead letters: %+v", dead)
		// 详细注释: }
	}

	// 详细注释: spans := tracer.Spans()
	spans := tracer.Spans()
	// 详细注释: if len(spans) != 2 {
	if len(spans) != 2 {
		// 详细注释: t.Fatalf("want 2 spans got %d", len(spans))
		t.Fatalf("want 2 spans got %d", len(spans))
		// 详细注释: }
	}
	// 详细注释: if spans[1].Error == "" {
	if spans[1].Error == "" {
		// 详细注释: t.Fatalf("failed job span should contain error")
		t.Fatalf("failed job span should contain error")
		// 详细注释: }
	}

	// 详细注释: events := audit.Events()
	events := audit.Events()
	// 详细注释: if len(events) != 2 {
	if len(events) != 2 {
		// 详细注释: t.Fatalf("want 2 audit events got %d", len(events))
		t.Fatalf("want 2 audit events got %d", len(events))
		// 详细注释: }
	}
	// 详细注释: }
}
