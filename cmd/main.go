package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abh1shekyadav/log-processing-pipeline/internal/app"
	"github.com/abh1shekyadav/log-processing-pipeline/internal/shutdown"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go shutdown.Listen(cancel)

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Pipeline stopped with error: %v", err)
	}
	fmt.Println("Pipeline finished.")
}
