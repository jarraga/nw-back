package customers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		CustomerID:    customerID,
		Type:          actionType,
		Comments:      strings.TrimSpace(request.Comments),
		InformantName: newInformantName(request.InformantName),
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

func (h *Handler) UpdateActionComments(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actionID, err := parseActionID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request updateActionCommentsRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	action, err := h.queries.UpdateCustomerActionComments(r.Context(), db.UpdateCustomerActionCommentsParams{
		ID:         actionID,
		CustomerID: customerID,
		Comments:   strings.TrimSpace(request.Comments),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "customer action not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update customer action comments", http.StatusInternalServerError)
		return
	}

	response := newActionResponse(action)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	customerID, err := parseCustomerID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actionID, err := parseActionID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.queries.DeleteCustomerAction(r.Context(), db.DeleteCustomerActionParams{
		ID:         actionID,
		CustomerID: customerID,
	})
	if err != nil {
		http.Error(w, "failed to delete customer action", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "customer action not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func newInformantName(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}

	informantName := strings.TrimSpace(*value)
	if informantName == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{
		String: informantName,
		Valid:  true,
	}
}
