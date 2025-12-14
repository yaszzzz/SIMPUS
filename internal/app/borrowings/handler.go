package borrowings

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"simpus/internal/app/books"
	"simpus/internal/app/members"
	"simpus/internal/middleware"
	"simpus/internal/models"
)

type Handler struct {
	service       *Service
	bookService   *books.Service
	memberService *members.Service
	templates     *template.Template
}

func NewHandler(
	service *Service,
	bookService *books.Service,
	memberService *members.Service,
	templates *template.Template,
) *Handler {
	return &Handler{
		service:       service,
		bookService:   bookService,
		memberService: memberService,
		templates:     templates,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	status := r.URL.Query().Get("status")

	filter := models.BorrowingFilter{
		Status: status,
		Page:   page,
		Limit:  10,
	}

	borrowings, total, err := h.service.GetBorrowings(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + 10 - 1) / 10

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":      "Manajemen Peminjaman - SIMPUS",
		"Borrowings": borrowings,
		"Total":      total,
		"Page":       page,
		"TotalPages": totalPages,
		"Status":     status,
		"User":       claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/borrowings/table.html", data)
		return
	}

	h.render(w, "admin/borrowings/index.html", data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	members, _, _ := h.memberService.GetMembers(1, 100, "")

	filter := models.BookFilter{Page: 1, Limit: 100, Available: true}
	books, _, _ := h.bookService.GetBooks(filter)

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":   "Tambah Peminjaman - SIMPUS",
		"Members": members,
		"Books":   books,
		"User":    claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/borrowings/form.html", data)
		return
	}

	h.render(w, "admin/borrowings/create.html", data)
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	memberID, _ := strconv.Atoi(r.FormValue("member_id"))
	bookID, _ := strconv.Atoi(r.FormValue("book_id"))
	borrowDays, _ := strconv.Atoi(r.FormValue("borrow_days"))

	claims := middleware.GetUserFromContext(r.Context())

	data := &models.BorrowingCreate{
		MemberID:   memberID,
		BookID:     bookID,
		BorrowDays: borrowDays,
		Notes:      r.FormValue("notes"),
	}

	_, err := h.service.CreateBorrowing(data, claims.UserID)
	if err != nil {
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Retarget", "#error-message")
			w.Write([]byte(`<div class="alert alert-error">` + err.Error() + `</div>`))
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/borrowings")
		return
	}

	http.Redirect(w, r, "/admin/borrowings", http.StatusSeeOther)
}

func (h *Handler) Return(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	borrowing, err := h.service.ReturnBook(id)
	if err != nil {
		if r.Header.Get("HX-Request") == "true" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		// Return success with fine info
		if borrowing.Fine > 0 {
			w.Write([]byte(`<div class="alert alert-warning">Buku berhasil dikembalikan. Denda: Rp ` + strconv.FormatFloat(borrowing.Fine, 'f', 0, 64) + `</div>`))
		} else {
			w.Write([]byte(`<div class="alert alert-success">Buku berhasil dikembalikan.</div>`))
		}
		w.Header().Set("HX-Trigger", "refreshTable")
		return
	}

	http.Redirect(w, r, "/admin/borrowings", http.StatusSeeOther)
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
