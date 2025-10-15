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

	go generator.GenerateLogs(ctx, jobs, 20)
	wp := workerpool.Workerpool{NumWorkers: 3}
	go wp.ProcessLogs(ctx, jobs, results)
	consumer.AggregateResults(ctx, results)

	<-ctx.Done()
	fmt.Println("[PIPELINE] Context canceled. Shutting down.")
	return nil
}
