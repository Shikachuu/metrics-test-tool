package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Shikachuu/metrics-test-tool/pkg/handler"
)

func main() {
	ctx := context.Background()
	if err := Run(ctx, os.Stdout, "8080", "text"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %v\n", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, w io.Writer, port string, logType string) error {
	var wg sync.WaitGroup
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := setLogger(logType, w)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      handler.NewServer(logger, http.NewServeMux()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	wg.Add(1)
	go func() {
		logger.InfoContext(ctx, "starting http server...", "port", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "failed to serve http", "error", err.Error())
			os.Exit(1)
		}
		wg.Done()
	}()

	go func() {
		<-ctx.Done()
		logger.InfoContext(ctx, "shutting down http server...")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.ErrorContext(ctx, "failed to shutdown http server", "error", err.Error())
			os.Exit(1)
		}
	}()

	wg.Wait()

	return nil
}

func setLogger(logType string, w io.Writer) *slog.Logger {
	var lh slog.Handler

	switch logType {
	case "json":
		lh = slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
	default:
		lh = slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
	}

	l := slog.New(lh)
	slog.SetDefault(l)

	return l
}
