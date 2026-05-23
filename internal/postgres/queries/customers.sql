-- name: ListCustomers :many
SELECT
  id,
  company_name,
  company_type,
  phone,
  email,
  monthly_fee,
  created_at
FROM customers
ORDER BY id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');
