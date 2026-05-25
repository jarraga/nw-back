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

-- name: ListCustomerActionsLastThreeMonths :many
SELECT
  id,
  customer_id,
  type,
  comments,
  action_date,
  created_at
FROM customer_actions
WHERE customer_id = sqlc.arg('customer_id')
  AND action_date >= NOW() - INTERVAL '3 months'
ORDER BY action_date DESC, id DESC;
