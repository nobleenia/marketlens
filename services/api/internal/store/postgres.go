package store

import (
	"context"
	"fmt"
	"marketlens/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres dsn: %w", err)
	}

	poolCfg.MaxConns = 5
	poolCfg.MinConns = 1
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	// quick ping to verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.pool.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.pool.Close()
}

func (pg *Postgres) InsertPriceObservation(ctx context.Context, obs models.PriceObservation) error {
	_, err := pg.pool.Exec(ctx, `
		INSERT INTO price_observations
		(crop_id, market_id, observed_at, price, currency, unit, source, reporter_id, notes, confidence_score)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, obs.CropID, obs.MarketID, obs.ObservedAt, obs.Price, obs.Currency, obs.Unit,
		obs.Source, obs.ReporterID, obs.Notes, obs.ConfidenceScore)
	if err != nil {
		return fmt.Errorf("InsertPriceObservation: %w", err)
	}
	return nil
}

func (pg *Postgres) LookupCropID(ctx context.Context, cropName string) (string, error) {
	var cropID string
	err := pg.pool.QueryRow(ctx, `
		SELECT id FROM crops
		WHERE LOWER(name) = LOWER($1)
	`, cropName).Scan(&cropID)
	if err != nil {
		return "", fmt.Errorf("LookupCropID: %w", err)
	}
	return cropID, nil
}

func (pg *Postgres) LookupMarketID(ctx context.Context, marketName, state string) (string, error) {
	var marketID string
	err := pg.pool.QueryRow(ctx, `
		SELECT id FROM markets
		WHERE LOWER(name) = LOWER($1) AND LOWER(state) = LOWER($2)
	`, marketName, state).Scan(&marketID)
	if err != nil {
		return "", fmt.Errorf("LookupMarketID: %w", err)
	}
	return marketID, nil
}

func (pg *Postgres) GetObservationsForAggregation(ctx context.Context, cropID, marketID string, periodStart, periodEnd time.Time) ([]models.PriceObservation, error) {
	rows, err := pg.pool.Query(ctx, `
		SELECT crop_id, market_id, observed_at, price, currency, unit, source, reporter_id, notes, confidence_score
		FROM price_observations
		WHERE crop_id = $1 AND market_id = $2 AND observed_at >= $3 AND observed_at <= $4
	`, cropID, marketID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("GetObservationsForAggregation: %w", err)
	}
	defer rows.Close()

	var observations []models.PriceObservation
	for rows.Next() {
		var obs models.PriceObservation
		if err := rows.Scan(&obs.CropID, &obs.MarketID, &obs.ObservedAt, &obs.Price,
			&obs.Currency, &obs.Unit, &obs.Source, &obs.ReporterID, &obs.Notes,
			&obs.ConfidenceScore); err != nil {
			return nil, fmt.Errorf("GetObservationsForAggregation scan: %w", err)
		}
		observations = append(observations, obs)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetObservationsForAggregation rows: %w", err)
	}

	return observations, nil
}

func (pg *Postgres) GetLatestAggregatedPrice(ctx context.Context, cropID, marketID, period string) (*models.AggregatedPrice, error) {
	var agg models.AggregatedPrice
	var confidenceScore float64

	err := pg.pool.QueryRow(ctx, `
		SELECT id, crop_id, market_id, period, period_start, period_end,
		       price_min, price_max, price_avg, price_median,
		       currency, unit, confidence_score, sample_size,
		       created_at, updated_at
		FROM aggregated_prices
		WHERE crop_id = $1 AND market_id = $2 AND period = $3
		ORDER BY period_end DESC
		LIMIT 1
	`, cropID, marketID, period).Scan(
		&agg.ID, &agg.CropID, &agg.MarketID, &agg.Period,
		&agg.PeriodStart, &agg.PeriodEnd, &agg.PriceMin,
		&agg.PriceMax, &agg.PriceMean, &agg.PriceMedian,
		&agg.Currency, &agg.Unit, &confidenceScore,
		&agg.SampleSize, &agg.CreatedAt, &agg.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // no aggregated price found, return nil without error
		}
		return nil, fmt.Errorf("GetLatestAggregatedPrice: %w", err)
	}
	agg.Confidence = models.ScoreToConfidenceLevel(confidenceScore)
	return &agg, nil
}

func (pg *Postgres) UpsertAggregatedPrice(ctx context.Context, agg models.AggregatedPrice) error {
	confidenceScore := agg.Confidence.ToScore()

	_, err := pg.pool.Exec(ctx, `
		INSERT INTO aggregated_prices (
			crop_id, market_id, period, period_start, period_end,
		    price_min, price_max, price_avg, price_median,
		    currency, unit, confidence_score, sample_size, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
		ON CONFLICT (crop_id, market_id, period, period_start, period_end)
		DO UPDATE SET 
			price_min = EXCLUDED.price_min,
		    price_max = EXCLUDED.price_max,
		    price_avg = EXCLUDED.price_avg,
		    price_median = EXCLUDED.price_median,
		    currency = EXCLUDED.currency,
		    unit = EXCLUDED.unit,
		    confidence_score = EXCLUDED.confidence_score,
		    sample_size = EXCLUDED.sample_size,
		    updated_at = NOW()
	`, agg.CropID, agg.MarketID, agg.Period, agg.PeriodStart, agg.PeriodEnd,
		agg.PriceMin, agg.PriceMax, agg.PriceMean, agg.PriceMedian,
		agg.Currency, agg.Unit, confidenceScore, agg.SampleSize,
	)
	if err != nil {
		return fmt.Errorf("UpsertAggregatedPrice: %w", err)
	}
	return nil
}

func (pg *Postgres) GetAllMarkets(ctx context.Context) ([]models.Market, error) {
	rows, err := pg.pool.Query(ctx, `
		SELECT id, name, state, country FROM markets
	`)
	if err != nil {
		return nil, fmt.Errorf("GetAllMarkets: %w", err)
	}
	defer rows.Close()

	var markets []models.Market
	for rows.Next() {
		var m models.Market
		if err := rows.Scan(&m.ID, &m.Name, &m.State, &m.Country); err != nil {
			return nil, fmt.Errorf("GetAllMarkets scan: %w", err)
		}
		markets = append(markets, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllMarkets rows: %w", err)
	}

	return markets, nil
}

func (pg *Postgres) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	rows, err := pg.pool.Query(ctx, `
		SELECT id, name, unit FROM crops
	`)
	if err != nil {
		return nil, fmt.Errorf("GetAllCrops: %w", err)
	}
	defer rows.Close()

	var crops []models.Crop
	for rows.Next() {
		var c models.Crop
		if err := rows.Scan(&c.ID, &c.Name, &c.Unit); err != nil {
			return nil, fmt.Errorf("GetAllCrops scan: %w", err)
		}
		crops = append(crops, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllCrops rows: %w", err)
	}

	return crops, nil
}

func (pg *Postgres) GetCropMarketCombinations(ctx context.Context, from time.Time, to time.Time) ([]models.CropMarketCombination, error) {
	rows, err := pg.pool.Query(ctx, `
		SELECT DISTINCT crop_id, market_id
		FROM price_observations
		WHERE observed_at >= $1 AND observed_at <= $2
	`, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetCropMarketCombinations: %w", err)
	}
	defer rows.Close()

	var combinations []models.CropMarketCombination
	for rows.Next() {
		var cm models.CropMarketCombination
		if err := rows.Scan(&cm.CropID, &cm.MarketID); err != nil {
			return nil, fmt.Errorf("GetCropMarketCombinations scan: %w", err)
		}
		combinations = append(combinations, cm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCropMarketCombinations rows: %w", err)
	}

	return combinations, nil
}
