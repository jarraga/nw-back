package customers

import (
	"encoding/json"
	"net/http"

	"nw-back/internal/postgres/db"
)

func (h *Handler) MonthlyDelinquency(w http.ResponseWriter, r *http.Request) {
	params, err := parseMonthlyDelinquencyParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := h.queries.GetMonthlyDelinquencyRate(r.Context(), db.GetMonthlyDelinquencyRateParams{
		Year:   int32(params.year),
		DueDay: int32(params.dueDay),
	})
	if err != nil {
		http.Error(w, "failed to get monthly delinquency", http.StatusInternalServerError)
		return
	}

	response := newMonthlyDelinquencyResponse(int32(params.year), int32(params.dueDay), rows)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) DelinquencyRate(w http.ResponseWriter, r *http.Request) {
	params, err := parseDebtParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row, err := h.queries.GetCustomerDelinquencyRate(r.Context(), int32(params.dueDay))
	if err != nil {
		http.Error(w, "failed to get customer delinquency rate", http.StatusInternalServerError)
		return
	}

	response := newDelinquencyRateResponse(int32(params.dueDay), row)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	params, err := parseDebtParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := h.queries.GetCustomerMetrics(r.Context(), int32(params.dueDay))
	if err != nil {
		http.Error(w, "failed to get customer metrics", http.StatusInternalServerError)
		return
	}

	response := newCustomerMetricsResponse(int32(params.dueDay), rows)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) ReviewedDebtorsPercentage(w http.ResponseWriter, r *http.Request) {
	percentage, err := h.queries.GetReviewedDebtorsPercentage(r.Context())
	if err != nil {
		http.Error(w, "failed to get reviewed debtors percentage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(percentage)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
