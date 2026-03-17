package cetegoriesbudget

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockCategoriesBudgetService struct {
	createFunc          func(ctx context.Context, categoryID, allocatedAmount, period string) error
	getByCategoryIDFunc func(ctx context.Context, categoryID string) (*CategoriesBudget, error)
}

func (m *mockCategoriesBudgetService) Create(ctx context.Context, categoryID, allocatedAmount, period string) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, categoryID, allocatedAmount, period)
	}

	return nil
}

func (m *mockCategoriesBudgetService) GetByCategoryID(ctx context.Context, categoryID string) (*CategoriesBudget, error) {
	if m.getByCategoryIDFunc != nil {
		return m.getByCategoryIDFunc(ctx, categoryID)
	}

	return nil, nil
}

func TestCreateCategoriesBudgetHandler_Success(t *testing.T) {
	mockService := &mockCategoriesBudgetService{
		createFunc: func(ctx context.Context, categoryID, allocatedAmount, period string) error {
			assert.Equal(t, "CAT-001", categoryID)
			assert.Equal(t, "1000000", allocatedAmount)
			assert.Equal(t, "", period)
			return nil
		},
	}

	handler := NewCategoriesBudgetHandler(mockService)

	body := `{
		"category_id": "CAT-001",
		"allocated_amount": "1000000"
	}`

	req := httptest.NewRequest(http.MethodPost, "/category-budgets/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Category budget created")
}

func TestCreateCategoriesBudgetHandler_InvalidJSON(t *testing.T) {
	mockService := &mockCategoriesBudgetService{}
	handler := NewCategoriesBudgetHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/category-budgets/create", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCategoriesBudgetHandler_ValidationError(t *testing.T) {
	mockService := &mockCategoriesBudgetService{}
	handler := NewCategoriesBudgetHandler(mockService)

	body := `{
		"allocated_amount": "abc"
	}`

	req := httptest.NewRequest(http.MethodPost, "/category-budgets/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Validation Failed")
	assert.Contains(t, rec.Body.String(), "category_id is required")
	assert.Contains(t, rec.Body.String(), "allocated_amount must be numeric")
}

func TestCreateCategoriesBudgetHandler_ServiceError(t *testing.T) {
	mockService := &mockCategoriesBudgetService{
		createFunc: func(ctx context.Context, categoryID, allocatedAmount, period string) error {
			return errors.New("failed create category budget")
		},
	}

	handler := NewCategoriesBudgetHandler(mockService)

	body := `{
		"category_id": "CAT-001",
		"allocated_amount": "1000000"
	}`

	req := httptest.NewRequest(http.MethodPost, "/category-budgets/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed create category budget")
}

func TestGetCategoriesBudgetByCategoryIDHandler_Success(t *testing.T) {
	mockService := &mockCategoriesBudgetService{
		getByCategoryIDFunc: func(ctx context.Context, categoryID string) (*CategoriesBudget, error) {
			assert.Equal(t, "CAT-001", categoryID)

			return &CategoriesBudget{
				ID:              1,
				BUDGET_ID:       "BUD-001",
				USER_ID:         "USER-001",
				CATEGORY_ID:     "CAT-001",
				AllocatedAmount: "1000000",
				UsedAmount:      "250000",
				Period:          "monthly",
			}, nil
		},
	}

	handler := NewCategoriesBudgetHandler(mockService)

	r := chi.NewRouter()
	r.Get("/category-budgets/category/{id}", handler.GetByCategoryID)

	req := httptest.NewRequest(http.MethodGet, "/category-budgets/category/CAT-001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Category budget fetched")
	assert.Contains(t, rec.Body.String(), `"category_id":"CAT-001"`)
}

func TestGetCategoriesBudgetByCategoryIDHandler_NotFound(t *testing.T) {
	mockService := &mockCategoriesBudgetService{
		getByCategoryIDFunc: func(ctx context.Context, categoryID string) (*CategoriesBudget, error) {
			return nil, ErrCategoryBudgetNotFound
		},
	}

	handler := NewCategoriesBudgetHandler(mockService)

	r := chi.NewRouter()
	r.Get("/category-budgets/category/{id}", handler.GetByCategoryID)

	req := httptest.NewRequest(http.MethodGet, "/category-budgets/category/CAT-001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrCategoryBudgetNotFound.Error())
}

func TestGetCategoriesBudgetByCategoryIDHandler_Unauthorized(t *testing.T) {
	mockService := &mockCategoriesBudgetService{
		getByCategoryIDFunc: func(ctx context.Context, categoryID string) (*CategoriesBudget, error) {
			return nil, ErrInvalidCredentials
		},
	}

	handler := NewCategoriesBudgetHandler(mockService)

	r := chi.NewRouter()
	r.Get("/category-budgets/category/{id}", handler.GetByCategoryID)

	req := httptest.NewRequest(http.MethodGet, "/category-budgets/category/CAT-001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInvalidCredentials.Error())
}

func TestGetCategoriesBudgetByCategoryIDHandler_InternalError(t *testing.T) {
	mockService := &mockCategoriesBudgetService{
		getByCategoryIDFunc: func(ctx context.Context, categoryID string) (*CategoriesBudget, error) {
			return nil, errors.New("boom")
		},
	}

	handler := NewCategoriesBudgetHandler(mockService)

	r := chi.NewRouter()
	r.Get("/category-budgets/category/{id}", handler.GetByCategoryID)

	req := httptest.NewRequest(http.MethodGet, "/category-budgets/category/CAT-001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "boom")
}
