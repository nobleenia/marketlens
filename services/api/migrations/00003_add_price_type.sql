-- +goose Up
-- +goose StatementBegin

ALTER TABLE price_observations
    ADD COLUMN IF NOT EXISTS price_type TEXT NOT NULL DEFAULT 'unknown';

-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin

ALTER TABLE price_observations
    DROP COLUMN IF EXISTS price_type;

-- +goose StatementEnd