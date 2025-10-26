package consumer

import (
	"context"
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/metrics"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func AggregateResults(ctx context.Context, in <-chan pipeline.Log, m *metrics.PipelineMetrics) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[AGGREGATOR] Context canceled — stopping aggregation.")
			return ctx.Err()
		case l, ok := <-in:
			if !ok {
				fmt.Println("[AGGREGATOR] All logs processed.")
				return nil
			}
			if m != nil {
				m.IncAggregated()
			}
			fmt.Printf("[AGGREGATOR] %v at %v\n", l.Message, l.Timestamp.Format("15:04:05"))
		}
	}
}
