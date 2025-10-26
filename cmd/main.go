package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/orchestrator"
)

func main() {
	fmt.Println("[MAIN] Starting LogStream (final phase)")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle OS signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n[MAIN] signal received — canceling")
		cancel()
	}()

	// configure two pipelines
	configs := []orchestrator.PipelineConfig{
		{
			Name:       "access",
			NumWorkers: 3,
			GenCount:   60,
			GenRate:    5, // logs/sec
			BatchSize:  8,
			BatchWait:  2 * time.Second,
		},
		{
			Name:       "error",
			NumWorkers: 2,
			GenCount:   30,
			GenRate:    2,
			BatchSize:  4,
			BatchWait:  3 * time.Second,
		},
	}

	// run orchestrator, expose metrics on :9090
	if err := orchestrator.Run(ctx, configs, ":9090"); err != nil {
		log.Printf("[MAIN] orchestrator returned error: %v", err)
	} else {
		log.Println("[MAIN] orchestrator finished successfully")
	}
	// let logs flush
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[MAIN] exiting")
}
