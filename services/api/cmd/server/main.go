package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketlens/internal/cache"
	"marketlens/internal/config"
	"marketlens/internal/httpserver"
	"marketlens/internal/store"
)

func main() {
	cfg := config.FromEnv()

	// dependency initialization
	db, err := store.NewPostgresPool(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer db.Close()

	rdb, err := cache.NewRedisClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to initialize redis: %v", err)
	}
	defer func() {
		_ = rdb.Close()
	}()

	srv := httpserver.New(cfg)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("marketlens-api listening on: %s (env=%s)", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
