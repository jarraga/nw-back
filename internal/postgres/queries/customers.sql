-- name: ListCustomers :many
SELECT
  id,
  company_name,
  company_type,
  phone,
  email,
  monthly_fee,
  billing_started_at,
  created_at
FROM customers
ORDER BY id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateCustomer :one
INSERT INTO customers (
  company_name,
  company_type,
  phone,
  email,
  monthly_fee,
  billing_started_at
) VALUES (
  sqlc.arg('company_name'),
  sqlc.arg('company_type'),
  sqlc.arg('phone'),
  sqlc.arg('email'),
  sqlc.arg('monthly_fee'),
  sqlc.arg('billing_started_at')
)
RETURNING
  id,
  company_name,
  company_type,
  phone,
  email,
  monthly_fee,
  billing_started_at,
  created_at;
