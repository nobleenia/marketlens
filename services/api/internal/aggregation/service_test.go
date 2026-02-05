package aggregation

import (
	"marketlens/internal/models"
	"math"
	"testing"
	"time"
)

// Helper function to compare floats with tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// ===========================================
// Task 1.1.2: CalculateMedian
// ===========================================

func TestCalculateMedian_EmptySlice(t *testing.T) {
	prices := []float64{}

	got := CalculateMedian(prices)

	if got != 0 {
		t.Errorf("CalculateMedian(%v) = %v, want 0", prices, got)
	}
}

func TestCalculateMedian_SingleElement(t *testing.T) {
	prices := []float64{5.0}

	got := CalculateMedian(prices)

	if got != 5.0 {
		t.Errorf("CalculateMedian(%v) = %v, want 5.0", prices, got)
	}
}

func TestCalculateMedian_OddLength(t *testing.T) {
	prices := []float64{3.0, 1.0, 5.0} // Unsorted deliberately

	got := CalculateMedian(prices)

	if got != 3.0 {
		t.Errorf("CalculateMedian(%v) = %v, want 3.0 (middle after sorting)", prices, got)
	}
}

func TestCalculateMedian_EvenLength(t *testing.T) {
	prices := []float64{1.0, 2.0, 3.0, 4.0}

	got := CalculateMedian(prices)

	if got != 2.5 {
		t.Errorf("CalculateMedian(%v) = %v, want 2.5", prices, got)
	}
}

func TestCalculateMedian_EvenLengthUnsorted(t *testing.T) {
	prices := []float64{4.0, 1.0, 3.0, 2.0}

	got := CalculateMedian(prices)

	if got != 2.5 {
		t.Errorf("CalculateMedian(%v) = %v, want 2.5", prices, got)
	}
}

func TestCalculateMedian_DoesNotMutateInput(t *testing.T) {
	prices := []float64{3.0, 1.0, 2.0}
	original := make([]float64, len(prices))
	copy(original, prices)

	_ = CalculateMedian(prices)

	for i, v := range prices {
		if v != original[i] {
			t.Errorf("CalculateMedian mutated input: got %v, original was %v", prices, original)
			break
		}
	}
}

func TestCalculateMedian_LargerDataset(t *testing.T) {
	prices := []float64{15000, 18000, 17000, 16500, 19000, 17500, 16000}

	got := CalculateMedian(prices)

	if got != 17000 {
		t.Errorf("CalculateMedian(%v) = %v, want 17000", prices, got)
	}
}

// ===========================================
// Task 1.1.3: CalculateVariance
// ===========================================

func TestCalculateVariance_EmptySlice(t *testing.T) {
	prices := []float64{}

	got := CalculateVariance(prices)

	if got != 0 {
		t.Errorf("CalculateVariance(%v) = %v, want 0", prices, got)
	}
}

func TestCalculateVariance_SingleElement(t *testing.T) {
	prices := []float64{5.0}

	got := CalculateVariance(prices)

	if got != 0 {
		t.Errorf("CalculateVariance(%v) = %v, want 0 (no variance with one element)", prices, got)
	}
}

func TestCalculateVariance_UniformValues(t *testing.T) {
	prices := []float64{5.0, 5.0, 5.0}

	got := CalculateVariance(prices)

	if got != 0 {
		t.Errorf("CalculateVariance(%v) = %v, want 0 (no variance when all same)", prices, got)
	}
}

func TestCalculateVariance_KnownValues(t *testing.T) {
	prices := []float64{2.0, 4.0, 6.0}

	got := CalculateVariance(prices)

	expected := 8.0 / 3.0

	if !almostEqual(got, expected, 0.0001) {
		t.Errorf("CalculateVariance(%v) = %v, want %v", prices, got, expected)
	}
}

func TestCalculateVariance_AnotherKnownValues(t *testing.T) {
	prices := []float64{10.0, 20.0, 30.0, 40.0, 50.0}

	got := CalculateVariance(prices)

	expected := 200.0

	if !almostEqual(got, expected, 0.0001) {
		t.Errorf("CalculateVariance(%v) = %v, want %v", prices, got, expected)
	}
}

// ===========================================
// Task 1.1.4: CalculateCoefficientOfVariation
// ===========================================

func TestCalculateCoefficientOfVariation_EmptySlice(t *testing.T) {
	prices := []float64{}

	got := CalculateCoefficientOfVariation(prices)

	if got != 0 {
		t.Errorf("CalculateCoefficientOfVariation(%v) = %v, want 0", prices, got)
	}
}

func TestCalculateCoefficientOfVariation_ZeroMean(t *testing.T) {
	prices := []float64{0.0, 0.0, 0.0}

	got := CalculateCoefficientOfVariation(prices)

	if got != 0 {
		t.Errorf("CalculateCoefficientOfVariation(%v) = %v, want 0 (handle zero mean)", prices, got)
	}
}

