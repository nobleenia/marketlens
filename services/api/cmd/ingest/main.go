package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"marketlens/internal/config"
	"marketlens/internal/store"
)

type row struct {
	ObservedAt time.Time
	CropName   string
	MarketName string
	State      string
	Price      float64
	Currency   string
	Unit       string
	Source     string
	ReporterID string
	Notes      string
}

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

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("failed to read CSV file: %v", err)
	}
	if len(records) < 2 {
		log.Fatal("CSV file has no data rows")
	}

	// header row expected
	inserted := 0
	skipped := 0

	for i := 1; i < len(records); i++ {
		rec := records[i]
		// observed_at,crop_name,market_name,state,price,currency,unit,source,reporter_id,notes
		if len(rec) < 10 {
			log.Printf("skipping row %d: expected 10 columns, got %d", i+1, len(rec))
			skipped++
			continue
		}

		t, err := time.Parse(time.RFC3339, rec[0])
		if err != nil {
			log.Printf("skipping row %d: invalid observed_at: %v", i+1, err)
			skipped++
			continue
		}
		price, err := strconv.ParseFloat(rec[4], 64)
		if err != nil {
			log.Printf("skipping row %d: invalid price: %v", i+1, err)
			skipped++
			continue
		}

		row := &row{
			ObservedAt: t,
			CropName:   rec[1],
			MarketName: rec[2],
			State:      rec[3],
			Price:      price,
			Currency:   rec[5],
			Unit:       rec[6],
			Source:     rec[7],
			ReporterID: rec[8],
			Notes:      rec[9],
		}

		ok, err := upsertObservation(ctx, pg, row)
		if err != nil {
			log.Printf("failed to upsert row %d: %v", i+1, err)
			skipped++
			continue
		}
		if ok {
			inserted++
		} else {
			skipped++
		}
	}

	log.Printf("Ingestion complete. Inserted: %d, Skipped: %d\n", inserted, skipped)
}

func upsertObservation(ctx context.Context, pg *store.Postgres, r *row) (bool, error) {
	// Resolve IDs by names (case-insensitive)
	cropID, err := pg.LookupCropID(ctx, r.CropName)
	if err != nil {
		return false, fmt.Errorf("Crop lookup %q/%q: %w", r.CropName, err)
	}

	marketID, err := pg.LookupMarketID(ctx, r.MarketName, r.State)
	if err != nil {
		return false, fmt.Errorf("Market lookup %q, %q/%q: %w", r.MarketName, r.State, err)
	}

	err = pg.InsertPriceObservation(ctx, store.PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      r.ObservedAt,
		Price:           r.Price,
		Currency:        r.Currency,
		Unit:            r.Unit,
		Source:          r.Source,
		ReporterID:      r.ReporterID,
		Notes:           r.Notes,
		ConfidenceScore: 0.50,
	})
	if err != nil {
		return false, fmt.Errorf("InsertPriceObservation: %w", err)
	}

	return true, nil
}
