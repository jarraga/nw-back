-- name: CreateCustomerPayment :one
INSERT INTO customer_payments (
  customer_id,
  year,
  month,
  status,
  paid_at
) VALUES (
  sqlc.arg('customer_id'),
  sqlc.arg('year'),
  sqlc.arg('month'),
  sqlc.arg('status'),
  sqlc.arg('paid_at')
)
RETURNING
  id,
  customer_id,
  year,
  month,
  status,
  paid_at,
  created_at;
