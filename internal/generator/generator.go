package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func GenerateLogs(ctx context.Context, out chan<- pipeline.Log, count int) {
	defer close(out)
	for i := 1; i <= count; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("[GENERATOR] Context canceled — stopping log generation.")
			return
		default:
			log := pipeline.Log{
				ID:      i,
				Message: fmt.Sprintf("Log message #%d", i),
			}
			out <- log
			time.Sleep(200 * time.Millisecond)
		}
	}
}
