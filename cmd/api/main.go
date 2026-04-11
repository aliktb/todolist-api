package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ali/todolist-api/internal/config"
	"github.com/ali/todolist-api/internal/db"
	apihttp "github.com/ali/todolist-api/internal/http"
	"github.com/ali/todolist-api/internal/tasks"
)

// @title Todo List API
// @version 1.0
// @description REST API for managing todo tasks.
// @BasePath /
// @schemes http https
func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	taskRepo := tasks.NewRepository(pool)
	handler := apihttp.NewHandler(cfg, logger, taskRepo, pool)

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           handler.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("starting api server", "addr", cfg.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("api server exited", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down api server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
}
