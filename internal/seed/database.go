package seed

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"nw-back/internal/postgres"
)

const schemaPath = "internal/postgres/schema/*.sql"

func resetDatabase(ctx context.Context) error {
	err := dropDatabase(ctx)
	if err != nil {
		return err
	}

	err = createDatabase(ctx)
	if err != nil {
		return err
	}

	return nil
}

func dropDatabase(ctx context.Context) error {
	_, err := postgres.DB.Exec(ctx, `
		DROP TABLE IF EXISTS customer_payments;
		DROP TABLE IF EXISTS customers;
		DROP TYPE IF EXISTS payment_status;
		DROP TYPE IF EXISTS company_type;
	`)
	if err != nil {
		return err
	}

	return nil
}

func createDatabase(ctx context.Context) error {
	files, err := filepath.Glob(schemaPath)
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = postgres.DB.Exec(ctx, string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
