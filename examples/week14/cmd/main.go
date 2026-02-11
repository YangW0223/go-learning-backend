package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/yang/go-learning-backend/examples/week14"
)

func main() {
	attemptByJob := map[string]int{}
	worker := week14.NewAsyncWorker(2, func(_ context.Context, job week14.Job, _ int) error {
		attemptByJob[job.ID]++
		if job.ID == "job-2" {
			return errors.New("permanent error")
		}
		if job.ID == "job-1" && attemptByJob[job.ID] < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	tracer := week14.NewTracer()
	audit := week14.NewAuditLogger()

	results := worker.ProcessBatch(context.Background(), tracer, audit, []week14.Job{
		{ID: "job-1", Payload: "notify todo 1"},
		{ID: "job-2", Payload: "notify todo 2"},
	})

	fmt.Println("results:", results)
	fmt.Println("dead letters:", worker.DeadLetters())
	fmt.Println("spans:", tracer.Spans())
	fmt.Println("audit events:", audit.Events())
}
