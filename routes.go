package main

import (
	"net/http"

	"nw-back/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.Home())

	return router
}
