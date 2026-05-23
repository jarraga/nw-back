package main

import (
	"context"
	"log"
	"os"

	"nw-back/internal/seed"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using environment defaults")
	}

	err = seed.Run(context.Background())
	if err != nil {
		log.Printf("seed failed: %v", err)
		os.Exit(1)
	}
}
