-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS markets (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name                TEXT NOT NULL,
    state               TEXT NOT NULL,
    country             TEXT NOT NULL DEFAULT 'Nigeria',
    latitude            DOUBLE PRECISION,
    longitude           DOUBLE PRECISION,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (name, state, country)
);

CREATE TABLE IF NOT EXISTS crops (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name                TEXT NOT NULL UNIQUE,
    unit                TEXT NOT NULL DEFAULT 'kg',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Raw prices collected from markets (from enumerators, farmers, scrapers)
CREATE TABLE IF NOT EXISTS price_observations (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    crop_id             UUID NOT NULL REFERENCES crops(id) ON DELETE RESTRICT,
    market_id           UUID NOT NULL REFERENCES markets(id) ON DELETE RESTRICT,
    
    observed_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    price               NUMERIC(12, 2) NOT NULL,
    currency            TEXT NOT NULL DEFAULT 'NGN',
    unit                TEXT NOT NULL DEFAULT 'kg',

    source              TEXT NOT NULL DEFAULT 'manual',
    reporter_id         TEXT,
    notes               TEXT,

    confidence_score    NUMERIC(4, 3) NOT NULL DEFAULT 0.500, -- value between 0 and 1
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_price_observations_crop_market_time 
    ON price_observations (crop_id, market_id, observed_at DESC);

-- Aggregated daily/weekly prices used by USSD + apps for fast reads
CREATE TABLE IF NOT EXISTS aggregated_prices (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    crop_id             UUID NOT NULL REFERENCES crops(id) ON DELETE RESTRICT,
    market_id           UUID NOT NULL REFERENCES markets(id) ON DELETE RESTRICT,
    
    period              TEXT NOT NULL,
    period_start        DATE NOT NULL,
    period_end          DATE NOT NULL,

    price_min           NUMERIC(12, 2),
    price_max           NUMERIC(12, 2),
    price_avg           NUMERIC(12, 2),
    price_median        NUMERIC(12, 2),

    currency            TEXT NOT NULL DEFAULT 'NGN',
    unit                TEXT NOT NULL DEFAULT 'kg',

    confidence_score    NUMERIC(4, 3) NOT NULL DEFAULT 0.500, -- value between 0 and 1
    sample_size         INTEGER NOT NULL DEFAULT 0,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (crop_id, market_id, period, period_start, period_end)
);

CREATE INDEX IF NOT EXISTS idx_aggregated_prices_lookup 
    ON aggregated_prices (crop_id, market_id, period, period_start DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS aggregated_prices;
DROP TABLE IF EXISTS price_observations;
DROP TABLE IF EXISTS crops;
DROP TABLE IF EXISTS markets;
-- +goose StatementEnd