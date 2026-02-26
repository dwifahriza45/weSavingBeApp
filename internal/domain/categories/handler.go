package categories

import (
	"BE_WE_SAVING/internal/shared/utils"
	"encoding/json"
	"net/http"
	"strings"

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
