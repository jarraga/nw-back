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

-- name: CountCustomers :one
SELECT COUNT(*)::int
FROM customers;

-- name: SearchCustomersByCompanyName :many
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
WHERE company_name ILIKE '%' || sqlc.arg('company_name')::text || '%'
ORDER BY company_name ASC, id ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountCustomersByCompanyName :one
SELECT COUNT(*)::int
FROM customers
WHERE company_name ILIKE '%' || sqlc.arg('company_name')::text || '%';

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

-- name: ListCustomersDebt :many
WITH customer_debts AS (
  SELECT
    c.id,
    c.company_name,
    c.company_type,
    c.phone,
    c.email,
    c.monthly_fee,
    c.billing_started_at,
    COUNT(overdue_months.month_date)::int AS overdue_months,
    COALESCE(SUM(c.monthly_fee), 0)::numeric AS overdue_amount
  FROM customers c
  LEFT JOIN LATERAL (
    SELECT
      month_date
    FROM generate_series(
      date_trunc('month', c.billing_started_at)::date,
      date_trunc('month', CURRENT_DATE)::date,
      INTERVAL '1 month'
    ) AS month_date
    LEFT JOIN customer_payments cp
      ON cp.customer_id = c.id
     AND cp.year = EXTRACT(YEAR FROM month_date)::int
     AND cp.month = EXTRACT(MONTH FROM month_date)::int
    WHERE cp.id IS NULL
      AND (
        month_date::date +
        (
          LEAST(
            sqlc.arg('due_day')::int,
            EXTRACT(
              DAY FROM date_trunc('month', month_date)::date + INTERVAL '1 month - 1 day'
            )::int
          ) - 1
        ) * INTERVAL '1 day'
      )::date < CURRENT_DATE
  ) overdue_months ON true
  WHERE cardinality(sqlc.arg('company_types')::text[]) = 0
     OR c.company_type::text = ANY(sqlc.arg('company_types')::text[])
  GROUP BY c.id
)
SELECT
  id,
  company_name,
  company_type,
  phone,
  email,
  monthly_fee,
  billing_started_at,
  overdue_months,
  overdue_amount
FROM customer_debts
ORDER BY
  CASE WHEN sqlc.arg('sort_by')::text = 'amount' AND sqlc.arg('sort_direction')::text = 'asc' THEN overdue_amount END ASC,
  CASE WHEN sqlc.arg('sort_by')::text = 'amount' AND sqlc.arg('sort_direction')::text = 'desc' THEN overdue_amount END DESC,
  CASE WHEN sqlc.arg('sort_by')::text = 'months' AND sqlc.arg('sort_direction')::text = 'asc' THEN overdue_months END ASC,
  CASE WHEN sqlc.arg('sort_by')::text = 'months' AND sqlc.arg('sort_direction')::text = 'desc' THEN overdue_months END DESC,
  overdue_amount DESC,
  overdue_months DESC,
  id ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountCustomersDebt :one
SELECT COUNT(*)::int
FROM customers c
WHERE cardinality(sqlc.arg('company_types')::text[]) = 0
   OR c.company_type::text = ANY(sqlc.arg('company_types')::text[]);
