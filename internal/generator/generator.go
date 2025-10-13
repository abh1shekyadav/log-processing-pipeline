package generator

import (
	"fmt"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func GenerateLogs(out chan<- pipeline.Log, count int) {
	for i := 1; i <= count; i++ {
		log := pipeline.Log{
			ID:      i,
			Message: fmt.Sprintf("Log message #%d", i),
		}
		out <- log
		time.Sleep(200 * time.Millisecond)
	}
	close(out)
}
