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
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
