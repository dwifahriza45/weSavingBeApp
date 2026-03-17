package app

import (
	"BE_WE_SAVING/internal/domain/categories"

	"github.com/go-chi/chi/v5"
)

func registerCategoryRoutes(r chi.Router, categoriesHandler *categories.CategoriesHandler) {
	r.Route("/categories", func(r chi.Router) {
		r.Post("/create", categoriesHandler.Create)
		r.Get("/all", categoriesHandler.GetAll)
		r.Get("/{id}", categoriesHandler.GetByCategoryID)
		r.Put("/{id}", categoriesHandler.Update)
		r.Delete("/{id}", categoriesHandler.Delete)
	})
}
