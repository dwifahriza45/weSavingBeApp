package app

import (
	"BE_WE_SAVING/internal/app/middleware"
	"BE_WE_SAVING/internal/shared/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) routes() {
	s.Router.Use(corsMiddleware())

	authMd := middleware.Auth(s.JWTSecret)
	handlers := s.buildRouteHandlers()

	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			utils.JSON(w, http.StatusOK, "OK", "Test OK", false, map[string]string{
				"result": "ok",
			})
		})

		registerAuthRoutes(r, handlers.auth, authMd)

		r.Group(func(r chi.Router) {
			r.Use(authMd)

			registerCategoryRoutes(r, handlers.categories)
			registerCategoryBudgetRoutes(r, handlers.categoriesBudget)
			registerSalaryRoutes(r, handlers.salaries)
		})
	})
}

func corsMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
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
	})
}
