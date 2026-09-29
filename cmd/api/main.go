package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hexleturlshort/internal/apiapp"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	err := apiapp.Run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
