package customers

import (
	"encoding/json"
	"errors"
	"net/http"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	debtParams, err := parseDebtParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.queries.GetCustomer(r.Context(), customerID)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get customer", http.StatusInternalServerError)
		return
	}

	actions, err := h.queries.ListCustomerActionsLastThreeMonths(r.Context(), customerID)
	if err != nil {
		http.Error(w, "failed to list customer actions", http.StatusInternalServerError)
		return
	}

	payments, err := h.queries.ListCustomerPaymentsLastYear(r.Context(), customerID)
	if err != nil {
		http.Error(w, "failed to list customer payments", http.StatusInternalServerError)
		return
	}

	debt, err := h.queries.GetCustomerDebtSummary(r.Context(), db.GetCustomerDebtSummaryParams{
		DueDay:     int32(debtParams.dueDay),
		CustomerID: customerID,
	})
	if err != nil {
		http.Error(w, "failed to get customer debt", http.StatusInternalServerError)
		return
	}

	behavior, err := h.queries.GetCustomerPaymentBehavior(r.Context(), db.GetCustomerPaymentBehaviorParams{
		DueDay:     int32(debtParams.dueDay),
		CustomerID: customerID,
	})
	if err != nil {
		http.Error(w, "failed to get customer payment behavior", http.StatusInternalServerError)
		return
	}

	response, err := newDetailResponse(customer, actions, payments, debt, behavior, int32(debtParams.dueDay))
	if err != nil {
		http.Error(w, "failed to build customer detail response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
