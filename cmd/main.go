package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	headliner "github.com/lmullen/chronam-headliner"
)

func main() {
	// Initialize the logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Create a context and listen for signals to gracefully shutdown the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the app first
	app := headliner.NewApp(ctx)

	slog.Debug("starting the app")

	// Set up signal handling after app creation
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(quit)

	// Listen for shutdown signals in a go-routine
	go func() {
		select {
		case sig := <-quit:
			slog.Info("shutting down in response to signal", "signal", sig.String())

			// Create a timeout context for shutdown
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			done := make(chan struct{})
			go func() {
				app.Shutdown()
				close(done)
			}()

			select {
			case <-done:
				slog.Info("shutdown completed cleanly")
			case <-shutdownCtx.Done():
				slog.Error("shutdown timed out, forcing exit")
				os.Exit(1)
			}

			cancel()
		case <-ctx.Done():
			slog.Debug("context cancelled, shutting down signal handler")
		}
	}()

	// Run the application
	if err := app.Run(); err != nil {
		slog.Error("error running the application", "error", err)
		panic(nil)
	}
}
