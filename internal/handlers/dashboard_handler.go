package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"simpus/internal/middleware"
	"simpus/internal/services"
)

type DashboardHandler struct {
	bookService   *services.BookService
	memberService *services.MemberService
	borrowService *services.BorrowService
	templates     *template.Template
}

func NewDashboardHandler(
	bookService *services.BookService,
	memberService *services.MemberService,
	borrowService *services.BorrowService,
	templates *template.Template,
) *DashboardHandler {
	return &DashboardHandler{
		bookService:   bookService,
		memberService: memberService,
		borrowService: borrowService,
		templates:     templates,
	}
}

func (h *DashboardHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	totalBooks, availableBooks, _ := h.bookService.GetStats()
	totalMembers, _ := h.memberService.GetMemberCount()
	activeBorrowings, _ := h.borrowService.GetActiveCount()
	overdueCount, _ := h.borrowService.GetOverdueCount()

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":            "Dashboard - SIMPUS",
		"TotalBooks":       totalBooks,
		"AvailableBooks":   availableBooks,
		"TotalMembers":     totalMembers,
		"ActiveBorrowings": activeBorrowings,
		"OverdueCount":     overdueCount,
		"User":             claims,
	}

	h.render(w, "admin/dashboard.html", data)
}

func (h *DashboardHandler) MemberDashboard(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())

	borrowings, _ := h.borrowService.GetMemberBorrowings(claims.UserID)

	// Count active and history
	activeBorrowings := 0
	for _, b := range borrowings {
		if b.Status == "dipinjam" {
			activeBorrowings++
		}
	}

	data := map[string]interface{}{
		"Title":            "Dashboard Anggota - SIMPUS",
		"Borrowings":       borrowings,
		"ActiveBorrowings": activeBorrowings,
		"TotalBorrowings":  len(borrowings),
		"User":             claims,
	}

	h.renderMember(w, "member/dashboard.html", data)
}

func (h *DashboardHandler) render(w http.ResponseWriter, name string, data interface{}) {
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

func (h *DashboardHandler) renderMember(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layouts", "member.html"),
		filepath.Join("templates", "components", "member-navbar.html"),
		filepath.Join("templates", name),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "member.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
