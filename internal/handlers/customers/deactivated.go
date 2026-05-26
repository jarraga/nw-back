package customers

import (
	"encoding/json"
	"errors"
	"net/http"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) UpdateDeactivated(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request updateCustomerDeactivatedRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	customer, err := h.queries.UpdateCustomerDeactivated(r.Context(), db.UpdateCustomerDeactivatedParams{
		ID:          customerID,
		Deactivated: request.Deactivated,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update customer", http.StatusInternalServerError)
		return
	}

	response, err := newCustomerResponse(customer)
	if err != nil {
		http.Error(w, "failed to build customer response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
