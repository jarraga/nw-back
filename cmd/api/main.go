package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"nw-back/internal/postgres"
	"nw-back/internal/postgres/db"
	"nw-back/internal/routes"

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

	queries := db.New(postgres.DB)
	router := routes.NewRouter(queries, postgres.DB)

	addr := ":" + envPort()
	log.Printf("northwind backend listening on %s", addr)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}

func envPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		return "8080"
	}

	return value
}
