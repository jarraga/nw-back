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

-- name: GetTotalCustomerDebt :one
WITH customer_months AS (
  SELECT
    c.id AS customer_id,
    c.monthly_fee,
    EXTRACT(YEAR FROM month_date)::int AS year,
    EXTRACT(MONTH FROM month_date)::int AS month,
    (
      month_date::date +
      (
        LEAST(
          sqlc.arg('due_day')::int,
          EXTRACT(
            DAY FROM date_trunc('month', month_date)::date + INTERVAL '1 month - 1 day'
          )::int
        ) - 1
      ) * INTERVAL '1 day'
    )::date AS due_date
  FROM customers c
  CROSS JOIN LATERAL generate_series(
    date_trunc('month', c.billing_started_at)::date,
    date_trunc('month', CURRENT_DATE)::date,
    INTERVAL '1 month'
  ) AS month_date
)
SELECT
  COALESCE(SUM(cm.monthly_fee), 0)::numeric AS total_debt
FROM customer_months cm
LEFT JOIN customer_payments cp
  ON cp.customer_id = cm.customer_id
 AND cp.year = cm.year
 AND cp.month = cm.month
WHERE cp.id IS NULL
  AND cm.due_date < CURRENT_DATE;
