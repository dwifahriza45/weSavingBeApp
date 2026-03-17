package app

import (
	"BE_WE_SAVING/internal/domain/salaries"

	"github.com/go-chi/chi/v5"
)

func registerSalaryRoutes(r chi.Router, salariesHandler *salaries.SalariesHandler) {
	r.Route("/salary", func(r chi.Router) {
		r.Post("/create", salariesHandler.Create)
		r.Get("/check", salariesHandler.CheckSalary)
		r.Get("/total", salariesHandler.GetTotalSalary)
	})
}
