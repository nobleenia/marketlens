package models

import "time"

type ConfidenceLevel string

const (
	ConfidenceLow    ConfidenceLevel = "low"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceHigh   ConfidenceLevel = "high"
)

type Trend string

const (
	TrendUp     Trend = "up"
	TrendDown   Trend = "down"
	TrendStable Trend = "stable"
)

type PriceObservation struct {
	ID              string    `json:"id"`
	CropID          string    `json:"crop_id"`
	MarketID        string    `json:"market_id"`
	ObservedAt      time.Time `json:"observed_at"`
	Price           float64   `json:"price"`
	Currency        string    `json:"currency"`
	Unit            string    `json:"unit"`
	PriceType       string    `json:"price_type"`
	Source          string    `json:"source"`
	ReporterID      string    `json:"reporter_id"`
	Notes           string    `json:"notes"`
	ConfidenceScore float64   `json:"confidence_score"`
	CreatedAt       time.Time `json:"created_at"`
}

type AggregatedPrice struct {
	ID          string
	CropID      string
	CropName    string
	MarketID    string
	MarketName  string
	Period      string // "daily", "weekly", "monthly"
	PeriodStart time.Time
	PeriodEnd   time.Time
	PriceMin    float64
	PriceMax    float64
	PriceMean   float64
	PriceMedian float64
	Currency    string
	Unit        string
	Confidence  ConfidenceLevel
	SampleSize  int
	Trend       Trend
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c ConfidenceLevel) ToScore() float64 {
	switch c {
	case ConfidenceHigh:
		return 0.9
	case ConfidenceMedium:
		return 0.6
	case ConfidenceLow:
		return 0.3
	default:
		return 0.5
	}
}

func ScoreToConfidenceLevel(score float64) ConfidenceLevel {
	if score >= 0.8 {
		return ConfidenceHigh
	} else if score >= 0.5 {
		return ConfidenceMedium
	} else {
		return ConfidenceLow
	}
}

type Market struct {
	ID      string
	Name    string
	State   string
	Country string
}

type Crop struct {
	ID   string
	Name string
	Unit string
}

type CropMarketCombination struct {
	CropID   string
	MarketID string
}
