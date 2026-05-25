-- name: CreateCustomerAction :one
INSERT INTO customer_actions (
  customer_id,
  type,
  comments
) VALUES (
  sqlc.arg('customer_id'),
  sqlc.arg('type'),
  sqlc.arg('comments')
)
RETURNING
  id,
  customer_id,
  type,
  comments,
  action_date,
  created_at;
