-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE payment_status AS ENUM (
    'pending',
    'paid'
  );
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS customer_payments (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  year INT NOT NULL,
  month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
  status payment_status NOT NULL DEFAULT 'pending',
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (customer_id, year, month)
);

CREATE INDEX IF NOT EXISTS idx_customer_payments_year_customer
ON customer_payments (year, customer_id);

-- +goose Down
DROP INDEX IF EXISTS idx_customer_payments_year_customer;
DROP TABLE IF EXISTS customer_payments;
DROP TYPE IF EXISTS payment_status;
