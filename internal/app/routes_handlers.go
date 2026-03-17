package app

import (
	"BE_WE_SAVING/internal/domain/auth"
	"BE_WE_SAVING/internal/domain/categories"
	cetegoriesbudget "BE_WE_SAVING/internal/domain/cetegories_budget"
	"BE_WE_SAVING/internal/domain/salaries"
	"BE_WE_SAVING/internal/domain/users"
)

type routeHandlers struct {
	auth             *auth.AuthHandler
	categories       *categories.CategoriesHandler
	categoriesBudget *cetegoriesbudget.CategoriesBudgetHandler
	salaries         *salaries.SalariesHandler
}

func (s *Server) buildRouteHandlers() routeHandlers {
	authRepo := users.NewUserRepository(s.DB.DB)
	authService := auth.NewAuthService(authRepo, s.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	categoriesRepo := categories.NewCategoriesRepository(s.DB.DB)
	categoriesService := categories.NewCategoriesService(categoriesRepo)
	categoriesHandler := categories.NewCategoriesHandler(categoriesService)

	categoriesBudgetRepo := cetegoriesbudget.NewCategoriesBudgetRepository(s.DB.DB)
	categoriesBudgetService := cetegoriesbudget.NewCategoriesBudgetService(categoriesBudgetRepo)
	categoriesBudgetHandler := cetegoriesbudget.NewCategoriesBudgetHandler(categoriesBudgetService)

	salariesRepo := salaries.NewSalariesRepository(s.DB.DB)
	salariesService := salaries.NewSalariesService(salariesRepo)
	salariesHandler := salaries.NewSalariesHandler(salariesService)

	return routeHandlers{
		auth:             authHandler,
		categories:       categoriesHandler,
		categoriesBudget: categoriesBudgetHandler,
		salaries:         salariesHandler,
	}
}
