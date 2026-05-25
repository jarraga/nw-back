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
