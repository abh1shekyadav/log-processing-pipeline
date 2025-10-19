package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/metrics"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func GenerateLogs(ctx context.Context, out chan<- pipeline.Log, count int) error {
	for i := 1; i <= count; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log := pipeline.Log{
				ID:      i,
				Message: fmt.Sprintf("Log message #%d", i),
			}
			out <- log
			metrics.IncGenerated()
			time.Sleep(200 * time.Millisecond)
		}
	}
	return nil
}
