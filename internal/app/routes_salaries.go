package app

import (
	"BE_WE_SAVING/internal/domain/salaries"

	"github.com/go-chi/chi/v5"
)

func registerSalaryRoutes(r chi.Router, salariesHandler *salaries.SalariesHandler) {
	r.Route("/salary", func(r chi.Router) {
		r.Post("/create", salariesHandler.Create)
		r.Get("/all", salariesHandler.GetAllByUserID)
		r.Get("/check", salariesHandler.CheckSalary)
		r.Get("/total", salariesHandler.GetTotalSalary)
		r.Get("/{id}", salariesHandler.GetBySalaryID)
		r.Put("/{id}", salariesHandler.Update)
		r.Delete("/{id}", salariesHandler.Delete)
	})
}
