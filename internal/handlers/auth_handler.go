package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"simpus/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	templates   *template.Template
}

func NewAuthHandler(authService *services.AuthService, templates *template.Template) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		templates:   templates,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login Admin - SIMPUS",
		"Error": r.URL.Query().Get("error"),
	}
	h.render(w, "auth/login.html", data)
}

func (h *AuthHandler) MemberLoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login Anggota - SIMPUS",
		"Error": r.URL.Query().Get("error"),
	}
	h.render(w, "auth/login-member.html", data)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error=Form tidak valid", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, token, err := h.authService.LoginAdmin(username, password)
	if err != nil {
		http.Redirect(w, r, "/login?error="+err.Error(), http.StatusSeeOther)
		return
	}

	h.setTokenCookie(w, token)

	_ = user // Can be used for logging
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) MemberLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login/member?error=Form tidak valid", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	member, token, err := h.authService.LoginMember(email, password)
	if err != nil {
		http.Redirect(w, r, "/login/member?error="+err.Error(), http.StatusSeeOther)
		return
	}

	h.setTokenCookie(w, token)

	_ = member // Can be used for logging
	http.Redirect(w, r, "/member/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(
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
