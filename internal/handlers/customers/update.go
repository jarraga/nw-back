package customers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request updateCustomerRequest

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	err = decoder.Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params, err := newUpdateCustomerParams(customerID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.queries.UpdateCustomerContact(r.Context(), params)
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

func newUpdateCustomerParams(customerID int64, request updateCustomerRequest) (db.UpdateCustomerContactParams, error) {
	phone := strings.TrimSpace(request.Phone)
	email := strings.TrimSpace(request.Email)

	monthlyFee, err := parseMonthlyFee(request.MonthlyFee)
	if err != nil {
		return db.UpdateCustomerContactParams{}, err
	}

	return db.UpdateCustomerContactParams{
		ID:         customerID,
		Phone:      phone,
		Email:      email,
		MonthlyFee: monthlyFee,
	}, nil
}
