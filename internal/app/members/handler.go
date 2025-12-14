package members

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"simpus/internal/middleware"
	"simpus/internal/models"
)

type Handler struct {
	service   *Service
	templates *template.Template
}

func NewHandler(service *Service, templates *template.Template) *Handler {
	return &Handler{
		service:   service,
		templates: templates,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	search := r.URL.Query().Get("search")

	members, total, err := h.service.GetMembers(page, 10, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + 10 - 1) / 10

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":      "Manajemen Anggota - SIMPUS",
		"Members":    members,
		"Total":      total,
		"Page":       page,
		"TotalPages": totalPages,
		"Search":     search,
		"User":       claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/members/table.html", data)
		return
	}

	h.render(w, "admin/members/index.html", data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title": "Tambah Anggota - SIMPUS",
		"User":  claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/members/form.html", data)
		return
	}

	h.render(w, "admin/members/create.html", data)
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	data := &models.MemberCreate{
		Name:       r.FormValue("name"),
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		Phone:      r.FormValue("phone"),
		MemberType: r.FormValue("member_type"),
		Address:    r.FormValue("address"),
	}

	_, err := h.service.CreateMember(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/members")
		return
	}

	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	member, err := h.service.GetMember(id)
	if err != nil {
		http.Error(w, "Anggota tidak ditemukan", http.StatusNotFound)
		return
	}

	claims := middleware.GetUserFromContext(r.Context())

	data := map[string]interface{}{
		"Title":  "Edit Anggota - SIMPUS",
		"Member": member,
		"User":   claims,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderPartial(w, "admin/members/form.html", data)
		return
	}

	h.render(w, "admin/members/edit.html", data)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	isActive := r.FormValue("is_active") == "on" || r.FormValue("is_active") == "true"

	data := &models.MemberUpdate{
		Name:       r.FormValue("name"),
		Email:      r.FormValue("email"),
		Phone:      r.FormValue("phone"),
		MemberType: r.FormValue("member_type"),
		Address:    r.FormValue("address"),
		IsActive:   isActive,
	}

	err := h.service.UpdateMember(id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/members")
		return
	}

	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	err := h.service.DeleteMember(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
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
