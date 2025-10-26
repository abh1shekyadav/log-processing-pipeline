package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/orchestrator"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/shutdown"
)

func main() {

	fmt.Println("[MAIN] Starting LogStream (orchestrator)")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go shutdown.Listen(cancel)

	configs := []orchestrator.PipelineConfig{
		{Name: "access-log", NumWorkers: 3, GenCount: 50, GenRate: 5},
		{Name: "error-log", NumWorkers: 2, GenCount: 20, GenRate: 2},
	}
	if err := orchestrator.Run(ctx, configs); err != nil {
		log.Printf("[MAIN] orchestrator finished with error: %v", err)
	} else {
		log.Println("[MAIN] orchestrator finished successfully")
	}
	time.Sleep(100 * time.Millisecond)
}
