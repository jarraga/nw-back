-- name: CreateCustomerAction :one
INSERT INTO customer_actions (
  customer_id,
  type,
  comments,
  informant_name,
  action_date
) VALUES (
  sqlc.arg('customer_id'),
  sqlc.arg('type'),
  sqlc.arg('comments'),
  sqlc.narg('informant_name'),
  COALESCE(sqlc.narg('action_date')::timestamptz, NOW())
)
RETURNING
  id,
  customer_id,
  type,
  comments,
  informant_name,
  action_date,
  created_at;

-- name: UpdateCustomerActionComments :one
UPDATE customer_actions
SET comments = sqlc.arg('comments')
WHERE id = sqlc.arg('id')
  AND customer_id = sqlc.arg('customer_id')
RETURNING
  id,
  customer_id,
  type,
  comments,
  informant_name,
  action_date,
  created_at;

-- name: DeleteCustomerAction :execrows
DELETE FROM customer_actions
WHERE id = sqlc.arg('id')
  AND customer_id = sqlc.arg('customer_id');

-- name: ListCustomerActionsLastThreeMonths :many
SELECT
  id,
  customer_id,
  type,
  comments,
  informant_name,
  action_date,
  created_at
FROM customer_actions
WHERE customer_id = sqlc.arg('customer_id')
  AND action_date >= NOW() - INTERVAL '3 months'
ORDER BY action_date DESC, id DESC;
