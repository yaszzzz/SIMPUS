package books

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"simpus/internal/middleware"
	"simpus/internal/models"
)

type CategoryHandler struct {
	service   *Service
	templates *template.Template
}

func NewCategoryHandler(service *Service, templates *template.Template) *CategoryHandler {
	return &CategoryHandler{
		service:   service,
		templates: templates,
	}
}

func (h *CategoryHandler) Index(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":      "Manajemen Kategori - SIMPUS",
		"Categories": categories,
		"User":       claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/categories/table.html", data)
		return
	}

	h.render(w, "admin/categories/index.html", data)
}

func (h *CategoryHandler) Store(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	data := &models.CategoryCreate{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	_, err := h.service.CreateCategory(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", "refreshTable")
		w.WriteHeader(http.StatusCreated)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	data := &models.CategoryCreate{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	err := h.service.UpdateCategory(id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", "refreshTable")
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	err := h.service.DeleteCategory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoryHandler) render(w http.ResponseWriter, name string, data interface{}) {
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

func (h *CategoryHandler) renderPartial(w http.ResponseWriter, name string, data interface{}) {
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
