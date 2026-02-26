package categories

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCategoriesService struct {
	createFunc func(ctx context.Context, name, description string) error
}

func (m *mockCategoriesService) Create(ctx context.Context, name, description string) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, name, description)
	}
	return nil
}

func TestCreateHandler_Success(t *testing.T) {
	mockService := &mockCategoriesService{
		createFunc: func(ctx context.Context, name, description string) error {
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
	mockService := &mockCategoriesService{}
	handler := NewCategoriesHandler(mockService)

	body := `{
		"name": ""
	}`

	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Validation Failed")
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
