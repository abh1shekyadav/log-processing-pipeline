package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/metrics"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func GenerateLogs(ctx context.Context, out chan<- pipeline.Log, count int, rate int, m *metrics.PipelineMetrics) error {
	if count <= 0 {
		return fmt.Errorf("GenerateLogs: count must be > 0")
	}
	var ticker *time.Ticker
	if rate > 0 {
		interval := time.Second / time.Duration(rate)
		ticker = time.NewTicker(interval)
		defer ticker.Stop()
	}
	for i := 1; i <= count; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if ticker != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
				}
			}
			log := pipeline.Log{
				ID:      i,
				Message: fmt.Sprintf("Log message #%d", i),
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- log:
			}
			if m != nil {
				m.IncGenerated()
			}
		}
	}
	return nil
}
