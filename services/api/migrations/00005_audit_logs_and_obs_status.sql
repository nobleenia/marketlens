-- +goose Up
-- +goose StatementBegin

-- Add review status to price_observations so admins can approve/flag/reject
ALTER TABLE price_observations
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'pending';

CREATE INDEX IF NOT EXISTS idx_price_observations_status
    ON price_observations (status);

-- Audit trail for all admin actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id    TEXT NOT NULL,
    action      TEXT NOT NULL,            -- 'approve', 'flag', 'reject', 'override'
    entity_type TEXT NOT NULL,            -- 'price_observation', 'aggregated_price'
    entity_id   UUID NOT NULL,
    old_value   JSONB,                    -- snapshot before change
    new_value   JSONB,                    -- snapshot after change
    reason      TEXT NOT NULL DEFAULT '', -- justification text
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_entity
    ON audit_logs (entity_type, entity_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_admin
    ON audit_logs (admin_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_time
    ON audit_logs (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS audit_logs;
ALTER TABLE price_observations DROP COLUMN IF EXISTS status;

-- +goose StatementEnd