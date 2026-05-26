package migrate

import (
	"context"
	"fmt"

	"nw-back/internal/postgres"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func Run(ctx context.Context, command string) error {
	config, err := postgres.Config()
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	switch command {
	case "up":
		return goose.UpContext(ctx, db, migrationsDir)
	case "status":
		return goose.StatusContext(ctx, db, migrationsDir)
	default:
		return fmt.Errorf("unsupported migrate command %q", command)
	}
}

func Up(ctx context.Context) error {
	return Run(ctx, "up")
}
