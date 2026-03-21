package cetegoriesbudget

import (
	"BE_WE_SAVING/internal/shared/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CategoriesBudgetHandler struct {
	categoriesBudgetService CategoriesBudgetService
}

func NewCategoriesBudgetHandler(categoriesBudgetService CategoriesBudgetService) *CategoriesBudgetHandler {
	return &CategoriesBudgetHandler{
		categoriesBudgetService: categoriesBudgetService,
	}
}

type categoriesBudgetRequest struct {
	CategoryID      string `json:"category_id" validate:"required"`
	AllocatedAmount string `json:"allocated_amount" validate:"required,numeric"`
}

type updateCategoriesBudgetRequest struct {
	AllocatedAmount string `json:"allocated_amount" validate:"required,numeric"`
}

func (h *CategoriesBudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req categoriesBudgetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "Invalid JSON body", true)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		errMap := make(map[string]string)

		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			switch e.Field() {
			case "CategoryID":
				field = "category_id"
			case "AllocatedAmount":
				field = "allocated_amount"
			}

			switch e.Tag() {
			case "required":
				errMap[field] = field + " is required"
			case "numeric":
				errMap[field] = field + " must be numeric"
			default:
				errMap[field] = "invalid value"
			}
		}

		utils.JSONErrorWithData(
			w,
			http.StatusBadRequest,
			"NOK",
			"Validation Failed",
			true,
			errMap,
		)
		return
	}

	err := h.categoriesBudgetService.Create(
		r.Context(),
		req.CategoryID,
		req.AllocatedAmount,
	)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		return
	}

	utils.JSON(w, http.StatusCreated, "OK", "Category budget created", false, nil)
}

func (h *CategoriesBudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateCategoriesBudgetRequest

	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "budget id is required", true)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "Invalid JSON body", true)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		errMap := make(map[string]string)

		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			if e.Field() == "AllocatedAmount" {
				field = "allocated_amount"
			}

			switch e.Tag() {
			case "required":
				errMap[field] = field + " is required"
			case "numeric":
				errMap[field] = field + " must be numeric"
			default:
				errMap[field] = "invalid value"
			}
		}

		utils.JSONErrorWithData(
			w,
			http.StatusBadRequest,
			"NOK",
			"Validation Failed",
			true,
			errMap,
		)
		return
	}

	err := h.categoriesBudgetService.Update(r.Context(), budgetID, req.AllocatedAmount)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case ErrInvalidAllocatedAmount:
			utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		case ErrCategoryBudgetNotFound:
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Category budget updated", false, nil)
}

func (h *CategoriesBudgetHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "category id is required", true)
		return
	}

	categoryBudget, err := h.categoriesBudgetService.GetByCategoryID(r.Context(), categoryID)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case ErrCategoryBudgetNotFound:
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	response := CategoriesBudget{
		ID:              categoryBudget.ID,
		BUDGET_ID:       categoryBudget.BUDGET_ID,
		USER_ID:         categoryBudget.USER_ID,
		CATEGORY_ID:     categoryBudget.CATEGORY_ID,
		AllocatedAmount: categoryBudget.AllocatedAmount,
		UsedAmount:      categoryBudget.UsedAmount,
	}

	utils.JSON(w, http.StatusOK, "OK", "Category budget fetched", false, response)
}

func (h *CategoriesBudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "budget id is required", true)
		return
	}

	err := h.categoriesBudgetService.Delete(r.Context(), budgetID)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case ErrCategoryBudgetNotFound:
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Category budget deleted", false, nil)
}
