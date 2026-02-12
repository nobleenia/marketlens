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
		(crop_id, market_id, observed_at, price, currency, unit, price_type, source, reporter_id, notes, confidence_score)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, obs.CropID, obs.MarketID, obs.ObservedAt, obs.Price, obs.Currency, obs.Unit,
		obs.PriceType, obs.Source, obs.ReporterID, obs.Notes, obs.ConfidenceScore)
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
		SELECT id, name, state, country, COALESCE(latitude, 0), COALESCE(longitude, 0) FROM markets
	`)
	if err != nil {
		return nil, fmt.Errorf("GetAllMarkets: %w", err)
	}
	defer rows.Close()

	var markets []models.Market
	for rows.Next() {
		var m models.Market
		if err := rows.Scan(&m.ID, &m.Name, &m.State, &m.Country, &m.Latitude, &m.Longitude); err != nil {
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

// LookupMarketIDByName finds a market by name alone (case-insensitive).
// Use this when state is not available (e.g., URL path).
func (pg *Postgres) LookupMarketIDByName(ctx context.Context, marketName string) (string, error) {
	var marketID string
	err := pg.pool.QueryRow(ctx, `
        SELECT id FROM markets
        WHERE LOWER(name) = LOWER($1)
        LIMIT 1
    `, marketName).Scan(&marketID)
	if err != nil {
		return "", fmt.Errorf("LookupMarketIDByName: %w", err)
	}
	return marketID, nil
}

// GetAggregatedPrices returns aggregated prices with optional filters.
// Pass empty strings / zero time to omit a filter.
func (pg *Postgres) GetAggregatedPrices(ctx context.Context, cropName, marketName string, from, to time.Time) ([]models.AggregatedPrice, error) {
	query := `
        SELECT ap.id, ap.crop_id, c.name, ap.market_id, m.name,
               ap.period, ap.period_start, ap.period_end,
               ap.price_min, ap.price_max, ap.price_avg, ap.price_median,
               ap.currency, ap.unit, ap.confidence_score, ap.sample_size,
               ap.created_at, ap.updated_at
        FROM aggregated_prices ap
        JOIN crops c ON c.id = ap.crop_id
        JOIN markets m ON m.id = ap.market_id
        WHERE 1=1
    `
	args := []interface{}{}
	argIdx := 1

	if cropName != "" {
		query += fmt.Sprintf(" AND LOWER(c.name) = LOWER($%d)", argIdx)
		args = append(args, cropName)
		argIdx++
	}
	if marketName != "" {
		query += fmt.Sprintf(" AND LOWER(m.name) = LOWER($%d)", argIdx)
		args = append(args, marketName)
		argIdx++
	}
	if !from.IsZero() {
		query += fmt.Sprintf(" AND ap.period_start >= $%d", argIdx)
		args = append(args, from)
		argIdx++
	}
	if !to.IsZero() {
		query += fmt.Sprintf(" AND ap.period_end <= $%d", argIdx)
		args = append(args, to)
		argIdx++
	}

	query += " ORDER BY ap.period_end DESC"

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAggregatedPrices: %w", err)
	}
	defer rows.Close()

	results := make([]models.AggregatedPrice, 0)
	for rows.Next() {
		var agg models.AggregatedPrice
		var confidenceScore float64

		if err := rows.Scan(
			&agg.ID, &agg.CropID, &agg.CropName, &agg.MarketID, &agg.MarketName,
			&agg.Period, &agg.PeriodStart, &agg.PeriodEnd,
			&agg.PriceMin, &agg.PriceMax, &agg.PriceMean, &agg.PriceMedian,
			&agg.Currency, &agg.Unit, &confidenceScore,
			&agg.SampleSize, &agg.CreatedAt, &agg.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetAggregatedPrices scan: %w", err)
		}
		agg.Confidence = models.ScoreToConfidenceLevel(confidenceScore)
		results = append(results, agg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAggregatedPrices rows: %w", err)
	}

	return results, nil
}

// GetAggregatedPriceByCropAndMarket returns the single latest aggregated price
// for a crop+market pair (looked up by name, case-insensitive).
func (pg *Postgres) GetAggregatedPriceByCropAndMarket(ctx context.Context, cropName, marketName string) (*models.AggregatedPrice, error) {
	var agg models.AggregatedPrice
	var confidenceScore float64

	err := pg.pool.QueryRow(ctx, `
        SELECT ap.id, ap.crop_id, c.name, ap.market_id, m.name,
               ap.period, ap.period_start, ap.period_end,
               ap.price_min, ap.price_max, ap.price_avg, ap.price_median,
               ap.currency, ap.unit, ap.confidence_score, ap.sample_size,
               ap.created_at, ap.updated_at
        FROM aggregated_prices ap
        JOIN crops c ON c.id = ap.crop_id
        JOIN markets m ON m.id = ap.market_id
        WHERE LOWER(c.name) = LOWER($1) AND LOWER(m.name) = LOWER($2)
        ORDER BY ap.period_end DESC
        LIMIT 1
    `, cropName, marketName).Scan(
		&agg.ID, &agg.CropID, &agg.CropName, &agg.MarketID, &agg.MarketName,
		&agg.Period, &agg.PeriodStart, &agg.PeriodEnd,
		&agg.PriceMin, &agg.PriceMax, &agg.PriceMean, &agg.PriceMedian,
		&agg.Currency, &agg.Unit, &confidenceScore,
		&agg.SampleSize, &agg.CreatedAt, &agg.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetAggregatedPriceByCropAndMarket: %w", err)
	}
	agg.Confidence = models.ScoreToConfidenceLevel(confidenceScore)
	return &agg, nil
}

// InsertCrop creates a new crop and returns its ID.
func (pg *Postgres) InsertCrop(ctx context.Context, name, unit string) (string, error) {
	var id string
	err := pg.pool.QueryRow(ctx, `
        INSERT INTO crops (name, unit)
        VALUES ($1, $2)
        ON CONFLICT (name) DO UPDATE SET unit = EXCLUDED.unit
        RETURNING id
    `, name, unit).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("InsertCrop: %w", err)
	}
	return id, nil
}

// GetDistinctStates returns all unique state names from the markets table, sorted alphabetically.
func (pg *Postgres) GetDistinctStates(ctx context.Context) ([]string, error) {
	rows, err := pg.pool.Query(ctx, `
        SELECT DISTINCT state FROM markets ORDER BY state
    `)
	if err != nil {
		return nil, fmt.Errorf("GetDistinctStates: %w", err)
	}
	defer rows.Close()

	var states []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("GetDistinctStates scan: %w", err)
		}
		states = append(states, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetDistinctStates rows: %w", err)
	}
	return states, nil
}

// GetMarketsByState returns all markets in a given state (case-insensitive).
func (pg *Postgres) GetMarketsByState(ctx context.Context, state string) ([]models.Market, error) {
	rows, err := pg.pool.Query(ctx, `
        SELECT id, name, state, country, COALESCE(latitude, 0), COALESCE(longitude, 0)
        FROM markets
        WHERE LOWER(state) = LOWER($1)
        ORDER BY name
    `, state)
	if err != nil {
		return nil, fmt.Errorf("GetMarketsByState: %w", err)
	}
	defer rows.Close()

	var markets []models.Market
	for rows.Next() {
		var m models.Market
		if err := rows.Scan(&m.ID, &m.Name, &m.State, &m.Country, &m.Latitude, &m.Longitude); err != nil {
			return nil, fmt.Errorf("GetMarketsByState scan: %w", err)
		}
		markets = append(markets, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetMarketsByState rows: %w", err)
	}
	return markets, nil
}

// ── Admin: Observations listing ────────────────────────────────────

// ListObservations returns price_observations with optional filters and pagination.
func (pg *Postgres) ListObservations(ctx context.Context, cropName, marketName, status string, limit, offset int) ([]models.PriceObservation, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if cropName != "" {
		where += fmt.Sprintf(" AND LOWER(c.name) = LOWER($%d)", argIdx)
		args = append(args, cropName)
		argIdx++
	}
	if marketName != "" {
		where += fmt.Sprintf(" AND LOWER(m.name) = LOWER($%d)", argIdx)
		args = append(args, marketName)
		argIdx++
	}
	if status != "" {
		where += fmt.Sprintf(" AND po.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count total
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*) FROM price_observations po
        JOIN crops c ON c.id = po.crop_id
        JOIN markets m ON m.id = po.market_id
        %s`, where)

	var total int
	if err := pg.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ListObservations count: %w", err)
	}

	// Fetch page
	if limit <= 0 {
		limit = 50
	}
	dataQuery := fmt.Sprintf(`
        SELECT po.id, po.crop_id, c.name, po.market_id, m.name,
               po.observed_at, po.price, po.currency, po.unit,
               po.source, COALESCE(po.reporter_id, ''), COALESCE(po.notes, ''),
               po.confidence_score, po.status, po.created_at
        FROM price_observations po
        JOIN crops c ON c.id = po.crop_id
        JOIN markets m ON m.id = po.market_id
        %s
        ORDER BY po.created_at DESC
        LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := pg.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ListObservations: %w", err)
	}
	defer rows.Close()

	var results []models.PriceObservation
	for rows.Next() {
		var o models.PriceObservation
		var cropName, marketName string
		if err := rows.Scan(
			&o.ID, &o.CropID, &cropName, &o.MarketID, &marketName,
			&o.ObservedAt, &o.Price, &o.Currency, &o.Unit,
			&o.Source, &o.ReporterID, &o.Notes,
			&o.ConfidenceScore, &o.Status, &o.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("ListObservations scan: %w", err)
		}
		// Attach names for JSON response convenience
		o.CropName = cropName
		o.MarketName = marketName
		results = append(results, o)
	}
	return results, total, rows.Err()
}

// UpdateObservationStatus sets the status of a price_observation.
func (pg *Postgres) UpdateObservationStatus(ctx context.Context, id, status string) error {
	tag, err := pg.pool.Exec(ctx, `
        UPDATE price_observations SET status = $1 WHERE id = $2
    `, status, id)
	if err != nil {
		return fmt.Errorf("UpdateObservationStatus: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateObservationStatus: observation not found")
	}
	return nil
}

// GetObservationByID returns a single price_observation by ID.
func (pg *Postgres) GetObservationByID(ctx context.Context, id string) (*models.PriceObservation, error) {
	var o models.PriceObservation
	var cropName, marketName string
	err := pg.pool.QueryRow(ctx, `
        SELECT po.id, po.crop_id, c.name, po.market_id, m.name,
               po.observed_at, po.price, po.currency, po.unit,
               po.source, COALESCE(po.reporter_id, ''), COALESCE(po.notes, ''),
               po.confidence_score, po.status, po.created_at
        FROM price_observations po
        JOIN crops c ON c.id = po.crop_id
        JOIN markets m ON m.id = po.market_id
        WHERE po.id = $1
    `, id).Scan(
		&o.ID, &o.CropID, &cropName, &o.MarketID, &marketName,
		&o.ObservedAt, &o.Price, &o.Currency, &o.Unit,
		&o.Source, &o.ReporterID, &o.Notes,
		&o.ConfidenceScore, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetObservationByID: %w", err)
	}
	o.CropName = cropName
	o.MarketName = marketName
	return &o, nil
}

// ── Admin: Audit logs ──────────────────────────────────────────────

// InsertAuditLog records an admin action.
func (pg *Postgres) InsertAuditLog(ctx context.Context, log models.AuditLog) error {
	_, err := pg.pool.Exec(ctx, `
        INSERT INTO audit_logs (admin_id, action, entity_type, entity_id, old_value, new_value, reason)
        VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7)
    `, log.AdminID, log.Action, log.EntityType, log.EntityID, log.OldValue, log.NewValue, log.Reason)
	if err != nil {
		return fmt.Errorf("InsertAuditLog: %w", err)
	}
	return nil
}

// ListAuditLogs returns audit logs with optional entity filter and pagination.
func (pg *Postgres) ListAuditLogs(ctx context.Context, entityType, entityID string, limit, offset int) ([]models.AuditLog, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if entityType != "" {
		where += fmt.Sprintf(" AND entity_type = $%d", argIdx)
		args = append(args, entityType)
		argIdx++
	}
	if entityID != "" {
		where += fmt.Sprintf(" AND entity_id = $%d", argIdx)
		args = append(args, entityID)
		argIdx++
	}

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", where)
	if err := pg.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ListAuditLogs count: %w", err)
	}

	if limit <= 0 {
		limit = 50
	}
	dataQ := fmt.Sprintf(`
        SELECT id, admin_id, action, entity_type, entity_id,
               COALESCE(old_value::text, ''), COALESCE(new_value::text, ''),
               reason, created_at
        FROM audit_logs %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := pg.pool.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ListAuditLogs: %w", err)
	}
	defer rows.Close()

	var results []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(
			&l.ID, &l.AdminID, &l.Action, &l.EntityType, &l.EntityID,
			&l.OldValue, &l.NewValue, &l.Reason, &l.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("ListAuditLogs scan: %w", err)
		}
		results = append(results, l)
	}
	return results, total, rows.Err()
}
