package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/metrics"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

type Workerpool struct {
	NumWorkers int
}

func (wp *Workerpool) ProcessLogs(ctx context.Context, in <-chan pipeline.Log, out chan<- pipeline.Log) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
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
					select {
					case out <- processed:
					case <-ctx.Done():
						errCh <- ctx.Err() // Capture cancellation error
						return
					}
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errCh)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func processLog(workerID int, log pipeline.Log) pipeline.Log {
	time.Sleep(200 * time.Millisecond)
	log.Message = fmt.Sprintf("[Worker-%d] %s", workerID, log.Message)
	log.Processed = true
	log.Timestamp = time.Now()
	metrics.IncProcessed()
	return log
}
