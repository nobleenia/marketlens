package aggregation

import (
	"context"
	"fmt"
	"log"
	"marketlens/internal/models"
	"marketlens/internal/store"
	"time"
)

// Confidence thresholds (as percentages)
const (
	CVThresholdHigh   = 15.0 // CV must be <= 15% for high confidence
	CVThresholdMedium = 25.0 // CV must be <= 25% for medium confidence
	MinSamplesHigh    = 5    // Need at least 5 samples for high confidence
	MinSamplesMedium  = 2    // Need at least 2 samples for medium confidence
)

// DetermineConfidenceLevel determines confidence based on sample size and CV.
// CV is expected as a percentage (e.g., 10.0 means 10%).
func DetermineConfidenceLevel(sampleSize int, coefficientOfVariation float64) models.ConfidenceLevel {
	if sampleSize < MinSamplesMedium || coefficientOfVariation > CVThresholdMedium {
		return models.ConfidenceLow
	}
	if sampleSize >= MinSamplesHigh && coefficientOfVariation <= CVThresholdHigh {
		return models.ConfidenceHigh
	}
	return models.ConfidenceMedium
}

// AggregatePrices aggregates a slice of price observations into a single result.
// Returns nil if there are no observations.
func AggregatePrices(observations []models.PriceObservation, cropID, marketID, period, currency, unit string, periodStart, periodEnd time.Time) *models.AggregatedPrice {
	if len(observations) == 0 {
		return nil
	}

	// Extract price values
	prices := make([]float64, len(observations))
	for i, obs := range observations {
		prices[i] = obs.Price
	}

	// Calculate statistics
	mean := CalculateMean(prices)
	cv := CalculateCoefficientOfVariation(prices)
	median := CalculateMedian(prices)
	minPrice, maxPrice := findMinAndMax(prices)

	// Determine confidence
	confidence := DetermineConfidenceLevel(len(prices), cv)

	return &models.AggregatedPrice{
		CropID:      cropID,
		MarketID:    marketID,
		Period:      period,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		PriceMin:    minPrice,
		PriceMax:    maxPrice,
		PriceMean:   mean,
		PriceMedian: median,
		Currency:    currency,
		Unit:        unit,
		Confidence:  confidence,
		SampleSize:  len(prices),
	}
}

// CalculateTrend compares current price to previous price.
// Returns TrendUp if increase > threshold%, TrendDown if decrease > threshold%, else TrendStable.
func CalculateTrend(currentPrice, previousPrice float64, thresholdPercent float64) models.Trend {
	if previousPrice == 0 {
		return models.TrendStable
	}

	changePercent := ((currentPrice - previousPrice) / previousPrice) * 100

	if changePercent >= thresholdPercent {
		return models.TrendUp
	} else if changePercent <= -thresholdPercent {
		return models.TrendDown
	}
	return models.TrendStable
}

type Service struct {
	store *store.Postgres
}

func NewService(store *store.Postgres) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) RunDailyAggregation(ctx context.Context) error {
	periodStart := time.Now().Add(-24 * time.Hour)
	periodEnd := time.Now()

	// 1. Get all unique crop-market combinations that have observations in the last day
	cropMarkets, err := s.store.GetCropMarketCombinations(ctx, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("failed to get crop-market combinations: %w", err)
	}

	for _, cm := range cropMarkets {
		// 2. For each combination, get observations for the last day
		observations, err := s.store.GetObservationsForAggregation(ctx, cm.CropID, cm.MarketID, periodStart, periodEnd)
		if err != nil {
			log.Printf("failed to get observations for crop %s and market %s: %v", cm.CropID, cm.MarketID, err)
			continue
		}

		if len(observations) == 0 {
			continue // no observations to aggregate
		}

		// 3. Aggregate the prices
		agg := AggregatePrices(observations, cm.CropID, cm.MarketID, "daily", observations[0].Currency, observations[0].Unit, periodStart, periodEnd)
		if agg == nil {
			continue // no valid aggregation
		}

		// 4. Get the latest aggregated price for trend calculation
		prevAgg, err := s.store.GetLatestAggregatedPrice(ctx, cm.CropID, cm.MarketID, "daily")
		if err != nil {
			return fmt.Errorf("failed to get latest aggregated price for crop %s and market %s: %w", cm.CropID, cm.MarketID, err)
		}

		// 5. Calculate trend if previous aggregation exists
		if prevAgg != nil {
			agg.Trend = CalculateTrend(agg.PriceMean, prevAgg.PriceMean, 5.0) // using 5% threshold for trend
		} else {
			agg.Trend = models.TrendStable
		}

		// 6. Upsert the aggregated price into the database
		err = s.store.UpsertAggregatedPrice(ctx, *agg)
		if err != nil {
			return fmt.Errorf("failed to upsert aggregated price for crop %s and market %s: %w", cm.CropID, cm.MarketID, err)
		}
	}

	return nil
}
