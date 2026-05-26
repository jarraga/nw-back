CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE company_type AS ENUM (
  'enterprise',
  'pyme',
  'startup'
);

CREATE TABLE customers (
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

CREATE INDEX idx_customers_company_name_trgm
ON customers USING GIN (company_name gin_trgm_ops);

CREATE INDEX idx_customers_active
ON customers (id)
WHERE deactivated = FALSE;
