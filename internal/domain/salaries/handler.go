package salaries

import (
	"BE_WE_SAVING/internal/shared/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (h *SalariesHandler) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	salariesData, err := h.salariesService.GetAllByUserID(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	response := make([]Salaries, 0, len(salariesData))

	for _, salary := range salariesData {
		response = append(response, Salaries{
			ID:          salary.ID,
			SalaryID:    salary.SalaryID,
			UserID:      salary.UserID,
			Amount:      salary.Amount,
			Source:      salary.Source,
			Description: salary.Description,
			ReceivedAt:  salary.ReceivedAt,
		})
	}

	utils.JSON(w, http.StatusOK, "OK", "Salaries fetched", false, response)
}

func (h *SalariesHandler) GetBySalaryID(w http.ResponseWriter, r *http.Request) {
	salaryID := chi.URLParam(r, "id")
	if salaryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "salary id is required", true)
		return
	}

	salary, err := h.salariesService.GetBySalaryID(r.Context(), salaryID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case errors.Is(err, ErrSalaryNotFound):
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	response := Salaries{
		ID:          salary.ID,
		SalaryID:    salary.SalaryID,
		UserID:      salary.UserID,
		Amount:      salary.Amount,
		Source:      salary.Source,
		Description: salary.Description,
		ReceivedAt:  salary.ReceivedAt,
	}

	utils.JSON(w, http.StatusOK, "OK", "Salary fetched", false, response)
}

func (h *SalariesHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req salariesRequest

	salaryID := chi.URLParam(r, "id")
	if salaryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "salary id is required", true)
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

	err := h.salariesService.Update(
		r.Context(),
		salaryID,
		req.Amount,
		req.Source,
		req.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case errors.Is(err, ErrInvalidAmount):
			utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		case errors.Is(err, ErrSalaryNotFound):
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Salary Updated", false, nil)
}

func (h *SalariesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	salaryID := chi.URLParam(r, "id")
	if salaryID == "" {
		utils.JSONError(w, http.StatusBadRequest, "NOK", "salary id is required", true)
		return
	}

	err := h.salariesService.Delete(r.Context(), salaryID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		case errors.Is(err, ErrSalaryNotFound):
			utils.JSONError(w, http.StatusNotFound, "NOK", err.Error(), true)
		default:
			utils.JSONError(w, http.StatusInternalServerError, "NOK", err.Error(), true)
		}
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Salary Deleted", false, nil)
}
