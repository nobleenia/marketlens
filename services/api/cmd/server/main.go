package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketlens/internal/aggregation"
	"marketlens/internal/cache"
	"marketlens/internal/config"
	"marketlens/internal/httpserver"
	"marketlens/internal/store"
)

func main() {
	cfg := config.FromEnv()

	// dependency initialization
	db, err := store.NewPostgres(context.Background(), cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer db.Close()

	// Initialize Redis cache (optional, can be used for caching aggregated results or other data)
	rdb, err := cache.NewRedisClient(context.Background(), cfg)
	if err != nil {
		log.Printf("failed to initialize redis: %v", err)
		rdb = nil // proceed without Redis
	} else {
		defer func() {
			_ = rdb.Close()
		}()
	}

	// initialize and run server
	srv := httpserver.New(cfg, db)

	// Start background aggregation job
	go runScheduledAggregation(db)

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

func runScheduledAggregation(db *store.Postgres) {
	svc := aggregation.NewService(db)

	// Run immediately on startup
	log.Println("Running initial aggregation...")
	if err := svc.RunDailyAggregation(context.Background()); err != nil {
		log.Printf("initial aggregation error: %v", err)
	} else {
		log.Println("Initial aggregation completed successfully")
	}

	// Schedule to run every 1 hour (can be adjusted as needed)
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running scheduled aggregation...")
		if err := svc.RunDailyAggregation(context.Background()); err != nil {
			log.Printf("scheduled aggregation error: %v", err)
		} else {
			log.Println("Scheduled aggregation completed successfully")
		}
	}
}
