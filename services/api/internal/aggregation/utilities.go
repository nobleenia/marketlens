package aggregation

import (
	"math"
	"sort"
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
