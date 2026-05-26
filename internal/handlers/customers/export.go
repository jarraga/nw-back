package customers

import (
	"fmt"
	"net/http"
	"time"

	"nw-back/internal/xls"
)

func (h *Handler) ExportXLS(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("customers-%s.xlsx", time.Now().Format("2006-01-02-150405"))

	file, err := xls.CustomerWorkbook(r.Context(), h.pool)
	if err != nil {
		http.Error(w, "failed to export customers xls", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	err = file.Write(w)
	if err != nil {
		http.Error(w, "failed to export customers xls", http.StatusInternalServerError)
	}
}
