package week14

import (
	"context"
	"errors"
	"testing"
)

// TestProcessBatchRetriesAndDeadLetter 验证重试和死信行为。
func TestProcessBatchRetriesAndDeadLetter(t *testing.T) {
	attemptByJob := map[string]int{}
	handler := func(_ context.Context, job Job, _ int) error {
		attemptByJob[job.ID]++
		if job.ID == "job-success-after-retry" && attemptByJob[job.ID] < 2 {
			return errors.New("temporary failure")
		}
		if job.ID == "job-always-fail" {
			return errors.New("permanent failure")
		}
		return nil
	}

	worker := NewAsyncWorker(2, handler)
	tracer := NewTracer()
	audit := NewAuditLogger()

	results := worker.ProcessBatch(context.Background(), tracer, audit, []Job{
		{ID: "job-success-after-retry", Payload: "todo-1"},
		{ID: "job-always-fail", Payload: "todo-2"},
	})

	if len(results) != 2 {
		t.Fatalf("want 2 results got %d", len(results))
	}
	if !results[0].Success || results[0].Attempts != 2 {
		t.Fatalf("unexpected first result: %+v", results[0])
	}
	if results[1].Success || results[1].Attempts != 3 {
		t.Fatalf("unexpected second result: %+v", results[1])
	}

	dead := worker.DeadLetters()
	if len(dead) != 1 || dead[0].JobID != "job-always-fail" {
		t.Fatalf("unexpected dead letters: %+v", dead)
	}

	spans := tracer.Spans()
	if len(spans) != 2 {
		t.Fatalf("want 2 spans got %d", len(spans))
	}
	if spans[1].Error == "" {
		t.Fatalf("failed job span should contain error")
	}

	events := audit.Events()
	if len(events) != 2 {
		t.Fatalf("want 2 audit events got %d", len(events))
	}
}
