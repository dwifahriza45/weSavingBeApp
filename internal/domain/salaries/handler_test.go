package salaries

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

type mockSalariesService struct {
	createFunc      func(ctx context.Context, amount, source, description string) error
	checkSalaryFunc func(ctx context.Context) (int, error)
	getTotalFunc    func(ctx context.Context) (int64, error)
	getAllFunc      func(ctx context.Context) ([]*Salaries, error)
	updateFunc      func(ctx context.Context, salaryID, amount, source, description string) error
	deleteFunc      func(ctx context.Context, salaryID string) error
}

func (m *mockSalariesService) Create(ctx context.Context, amount, source, description string) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, amount, source, description)
	}

	return nil
}

func (m *mockSalariesService) CheckSalary(ctx context.Context) (int, error) {
	if m.checkSalaryFunc != nil {
		return m.checkSalaryFunc(ctx)
	}

	return 0, nil
}

func (m *mockSalariesService) GetTotalSalary(ctx context.Context) (int64, error) {
	if m.getTotalFunc != nil {
		return m.getTotalFunc(ctx)
	}

	return 0, nil
}

func (m *mockSalariesService) GetAllByUserID(ctx context.Context) ([]*Salaries, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}

	return []*Salaries{}, nil
}

func (m *mockSalariesService) Update(ctx context.Context, salaryID, amount, source, description string) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, salaryID, amount, source, description)
	}

	return nil
}

func (m *mockSalariesService) Delete(ctx context.Context, salaryID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, salaryID)
	}

	return nil
}

func TestCreateSalaryHandler_Success(t *testing.T) {
	mockService := &mockSalariesService{
		createFunc: func(ctx context.Context, amount, source, description string) error {
			assert.Equal(t, "5000000", amount)
			assert.Equal(t, "main job", source)
			assert.Equal(t, "monthly salary", description)
			return nil
		},
	}

	handler := NewSalariesHandler(mockService)

	body := `{
		"amount": "5000000",
		"source": "main job",
		"description": "monthly salary"
	}`

	req := httptest.NewRequest(http.MethodPost, "/salary/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Salary Created")
}

func TestCreateSalaryHandler_InvalidJSON(t *testing.T) {
	mockService := &mockSalariesService{}
	handler := NewSalariesHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/salary/create", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateSalaryHandler_ValidationError(t *testing.T) {
	mockService := &mockSalariesService{}
	handler := NewSalariesHandler(mockService)

	body := `{
		"amount": "abc"
	}`

	req := httptest.NewRequest(http.MethodPost, "/salary/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Validation Failed")
	assert.Contains(t, rec.Body.String(), "amount must be numeric")
}

func TestCreateSalaryHandler_ServiceError(t *testing.T) {
	mockService := &mockSalariesService{
		createFunc: func(ctx context.Context, amount, source, description string) error {
			return errors.New("failed create salary")
		},
	}

	handler := NewSalariesHandler(mockService)

	body := `{
		"amount": "5000000",
		"source": "main job"
	}`

	req := httptest.NewRequest(http.MethodPost, "/salary/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed create salary")
}

