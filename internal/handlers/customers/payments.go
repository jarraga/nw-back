package customers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request createPaymentRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params, err := newCreatePaymentParams(customerID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := h.queries.CreateCustomerPayment(r.Context(), params)
	if isUniqueViolation(err) {
		http.Error(w, "customer payment already exists", http.StatusConflict)
		return
	}
	if isForeignKeyViolation(err) {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to create customer payment", http.StatusInternalServerError)
		return
	}

	response := newPaymentResponse(payment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func newCreatePaymentParams(customerID int64, request createPaymentRequest) (db.CreateCustomerPaymentParams, error) {
	if request.Year <= 0 {
		return db.CreateCustomerPaymentParams{}, fmt.Errorf("year is required")
	}

	if request.Month < 1 || request.Month > 12 {
		return db.CreateCustomerPaymentParams{}, fmt.Errorf("month must be between 1 and 12")
	}

	status, err := parsePaymentStatus(request.Status)
	if err != nil {
		return db.CreateCustomerPaymentParams{}, err
	}

	paidAt, err := parsePaidAt(request.PaidAt)
	if err != nil {
		return db.CreateCustomerPaymentParams{}, err
	}

	return db.CreateCustomerPaymentParams{
		CustomerID: customerID,
		Year:       request.Year,
		Month:      request.Month,
		Status:     status,
		PaidAt:     paidAt,
	}, nil
}

func parsePaidAt(value *string) (pgtype.Timestamptz, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return pgtype.Timestamptz{}, fmt.Errorf("paidAt is required")
	}

	paidAt, err := time.Parse(time.RFC3339, strings.TrimSpace(*value))
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("paidAt must use RFC3339 format")
	}

	return pgtype.Timestamptz{
		Time:  paidAt,
		Valid: true,
	}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505"
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23503"
}
