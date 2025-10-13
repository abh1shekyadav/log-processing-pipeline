package main

import (
	"fmt"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/consumer"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/generator"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/pipeline"
)

func main() {
	fmt.Println("Starting LogStream")
	logCh := make(chan pipeline.Log)
	go generator.GenerateLogs(logCh, 10)

	consumer.ConsumeLogs(logCh)
	fmt.Println("Pipeline finished.")
}
