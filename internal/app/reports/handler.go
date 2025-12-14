package reports

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"simpus/internal/app/borrowings"
	"simpus/internal/middleware"
	"simpus/internal/models"
)

type Handler struct {
	borrowService *borrowings.Service
	templates     *template.Template
}

func NewHandler(
	borrowService *borrowings.Service,
	templates *template.Template,
) *Handler {
	return &Handler{
		borrowService: borrowService,
		templates:     templates,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	filter := models.BorrowingFilter{
		Page:  1,
		Limit: 1000,
	}

	if fromDate != "" {
		t, _ := time.Parse("2006-01-02", fromDate)
		filter.FromDate = t
	}
	if toDate != "" {
		t, _ := time.Parse("2006-01-02", toDate)
		filter.ToDate = t
	}

	borrowings, total, _ := h.borrowService.GetBorrowings(filter)

	claims := middleware.GetUserFromContext(r.Context())

	// Calculate stats
	var totalFine float64
	returnedCount := 0
	overdueCount := 0
	for _, b := range borrowings {
		totalFine += b.Fine
		if b.Status == "dikembalikan" {
			returnedCount++
		}
		if b.Status == "terlambat" {
			overdueCount++
		}
	}

	data := map[string]interface{}{
		"Title":         "Laporan Transaksi - SIMPUS",
		"Borrowings":    borrowings,
		"Total":         total,
		"TotalFine":     totalFine,
		"ReturnedCount": returnedCount,
		"OverdueCount":  overdueCount,
		"FromDate":      fromDate,
		"ToDate":        toDate,
		"User":          claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/reports/table.html", data)
		return
	}

	h.render(w, "admin/reports/index.html", data)
}

func (h *Handler) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layouts", "admin.html"),
		filepath.Join("templates", "components", "sidebar.html"),
		filepath.Join("templates", "components", "navbar.html"),
		filepath.Join("templates", name),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "admin.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) renderPartial(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
