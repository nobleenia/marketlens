package ingestion

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"marketlens/internal/models"
	"marketlens/internal/store"
	"strconv"
	"strings"
	"time"
)

type WFPIngestor struct {
	pg     *store.Postgres
	logger *log.Logger

	// Caches to avoid repeated lookups
	cropCache   map[string]string // lowercase crop name -> crop ID
	marketCache map[string]string // lowercase market name -> market ID
}

func NewWFPIngestor(pg *store.Postgres, logger *log.Logger) *WFPIngestor {
	if logger == nil {
		logger = log.Default()
	}
	return &WFPIngestor{
		pg:          pg,
		logger:      logger,
		cropCache:   make(map[string]string),
		marketCache: make(map[string]string),
	}
}

// IngestWFPCSV reads the WFP food prices QC CSV format:
//
//	date,code,usdprice
//	code = {State}-{Admin2}-{Market}-{Commodity}-{Quantity Unit}-{PriceType}-{Currency}
//
// It skips the HXL tag row (starting with #), auto-creates unknown crops,
// and resolves markets by name (must be pre-seeded).
func (w *WFPIngestor) IngestWFPCSV(ctx context.Context, r io.Reader) (inserted, skipped int, err error) {
	csvReader := csv.NewReader(r)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to read CSV: %w", err)
	}

	seen := make(map[string]bool)

	for rowIdx := 1; rowIdx < len(records); rowIdx++ {
		rec := records[rowIdx]

		// Skip HXL tag rows (start with #)
		// var skipped int
		if len(rec) > 0 && strings.HasPrefix(rec[0], "#") {
			continue
		}

		if len(rec) < 3 {
			w.logger.Printf("row %d: expected 3 columns, got %d --- skipping", rowIdx+1, len(rec))
			// skipped++
			continue
		}

		dateStr := strings.TrimSpace(rec[0])
		code := strings.TrimSpace(rec[1])
		priceStr := strings.TrimSpace(rec[2])

		// Dedup: skip if we've already seen this exact row
		dedup := dateStr + "|" + code + "|" + priceStr
		if seen[dedup] {
			skipped++
			continue
		}
		seen[dedup] = true

		// // Parse date (YYYY-MM-DD)
		observedAt, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			w.logger.Printf("row %d: invalid date %q — skipping", rowIdx+1, dateStr)
			skipped++
			continue
		}

		// Parse price
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || price <= 0 {
			w.logger.Printf("row %d: invalid or non-positive price %q — skipping", rowIdx+1, priceStr)
			skipped++
			continue
		}

		// Parse code: State-Admin2-Market-Commodity-QuantityUnit-PriceType-Currency
		// Example: Yobe-Damaturu-Damaturu-Groundnuts (shelled)-100 KG-Wholesale-NGN
		parts := strings.SplitN(code, "-", 4) // split into State, Admin2, Market, rest
		if len(parts) < 4 {
			w.logger.Printf("row %d: cannot parse code %q — skipping", rowIdx+1, code)
			skipped++
			continue
		}

		// The "rest" part has: Commodity-QuantityUnit-PriceType-Currency
		// We need to split from the RIGHT since commodity can contain hyphens (unlikely but safe)
		rest := parts[3]
		marketName := strings.TrimSpace(parts[2])
		state := strings.TrimSpace(parts[0])

		// Split rest from the right: last segment = Currency, second-to-last = PriceType,
		// third-to-last = QuantityUnit, everything before = Commodity
		restParts := strings.Split(rest, "-")
		if len(restParts) < 4 {
			w.logger.Printf("row %d: cannot parse rest-of-code %q — skipping", rowIdx+1, rest)
			skipped++
			continue
		}

		currency := strings.TrimSpace(restParts[len(restParts)-1])
		priceType := strings.TrimSpace(restParts[len(restParts)-2])
		unit := strings.TrimSpace(restParts[len(restParts)-3])
		commodity := strings.TrimSpace(strings.Join(restParts[:len(restParts)-3], "-"))

		// Resolve market ID
		marketID, err := w.pg.LookupMarketID(ctx, marketName, state)
		if err != nil {
			w.logger.Printf("row %d: market %q (state=%s) not found — skipping", rowIdx+1, marketName, state)
			skipped++
			continue
		}

		// Resolve crop ID (auto-create if unknown)
		cropID, err := w.pg.LookupCropID(ctx, commodity)
		if err != nil {
			// Auto-create the crop
			cropID, err = w.pg.InsertCrop(ctx, commodity, unit)
			if err != nil {
				w.logger.Printf("row %d: failed to create crop %q — skipping: %v", rowIdx+1, commodity, err)
				skipped++
				continue
			}
			w.logger.Printf("auto-created crop: %s (unit=%s, id=%s)", commodity, unit, cropID)
		}

		obs := models.PriceObservation{
			CropID:          cropID,
			MarketID:        marketID,
			ObservedAt:      observedAt,
			Price:           price,
			Currency:        currency,
			Unit:            unit,
			PriceType:       strings.ToLower(priceType),
			Source:          "wfp",
			Notes:           fmt.Sprintf("WFP QC import; original unit: %s", unit),
			ConfidenceScore: 0.70, // WFP data is fairly reliable
		}

		if err := w.pg.InsertPriceObservation(ctx, obs); err != nil {
			w.logger.Printf("row %d: insert failed: %v — skipping", rowIdx+1, err)
			skipped++
			continue
		}

		inserted++
	}
	return inserted, skipped, nil
}
