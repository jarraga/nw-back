CREATE TYPE customer_action_type AS ENUM (
  'call',
  'email',
  'personal_visit',
  'other'
);

CREATE TABLE customer_actions (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  type customer_action_type NOT NULL,
  comments TEXT NOT NULL DEFAULT '',
  informant_name TEXT,
  action_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customer_actions_customer_date
ON customer_actions (customer_id, action_date DESC);
