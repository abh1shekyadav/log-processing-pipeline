package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

type Workerpool struct {
	NumWorkers int
}

func (wp *Workerpool) ProcessLogs(ctx context.Context, in <-chan pipeline.Log, out chan<- pipeline.Log) {
	var wg sync.WaitGroup

	for i := 1; i <= wp.NumWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("[Worker-%d] Context canceled — exiting.\n", workerID)
					return
				case log, ok := <-in:
					if !ok {
						return
					}
					processed := processLog(workerID, log)
					out <- processed
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func processLog(workerID int, log pipeline.Log) pipeline.Log {
	time.Sleep(200 * time.Millisecond)
	log.Message = fmt.Sprintf("[Worker-%d] %s", workerID, log.Message)
	log.Processed = true
	log.Timestamp = time.Now()
	return log
}
