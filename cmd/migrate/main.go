package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"nw-back/internal/postgres"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using environment defaults")
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	err = run(command)
	if err != nil {
		log.Printf("migration failed: %v", err)
		os.Exit(1)
	}
}

func run(command string) error {
	ctx := context.Background()

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
