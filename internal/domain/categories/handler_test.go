package categories

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockCategoriesService struct {
	countByDateFunc     func(ctx context.Context, date string) (int, error)
	createFunc          func(ctx context.Context, name, description string) error
	getAllByUserIDFunc  func(ctx context.Context) ([]*Categories, error)
	getByCategoryIDFunc func(ctx context.Context, categoryID string) (*Categories, error)
	updateFunc          func(ctx context.Context, name, description, categoryID string) error
	deleteFunc          func(ctx context.Context, categoryID string) error
}

func (m *mockCategoriesService) Create(ctx context.Context, name, description string) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, name, description)
	}
	return nil
}

func (m *mockCategoriesService) GetAllByUserID(ctx context.Context) ([]*Categories, error) {
	if m.getAllByUserIDFunc != nil {
		return m.getAllByUserIDFunc(ctx)
	}
	return []*Categories{}, nil
}

func (m *mockCategoriesService) GetByCategoryID(ctx context.Context, categoryID string) (*Categories, error) {
	if m.getByCategoryIDFunc != nil {
		return m.getByCategoryIDFunc(ctx, categoryID)
	}
	return nil, nil
}

func (m *mockCategoriesService) Update(ctx context.Context, name, description, categoryID string) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, name, description, categoryID)
	}
	return nil
}

func (m *mockCategoriesService) Delete(ctx context.Context, categoryID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, categoryID)
	}
	return nil
}

func TestCreateHandler_Success(t *testing.T) {
	called := false

	mockService := &mockCategoriesService{
		createFunc: func(ctx context.Context, name, description string) error {
			called = true
			assert.Equal(t, "Test", name)
			assert.Equal(t, "Desc Test", description)
			return nil
		},
	}

	handler := NewCategoriesHandler(mockService)

	body := `{
		"name": "Test",
		"description": "Desc Test"
	}`

	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// inject userID ke context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Categories Created")
	assert.True(t, called)
}

func TestCreateHandler_InvalidJSON(t *testing.T) {
	mockService := &mockCategoriesService{}
	handler := NewCategoriesHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateHandler_ValidationError(t *testing.T) {
	called := false

	mockService := &mockCategoriesService{
		createFunc: func(ctx context.Context, name, description string) error {
			called = true
			return nil
		},
	}

	handler := NewCategoriesHandler(mockService)

	body := `{
		"name": ""
	}`

	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestCreateHandler_ServiceError(t *testing.T) {
	mockService := &mockCategoriesService{
		createFunc: func(ctx context.Context, name, description string) error {
			return errors.New("Email already exists")
		},
	}

	handler := NewCategoriesHandler(mockService)

	body := `{
		"name": "test",
		"description": "test"
	}`

	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email already exists")
}

func TestDeleteHandler_Success(t *testing.T) {
	called := false

	mockService := &mockCategoriesService{
		deleteFunc: func(ctx context.Context, categoryID string) error {
			called = true
			assert.Equal(t, "CAT-001", categoryID)
			return nil
		},
	}

	handler := NewCategoriesHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/CAT-001", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Category Deleted")
	assert.True(t, called)
}

func TestDeleteHandler_MissingID(t *testing.T) {
	mockService := &mockCategoriesService{}
	handler := NewCategoriesHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteHandler_ServiceError(t *testing.T) {
	mockService := &mockCategoriesService{
		deleteFunc: func(ctx context.Context, categoryID string) error {
			return errors.New("Something went wrong")
		},
	}

	handler := NewCategoriesHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/CAT-001", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Something went wrong")
}

func TestDeleteHandler_CategoryHasBudget(t *testing.T) {
	mockService := &mockCategoriesService{
		deleteFunc: func(ctx context.Context, categoryID string) error {
			return ErrCategoryHasBudget
		},
	}

	handler := NewCategoriesHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/CAT-001", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrCategoryHasBudget.Error())
}
