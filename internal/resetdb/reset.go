package resetdb

import (
	"context"

	"nw-back/internal/migrate"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Reset(ctx context.Context, pool *pgxpool.Pool) error {
	err := drop(ctx, pool)
	if err != nil {
		return err
	}

	return migrate.Up(ctx)
}

func drop(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		DROP TABLE IF EXISTS customer_actions;
		DROP TABLE IF EXISTS customer_payments;
		DROP TABLE IF EXISTS customers;
		DROP TABLE IF EXISTS goose_db_version;
		DROP TYPE IF EXISTS customer_action_type;
		DROP TYPE IF EXISTS payment_status;
		DROP TYPE IF EXISTS company_type;
	`)
	return err
}
