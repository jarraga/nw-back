package customers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) Review(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request reviewCustomerRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = parseReviewDays(request.Days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.queries.MarkCustomerReviewed(r.Context(), db.MarkCustomerReviewedParams{
		ID:         customerID,
		Days:       request.Days,
		ReviewedBy: strings.TrimSpace(request.ReviewedBy),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to review customer", http.StatusInternalServerError)
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

func (h *Handler) ClearReview(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.queries.ClearCustomerReview(r.Context(), customerID)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to clear customer review", http.StatusInternalServerError)
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

func parseReviewDays(days int32) error {
	if days <= 0 {
		return fmt.Errorf("days must be greater than 0")
	}

	return nil
}
