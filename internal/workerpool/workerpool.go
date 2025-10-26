package workerpool

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/metrics"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

type Workerpool struct {
	NumWorkers int
	BatchSize  int
	BatchWait  time.Duration
}

func (wp *Workerpool) ProcessLogs(ctx context.Context, in <-chan pipeline.Log, out chan<- pipeline.Log, m *metrics.PipelineMetrics) error {
	defer close(out)

	if wp.BatchSize <= 0 {
		wp.BatchSize = 10
	}
	if wp.BatchWait <= 0 {
		wp.BatchWait = 1 * time.Second
	}
	if wp.NumWorkers <= 0 {
		wp.NumWorkers = 1
	}

	batch := make([]pipeline.Log, 0, wp.BatchSize)
	timer := time.NewTimer(wp.BatchWait)
	defer timer.Stop()

	flush := func(ctx context.Context, toProcess []pipeline.Log) error {
		if len(toProcess) == 0 {
			return nil
		}
		if m != nil {
			m.IncBatches()
		}

		var wg sync.WaitGroup
		sem := make(chan struct{}, wp.NumWorkers)
		for _, entry := range toProcess {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			wg.Add(1)
			sem <- struct{}{}
			go func(e pipeline.Log) {
				defer wg.Done()
				time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
				e.Processed = true
				e.Timestamp = time.Now()
				e.Message = fmt.Sprintf("[processed] %s", e.Message)
				select {
				case <-ctx.Done():
				case out <- e:
					if m != nil {
						m.IncProcessed()
					}
				}
				<-sem
			}(entry)
		}
		wg.Wait()
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			_ = flush(ctx, batch)
			return ctx.Err()
		case l, ok := <-in:
			if !ok {
				_ = flush(ctx, batch)
				return nil
			}
			batch = append(batch, l)
			if len(batch) >= wp.BatchSize {
				toProcess := batch
				batch = make([]pipeline.Log, 0, wp.BatchSize)
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(wp.BatchWait)
				if err := flush(ctx, toProcess); err != nil {
					return err
				}
			}
		case <-timer.C:
			if len(batch) > 0 {
				toProcess := batch
				batch = make([]pipeline.Log, 0, wp.BatchSize)
				timer.Reset(wp.BatchWait)
				if err := flush(ctx, toProcess); err != nil {
					return err
				}
			} else {
				timer.Reset(wp.BatchWait)
			}
		}
	}
}
