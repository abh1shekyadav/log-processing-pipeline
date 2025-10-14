package consumer

import (
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func ConsumeLogs(in <-chan pipeline.Log) {
	for log := range in {
		fmt.Printf("[CONSUMER] Received: %v\n", log)
	}
	fmt.Println("[CONSUMER] All logs consumed.")
}
