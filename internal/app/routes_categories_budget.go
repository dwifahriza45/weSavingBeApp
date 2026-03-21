package app

import (
	cetegoriesbudget "BE_WE_SAVING/internal/domain/cetegories_budget"

	"github.com/go-chi/chi/v5"
)

func registerCategoryBudgetRoutes(r chi.Router, categoriesBudgetHandler *cetegoriesbudget.CategoriesBudgetHandler) {
	r.Route("/category-budgets", func(r chi.Router) {
		r.Post("/create", categoriesBudgetHandler.Create)
		r.Get("/category/{id}", categoriesBudgetHandler.GetByCategoryID)
		r.Put("/budget/{id}", categoriesBudgetHandler.Update)
		r.Delete("/budget/{id}", categoriesBudgetHandler.Delete)
	})
}
