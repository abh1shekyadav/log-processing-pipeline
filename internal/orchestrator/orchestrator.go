package orchestrator

import (
	"context"
	"fmt"
	"sync"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/consumer"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/generator"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/workerpool"
)

type PipelineConfig struct {
	Name       string
	NumWorkers int
	GenCount   int
	GenRate    int
}

func Run(ctx context.Context, configs []PipelineConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("no pipeline configs provided")
	}

	errCh := make(chan error, len(configs))

	var wg sync.WaitGroup
	wg.Add(len(configs))

	for _, cfg := range configs {
		cfg := cfg
		go func() {
			defer wg.Done()
			if err := runOnePipeline(ctx, cfg, errCh); err != nil {
				select {
				case errCh <- fmt.Errorf("%s pipeline error: %w", cfg.Name, err):
				default:
				}
			}
		}()
	}
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-waitCh:
		return nil
	}
}

func runOnePipeline(ctx context.Context, cfg PipelineConfig, errCh chan<- error) error {
	prefix := fmt.Sprintf("[%s]", cfg.Name)
	fmt.Println(prefix, "starting")

	jobs := make(chan pipeline.Log, 32)
	results := make(chan pipeline.Log, 32)

	genDone := make(chan error, 1)

	go func() {
		err := generator.GenerateLogs(ctx, jobs, cfg.GenCount)
		genDone <- err
		close(jobs)
	}()

	wp := workerpool.Workerpool{NumWorkers: cfg.NumWorkers}

	wpDone := make(chan error, 1)
	go func() {
		err := wp.ProcessLogs(ctx, jobs, results)
		wpDone <- err
	}()

	aggDone := make(chan error, 1)
	go func() {
		aggDone <- consumer.AggregateResults(ctx, results)
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(prefix, "context canceled; stopping pipeline")
			return ctx.Err()
		case err := <-genDone:
			if err != nil {
				return fmt.Errorf("generator failed: %w", err)
			}
			genDone = nil
			fmt.Println(prefix, "generator finished")
		case err := <-wpDone:
			if err != nil {
				return fmt.Errorf("workerpool failed: %w", err)
			}
			wpDone = nil
			fmt.Println(prefix, "workerpool finished")
		case err := <-aggDone:
			if err != nil {
				return fmt.Errorf("aggregator failed: %w", err)
			}
			aggDone = nil
			fmt.Println(prefix, "aggregator finished; pipeline complete")
			return nil
		}
	}
}
