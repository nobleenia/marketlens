package store

import (
	"context"
	"fmt"
	"time"

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

type PriceObservation struct {
	CropID          string
	MarketID        string
	ObservedAt      time.Time
	Price           float64
	Currency        string
	Unit            string
	Source          string
	ReporterID      string
	Notes           string
	ConfidenceScore float64
}

func (pg *Postgres) InsertPriceObservation(ctx context.Context, obs PriceObservation) error {
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

func (pg *Postgres) GetObservationsForAggregation(ctx context.Context, cropID, marketID string, from time.Time, to time.Time) ([]PriceObservation, error) {
	rows, err := pg.pool.Query(ctx, `
		SELECT crop_id, market_id, observed_at, price, currency, unit, source, reporter_id, notes, confidence_score
		FROM price_observations
		WHERE crop_id = $1 AND market_id = $2 AND observed_at >= $3 AND observed_at <= $4
	`, cropID, marketID, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetObservationsForAggregation: %w", err)
	}
	defer rows.Close()

	var observations []PriceObservation
	for rows.Next() {
		var obs PriceObservation
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

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.pool.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.pool.Close()
}
