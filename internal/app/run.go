package app

import (
	"context"
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/consumer"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/generator"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/workerpool"
)

func Run(ctx context.Context) error {
	fmt.Println("[PIPELINE] Starting...")
	jobs := make(chan pipeline.Log, 10)
	results := make(chan pipeline.Log, 10)
	errCh := make(chan error, 1)

	go func() {
		defer close(jobs)
		if err := generator.GenerateLogs(ctx, jobs, 20); err != nil {
			errCh <- fmt.Errorf("generator error: %w", err)
		}
	}()
	wp := workerpool.Workerpool{NumWorkers: 3}
	go func() {
		if err := wp.ProcessLogs(ctx, jobs, results); err != nil {
			errCh <- fmt.Errorf("workerpool error: %w", err)
		}
	}()
	go func() {
		if err := consumer.AggregateResults(ctx, results); err != nil {
			errCh <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("[APP] Context canceled. Stopping pipeline...")
	case err := <-errCh:
		return err
	}
	fmt.Println("[APP] Pipeline stopped gracefully.")
	return nil
}
