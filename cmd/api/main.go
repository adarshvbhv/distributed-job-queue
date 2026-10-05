package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adarshvbhv/distributed-job-queue/internal/config"
	apphttp "github.com/adarshvbhv/distributed-job-queue/internal/http"
	"github.com/adarshvbhv/distributed-job-queue/internal/job"
	"github.com/adarshvbhv/distributed-job-queue/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db1, err := storage.NewPostgres(ctx, logger, cfg.DatabaseURL)
	
	if err != nil || db1 == nil {
		logger.Error("error creating new postgres", "error", err)
		return
	}

	defer db1.Pool.Close()

	repo:= job.NewPostgresRepository(db1)


	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(ch)

	server := apphttp.NewServer(logger, &cfg, repo)

	go func() {

		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("unable to start server:", "error", err)
		}

	}()

	<-ch

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error shutting down server")
	}

}
