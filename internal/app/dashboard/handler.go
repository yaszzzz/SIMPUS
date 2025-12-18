package dashboard

import (
	"html/template"
	"net/http"
	"path/filepath"

	"simpus/internal/app/books"
	"simpus/internal/app/borrowings"
	"simpus/internal/app/members"
	"simpus/internal/middleware"
)

type Handler struct {
	bookService   *books.Service
	memberService *members.Service
	borrowService *borrowings.Service
	templates     *template.Template
}

func NewHandler(
	bookService *books.Service,
	memberService *members.Service,
	borrowService *borrowings.Service,
	templates *template.Template,
) *Handler {
	return &Handler{
		bookService:   bookService,
		memberService: memberService,
		borrowService: borrowService,
		templates:     templates,
	}
}

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Get stats
	activeBorrowings, _ := h.borrowService.GetActiveCount()
	overdueBorrowings, _ := h.borrowService.GetOverdueCount()
	totalBooks, _ := h.bookService.GetBookCount()
	totalMembers, _ := h.memberService.GetMemberCount()

	// Check overdue and create notifications
	h.borrowService.CheckAndCreateOverdueNotifications()

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":             "Dashboard - SIMPUS",
		"ActiveBorrowings":  activeBorrowings,
		"OverdueBorrowings": overdueBorrowings,
		"TotalBooks":        totalBooks,
		"TotalMembers":      totalMembers,
		"User":              claims,
	}

	h.render(w, "admin/dashboard.html", data)
}

func (h *Handler) MemberDashboard(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())

	// Get member borrowings
	borrowings, _ := h.borrowService.GetMemberBorrowings(claims.UserID)

	data := map[string]interface{}{
		"Title":      "Dashboard - SIMPUS",
		"Borrowings": borrowings,
		"User":       claims,
	}

	h.renderMember(w, "member/dashboard.html", data)
}

func (h *Handler) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := h.templates.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.ParseFiles(
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

func (h *Handler) renderMember(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := h.templates.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.ParseFiles(
		filepath.Join("templates", "layouts", "member.html"),
		filepath.Join("templates", "components", "member-sidebar.html"),
		filepath.Join("templates", "components", "member-header.html"),
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
