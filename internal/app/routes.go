package app

import (
	"BE_WE_SAVING/internal/app/middleware"
	"BE_WE_SAVING/internal/domain/auth"
	"BE_WE_SAVING/internal/domain/categories"
	"BE_WE_SAVING/internal/domain/users"
	"BE_WE_SAVING/internal/shared/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) routes() {
	r := s.Router

	// =====================
	// CORS FOR WEB
	// =====================
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // Vite
			"http://localhost:3000", // React CRA
			// "https://domain-production.com",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300, // cache preflight (5 menit)
	}))

	// middleware
	authMd := middleware.Auth(s.JWTSecret)

	// Dependency Auth
	authRepo := users.NewUserRepository(s.DB.DB)
	authService := auth.NewAuthService(authRepo, s.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// Dependency Categories
	categoriesRepo := categories.NewCategoriesRepository(s.DB.DB)
	categoriesService := categories.NewCategoriesService(categoriesRepo)
	categoriesHandler := categories.NewCategoriesHandler(categoriesService)

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			utils.JSON(w, http.StatusOK, "OK", "Test OK", false, map[string]string{
				"result": "ok",
			})
		})

		// =====================
		// AUTH (public + protected)
		// =====================
		r.Route("/auth", func(r chi.Router) {

			// public
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)
			r.With(authMd).Get("/me", authHandler.Me)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMd)

			r.Route("/categories", func(r chi.Router) {
				r.Post("/create", categoriesHandler.Create)
			})
		})
	})

}
