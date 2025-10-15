package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Listen(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	fmt.Printf("\nCaught signal: %v — shutting down...\n", sig)
	cancel()
}
