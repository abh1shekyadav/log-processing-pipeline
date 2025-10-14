package main

import (
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/consumer"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/generator"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/workerpool"
)

func main() {
	fmt.Println("Starting LogStream")

	jobs := make(chan pipeline.Log, 10)
	results := make(chan pipeline.Log, 10)

	go generator.GenerateLogs(jobs, 15)

	wp := workerpool.Workerpool{NumWorkers: 3}
	go wp.ProcessLogs(jobs, results)

	consumer.AggregateResults(results)
	fmt.Println("Pipeline finished.")
}
