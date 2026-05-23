package routes

import (
	"net/http"

	"nw-back/internal/handlers"
	"nw-back/internal/handlers/customers"
	"nw-back/internal/postgres/db"

	"github.com/go-chi/chi/v5"
)

func NewRouter(queries *db.Queries) http.Handler {
	router := chi.NewRouter()
	customersHandler := customers.NewHandler(queries)

	router.Get("/", handlers.Home())
	router.Get("/customers", customersHandler.List)

	return router
}