func TestCalculateCoefficientOfVariation_UniformValues(t *testing.T) {
	prices := []float64{100.0, 100.0, 100.0}

	got := CalculateCoefficientOfVariation(prices)

	if got != 0 {
		t.Errorf("CalculateCoefficientOfVariation(%v) = %v, want 0", prices, got)
	}
}

func TestCalculateCoefficientOfVariation_KnownValues(t *testing.T) {
	prices := []float64{10.0, 20.0, 30.0}

	got := CalculateCoefficientOfVariation(prices)

	// Mean: 20.0
	// Variance: 66.666...
	// StdDev: 8.165...
	// CV: (8.165 / 20) * 100 = 40.825%
	expected := 40.8248

	if !almostEqual(got, expected, 0.01) {
		t.Errorf("CalculateCoefficientOfVariation(%v) = %v, want approximately %v", prices, got, expected)
	}
}

func TestCalculateCoefficientOfVariation_LowVariation(t *testing.T) {
	prices := []float64{100.0, 102.0, 98.0, 101.0, 99.0}

	got := CalculateCoefficientOfVariation(prices)

	if got >= 5.0 {
		t.Errorf("CalculateCoefficientOfVariation(%v) = %v, expected low CV (< 5%%)", prices, got)
	}
}

// ===========================================
// Task 1.1.5: DetermineConfidenceLevel
// ===========================================

func TestDetermineConfidence_LowSampleSize(t *testing.T) {
	sampleSize := 1
	cv := 10.0

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceLow {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v", sampleSize, cv, got, models.ConfidenceLow)
	}
}

func TestDetermineConfidence_ZeroSampleSize(t *testing.T) {
	sampleSize := 0
	cv := 5.0

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceLow {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v", sampleSize, cv, got, models.ConfidenceLow)
	}
}

func TestDetermineConfidence_HighVariance(t *testing.T) {
	sampleSize := 10
	cv := 30.0

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceLow {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v (high variance should be low confidence)", sampleSize, cv, got, models.ConfidenceLow)
	}
}

func TestDetermineConfidence_MediumConfidence(t *testing.T) {
	sampleSize := 3
	cv := 20.0 // Within medium threshold, but not enough samples for high

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceMedium {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v", sampleSize, cv, got, models.ConfidenceMedium)
	}
}

func TestDetermineConfidence_HighConfidence(t *testing.T) {
	sampleSize := 5
	cv := 10.0

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceHigh {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v", sampleSize, cv, got, models.ConfidenceHigh)
	}
}

func TestDetermineConfidence_HighSampleLowVariance(t *testing.T) {
	sampleSize := 10
	cv := 5.0

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceHigh {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v", sampleSize, cv, got, models.ConfidenceHigh)
	}
}

func TestDetermineConfidence_EdgeCaseMediumSampleMediumVariance(t *testing.T) {
	sampleSize := 4
	cv := 20.0

	got := DetermineConfidenceLevel(sampleSize, cv)

	if got != models.ConfidenceMedium {
		t.Errorf("DetermineConfidenceLevel(%d, %v) = %v, want %v", sampleSize, cv, got, models.ConfidenceMedium)
	}
}

// ===========================================
// Task 1.1.6: AggregatePrices
// ===========================================

func TestAggregatePrice_EmptyObservations(t *testing.T) {
	observations := []models.PriceObservation{}
	periodStart := time.Now().Add(-24 * time.Hour)
	periodEnd := time.Now()

	got := AggregatePrices(observations, "", "", "daily", "NGN", "kg", periodStart, periodEnd)

	if got != nil {
		t.Errorf("AggregatePrices with empty observations = %v, want nil", got)
	}
}

func TestAggregatePrice_SingleObservation(t *testing.T) {
	cropID := "crop-123"
	marketID := "market-456"
	observations := []models.PriceObservation{
		{
			CropID:   cropID,
			MarketID: marketID,
			Price:    100.0,
		},
	}
	periodStart := time.Now().Add(-24 * time.Hour)
	periodEnd := time.Now()

	got := AggregatePrices(observations, cropID, marketID, "daily", "NGN", "kg", periodStart, periodEnd)

	if got == nil {
		t.Fatal("AggregatePrices with single observation returned nil")
	}

	if got.CropID != cropID {
		t.Errorf("CropID = %v, want %v", got.CropID, cropID)
	}
	if got.MarketID != marketID {
		t.Errorf("MarketID = %v, want %v", got.MarketID, marketID)
	}
	if got.PriceMin != 100.0 {
		t.Errorf("PriceMin = %v, want 100.0", got.PriceMin)
	}
	if got.PriceMax != 100.0 {
		t.Errorf("PriceMax = %v, want 100.0", got.PriceMax)
	}
	if got.PriceMean != 100.0 {
		t.Errorf("PriceMean = %v, want 100.0", got.PriceMean)
	}
	if got.PriceMedian != 100.0 {
		t.Errorf("PriceMedian = %v, want 100.0", got.PriceMedian)
	}
	if got.SampleSize != 1 {
		t.Errorf("SampleSize = %v, want 1", got.SampleSize)
	}
	if got.Confidence != models.ConfidenceLow {
		t.Errorf("Confidence = %v, want %v (single observation)", got.Confidence, models.ConfidenceLow)
	}
}

