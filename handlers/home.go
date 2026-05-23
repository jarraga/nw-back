package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type homeResponse struct {
	Message string `json:"message"`
	UpTime  string `json:"upTime"`
}

func Home() http.HandlerFunc {
	startedAt := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := homeResponse{
			Message: "Northwind backend",
			UpTime:  time.Since(startedAt).Round(time.Second).String(),
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
