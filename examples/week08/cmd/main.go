// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "time"
	"time"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week08"
	"github.com/yang/go-learning-backend/examples/week08"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: source := &week08.SlowSource{
	source := &week08.SlowSource{
		// 详细注释: Delay: 50 * time.Millisecond,
		Delay: 50 * time.Millisecond,
		// 详细注释: Data: map[string][]week08.Todo{
		Data: map[string][]week08.Todo{
			// 详细注释: "u1": {
			"u1": {
				// 详细注释: {ID: "1", Title: "cache todo 1"},
				{ID: "1", Title: "cache todo 1"},
				// 详细注释: {ID: "2", Title: "cache todo 2"},
				{ID: "2", Title: "cache todo 2"},
				// 详细注释: },
			},
			// 详细注释: },
		},
		// 详细注释: }
	}
	// 详细注释: service := week08.NewCachedTodoService(source, week08.NewTTLCache(), 10*time.Second)
	service := week08.NewCachedTodoService(source, week08.NewTTLCache(), 10*time.Second)

	// 详细注释: d1, _ := week08.MeasureLatency(func() error {
	d1, _ := week08.MeasureLatency(func() error {
		// 详细注释: _, err := service.List(context.Background(), "u1", 1, 10)
		_, err := service.List(context.Background(), "u1", 1, 10)
		// 详细注释: return err
		return err
		// 详细注释: })
	})
	// 详细注释: d2, _ := week08.MeasureLatency(func() error {
	d2, _ := week08.MeasureLatency(func() error {
		// 详细注释: _, err := service.List(context.Background(), "u1", 1, 10)
		_, err := service.List(context.Background(), "u1", 1, 10)
		// 详细注释: return err
		return err
		// 详细注释: })
	})
	// 详细注释: hit, miss := service.Stats()
	hit, miss := service.Stats()

	// 详细注释: fmt.Printf("first call latency=%s\n", d1)
	fmt.Printf("first call latency=%s\n", d1)
	// 详细注释: fmt.Printf("second call latency=%s\n", d2)
	fmt.Printf("second call latency=%s\n", d2)
	// 详细注释: fmt.Printf("cache stats hit=%d miss=%d\n", hit, miss)
	fmt.Printf("cache stats hit=%d miss=%d\n", hit, miss)
	// 详细注释: }
}
