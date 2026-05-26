package main

import (
	"context"
	"flag"
	"log"
	"os"

	"nw-back/internal/postgres"
	"nw-back/internal/xls"

	"github.com/joho/godotenv"
)

func main() {
	outputPath := flag.String("out", "", "xlsx output path")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using environment defaults")
	}

	ctx := context.Background()

	err = postgres.Connect(ctx)
	if err != nil {
		log.Printf("postgres connection failed: %v", err)
		os.Exit(1)
	}
	defer postgres.Close()

	path, err := xls.ExportCustomers(ctx, postgres.DB, *outputPath)
	if err != nil {
		log.Printf("xls export failed: %v", err)
		os.Exit(1)
	}

	log.Printf("xls exported: %s", path)
}
