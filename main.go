package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/svrakitin/breitbandmessung/cmd/breitbandmessung"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := breitbandmessung.Execute(ctx); err != nil {
		os.Exit(1)
	}
}
