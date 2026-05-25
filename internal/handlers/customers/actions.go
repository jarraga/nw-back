package customers

import (
	"encoding/json"
	"net/http"
	"strings"

	"nw-back/internal/postgres/db"
)

func (h *Handler) CreateAction(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request createActionRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	actionType, err := parseCustomerActionType(request.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	action, err := h.queries.CreateCustomerAction(r.Context(), db.CreateCustomerActionParams{
		CustomerID: customerID,
		Type:       actionType,
		Comments:   strings.TrimSpace(request.Comments),
	})
	if err != nil {
		http.Error(w, "failed to create customer action", http.StatusInternalServerError)
		return
	}

	response := newActionResponse(action)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
