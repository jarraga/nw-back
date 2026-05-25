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
  COALESCE(SUM(cm.monthly_fee), 0)::bigint AS total_debt
FROM customer_months cm
LEFT JOIN customer_payments cp
  ON cp.customer_id = cm.customer_id
 AND cp.year = cm.year
 AND cp.month = cm.month
WHERE cp.id IS NULL
  AND cm.due_date < CURRENT_DATE;

-- name: GetCustomerDebtSummary :one
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
  WHERE c.id = sqlc.arg('customer_id')
)
SELECT
  COUNT(customer_months.month)::int AS overdue_months,
  COALESCE(SUM(customer_months.monthly_fee), 0)::bigint AS overdue_amount
FROM customer_months
LEFT JOIN customer_payments cp
  ON cp.customer_id = customer_months.customer_id
 AND cp.year = customer_months.year
 AND cp.month = customer_months.month
WHERE customer_months.due_date < CURRENT_DATE
  AND (
    cp.id IS NULL
    OR cp.status <> 'paid'
  );

-- name: GetCustomerDelinquencyRate :one
WITH active_customers AS (
  SELECT
    id,
    billing_started_at
  FROM customers
  WHERE date_trunc('month', billing_started_at)::date <= date_trunc('month', CURRENT_DATE)::date
),
overdue_customers AS (
  SELECT DISTINCT
    active_customers.id
  FROM active_customers
  CROSS JOIN LATERAL (
    SELECT
      month_date,
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
    FROM generate_series(
      date_trunc('month', active_customers.billing_started_at)::date,
      date_trunc('month', CURRENT_DATE)::date,
      INTERVAL '1 month'
    ) AS month_date
  ) customer_months
  LEFT JOIN customer_payments cp
    ON cp.customer_id = active_customers.id
   AND cp.year = customer_months.year
   AND cp.month = customer_months.month
  WHERE customer_months.due_date < CURRENT_DATE
    AND (
      cp.id IS NULL
      OR cp.status <> 'paid'
    )
)
SELECT
  COUNT(active_customers.id)::int AS total_customers,
  COUNT(overdue_customers.id)::int AS overdue_customers,
  COALESCE(
    COUNT(overdue_customers.id)::double precision * 100 / NULLIF(COUNT(active_customers.id), 0),
    0
  )::double precision AS delinquency_percentage
FROM active_customers
LEFT JOIN overdue_customers
  ON overdue_customers.id = active_customers.id;

-- name: ListCustomerPaymentsLastYear :many
SELECT
  id,
  customer_id,
  year,
  month,
  status,
  paid_at,
  created_at
FROM customer_payments
WHERE customer_id = sqlc.arg('customer_id')
  AND make_date(year, month, 1) >= date_trunc('month', CURRENT_DATE)::date - INTERVAL '11 months'
  AND make_date(year, month, 1) <= date_trunc('month', CURRENT_DATE)::date
ORDER BY year DESC, month DESC;

-- name: GetMonthlyDelinquencyRate :many
WITH months AS (
  SELECT generate_series(1, 12)::int AS month
),
customer_months AS (
  SELECT
    months.month,
    c.id AS customer_id,
    (
      make_date(sqlc.arg('year')::int, months.month, 1) +
      (
        LEAST(
          sqlc.arg('due_day')::int,
          EXTRACT(
            DAY FROM make_date(sqlc.arg('year')::int, months.month, 1) + INTERVAL '1 month - 1 day'
          )::int
        ) - 1
      ) * INTERVAL '1 day'
    )::date AS due_date
  FROM months
  JOIN customers c
    ON date_trunc('month', c.billing_started_at)::date <= make_date(sqlc.arg('year')::int, months.month, 1)
)
SELECT
  months.month,
  COUNT(customer_months.customer_id)::int AS total_customers,
  COUNT(customer_months.customer_id) FILTER (
    WHERE customer_payments.id IS NULL
       OR customer_payments.status <> 'paid'
       OR customer_payments.paid_at::date > customer_months.due_date
  )::int AS overdue_customers,
  COALESCE(
    COUNT(customer_months.customer_id) FILTER (
      WHERE customer_payments.id IS NULL
         OR customer_payments.status <> 'paid'
         OR customer_payments.paid_at::date > customer_months.due_date
    )::double precision * 100 / NULLIF(COUNT(customer_months.customer_id), 0),
    0
  )::double precision AS delinquency_percentage
FROM months
LEFT JOIN customer_months
  ON customer_months.month = months.month
 AND customer_months.due_date <= CURRENT_DATE
LEFT JOIN customer_payments
  ON customer_payments.customer_id = customer_months.customer_id
 AND customer_payments.year = sqlc.arg('year')::int
 AND customer_payments.month = months.month
GROUP BY months.month
ORDER BY months.month;
