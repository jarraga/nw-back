-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE company_type AS ENUM (
    'enterprise',
    'pyme',
    'startup'
  );
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS customers (
  id BIGSERIAL PRIMARY KEY,
  company_name TEXT NOT NULL,
  company_type company_type NOT NULL,
  phone TEXT NOT NULL,
  email TEXT NOT NULL,
  monthly_fee INT NOT NULL,
  billing_started_at DATE NOT NULL,
  comments TEXT NOT NULL DEFAULT '',
  deactivated BOOLEAN NOT NULL DEFAULT FALSE,
  reviewed_at TIMESTAMPTZ,
  reviewed_until TIMESTAMPTZ,
  reviewed_by TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_customers_company_name_trgm
ON customers USING GIN (company_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_customers_active
ON customers (id)
WHERE deactivated = FALSE;

-- +goose Down
DROP INDEX IF EXISTS idx_customers_active;
DROP INDEX IF EXISTS idx_customers_company_name_trgm;
DROP TABLE IF EXISTS customers;
DROP TYPE IF EXISTS company_type;
