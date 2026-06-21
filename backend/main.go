package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/app"
	"github.com/carkeeper/backend/internal/handler"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/service"
	"github.com/carkeeper/backend/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	initLogging(cfg)

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	slog.Info("database connection established")

	fileStore, err := storage.NewLocal(cfg.Storage.RootPath)
	if err != nil {
		log.Fatalf("Document storage: %v", err)
	}

	repos := repository.New(db)
	service.BootstrapAuthz(context.Background(), repos)
	services := service.New(repos, cfg, fileStore)
	handlers := handler.New(services, cfg)
	router := app.NewRouter(handlers, cfg, db)

	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 65 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "addr", cfg.Server.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("server exited")
}

func initLogging(cfg *config.Config) {
	var h slog.Handler
	if cfg.Env == "production" {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(h))
}
