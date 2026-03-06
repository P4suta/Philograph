package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cmd := newRootCmd()

	// For server mode, block until signal
	originalRunE := cmd.RunE
	cmd.RunE = func(c *cobra.Command, args []string) error {
		err := originalRunE(c, args)
		if err != nil {
			return err
		}

		// If not JSON mode, wait for signal
		if !flagJSON {
			fmt.Fprintln(os.Stderr, "Press Ctrl+C to stop")
			<-ctx.Done()
			fmt.Fprintln(os.Stderr, "\nShutting down...")
		}

		return nil
	}

	return cmd.ExecuteContext(ctx)
}
