// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "fmt"
	"fmt"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week14"
	"github.com/yang/go-learning-backend/examples/week14"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: attemptByJob := map[string]int{}
	attemptByJob := map[string]int{}
	// 详细注释: worker := week14.NewAsyncWorker(2, func(_ context.Context, job week14.Job, _ int) error {
	worker := week14.NewAsyncWorker(2, func(_ context.Context, job week14.Job, _ int) error {
		// 详细注释: attemptByJob[job.ID]++
		attemptByJob[job.ID]++
		// 详细注释: if job.ID == "job-2" {
		if job.ID == "job-2" {
			// 详细注释: return errors.New("permanent error")
			return errors.New("permanent error")
			// 详细注释: }
		}
		// 详细注释: if job.ID == "job-1" && attemptByJob[job.ID] < 2 {
		if job.ID == "job-1" && attemptByJob[job.ID] < 2 {
			// 详细注释: return errors.New("temporary error")
			return errors.New("temporary error")
			// 详细注释: }
		}
		// 详细注释: return nil
		return nil
		// 详细注释: })
	})

	// 详细注释: tracer := week14.NewTracer()
	tracer := week14.NewTracer()
	// 详细注释: audit := week14.NewAuditLogger()
	audit := week14.NewAuditLogger()

	// 详细注释: results := worker.ProcessBatch(context.Background(), tracer, audit, []week14.Job{
	results := worker.ProcessBatch(context.Background(), tracer, audit, []week14.Job{
		// 详细注释: {ID: "job-1", Payload: "notify todo 1"},
		{ID: "job-1", Payload: "notify todo 1"},
		// 详细注释: {ID: "job-2", Payload: "notify todo 2"},
		{ID: "job-2", Payload: "notify todo 2"},
		// 详细注释: })
	})

	// 详细注释: fmt.Println("results:", results)
	fmt.Println("results:", results)
	// 详细注释: fmt.Println("dead letters:", worker.DeadLetters())
	fmt.Println("dead letters:", worker.DeadLetters())
	// 详细注释: fmt.Println("spans:", tracer.Spans())
	fmt.Println("spans:", tracer.Spans())
	// 详细注释: fmt.Println("audit events:", audit.Events())
	fmt.Println("audit events:", audit.Events())
	// 详细注释: }
}
