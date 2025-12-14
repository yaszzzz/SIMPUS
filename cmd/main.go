package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"simpus/config"
	"simpus/database"
	"simpus/internal/handlers"
	authMiddleware "simpus/internal/middleware"
	"simpus/internal/repository"
	"simpus/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Connected to database successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	memberRepo := repository.NewMemberRepository(database.DB)
	categoryRepo := repository.NewCategoryRepository(database.DB)
	authorRepo := repository.NewAuthorRepository(database.DB)
	bookRepo := repository.NewBookRepository(database.DB)
	borrowRepo := repository.NewBorrowingRepository(database.DB)
	notifRepo := repository.NewNotificationRepository(database.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, memberRepo, cfg)
	bookService := services.NewBookService(bookRepo, categoryRepo, authorRepo)
	memberService := services.NewMemberService(memberRepo)
	borrowService := services.NewBorrowService(borrowRepo, bookRepo, memberRepo, notifRepo)
	notifService := services.NewNotificationService(notifRepo)

	// Initialize template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"slice": func(s string, start, end int) string {
			if start >= len(s) {
				return ""
			}
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"deref": func(i *int) int {
			if i == nil {
				return 0
			}
			return *i
		},
	}
	templates := template.New("").Funcs(funcMap)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, templates)
	dashboardHandler := handlers.NewDashboardHandler(bookService, memberService, borrowService, templates)
	bookHandler := handlers.NewBookHandler(bookService, templates)
	memberHandler := handlers.NewMemberHandler(memberService, templates)
	borrowHandler := handlers.NewBorrowHandler(borrowService, bookService, memberService, templates)
	categoryHandler := handlers.NewCategoryHandler(bookService, templates)
	authorHandler := handlers.NewAuthorHandler(bookService, templates)

	// Initialize middleware
	authMw := authMiddleware.NewAuthMiddleware(authService)

	// Unused variable fix
	_ = notifService

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// Static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)
	r.Get("/login/member", authHandler.MemberLoginPage)
	r.Post("/login/member", authHandler.MemberLogin)
	r.Get("/logout", authHandler.Logout)

	// Admin routes (protected)
	r.Route("/admin", func(r chi.Router) {
		r.Use(authMw.RequireAuth)
		r.Use(authMw.RequireAdmin)

		// Dashboard
		r.Get("/dashboard", dashboardHandler.AdminDashboard)

		// Books
		r.Get("/books", bookHandler.Index)
		r.Get("/books/create", bookHandler.Create)
		r.Post("/books", bookHandler.Store)
		r.Get("/books/{id}/edit", bookHandler.Edit)
		r.Post("/books/{id}", bookHandler.Update)
		r.Delete("/books/{id}", bookHandler.Delete)

		// Categories
		r.Get("/categories", categoryHandler.Index)
		r.Post("/categories", categoryHandler.Store)
		r.Post("/categories/{id}", categoryHandler.Update)
		r.Delete("/categories/{id}", categoryHandler.Delete)

		// Authors
		r.Get("/authors", authorHandler.Index)
		r.Post("/authors", authorHandler.Store)
		r.Post("/authors/{id}", authorHandler.Update)
		r.Delete("/authors/{id}", authorHandler.Delete)

		// Members
		r.Get("/members", memberHandler.Index)
		r.Get("/members/create", memberHandler.Create)
		r.Post("/members", memberHandler.Store)
		r.Get("/members/{id}/edit", memberHandler.Edit)
		r.Post("/members/{id}", memberHandler.Update)
		r.Delete("/members/{id}", memberHandler.Delete)

		// Borrowings
		r.Get("/borrowings", borrowHandler.Index)
		r.Get("/borrowings/create", borrowHandler.Create)
		r.Post("/borrowings", borrowHandler.Store)
		r.Post("/borrowings/{id}/return", borrowHandler.Return)

		// Reports
		r.Get("/reports", borrowHandler.Report)
	})

	// Member routes (protected)
	r.Route("/member", func(r chi.Router) {
		r.Use(authMw.RequireAuth)
		r.Use(authMw.RequireMember)

		r.Get("/dashboard", dashboardHandler.MemberDashboard)
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 SIMPUS server running at http://%s", addr)
	log.Printf("📚 Login dengan username: admin, password: admin123")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
