package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"marketlens/internal/auth"
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

// CORS middleware:
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
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
		Handler:           corsMiddleware(s.mux),
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

	// Admin endpoints
	adminMux := http.NewServeMux()
	s.mux.HandleFunc("/v1/admin/observations", s.handleListObservations)
	s.mux.HandleFunc("/v1/admin/observations/", s.handleAdminObservation) // PATCH /v1/admin/observations/{id}
	s.mux.HandleFunc("/v1/admin/audit", s.handleListAuditLogs)

	adminHandler := auth.APIKeyAuth(s.cfg.AdminAPIKey)(adminMux)
	s.mux.Handle("/v1/admin/", adminHandler)

	// USSD API Endpoints
	ussdHandler := auth.RateLimiter(s.cfg.USSDRateLimit)(http.HandlerFunc(s.ussdHandler.ServeUSSD))
	s.mux.Handle("/ussd", ussdHandler)
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

func (s *Server) handleListObservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	crop := q.Get("crop")
	market := q.Get("market")
	status := q.Get("status")
	limit := intParam(q.Get("limit"), 50)
	offset := intParam(q.Get("offset"), 0)

	obs, total, err := s.db.ListObservations(r.Context(), crop, market, status, limit, offset)
	if err != nil {
		log.Printf("handleListObservations error: %v", err)
		http.Error(w, "failed to list observations", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{
		"data":   obs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (s *Server) handleAdminObservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/admin/observations/")
	if id == "" {
		http.Error(w, "observation ID required", http.StatusBadRequest)
		return
	}

	var req struct {
		Status  string `json:"status"` // approve, flag, reject
		Reason  string `json:"reason"`
		AdminID string `json:"admin_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := map[string]bool{"approved": true, "flagged": true, "rejected": true}
	if !validStatuses[req.Status] {
		http.Error(w, "status must be 'approved', 'flagged', or 'rejected'", http.StatusBadRequest)
		return
	}
	if req.AdminID == "" {
		req.AdminID = "admin" // default for MVP
	}

	ctx := r.Context()

	// Get current observation for audit snapshot
	obs, err := s.db.GetObservationByID(ctx, id)
	if err != nil {
		log.Printf("handleAdminObservation GetObservationByID error: %v", err)
		http.Error(w, "observation not found", http.StatusNotFound)
		return
	}

	oldStatus := obs.Status

	// Update status
	if err := s.db.UpdateObservationStatus(ctx, id, req.Status); err != nil {
		log.Printf("handleAdminObservation UpdateObservationStatus error: %v", err)
		http.Error(w, "failed to update observation", http.StatusInternalServerError)
		return
	}

	// Record audit log
	oldJSON, _ := json.Marshal(map[string]string{"status": oldStatus})
	newJSON, _ := json.Marshal(map[string]string{"status": req.Status})

	auditLog := models.AuditLog{
		AdminID:    req.AdminID,
		Action:     req.Status, // "approved", "flagged", "rejected"
		EntityType: "price_observation",
		EntityID:   id,
		OldValue:   string(oldJSON),
		NewValue:   string(newJSON),
		Reason:     req.Reason,
	}
	if err := s.db.InsertAuditLog(ctx, auditLog); err != nil {
		log.Printf("handleAdminObservation InsertAuditLog error: %v", err)
		// Non-fatal — the status was already updated
	}

	writeJSON(w, map[string]string{"status": "updated", "new_status": req.Status})
}

func (s *Server) handleListAuditLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	entityType := q.Get("entity_type")
	entityID := q.Get("entity_id")
	limit := intParam(q.Get("limit"), 50)
	offset := intParam(q.Get("offset"), 0)

	logs, total, err := s.db.ListAuditLogs(r.Context(), entityType, entityID, limit, offset)
	if err != nil {
		log.Printf("handleListAuditLogs error: %v", err)
		http.Error(w, "failed to list audit logs", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// intParam parses a query string int with a default fallback.
func intParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return defaultVal
	}
	return n
}
