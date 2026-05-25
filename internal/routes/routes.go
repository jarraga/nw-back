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
	router.Post("/customers", customersHandler.Create)
	router.Get("/customers/search", customersHandler.Search)
	router.Get("/customers/debt", customersHandler.Debt)
	router.Get("/customers/debt-list", customersHandler.DebtList)
	router.Get("/customers/monthly-delinquency", customersHandler.MonthlyDelinquency)
	router.Get("/customers/{customerID}", customersHandler.Detail)
	router.Patch("/customers/{customerID}", customersHandler.Update)
	router.Patch("/customers/{customerID}/comments", customersHandler.UpdateComments)
	router.Post("/customers/{customerID}/actions", customersHandler.CreateAction)
	router.Post("/customers/{customerID}/payments", customersHandler.CreatePayment)

	return router
}