func TestCheckSalaryHandler_Success_NotYetInserted(t *testing.T) {
	mockService := &mockSalariesService{
		checkSalaryFunc: func(ctx context.Context) (int, error) {
			return 0, nil
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/check", nil)
	rec := httptest.NewRecorder()

	handler.CheckSalary(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"data":0`)
}

func TestCheckSalaryHandler_Success_AlreadyInserted(t *testing.T) {
	mockService := &mockSalariesService{
		checkSalaryFunc: func(ctx context.Context) (int, error) {
			return 1, nil
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/check", nil)
	rec := httptest.NewRecorder()

	handler.CheckSalary(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"data":1`)
}

func TestCheckSalaryHandler_Unauthorized(t *testing.T) {
	mockService := &mockSalariesService{
		checkSalaryFunc: func(ctx context.Context) (int, error) {
			return 0, ErrInvalidCredentials
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/check", nil)
	rec := httptest.NewRecorder()

	handler.CheckSalary(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInvalidCredentials.Error())
}

func TestCheckSalaryHandler_InternalError(t *testing.T) {
	mockService := &mockSalariesService{
		checkSalaryFunc: func(ctx context.Context) (int, error) {
			return 0, errors.New("boom")
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/check", nil)
	rec := httptest.NewRecorder()

	handler.CheckSalary(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "boom")
}

func TestGetTotalSalaryHandler_Success(t *testing.T) {
	mockService := &mockSalariesService{
		getTotalFunc: func(ctx context.Context) (int64, error) {
			return 7500000, nil
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/total", nil)
	rec := httptest.NewRecorder()

	handler.GetTotalSalary(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"data":7500000`)
}

func TestGetTotalSalaryHandler_Unauthorized(t *testing.T) {
	mockService := &mockSalariesService{
		getTotalFunc: func(ctx context.Context) (int64, error) {
			return 0, ErrInvalidCredentials
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/total", nil)
	rec := httptest.NewRecorder()

	handler.GetTotalSalary(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInvalidCredentials.Error())
}

func TestGetTotalSalaryHandler_InternalError(t *testing.T) {
	mockService := &mockSalariesService{
		getTotalFunc: func(ctx context.Context) (int64, error) {
			return 0, errors.New("boom")
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/total", nil)
	rec := httptest.NewRecorder()

	handler.GetTotalSalary(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "boom")
}

func TestGetAllSalaryHandler_Success(t *testing.T) {
	mockService := &mockSalariesService{
		getAllFunc: func(ctx context.Context) ([]*Salaries, error) {
			return []*Salaries{
				{
					ID:          1,
					SalaryID:    "SAL-20260317-000001",
					UserID:      "USER-001",
					Amount:      "7500000",
					Source:      "main job",
					Description: "monthly salary",
					ReceivedAt:  "2026-03-17T09:30:00+07:00",
				},
			}, nil
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/all", nil)
	rec := httptest.NewRecorder()

	handler.GetAllByUserID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"salary_id":"SAL-20260317-000001"`)
}

func TestGetAllSalaryHandler_Unauthorized(t *testing.T) {
	mockService := &mockSalariesService{
		getAllFunc: func(ctx context.Context) ([]*Salaries, error) {
			return nil, ErrInvalidCredentials
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/all", nil)
	rec := httptest.NewRecorder()

	handler.GetAllByUserID(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInvalidCredentials.Error())
}

func TestGetAllSalaryHandler_InternalError(t *testing.T) {
	mockService := &mockSalariesService{
		getAllFunc: func(ctx context.Context) ([]*Salaries, error) {
			return nil, errors.New("boom")
		},
	}

	handler := NewSalariesHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/salary/all", nil)
	rec := httptest.NewRecorder()

	handler.GetAllByUserID(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "boom")
}

func TestUpdateSalaryHandler_Success(t *testing.T) {
	mockService := &mockSalariesService{
		updateFunc: func(ctx context.Context, salaryID, amount, source, description string) error {
			assert.Equal(t, "SAL-20260317-000001", salaryID)
			assert.Equal(t, "8000000", amount)
			assert.Equal(t, "main job", source)
			assert.Equal(t, "salary updated", description)
			return nil
		},
	}

	handler := NewSalariesHandler(mockService)
	r := chi.NewRouter()
	r.Put("/salary/{id}", handler.Update)

	body := `{
		"amount": "8000000",
		"source": "main job",
		"description": "salary updated"
	}`

	req := httptest.NewRequest(http.MethodPut, "/salary/SAL-20260317-000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Salary Updated")
}

func TestUpdateSalaryHandler_InvalidJSON(t *testing.T) {
	mockService := &mockSalariesService{}
	handler := NewSalariesHandler(mockService)

	r := chi.NewRouter()
	r.Put("/salary/{id}", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/salary/SAL-20260317-000001", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateSalaryHandler_ValidationError(t *testing.T) {
	mockService := &mockSalariesService{}
	handler := NewSalariesHandler(mockService)

	r := chi.NewRouter()
	r.Put("/salary/{id}", handler.Update)

	body := `{
		"amount": "abc"
	}`

	req := httptest.NewRequest(http.MethodPut, "/salary/SAL-20260317-000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Validation Failed")
}

func TestUpdateSalaryHandler_NotFound(t *testing.T) {
	mockService := &mockSalariesService{
		updateFunc: func(ctx context.Context, salaryID, amount, source, description string) error {
			return ErrSalaryNotFound
		},
	}

	handler := NewSalariesHandler(mockService)
	r := chi.NewRouter()
	r.Put("/salary/{id}", handler.Update)

	body := `{
		"amount": "8000000",
		"source": "main job",
		"description": "salary updated"
	}`

	req := httptest.NewRequest(http.MethodPut, "/salary/SAL-20260317-000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrSalaryNotFound.Error())
}

func TestDeleteSalaryHandler_Success(t *testing.T) {
	mockService := &mockSalariesService{
		deleteFunc: func(ctx context.Context, salaryID string) error {
			assert.Equal(t, "SAL-20260317-000001", salaryID)
			return nil
		},
	}

	handler := NewSalariesHandler(mockService)
	r := chi.NewRouter()
	r.Delete("/salary/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/salary/SAL-20260317-000001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Salary Deleted")
}

func TestDeleteSalaryHandler_NotFound(t *testing.T) {
	mockService := &mockSalariesService{
		deleteFunc: func(ctx context.Context, salaryID string) error {
			return ErrSalaryNotFound
		},
	}

	handler := NewSalariesHandler(mockService)
	r := chi.NewRouter()
	r.Delete("/salary/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/salary/SAL-20260317-000001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrSalaryNotFound.Error())
}

func TestDeleteSalaryHandler_Unauthorized(t *testing.T) {
	mockService := &mockSalariesService{
		deleteFunc: func(ctx context.Context, salaryID string) error {
			return ErrInvalidCredentials
		},
	}

	handler := NewSalariesHandler(mockService)
	r := chi.NewRouter()
	r.Delete("/salary/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/salary/SAL-20260317-000001", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInvalidCredentials.Error())
}
