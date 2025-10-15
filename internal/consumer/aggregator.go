package consumer

import (
	"context"
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func AggregateResults(ctx context.Context, in <-chan pipeline.Log) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[AGGREGATOR] Context canceled — stopping aggregation.")
			return
		case log, ok := <-in:
			if !ok {
				fmt.Println("[AGGREGATOR] All logs processed.")
				return
			}
			fmt.Printf("[AGGREGATOR] %v at %v\n", log.Message, log.Timestamp.Format("15:04:05"))
		}
	}
}
