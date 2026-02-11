package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yang/go-learning-backend/examples/week08"
)

func main() {
	source := &week08.SlowSource{
		Delay: 50 * time.Millisecond,
		Data: map[string][]week08.Todo{
			"u1": {
				{ID: "1", Title: "cache todo 1"},
				{ID: "2", Title: "cache todo 2"},
			},
		},
	}
	service := week08.NewCachedTodoService(source, week08.NewTTLCache(), 10*time.Second)

	d1, _ := week08.MeasureLatency(func() error {
		_, err := service.List(context.Background(), "u1", 1, 10)
		return err
	})
	d2, _ := week08.MeasureLatency(func() error {
		_, err := service.List(context.Background(), "u1", 1, 10)
		return err
	})
	hit, miss := service.Stats()

	fmt.Printf("first call latency=%s\n", d1)
	fmt.Printf("second call latency=%s\n", d2)
	fmt.Printf("cache stats hit=%d miss=%d\n", hit, miss)
}
