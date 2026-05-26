package routes

import (
	"net/http"

	"nw-back/internal/handlers"
	"nw-back/internal/handlers/customers"
	"nw-back/internal/handlers/presence"
	"nw-back/internal/postgres/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(queries *db.Queries) http.Handler {
	router := chi.NewRouter()
	customersHandler := customers.NewHandler(queries)
	presenceHub := presence.NewHub()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{"*"},
	}))

	router.Get("/", handlers.Home())
	router.Get("/ws/customer-viewers", presenceHub.Handle)
	router.Get("/customers", customersHandler.List)
	router.Post("/customers", customersHandler.Create)
	router.Get("/customers/debt", customersHandler.Debt)
	router.Get("/customers/debt-list", customersHandler.DebtList)
	router.Get("/customers/delinquency-rate", customersHandler.DelinquencyRate)
	router.Get("/customers/monthly-delinquency", customersHandler.MonthlyDelinquency)
	router.Get("/customers/metrics", customersHandler.Metrics)
	router.Get("/customers/{customerID:[0-9]+}", customersHandler.Detail)
	router.Patch("/customers/{customerID:[0-9]+}", customersHandler.Update)
	router.Patch("/customers/{customerID:[0-9]+}/deactivated", customersHandler.UpdateDeactivated)
	router.Patch("/customers/{customerID:[0-9]+}/comments", customersHandler.UpdateComments)
	router.Patch("/customers/{customerID:[0-9]+}/review", customersHandler.Review)
	router.Delete("/customers/{customerID:[0-9]+}/review", customersHandler.ClearReview)
	router.Post("/customers/{customerID:[0-9]+}/actions", customersHandler.CreateAction)
	router.Patch("/customers/{customerID:[0-9]+}/actions/{actionID:[0-9]+}/comments", customersHandler.UpdateActionComments)
	router.Delete("/customers/{customerID:[0-9]+}/actions/{actionID:[0-9]+}", customersHandler.DeleteAction)
	router.Post("/customers/{customerID:[0-9]+}/payments", customersHandler.CreatePayment)

	return router
}
