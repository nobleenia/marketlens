package ussd

import (
	"marketlens/internal/config"
	"marketlens/internal/store"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	cfg   config.Config
	rdb   *redis.Client
	store *store.Postgres
}

func NewHandler(cfg config.Config, rdb *redis.Client, store *store.Postgres) *Handler {
	return &Handler{cfg: cfg, rdb: rdb, store: store}
}
