package books

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"simpus/internal/middleware"
	"simpus/internal/models"
)

type AuthorHandler struct {
	service   *Service
	templates *template.Template
}

func NewAuthorHandler(service *Service, templates *template.Template) *AuthorHandler {
	return &AuthorHandler{
		service:   service,
		templates: templates,
	}
}

func (h *AuthorHandler) Index(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.GetAuthors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":   "Manajemen Penulis - SIMPUS",
		"Authors": authors,
		"User":    claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/authors/table.html", data)
		return
	}

	h.render(w, "admin/authors/index.html", data)
}

func (h *AuthorHandler) Store(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	data := &models.AuthorCreate{
		Name: r.FormValue("name"),
		Bio:  r.FormValue("bio"),
	}

	_, err := h.service.CreateAuthor(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", "refreshTable")
		w.WriteHeader(http.StatusCreated)
		return
	}

	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

func (h *AuthorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	data := &models.AuthorCreate{
		Name: r.FormValue("name"),
		Bio:  r.FormValue("bio"),
	}

	err := h.service.UpdateAuthor(id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", "refreshTable")
		return
	}

	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

func (h *AuthorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	err := h.service.DeleteAuthor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

func (h *AuthorHandler) render(w http.ResponseWriter, name string, data interface{}) {
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

func (h *AuthorHandler) renderPartial(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
