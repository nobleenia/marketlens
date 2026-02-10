package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"marketlens/internal/config"
	"marketlens/internal/models"
	"marketlens/internal/store"
	"marketlens/internal/ussd"
)

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Server struct {
	cfg         config.Config
	db          *store.Postgres
	ussdHandler *ussd.Handler
	mux         *http.ServeMux
}

func New(cfg config.Config, db *store.Postgres, ussdH *ussd.Handler) *http.Server {
	s := &Server{
		cfg:         cfg,
		db:          db,
		ussdHandler: ussdH,
		mux:         http.NewServeMux(),
	}
	s.routes()

	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.handleHealth)
	s.mux.HandleFunc("/", s.handleRoot)

	// API endpoints
	s.mux.HandleFunc("/v1/crops", s.handleGetCrops)
	s.mux.HandleFunc("/v1/markets", s.handleGetMarkets)
	s.mux.HandleFunc("/v1/prices", s.handleGetPrices)     // list + filters via query params
	s.mux.HandleFunc("/v1/prices/", s.handleGetPricePath) // /v1/prices/{crop}/{market}
	s.mux.HandleFunc("/v1/observations", s.handlePostObservation)

	// USSD API Endpoints
	s.mux.HandleFunc("/ussd", s.ussdHandler.ServeUSSD)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Welcome to MarketLens API!\n"))
}

func (s *Server) handleGetCrops(w http.ResponseWriter, r *http.Request) {
	crops, err := s.db.GetAllCrops(r.Context())
	if err != nil {
		log.Printf("handleGetCrops error: %v", err)
		http.Error(w, "failed to fetch crops", http.StatusInternalServerError)
		return
	}
	writeJSON(w, crops)
}

func (s *Server) handleGetMarkets(w http.ResponseWriter, r *http.Request) {
	markets, err := s.db.GetAllMarkets(r.Context())
	if err != nil {
		log.Printf("handleGetMarkets error: %v", err)
		http.Error(w, "failed to fetch markets", http.StatusInternalServerError)
		return
	}
	writeJSON(w, markets)
}

func (s *Server) handleGetPrices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	crop := q.Get("crop")
	market := q.Get("market")

	var from, to time.Time
	var err error
	if v := q.Get("from"); v != "" {
		from, err = time.Parse("2006-01-02", v)
		if err != nil {
			http.Error(w, "invalid 'from' date (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}
	if v := q.Get("to"); v != "" {
		to, err = time.Parse("2006-01-02", v)
		if err != nil {
			http.Error(w, "invalid 'to' date (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}

	prices, err := s.db.GetAggregatedPrices(r.Context(), crop, market, from, to)
	if err != nil {
		log.Printf("handleGetPrices error: %v", err)
		http.Error(w, "failed to fetch prices", http.StatusInternalServerError)
		return
	}
	writeJSON(w, prices)
}

func (s *Server) handleGetPricePath(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/prices/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(w, "use /v1/prices/{crop}/{market}", http.StatusBadRequest)
		return
	}

	cropName, _ := url.PathUnescape(parts[0])
	marketName, _ := url.PathUnescape(parts[1])

	price, err := s.db.GetAggregatedPriceByCropAndMarket(r.Context(), cropName, marketName)
	if err != nil {
		log.Printf("handleGetPricePath error: %v", err)
		http.Error(w, "failed to fetch price", http.StatusInternalServerError)
		return
	}
	if price == nil {
		http.Error(w, "price not found", http.StatusNotFound)
		return
	}
	writeJSON(w, price)
}

func (s *Server) handlePostObservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CropID     string  `json:"crop_id"`
		CropName   string  `json:"crop_name"`
		MarketID   string  `json:"market_id"`
		MarketName string  `json:"market_name"`
		ObservedAt string  `json:"observed_at"` // RFC3339 or YYYY-MM-DD
		Price      float64 `json:"price"`
		Currency   string  `json:"currency"`
		Unit       string  `json:"unit"`
		PriceType  string  `json:"price_type"`
		Source     string  `json:"source"`
		ReporterID string  `json:"reporter_id"`
		Notes      string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	// Basic Validation
	if req.Price <= 0 {
		http.Error(w, "Price must be > 0", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = "NGN"
	}
	if req.Unit == "" {
		req.Unit = "kg"
	}
	if req.Source == "" {
		req.Source = "api"
	}
	if req.PriceType == "" {
		req.PriceType = "unknown"
	}

	ctx := r.Context()

	// Resolve CropID
	cropID := req.CropID
	if cropID == "" {
		if req.CropName == "" {
			http.Error(w, "crop_id or crop_name required", http.StatusBadRequest)
			return
		}
		id, err := s.db.LookupCropID(ctx, req.CropName)
		if err != nil {
			http.Error(w, "Unknown crop", http.StatusBadRequest)
			return
		}
		cropID = id
	}

	// Resolve market ID (try LookupMarketIDByName if state not supplied)
	marketID := req.MarketID
	if marketID == "" {
		if req.MarketName == "" {
			http.Error(w, "market_id or market_name required", http.StatusBadRequest)
			return
		}
		id, err := s.db.LookupMarketIDByName(ctx, req.MarketName)
		if err != nil {
			http.Error(w, "Unknown market", http.StatusBadRequest)
			return
		}
		marketID = id
	}

	// Parse observed_at
	var observedAt time.Time
	var err error
	if req.ObservedAt != "" {
		observedAt, err = time.Parse(time.RFC3339, req.ObservedAt)
		if err != nil {
			observedAt, err = time.Parse("2006-01-02", req.ObservedAt)
			if err != nil {
				http.Error(w, "Invalid observed_at (use RFC3339 or YYYY-MM-DD)", http.StatusBadRequest)
				return
			}
		}
	} else {
		observedAt = time.Now().UTC()
	}

	obs := models.PriceObservation{
		CropID:          cropID,
		MarketID:        marketID,
		ObservedAt:      observedAt,
		Price:           req.Price,
		Currency:        req.Currency,
		Unit:            req.Unit,
		PriceType:       req.PriceType,
		Source:          req.Source,
		ReporterID:      req.ReporterID,
		Notes:           req.Notes,
		ConfidenceScore: 0.50,
	}

	if err := s.db.InsertPriceObservation(ctx, obs); err != nil {
		log.Printf("InsertPriceObservation error: %v", err)
		http.Error(w, "Failed to insert observation", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"status": "created"})
}

// writeJSON is a small helper to avoid repeating Content-Type + Encode.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
