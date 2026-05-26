package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(ctx context.Context) error {
	config, err := Config()
	if err != nil {
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return err
	}

	DB = pool
	return nil
}

func Config() (*pgxpool.Config, error) {
	return pgxpool.ParseConfig(databaseURL())
}

func Close() {
	if DB == nil {
		return
	}

	DB.Close()
}

func databaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		env("PGUSER", "postgres"),
		env("PGPASSWORD", ""),
		env("PGHOST", "localhost"),
		env("PGPORT", "5432"),
		env("PGDATABASE", "postgres"),
	)
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
