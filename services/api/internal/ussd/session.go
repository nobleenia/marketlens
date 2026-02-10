package ussd

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Session represents USSD session state.
type Session struct {
	State            string `json:"state"` // MAIN_MENU, SELECT_STATE, SELECT_CROP, SELECT_MARKET
	ChosenState      string `json:"chosen_state"`
	ChosenCropID     string `json:"chosen_crop_id"`
	ChosenCropName   string `json:"chosen_crop_name"`
	ChosenMarketID   string `json:"chosen_market_id"`
	ChosenMarketName string `json:"chosen_market_name"`
	Page             int    `json:"page"` // current page (0-indexed) for paginated lists
	Tries            int    `json:"tries"`
}

// Key: ussd:session:{sessionId}
// TTL: 5 minutes
const sessionTTL = 5 * time.Minute

// memEntry holds an in-memory session with its expiry time.
type memEntry struct {
	session   *Session
	expiresAt time.Time
}

// memStore is a simple in-memory session store used when Redis is unavailable.
type memStore struct {
	mu       sync.Mutex
	sessions map[string]memEntry
}

func newMemStore() *memStore {
	return &memStore{sessions: make(map[string]memEntry)}
}

func (m *memStore) get(key string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.sessions[key]
	if !ok || time.Now().After(entry.expiresAt) {
		delete(m.sessions, key)
		return nil
	}
	cp := *entry.session
	return &cp
}

func (m *memStore) set(key string, s *Session, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *s
	m.sessions[key] = memEntry{session: &cp, expiresAt: time.Now().Add(ttl)}
}

func (m *memStore) del(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, key)
}

// getSession loads the session from Redis, falling back to in-memory.
func (h *Handler) getSession(ctx context.Context, sessionID string) (*Session, error) {
	key := "ussd:session:" + sessionID

	if h.rdb != nil {
		val, err := h.rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			return nil, nil
		}
		if err != nil {
			// Redis error — fall through to memory
		} else {
			var s Session
			if err := json.Unmarshal([]byte(val), &s); err != nil {
				return nil, err
			}
			return &s, nil
		}
	}

	return h.mem.get(key), nil
}

// saveSession writes the session to Redis and in-memory.
func (h *Handler) saveSession(ctx context.Context, sessionID string, s *Session) error {
	key := "ussd:session:" + sessionID
	h.mem.set(key, s, sessionTTL)

	if h.rdb != nil {
		b, err := json.Marshal(s)
		if err != nil {
			return err
		}
		return h.rdb.Set(ctx, key, b, sessionTTL).Err()
	}
	return nil
}

// deleteSession removes the session from Redis and in-memory.
func (h *Handler) deleteSession(ctx context.Context, sessionID string) error {
	key := "ussd:session:" + sessionID
	h.mem.del(key)

	if h.rdb != nil {
		return h.rdb.Del(ctx, key).Err()
	}
	return nil
}
