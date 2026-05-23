package customers

import (
	"encoding/json"
	"net/http"

	"nw-back/internal/postgres/db"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params, err := parseListParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customers, err := h.queries.ListCustomers(r.Context(), db.ListCustomersParams{
		Limit:  int32(params.limit),
		Offset: int32(params.offset),
	})
	if err != nil {
		http.Error(w, "failed to list customers", http.StatusInternalServerError)
		return
	}

	response, err := newListResponse(customers)
	if err != nil {
		http.Error(w, "failed to build customers response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
