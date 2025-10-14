package consumer

import (
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func AggregateResults(in <-chan pipeline.Log) {
	for log := range in {
		fmt.Printf("[AGGREGATOR] %v at %v\n", log.Message, log.Timestamp.Format("15:04:05"))
	}
	fmt.Println("[AGGREGATOR] All logs processed.")

}
