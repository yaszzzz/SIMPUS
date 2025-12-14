package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"simpus/internal/middleware"
	"simpus/internal/models"
	"simpus/internal/services"
)

type MemberHandler struct {
	memberService *services.MemberService
	templates     *template.Template
}

func NewMemberHandler(memberService *services.MemberService, templates *template.Template) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
		templates:     templates,
	}
}

func (h *MemberHandler) Index(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	search := r.URL.Query().Get("search")

	members, total, err := h.memberService.GetMembers(page, 10, search)
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

func (h *MemberHandler) Create(w http.ResponseWriter, r *http.Request) {
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

func (h *MemberHandler) Store(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.memberService.CreateMember(data)
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

func (h *MemberHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	member, err := h.memberService.GetMember(id)
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

func (h *MemberHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	err := h.memberService.UpdateMember(id, data)
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

func (h *MemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	err := h.memberService.DeleteMember(id)
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

func (h *MemberHandler) render(w http.ResponseWriter, name string, data interface{}) {
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

func (h *MemberHandler) renderPartial(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
