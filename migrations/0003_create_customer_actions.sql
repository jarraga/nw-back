-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE customer_action_type AS ENUM (
    'call',
    'email',
    'personal_visit',
    'other'
  );
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS customer_actions (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  type customer_action_type NOT NULL,
  comments TEXT NOT NULL DEFAULT '',
  informant_name TEXT,
  action_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_customer_actions_customer_date
ON customer_actions (customer_id, action_date DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_customer_actions_customer_date;
DROP TABLE IF EXISTS customer_actions;
DROP TYPE IF EXISTS customer_action_type;
