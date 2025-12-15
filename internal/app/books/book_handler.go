package books

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"simpus/internal/middleware"
	"simpus/internal/models"
)

type BookHandler struct {
	service   *Service
	templates *template.Template
}

func NewBookHandler(service *Service, templates *template.Template) *BookHandler {
	return &BookHandler{
		service:   service,
		templates: templates,
	}
}

func (h *BookHandler) Index(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	search := r.URL.Query().Get("search")
	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category"))

	filter := models.BookFilter{
		Search:     search,
		CategoryID: categoryID,
		Page:       page,
		Limit:      10,
	}

	books, total, err := h.service.GetBooks(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, _ := h.service.GetCategories()

	totalPages := (total + filter.Limit - 1) / filter.Limit

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":      "Manajemen Buku - SIMPUS",
		"Books":      books,
		"Total":      total,
		"Page":       page,
		"TotalPages": totalPages,
		"Search":     search,
		"CategoryID": categoryID,
		"Categories": categories,
		"User":       claims,
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/books/table.html", data)
		return
	}

	h.render(w, "admin/books/index.html", data)
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.service.GetCategories()
	authors, _ := h.service.GetAuthors()

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":      "Tambah Buku - SIMPUS",
		"Categories": categories,
		"Authors":    authors,
		"User":       claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/books/form.html", data)
		return
	}

	h.render(w, "admin/books/create.html", data)
}

func (h *BookHandler) Store(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
	authorID, _ := strconv.Atoi(r.FormValue("author_id"))
	publishYear, _ := strconv.Atoi(r.FormValue("publish_year"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	data := &models.BookCreate{
		ISBN:        r.FormValue("isbn"),
		Title:       r.FormValue("title"),
		CategoryID:  categoryID,
		AuthorID:    authorID,
		Publisher:   r.FormValue("publisher"),
		PublishYear: publishYear,
		Stock:       stock,
		Description: r.FormValue("description"),
	}

	_, err := h.service.CreateBook(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/books")
		return
	}

	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

func (h *BookHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	book, err := h.service.GetBook(id)
	if err != nil {
		http.Error(w, "Buku tidak ditemukan", http.StatusNotFound)
		return
	}

	categories, _ := h.service.GetCategories()
	authors, _ := h.service.GetAuthors()

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":      "Edit Buku - SIMPUS",
		"Book":       book,
		"Categories": categories,
		"Authors":    authors,
		"User":       claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/books/form.html", data)
		return
	}

	h.render(w, "admin/books/edit.html", data)
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
	authorID, _ := strconv.Atoi(r.FormValue("author_id"))
	publishYear, _ := strconv.Atoi(r.FormValue("publish_year"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	data := &models.BookUpdate{
		ISBN:        r.FormValue("isbn"),
		Title:       r.FormValue("title"),
		CategoryID:  categoryID,
		AuthorID:    authorID,
		Publisher:   r.FormValue("publisher"),
		PublishYear: publishYear,
		Stock:       stock,
		Description: r.FormValue("description"),
	}

	err := h.service.UpdateBook(id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/books")
		return
	}

	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	err := h.service.DeleteBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

func (h *BookHandler) render(w http.ResponseWriter, name string, data interface{}) {
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

func (h *BookHandler) renderPartial(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := h.templates.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.ParseFiles(filepath.Join("templates", name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
