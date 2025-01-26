package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	headliner "github.com/lmullen/chronam-headliner"
)

func main() {

	// Initialize the logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	var server *headliner.Server

	// Create a context and listen for signals to gracefully shutdown the application
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Clean up function that will be called at program end no matter what
	defer func() {
		signal.Stop(quit)
		cancel()
	}()

	// Listen for shutdown signals in a go-routine and cancel context then
	go func() {
		select {
		case <-quit:
			slog.Info("shutdown signal received, quitting gracefully")
			cancel()
			server.Shutdown()
		case <-ctx.Done():
		}
	}()

	server = headliner.NewServer(ctx)
	err := server.Run()
	if err != nil {
		slog.Error("error running the server", "error", err)
	}

}
