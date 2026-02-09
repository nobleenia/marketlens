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
	ID          string          `json:"id"`
	CropID      string          `json:"crop_id"`
	CropName    string          `json:"crop_name"`
	MarketID    string          `json:"market_id"`
	MarketName  string          `json:"market_name"`
	Period      string          `json:"period"` // "daily", "weekly", "monthly"
	PeriodStart time.Time       `json:"period_start"`
	PeriodEnd   time.Time       `json:"period_end"`
	PriceMin    float64         `json:"price_min"`
	PriceMax    float64         `json:"price_max"`
	PriceMean   float64         `json:"price_mean"`
	PriceMedian float64         `json:"price_median"`
	Currency    string          `json:"currency"`
	Unit        string          `json:"unit"`
	Confidence  ConfidenceLevel `json:"confidence"`
	SampleSize  int             `json:"sample_size"`
	Trend       Trend           `json:"trend"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
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
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Crop struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type CropMarketCombination struct {
	CropID   string
	MarketID string
}
