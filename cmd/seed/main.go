package main

import (
	"context"
	"log"
	"os"

	"nw-back/internal/postgres"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using environment defaults")
	}

	err = postgres.Connect(context.Background())
	if err != nil {
		log.Printf("postgres connection failed: %v", err)
		os.Exit(1)
	}
	log.Println("postgres connection OK")
	defer postgres.Close()

	gofakeit.Seed(0)
	log.Println("seed data generator ready")
}
