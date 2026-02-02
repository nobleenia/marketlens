package ingestion

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"marketlens/internal/store"
)

type Record struct {
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

type Ingestor struct {
	pg     *store.Postgres
	logger *log.Logger
}

func NewIngestor(pg *store.Postgres, logger *log.Logger) *Ingestor {
	if logger == nil {
		logger = log.Default()
	}
	return &Ingestor{
		pg:     pg,
		logger: logger,
	}
}

// IngestCSV ingests records from a CSV reader into the database.
// oberved_at,crop_name,market_name,state,price,currency,unit,source,reporter_id,notes
func (ing *Ingestor) IngestCSV(ctx context.Context, r io.Reader) (inserted int, skipped int, err error) {
	csvReader := csv.NewReader(r)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read CSV: %w", err)
	}
	if len(records) < 2 {
		return 0, 0, fmt.Errorf("CSV has no data rows")
	}

	for rowIdx := 1; rowIdx < len(records); rowIdx++ {
		rec := records[rowIdx]

		if len(rec) < 10 {
			ing.logger.Printf("skipping row %d: expected 10 columns, got %d", rowIdx+1, len(rec))
			skipped++
			continue
		}

		// Trim values to avoid issues with extra spaces (helps when CSV contains spaces around commas)
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}

		t, err := time.Parse(time.RFC3339, rec[0])
		if err != nil {
			ing.logger.Printf("skipping row %d: invalid observed_at: %v", rowIdx+1, err)
			skipped++
			continue
		}

		price, err := strconv.ParseFloat(rec[4], 64)
		if err != nil {
			ing.logger.Printf("skipping row %d: invalid price: %v", rowIdx+1, err)
			skipped++
			continue
		}

		record := &Record{
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

		ok, err := ing.upsertObservation(ctx, record)
		if err != nil {
			ing.logger.Printf("failed to upsert row %d: %v", rowIdx+1, err)
			skipped++
			continue
		}
		if ok {
			inserted++
		} else {
			skipped++
		}
	}

	return inserted, skipped, nil
}

func (ing *Ingestor) upsertObservation(ctx context.Context, rec *Record) (bool, error) {
	// Resolve IDs by names (case-insensitive)
	cropID, err := ing.pg.LookupCropID(ctx, rec.CropName)
	if err != nil {
		return false, fmt.Errorf("lookup crop ID: %w", err)
	}

	marketID, err := ing.pg.LookupMarketID(ctx, rec.MarketName, rec.State)
	if err != nil {
		return false, fmt.Errorf("lookup market ID: %w", err)
	}

	err = ing.pg.InsertPriceObservation(ctx, store.PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      rec.ObservedAt,
		Price:           rec.Price,
		Currency:        rec.Currency,
		Unit:            rec.Unit,
		Source:          rec.Source,
		ReporterID:      rec.ReporterID,
		Notes:           rec.Notes,
		ConfidenceScore: 0.50,
	})
	if err != nil {
		return false, fmt.Errorf("insert price observation: %w", err)
	}

	return true, nil
}
