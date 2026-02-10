package ussd

import (
	"fmt"
	"log"
	"marketlens/internal/config"
	"marketlens/internal/store"
	"net/http"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

// perPage is the number of items shown per USSD screen.
const perPage = 7

type Handler struct {
	cfg   config.Config
	rdb   *redis.Client
	store *store.Postgres
	mem   *memStore
}

func NewHandler(cfg config.Config, rdb *redis.Client, store *store.Postgres) *Handler {
	return &Handler{
		cfg:   cfg,
		rdb:   rdb,
		store: store,
		mem:   newMemStore(),
	}
}

// ServeUSSD handles the Africa's Talking USSD webhook POST.
//
// Flow: MAIN_MENU → SELECT_STATE → SELECT_CROP → SELECT_MARKET → SHOW_PRICE
func (h *Handler) ServeUSSD(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	sessionID := r.FormValue("sessionId")
	if sessionID == "" {
		sessionID = r.FormValue("phoneNumber")
	}
	phone := r.FormValue("phoneNumber")
	text := r.FormValue("text")
	ctx := r.Context()

	// Load or create session
	session, err := h.getSession(ctx, sessionID)
	if err != nil {
		log.Printf("USSD getSession error: sessionId=%s err=%v", sessionID, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if session == nil {
		session = &Session{}
	}

	log.Printf("USSD req: sessionId=%s phone=%s text=%q state=%s page=%d tries=%d",
		sessionID, phone, text, session.State, session.Page, session.Tries)

	// New session with empty text → show main menu
	if strings.TrimSpace(text) == "" && session.State == "" {
		session.State = "MAIN_MENU"
		_ = h.saveSession(ctx, sessionID, session)
		h.respond(w, renderMainMenu())
		return
	}

	// Extract last input token (AT sends cumulative *-separated text)
	last := ""
	if trimmed := strings.TrimSpace(text); trimmed != "" {
		tokens := strings.Split(trimmed, "*")
		last = strings.TrimSpace(tokens[len(tokens)-1])
	}

	// Ensure state exists for mid-session joins
	if session.State == "" {
		session.State = "MAIN_MENU"
	}

	// Too many bad inputs → reset
	if session.Tries >= 3 {
		session.Tries = 0
		session.Page = 0
		session.State = "MAIN_MENU"
		_ = h.saveSession(ctx, sessionID, session)
		h.respond(w, renderMainMenu())
		return
	}

	var response string

	switch session.State {

	// ── MAIN MENU ─────────────────────────────────────────────
	case "MAIN_MENU":
		switch last {
		case "1":
			states, err := h.store.GetDistinctStates(ctx)
			if err != nil || len(states) == 0 {
				response = "END No states available now."
				_ = h.deleteSession(ctx, sessionID)
				break
			}
			session.State = "SELECT_STATE"
			session.Page = 0
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectState(states, session.Page, perPage)
		case "2":
			response = renderHelp()
			_ = h.deleteSession(ctx, sessionID)
		case "3":
			response = renderGoodbye()
			_ = h.deleteSession(ctx, sessionID)
		default:
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
		}

	// ── SELECT STATE ──────────────────────────────────────────
	case "SELECT_STATE":
		states, _ := h.store.GetDistinctStates(ctx)

		if last == "0" {
			session.State = "MAIN_MENU"
			session.Page = 0
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderMainMenu()
			break
		}
		if last == "00" { // Next page
			session.Page++
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectState(states, session.Page, perPage)
			break
		}
		if last == "98" && session.Page > 0 { // Previous page
			session.Page--
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectState(states, session.Page, perPage)
			break
		}

		idx, err := strconv.Atoi(last)
		if err != nil || idx <= 0 {
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
			break
		}
		// Map page-relative index to absolute index
		absIdx := session.Page*perPage + idx - 1
		if absIdx >= len(states) {
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
			break
		}
		session.ChosenState = states[absIdx]
		session.Tries = 0
		session.Page = 0
		session.State = "SELECT_CROP"
		_ = h.saveSession(ctx, sessionID, session)

		crops, _ := h.store.GetAllCrops(ctx)
		response = renderSelectCrop(crops, session.Page, perPage)

	// ── SELECT CROP ───────────────────────────────────────────
	case "SELECT_CROP":
		crops, _ := h.store.GetAllCrops(ctx)

		if last == "0" {
			states, _ := h.store.GetDistinctStates(ctx)
			session.State = "SELECT_STATE"
			session.Page = 0
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectState(states, session.Page, perPage)
			break
		}
		if last == "00" {
			session.Page++
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectCrop(crops, session.Page, perPage)
			break
		}
		if last == "98" && session.Page > 0 {
			session.Page--
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectCrop(crops, session.Page, perPage)
			break
		}

		idx, err := strconv.Atoi(last)
		if err != nil || idx <= 0 {
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
			break
		}
		absIdx := session.Page*perPage + idx - 1
		if absIdx >= len(crops) {
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
			break
		}
		chosen := crops[absIdx]
		session.ChosenCropID = chosen.ID
		session.ChosenCropName = chosen.Name
		session.Tries = 0
		session.Page = 0
		session.State = "SELECT_MARKET"
		_ = h.saveSession(ctx, sessionID, session)

		markets, _ := h.store.GetMarketsByState(ctx, session.ChosenState)
		response = renderSelectMarket(markets, session.Page, perPage)

	// ── SELECT MARKET ─────────────────────────────────────────
	case "SELECT_MARKET":
		markets, _ := h.store.GetMarketsByState(ctx, session.ChosenState)

		if last == "0" {
			crops, _ := h.store.GetAllCrops(ctx)
			session.State = "SELECT_CROP"
			session.Page = 0
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectCrop(crops, session.Page, perPage)
			break
		}
		if last == "00" {
			session.Page++
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectMarket(markets, session.Page, perPage)
			break
		}
		if last == "98" && session.Page > 0 {
			session.Page--
			session.Tries = 0
			_ = h.saveSession(ctx, sessionID, session)
			response = renderSelectMarket(markets, session.Page, perPage)
			break
		}

		idx, err := strconv.Atoi(last)
		if err != nil || idx <= 0 {
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
			break
		}
		absIdx := session.Page*perPage + idx - 1
		if absIdx >= len(markets) {
			session.Tries++
			_ = h.saveSession(ctx, sessionID, session)
			response = renderError()
			break
		}
		chosenMarket := markets[absIdx]
		session.ChosenMarketID = chosenMarket.ID
		session.ChosenMarketName = chosenMarket.Name

		agg, _ := h.store.GetAggregatedPriceByCropAndMarket(ctx, session.ChosenCropName, session.ChosenMarketName)
		response = renderShowPrice(agg)
		_ = h.deleteSession(ctx, sessionID)

	// ── UNKNOWN → RESET ───────────────────────────────────────
	default:
		session.State = "MAIN_MENU"
		session.Page = 0
		_ = h.saveSession(ctx, sessionID, session)
		response = renderMainMenu()
	}

	log.Printf("USSD resp: sessionId=%s state=%s page=%d", sessionID, session.State, session.Page)
	h.respond(w, response)
}

// respond writes a plain-text USSD response.
func (h *Handler) respond(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, body)
}