func TestAggregatePrice_MultipleObservations(t *testing.T) {
	cropID := "crop-123"
	marketID := "market-456"
	observations := []models.PriceObservation{
		{CropID: cropID, MarketID: marketID, Price: 100.0},
		{CropID: cropID, MarketID: marketID, Price: 120.0},
		{CropID: cropID, MarketID: marketID, Price: 110.0},
		{CropID: cropID, MarketID: marketID, Price: 115.0},
		{CropID: cropID, MarketID: marketID, Price: 105.0},
	}
	periodStart := time.Now().Add(-24 * time.Hour)
	periodEnd := time.Now()

	got := AggregatePrices(observations, cropID, marketID, "daily", "NGN", "kg", periodStart, periodEnd)

	if got == nil {
		t.Fatal("AggregatePrices with multiple observations returned nil")
	}

	if got.PriceMin != 100.0 {
		t.Errorf("PriceMin = %v, want 100.0", got.PriceMin)
	}
	if got.PriceMax != 120.0 {
		t.Errorf("PriceMax = %v, want 120.0", got.PriceMax)
	}
	if got.PriceMean != 110.0 {
		t.Errorf("PriceMean = %v, want 110.0", got.PriceMean)
	}
	if got.PriceMedian != 110.0 {
		t.Errorf("PriceMedian = %v, want 110.0", got.PriceMedian)
	}
	if got.SampleSize != 5 {
		t.Errorf("SampleSize = %v, want 5", got.SampleSize)
	}
	if got.Confidence == models.ConfidenceLow {
		t.Errorf("Confidence = %v, expected medium or high for 5 samples with low variance", got.Confidence)
	}
}

func TestAggregatePrice_SetsCorrectPeriod(t *testing.T) {
	observations := []models.PriceObservation{
		{CropID: "crop-1", MarketID: "market-1", Price: 100.0},
	}
	periodStart := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	got := AggregatePrices(observations, "crop-1", "market-1", "daily", "NGN", "kg", periodStart, periodEnd)

	if got == nil {
		t.Fatal("AggregatePrices returned nil")
	}

	if !got.PeriodStart.Equal(periodStart) {
		t.Errorf("PeriodStart = %v, want %v", got.PeriodStart, periodStart)
	}
	if !got.PeriodEnd.Equal(periodEnd) {
		t.Errorf("PeriodEnd = %v, want %v", got.PeriodEnd, periodEnd)
	}
}

func TestAggregatePrice_HighVarianceData(t *testing.T) {
	cropID := "crop-123"
	marketID := "market-456"
	observations := []models.PriceObservation{
		{CropID: cropID, MarketID: marketID, Price: 50.0},
		{CropID: cropID, MarketID: marketID, Price: 150.0},
		{CropID: cropID, MarketID: marketID, Price: 200.0},
	}
	periodStart := time.Now().Add(-24 * time.Hour)
	periodEnd := time.Now()

	got := AggregatePrices(observations, cropID, marketID, "daily", "NGN", "kg", periodStart, periodEnd)

	if got == nil {
		t.Fatal("AggregatePrices returned nil")
	}

	if got.Confidence == models.ConfidenceHigh {
		t.Errorf("Confidence = %v, expected low or medium for high variance data", got.Confidence)
	}
}

// ===========================================
// Task 1.2.3: CalculateTrend
// ===========================================

func TestCalculateTrend_NoPreviousPrice(t *testing.T) {
	got := CalculateTrend(100.0, 0, 5.0)

	if got != models.TrendStable {
		t.Errorf("CalculateTrend(100, 0, 5) = %v, want %v", got, models.TrendStable)
	}
}

func TestCalculateTrend_PriceIncrease(t *testing.T) {
	// 10% increase, threshold 5%
	got := CalculateTrend(110.0, 100.0, 5.0)

	if got != models.TrendUp {
		t.Errorf("CalculateTrend(110, 100, 5) = %v, want %v", got, models.TrendUp)
	}
}

func TestCalculateTrend_PriceDecrease(t *testing.T) {
	// 10% decrease, threshold 5%
	got := CalculateTrend(90.0, 100.0, 5.0)

	if got != models.TrendDown {
		t.Errorf("CalculateTrend(90, 100, 5) = %v, want %v", got, models.TrendDown)
	}
}

func TestCalculateTrend_SmallChange(t *testing.T) {
	// 2% increase, threshold 5%
	got := CalculateTrend(102.0, 100.0, 5.0)

	if got != models.TrendStable {
		t.Errorf("CalculateTrend(102, 100, 5) = %v, want %v (small change)", got, models.TrendStable)
	}
}

func TestCalculateTrend_ExactThreshold(t *testing.T) {
	// Exactly 5% increase, threshold 5%
	got := CalculateTrend(105.0, 100.0, 5.0)

	if got != models.TrendUp {
		t.Errorf("CalculateTrend(105, 100, 5) = %v, want %v (at threshold)", got, models.TrendUp)
	}
}
