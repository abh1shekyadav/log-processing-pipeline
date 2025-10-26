package orchestrator

import (
	"context"
	"fmt"
	"sync"

	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/consumer"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/generator"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/metrics"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/workerpool"
)

// PipelineConfig holds per-pipeline configuration.
type PipelineConfig struct {
	Name       string
	NumWorkers int
	GenCount   int
	GenRate    int
	BatchSize  int
	BatchWait  time.Duration
}

// Run starts pipelines and returns when all complete or an error occurs.
func Run(ctx context.Context, configs []PipelineConfig, httpMetricsAddr string) error {
	if len(configs) == 0 {
		return fmt.Errorf("no pipeline configs provided")
	}

	// Build metrics list and register HTTP endpoint if requested.
	metricsList := make([]*metrics.PipelineMetrics, 0, len(configs))
	for _, c := range configs {
		metricsList = append(metricsList, metrics.NewPipelineMetrics(c.Name))
	}
	if httpMetricsAddr != "" {
		metrics.StartHTTPReporter(httpMetricsAddr, metricsList)
	}

	// central error channel
	errCh := make(chan error, len(configs))
	var wg sync.WaitGroup
	wg.Add(len(configs))

	for i := range configs {
		cfg := configs[i]
		m := metricsList[i]
		go func(cfg PipelineConfig, m *metrics.PipelineMetrics) {
			defer wg.Done()
			if err := runOnePipeline(ctx, cfg, m); err != nil {
				select {
				case errCh <- fmt.Errorf("%s pipeline error: %w", cfg.Name, err):
				default:
				}
			}
		}(cfg, m)
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

// runOnePipeline composes generator -> workerpool -> consumer for a single pipeline.
func runOnePipeline(ctx context.Context, cfg PipelineConfig, m *metrics.PipelineMetrics) error {
	prefix := fmt.Sprintf("[%s]", cfg.Name)
	fmt.Println(prefix, "starting")

	// per-pipeline channels (owned by this function)
	jobs := make(chan pipeline.Log, 128)
	results := make(chan pipeline.Log, 128)

	// generator: generate and return; caller closes jobs
	genDone := make(chan error, 1)
	go func() {
		err := generator.GenerateLogs(ctx, jobs, cfg.GenCount, cfg.GenRate, m)
		genDone <- err
		close(genDone)
		// caller will close jobs
	}()

	// Close jobs when generator returns (so workerpool sees end-of-stream)
	go func() {
		err := <-genDone
		_ = err // error handled below via genDone value
		close(jobs)
	}()

	// workerpool: processes and closes results when done
	wp := workerpool.Workerpool{
		NumWorkers: cfg.NumWorkers,
		BatchSize:  cfg.BatchSize,
		BatchWait:  cfg.BatchWait,
	}
	wpDone := make(chan error, 1)
	go func() {
		err := wp.ProcessLogs(ctx, jobs, results, m)
		wpDone <- err
		close(wpDone)
	}()

	// consumer: read results (its error will signal completion)
	aggDone := make(chan error, 1)
	go func() {
		err := consumer.AggregateResults(ctx, results, m)
		aggDone <- err
		close(aggDone)
	}()

	// monitor
	for {
		select {
		case <-ctx.Done():
			fmt.Println(prefix, "context canceled, returning")
			return ctx.Err()
		case err := <-genDone:
			if err != nil {
				return fmt.Errorf("generator failed: %w", err)
			}
			fmt.Println(prefix, "generator finished")
			// continue waiting for wp/agg
		case err := <-wpDone:
			if err != nil {
				return fmt.Errorf("workerpool failed: %w", err)
			}
			fmt.Println(prefix, "workerpool finished")
		case err := <-aggDone:
			if err != nil {
				return fmt.Errorf("aggregator failed: %w", err)
			}
			fmt.Println(prefix, "aggregator finished; pipeline complete")
			return nil
		}
	}
}
