package aggregation

import (
	"marketlens/internal/models"
	"math"
	"sort"
	"time"
)

// CalculateMean calculates the arithmetic mean of a slice of prices.
func CalculateMean(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}

	var sum float64
	for _, price := range prices {
		sum += price
	}
	return sum / float64(n)
}

// CalculateMedian calculates the median of a slice of prices.
// It does NOT mutate the input slice.
func CalculateMedian(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}

	// Create a copy to avoid mutating the input
	sorted := make([]float64, n)
	copy(sorted, prices)
	sort.Float64s(sorted)

	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// CalculateVariance calculates the population variance of a slice of prices.
// Variance = Σ(x - mean)² / n
func CalculateVariance(prices []float64) float64 {
	n := len(prices)
	if n <= 1 {
		return 0
	}

	mean := CalculateMean(prices)

	var sumSquares float64
	for _, price := range prices {
		diff := price - mean
		sumSquares += diff * diff
	}
	return sumSquares / float64(n)
}

// CalculateStandardDeviation calculates the standard deviation (sqrt of variance).
func CalculateStandardDeviation(prices []float64) float64 {
	variance := CalculateVariance(prices)
	return math.Sqrt(variance)
}

// CalculateCoefficientOfVariation calculates CV as a percentage.
// CV = (stdDev / mean) * 100
func CalculateCoefficientOfVariation(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}

	mean := CalculateMean(prices)
	if mean == 0 {
		return 0
	}

	stdDev := CalculateStandardDeviation(prices)
	return (stdDev / mean) * 100
}

// findMinAndMax returns the minimum and maximum values in a slice.
func findMinAndMax(prices []float64) (min float64, max float64) {
	if len(prices) == 0 {
		return 0, 0
	}

	min = prices[0]
	max = prices[0]

	for _, price := range prices {
		if price < min {
			min = price
		}
		if price > max {
			max = price
		}
	}
	return min, max
}

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
