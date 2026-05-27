package admin

import (
	"encoding/json"
	"net/http"
	"sync"

	"nw-back/internal/resetdb"
	"nw-back/internal/seed"
	"nw-back/internal/xls"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	mu      sync.Mutex
	running bool
	pool    *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
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

func (h *Handler) ResetData(w http.ResponseWriter, r *http.Request) {
	if !h.lock() {
		http.Error(w, "reset data already running", http.StatusConflict)
		return
	}
	defer h.unlock()

	err := resetdb.Reset(r.Context(), h.pool)
	if err != nil {
		http.Error(w, "failed to reset data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) ImportXLS(w http.ResponseWriter, r *http.Request) {
	if !h.lock() {
		http.Error(w, "import xls already running", http.StatusConflict)
		return
	}
	defer h.unlock()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := xls.ImportCustomers(r.Context(), h.pool, file)
	if err != nil {
		http.Error(w, "failed to import xls", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(result)
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
