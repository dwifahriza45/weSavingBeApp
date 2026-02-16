package auth

import (
	"BE_WE_SAVING/internal/app/middleware"
	"BE_WE_SAVING/internal/shared/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

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
			case "email":
				errMap[field] = "email format is not valid"
			case "min":
				errMap[field] = field + " must be at least " + e.Param() + " characters"
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

	err := h.authService.Register(
		r.Context(),
		req.Username,
		req.Fullname,
		req.Email,
		req.Password,
	)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "NOK", err.Error(), true)
		return
	}

	utils.JSON(w, http.StatusCreated, "OK", "User Registered Successfully", false, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

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
			case "min":
				errMap[field] = field + " must be at least " + e.Param() + " characters"
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

	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, "NOK", err.Error(), true)
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "Login Success", false, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "NOK", "Unauthorized", true)
		return
	}

	user, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "NOK", "User Not Found", true)
		return
	}

	utils.JSON(w, http.StatusOK, "OK", "success", false, map[string]interface{}{
		"user_id":  user.USER_ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
