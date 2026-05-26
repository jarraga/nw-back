package seed

import (
	"context"

	"nw-back/internal/migrate"
	"nw-back/internal/postgres"
)

func resetDatabase(ctx context.Context) error {
	err := dropDatabase(ctx)
	if err != nil {
		return err
	}

	err = migrate.Up(ctx)
	if err != nil {
		return err
	}

	return nil
}

func dropDatabase(ctx context.Context) error {
	_, err := postgres.DB.Exec(ctx, `
		DROP TABLE IF EXISTS customer_actions;
		DROP TABLE IF EXISTS customer_payments;
		DROP TABLE IF EXISTS customers;
		DROP TABLE IF EXISTS goose_db_version;
		DROP TYPE IF EXISTS customer_action_type;
		DROP TYPE IF EXISTS payment_status;
		DROP TYPE IF EXISTS company_type;
	`)
	if err != nil {
		return err
	}

	return nil
}
