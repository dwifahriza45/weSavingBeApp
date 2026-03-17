package app

import (
	"BE_WE_SAVING/internal/domain/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerAuthRoutes(r chi.Router, authHandler *auth.AuthHandler, authMd func(http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.With(authMd).Get("/me", authHandler.Me)
	})
}
