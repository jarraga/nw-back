package admin

import (
	"encoding/json"
	"net/http"
	"sync"

	"nw-back/internal/seed"
)

type Handler struct {
	mu      sync.Mutex
	running bool
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ResetDemoData(w http.ResponseWriter, r *http.Request) {
	if !h.lock() {
		http.Error(w, "reset demo data already running", http.StatusConflict)
		return
	}
	defer h.unlock()

	err := seed.Run(r.Context())
	if err != nil {
		http.Error(w, "failed to reset demo data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) lock() bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return false
	}

	h.running = true
	return true
}

func (h *Handler) unlock() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.running = false
}
