package store

import (
	"context"
	"os"
	"testing"
	"time"
)

// getTestDB returns a Postgres connection for testing.
// Skips the test if database is not available.
func getTestDB(t *testing.T) *Postgres {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "marketlens"
		}
		password := os.Getenv("POSTGRES_PASSWORD")
		if password == "" {
			password = "marketlens"
		}
		dbname := os.Getenv("POSTGRES_DB")
		if dbname == "" {
			dbname = "marketlens"
		}
		dsn = "postgresql://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pg, err := NewPostgres(ctx, dsn)
	if err != nil {
		t.Skipf("Skipping database test: could not connect to postgres: %v", err)
	}

	return pg
}

// ===========================================
// GetObservationsForAggregation Tests
// ===========================================

func TestGetObservationsForAggregation_NoResults(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	observations, err := pg.GetObservationsForAggregation(
		ctx,
		"00000000-0000-0000-0000-000000000000",
		"00000000-0000-0000-0000-000000000000",
		time.Now().Add(-24*time.Hour),
		time.Now(),
	)

	if err != nil {
		t.Fatalf("GetObservationsForAggregation returned error: %v", err)
	}

	if len(observations) != 0 {
		t.Errorf("Expected empty slice, got %d observations", len(observations))
	}
}

func TestGetObservationsForAggregation_FiltersByCropAndMarket(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	cropID, err := pg.LookupCropID(ctx, "Tomato")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded (no Tomato crop): %v", err)
	}

	marketID, err := pg.LookupMarketID(ctx, "Mile 12", "Lagos")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded (no Mile 12 market): %v", err)
	}

	testObs := PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      time.Now(),
		Price:           15000.0,
		Currency:        "NGN",
		Unit:            "basket",
		Source:          "test",
		ReporterID:      "test-reporter",
		Notes:           "integration test",
		ConfidenceScore: 0.5,
	}

	err = pg.InsertPriceObservation(ctx, testObs)
	if err != nil {
		t.Fatalf("Failed to insert test observation: %v", err)
	}

	observations, err := pg.GetObservationsForAggregation(
		ctx,
		cropID,
		marketID,
		time.Now().Add(-1*time.Hour),
		time.Now().Add(1*time.Hour),
	)

	if err != nil {
		t.Fatalf("GetObservationsForAggregation returned error: %v", err)
	}

	if len(observations) == 0 {
		t.Error("Expected at least one observation, got 0")
	}

	for _, obs := range observations {
		if obs.CropID != cropID {
			t.Errorf("Observation has wrong CropID: got %v, want %v", obs.CropID, cropID)
		}
		if obs.MarketID != marketID {
			t.Errorf("Observation has wrong MarketID: got %v, want %v", obs.MarketID, marketID)
		}
	}
}

func TestGetObservationsForAggregation_FiltersByTimeWindow(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	cropID, err := pg.LookupCropID(ctx, "Maize")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	marketID, err := pg.LookupMarketID(ctx, "Bodija", "Oyo")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	now := time.Now()

	// Insert old and recent observations
	oldObs := PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      now.Add(-48 * time.Hour),
		Price:           800.0,
		Currency:        "NGN",
		Unit:            "kg",
		Source:          "test",
		ReporterID:      "test-old",
		Notes:           "old observation",
		ConfidenceScore: 0.5,
	}
	_ = pg.InsertPriceObservation(ctx, oldObs)

	recentObs := PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      now.Add(-1 * time.Hour),
		Price:           850.0,
		Currency:        "NGN",
		Unit:            "kg",
		Source:          "test",
		ReporterID:      "test-recent",
		Notes:           "recent observation",
		ConfidenceScore: 0.5,
	}
	_ = pg.InsertPriceObservation(ctx, recentObs)

	// Query last 24 hours only
	observations, err := pg.GetObservationsForAggregation(
		ctx,
		cropID,
		marketID,
		now.Add(-24*time.Hour),
		now,
	)

	if err != nil {
		t.Fatalf("GetObservationsForAggregation returned error: %v", err)
	}

	for _, obs := range observations {
		if obs.ObservedAt.Before(now.Add(-24 * time.Hour)) {
			t.Errorf("Observation outside time window: %v", obs.ObservedAt)
		}
	}
}

