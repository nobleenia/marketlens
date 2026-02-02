package main

import (
	"context"
	"flag"
	"log"
	"os"

	"marketlens/internal/config"
	"marketlens/internal/ingestion"
	"marketlens/internal/store"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "", "Path to CSV file to ingest")
	flag.Parse()

	if file == "" {
		log.Fatal("Please provide a CSV file path using the -file flag")
	}

	cfg := config.FromEnv()
	ctx := context.Background()

	pg, err := store.NewPostgres(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	defer pg.Close()

	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	ingestor := ingestion.NewIngestor(pg, log.Default())

	inserted, skipped, err := ingestor.IngestCSV(ctx, f)
	if err != nil {
		log.Fatalf("failed to ingest CSV: %v", err)
	}

	log.Printf("Ingestion complete. Inserted: %d, Skipped: %d", inserted, skipped)
}
