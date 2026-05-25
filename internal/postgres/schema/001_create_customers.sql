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
  monthly_fee NUMERIC(12, 2) NOT NULL,
  billing_started_at DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_company_name_trgm
ON customers USING GIN (company_name gin_trgm_ops);
