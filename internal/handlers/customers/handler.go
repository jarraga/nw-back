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

	total, err := h.queries.CountCustomers(r.Context())
	if err != nil {
		http.Error(w, "failed to count customers", http.StatusInternalServerError)
		return
	}

	response, err := newPaginatedListResponse(customers, total)
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

func (h *Handler) Debt(w http.ResponseWriter, r *http.Request) {
	params, err := parseDebtParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalDebt, err := h.queries.GetTotalCustomerDebt(r.Context(), int32(params.dueDay))
	if err != nil {
		http.Error(w, "failed to calculate customer debt", http.StatusInternalServerError)
		return
	}

	response := newDebtResponse(totalDebt)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) DebtList(w http.ResponseWriter, r *http.Request) {
	params, err := parseDebtListParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customers, err := h.queries.ListCustomersDebt(r.Context(), db.ListCustomersDebtParams{
		DueDay:          int32(params.dueDay),
		SortBy:          params.sortBy,
		SortDirection:   params.sortDirection,
		CompanyName:     params.companyName,
		Limit:           int32(params.limit),
		Offset:          int32(params.offset),
		CompanyTypes:    params.companyTypes,
		IncludeReviewed: params.includeReviewed,
	})
	if err != nil {
		http.Error(w, "failed to list customers debt", http.StatusInternalServerError)
		return
	}

	total, err := h.queries.CountCustomersDebt(r.Context(), db.CountCustomersDebtParams{
		CompanyTypes:    params.companyTypes,
		CompanyName:     params.companyName,
		IncludeReviewed: params.includeReviewed,
	})
	if err != nil {
		http.Error(w, "failed to count customers debt", http.StatusInternalServerError)
		return
	}

	response, err := newPaginatedDebtListResponse(customers, total)
	if err != nil {
		http.Error(w, "failed to build customers debt response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
