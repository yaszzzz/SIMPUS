package notifications

import (
	"html/template"
	"net/http"
	"path/filepath"
	"simpus/internal/middleware"
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

func (h *Handler) MemberIndex(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	notifications, err := h.service.GetMemberNotifications(claims.UserID, 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":         "Notifikasi - SIMPUS",
		"Notifications": notifications,
		"User":          claims,
	}

	h.renderMember(w, "member/notifications/index.html", data)
}

func (h *Handler) renderMember(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := h.templates.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.ParseFiles(
		filepath.Join("templates", "layouts", "member.html"),
		filepath.Join("templates", "components", "member_navbar.html"),
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
