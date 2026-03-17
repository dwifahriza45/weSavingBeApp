package salaries

import (
	"BE_WE_SAVING/internal/shared/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type SalariesHandler struct {
	salariesService SalariesService
}

func NewSalariesHandler(salariesService SalariesService) *SalariesHandler {
	return &SalariesHandler{
		salariesService: salariesService,
	}
}

type salariesRequest struct {
	Amount      string `json:"amount" validate:"required,numeric"`
	Source      string `json:"source" validate:"required"`
	Description string `json:"description"`
}

func (h *SalariesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req salariesRequest

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

	err := h.salariesService.Create(
		r.Context(),
		req.Amount,
		req.Source,
		req.Description,
	)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		return
	}

	utils.JSON(w, http.StatusCreated, "OK", "Salary Created", false, nil)
}

func (h *SalariesHandler) CheckSalary(w http.ResponseWriter, r *http.Request) {
	result, err := h.salariesService.CheckSalary(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Salary checked", false, result)
}

func (h *SalariesHandler) GetTotalSalary(w http.ResponseWriter, r *http.Request) {
	total, err := h.salariesService.GetTotalSalary(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Total salary fetched", false, total)
}