func TestGetObservationsForAggregation_ReturnsCorrectFields(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	cropID, err := pg.LookupCropID(ctx, "Rice")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	marketID, err := pg.LookupMarketID(ctx, "Wuse Market", "FCT")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	now := time.Now()

	testObs := PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      now,
		Price:           1250.0,
		Currency:        "NGN",
		Unit:            "kg",
		Source:          "test-source",
		ReporterID:      "test-reporter-123",
		Notes:           "test notes here",
		ConfidenceScore: 0.75,
	}

	err = pg.InsertPriceObservation(ctx, testObs)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	observations, err := pg.GetObservationsForAggregation(
		ctx,
		cropID,
		marketID,
		now.Add(-1*time.Minute),
		now.Add(1*time.Minute),
	)

	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(observations) == 0 {
		t.Fatal("No observations returned")
	}

	// Find our test observation
	var found bool
	for _, obs := range observations {
		if obs.ReporterID == "test-reporter-123" {
			found = true
			if obs.Price != 1250.0 {
				t.Errorf("Price = %v, want 1250.0", obs.Price)
			}
			if obs.Currency != "NGN" {
				t.Errorf("Currency = %v, want NGN", obs.Currency)
			}
			if obs.Unit != "kg" {
				t.Errorf("Unit = %v, want kg", obs.Unit)
			}
			if obs.Source != "test-source" {
				t.Errorf("Source = %v, want test-source", obs.Source)
			}
			break
		}
	}

	if !found {
		t.Error("Test observation not found in results")
	}
}

// ===========================================
// Existing Function Tests
// ===========================================

func TestLookupCropID_ExistingCrop(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	cropID, err := pg.LookupCropID(ctx, "Tomato")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	if cropID == "" {
		t.Error("LookupCropID returned empty string")
	}
}

func TestLookupCropID_CaseInsensitive(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	id1, err1 := pg.LookupCropID(ctx, "tomato")
	id2, err2 := pg.LookupCropID(ctx, "TOMATO")
	id3, err3 := pg.LookupCropID(ctx, "Tomato")

	if err1 != nil || err2 != nil || err3 != nil {
		t.Skipf("Skipping: seed data not loaded")
	}

	if id1 != id2 || id2 != id3 {
		t.Errorf("LookupCropID not case-insensitive: %v, %v, %v", id1, id2, id3)
	}
}

func TestLookupCropID_NonExistent(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	_, err := pg.LookupCropID(ctx, "NonExistentCrop12345")

	if err == nil {
		t.Error("Expected error for non-existent crop")
	}
}

func TestLookupMarketID_ExistingMarket(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	marketID, err := pg.LookupMarketID(ctx, "Mile 12", "Lagos")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	if marketID == "" {
		t.Error("LookupMarketID returned empty string")
	}
}

func TestLookupMarketID_NonExistent(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	_, err := pg.LookupMarketID(ctx, "NonExistentMarket", "NoState")

	if err == nil {
		t.Error("Expected error for non-existent market")
	}
}

func TestInsertPriceObservation_ValidData(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	cropID, err := pg.LookupCropID(ctx, "Beans")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	marketID, err := pg.LookupMarketID(ctx, "Main Market", "Kano")
	if err != nil {
		t.Skipf("Skipping: seed data not loaded: %v", err)
	}

	obs := PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      time.Now(),
		Price:           950.0,
		Currency:        "NGN",
		Unit:            "kg",
		Source:          "test",
		ReporterID:      "integration-test",
		Notes:           "test insert",
		ConfidenceScore: 0.75,
	}

	err = pg.InsertPriceObservation(ctx, obs)

	if err != nil {
		t.Errorf("InsertPriceObservation failed: %v", err)
	}
}

func TestPing(t *testing.T) {
	pg := getTestDB(t)
	defer pg.Close()

	ctx := context.Background()

	err := pg.Ping(ctx)

	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}
