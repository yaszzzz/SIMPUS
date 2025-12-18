package auth

import (
	"html/template"
	"net/http"
	"path/filepath"
	"simpus/internal/models"
	"time"
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

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login Admin - SIMPUS",
		"Error": r.URL.Query().Get("error"),
	}
	h.render(w, "auth/login.html", data)
}

func (h *Handler) MemberLoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login Anggota - SIMPUS",
		"Error": r.URL.Query().Get("error"),
	}
	h.render(w, "auth/login-member.html", data)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error=Form tidak valid", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, token, err := h.service.LoginAdmin(username, password)
	if err != nil {
		http.Redirect(w, r, "/login?error="+err.Error(), http.StatusSeeOther)
		return
	}

	h.setTokenCookie(w, token)

	_ = user // Can be used for logging
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *Handler) MemberLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login/member?error=Form tidak valid", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	member, token, err := h.service.LoginMember(email, password)
	if err != nil {
		http.Redirect(w, r, "/login/member?error="+err.Error(), http.StatusSeeOther)
		return
	}

	h.setTokenCookie(w, token)

	_ = member // Can be used for logging
	http.Redirect(w, r, "/member/dashboard", http.StatusSeeOther)
}

func (h *Handler) RegisterMemberPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Registrasi Anggota - SIMPUS",
		"Error": r.URL.Query().Get("error"),
	}
	h.render(w, "auth/register-member.html", data)
}

func (h *Handler) RegisterMember(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/register?error=Form tidak valid", http.StatusSeeOther)
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

	if data.Name == "" || data.Email == "" || data.Password == "" {
		http.Redirect(w, r, "/register?error=Data tidak lengkap", http.StatusSeeOther)
		return
	}

	_, err := h.service.RegisterMember(data)
	if err != nil {
		http.Redirect(w, r, "/register?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login/member?success=Registrasi berhasil, silakan login", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := h.templates.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.ParseFiles(
		filepath.Join("templates", "layouts", "auth.html"),
		filepath.Join("templates", name),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "auth.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
