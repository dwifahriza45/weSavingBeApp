package app

import (
	"BE_WE_SAVING/internal/Infrastructures/config"
	"BE_WE_SAVING/internal/Infrastructures/database"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router    *chi.Mux
	DB        *database.Store
	JWTSecret string
}

func NewServer() *Server {
	cfg := config.Load()
	db := database.NewStore(cfg)

	s := &Server{
		Router:    chi.NewRouter(),
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	s.routes()
	return s
}
