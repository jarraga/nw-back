package seed

import (
	"context"

	"nw-back/internal/postgres"
	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5"
)

func copyCustomers(ctx context.Context, customers []db.Customer) error {
	if len(customers) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(customers))

	for _, customer := range customers {
		rows = append(rows, []any{
			customer.ID,
			customer.CompanyName,
			string(customer.CompanyType),
			customer.Phone,
			customer.Email,
			customer.MonthlyFee,
			customer.BillingStartedAt.Time,
			customer.Comments,
		})
	}

	_, err := postgres.DB.CopyFrom(
		ctx,
		pgx.Identifier{"customers"},
		[]string{
			"id",
			"company_name",
			"company_type",
			"phone",
			"email",
			"monthly_fee",
			"billing_started_at",
			"comments",
		},
		pgx.CopyFromRows(rows),
	)
	return err
}

func copyPayments(ctx context.Context, payments []db.CreateCustomerPaymentParams) error {
	if len(payments) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(payments))

	for _, payment := range payments {
		rows = append(rows, []any{
			payment.CustomerID,
			payment.Year,
			payment.Month,
			string(payment.Status),
			payment.PaidAt.Time,
		})
	}

	_, err := postgres.DB.CopyFrom(
		ctx,
		pgx.Identifier{"customer_payments"},
		[]string{
			"customer_id",
			"year",
			"month",
			"status",
			"paid_at",
		},
		pgx.CopyFromRows(rows),
	)
	return err
}

func copyActions(ctx context.Context, actions []db.CreateCustomerActionParams) error {
	if len(actions) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(actions))

	for _, action := range actions {
		var informantName any
		if action.InformantName.Valid {
			informantName = action.InformantName.String
		}

		rows = append(rows, []any{
			action.CustomerID,
			string(action.Type),
			action.Comments,
			informantName,
			action.ActionDate.Time,
		})
	}

	_, err := postgres.DB.CopyFrom(
		ctx,
		pgx.Identifier{"customer_actions"},
		[]string{
			"customer_id",
			"type",
			"comments",
			"informant_name",
			"action_date",
		},
		pgx.CopyFromRows(rows),
	)
	return err
}

func resetCustomerSequence(ctx context.Context) error {
	_, err := postgres.DB.Exec(ctx, `
		SELECT setval(
			pg_get_serial_sequence('customers', 'id'),
			COALESCE((SELECT MAX(id) FROM customers), 1),
			(SELECT COUNT(*) > 0 FROM customers)
		)
	`)
	return err
}
