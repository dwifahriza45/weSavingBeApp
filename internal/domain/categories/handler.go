package categories

import (
	"BE_WE_SAVING/internal/shared/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CategoriesHandler struct {
	categoriesService CategoriesService
}

func NewCategoriesHandler(categoriesService CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{
		categoriesService: categoriesService,
	}
}

type categoriesRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req categoriesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "Invalid JSON body", true)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		errMap := make(map[string]string)

		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())

			switch e.Tag() {
			case "required":
				errMap[field] = field + " is required"
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

	err := h.categoriesService.Create(
		r.Context(),
		req.Name,
		req.Description,
	)

	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		return
	}

	utils.JSON(w, http.StatusCreated, "OK", "Categories Created", false, nil)
}

func (h *CategoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req categoriesRequest

	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "category id is required", true)
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

			switch e.Tag() {
			case "required":
				errMap[field] = field + " is required"
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

	err := h.categoriesService.Update(
		r.Context(),
		req.Name,
		req.Description,
		categoryID,
	)

	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		return
	}

	utils.JSON(w, http.StatusCreated, "OK", "Categories Updated", false, nil)
}

func (h *CategoriesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoriesService.GetAllByUserID(r.Context())
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	response := make([]Categories, 0, len(categories))

	for _, c := range categories {
		response = append(response, Categories{
			ID:          c.ID,
			USER_ID:     c.USER_ID,
			CATEGORY_ID: c.CATEGORY_ID,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	utils.JSON(w, http.StatusOK, "OK", "Categories fetched", false, response)
}

func (h *CategoriesHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "category id is required", true)
		return
	}

	category, err := h.categoriesService.GetByCategoryID(r.Context(), categoryID)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)

		case ErrCategoryNotFound:
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)

		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	response := Categories{
		ID:          category.ID,
		USER_ID:     category.USER_ID,
		CATEGORY_ID: category.CATEGORY_ID,
		Name:        category.Name,
		Description: category.Description,
	}

	utils.JSON(w, http.StatusOK, "OK", "Categories fetched", false, response)
}

func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "category id is required", true)
		return
	}

	err := h.categoriesService.Delete(r.Context(), categoryID)
	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case errors.Is(err, ErrCategoryHasBudget):
			utils.JSONError(w, http.StatusConflict, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", "Something went wrong", true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Category Deleted", false, nil)
}
