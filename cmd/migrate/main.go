package main

import (
	"context"
	"log"
	"os"

	"nw-back/internal/migrate"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using environment defaults")
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	err = migrate.Run(context.Background(), command)
	if err != nil {
		log.Printf("migration failed: %v", err)
		os.Exit(1)
	}
}
