package auth

import (
	"BE_WE_SAVING/internal/app/middleware"
	"BE_WE_SAVING/internal/domain/users"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	registerFunc func(ctx context.Context, username, fullname, email, password string) error
	loginFunc    func(ctx context.Context, email, password string) (string, error)
	getMeFunc    func(ctx context.Context, userID string) (*users.User, error)
}

func (m *mockAuthService) Register(ctx context.Context, username, fullname, email, password string) error {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, username, fullname, email, password)
	}
	return nil
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, email, password)
	}
	return "", nil
}

func (m *mockAuthService) GetMe(ctx context.Context, userID string) (*users.User, error) {
	if m.getMeFunc != nil {
		return m.getMeFunc(ctx, userID)
	}
	return nil, nil
}

func TestRegisterHandler_Success(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, username, fullname, email, password string) error {
			return nil
		},
	}

	handler := NewAuthHandler(mockService)

	body := `{
		"username": "john",
		"fullname": "John Doe",
		"email": "john@mail.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "User Registered Successfully")
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_ValidationError(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	body := `{
		"username": "",
		"fullname": "",
		"email": "invalid-email",
		"password": "123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Validation Failed")
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	mockService := &mockAuthService{
		registerFunc: func(ctx context.Context, username, fullname, email, password string) error {
			return errors.New("Email already exists")
		},
	}

	handler := NewAuthHandler(mockService)

	body := `{
		"username": "john",
		"fullname": "John Doe",
		"email": "john@mail.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email already exists")
}

func TestMeHandler_Unauthorized(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_ValidationErrors(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedErrors []string
	}{
		{
			name: "username required",
			body: `{
				"username": "",
				"password": "password123"
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: []string{"username is required"},
		},
		{
			name: "password required",
			body: `{
				"username": "john",
				"password": ""
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: []string{"password is required"},
		},
		{
			name: "password min length",
			body: `{
				"username": "john",
				"password": "123"
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: []string{"password must be at least 6 characters"},
		},
		{
			name: "multiple validation errors",
			body: `{
				"username": "",
				"password": "123"
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: []string{
				"username is required",
				"password must be at least 6 characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAuthService{}
			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Login(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			for _, errMsg := range tt.expectedErrors {
				assert.Contains(t, rec.Body.String(), errMsg)
			}
		})
	}
}

func TestLoginHandler_ServiceError(t *testing.T) {
	mockService := &mockAuthService{
		loginFunc: func(ctx context.Context, username, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}

	handler := NewAuthHandler(mockService)

	body := `{
		"username": "john",
		"password": "password123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid credentials")
}

func TestLoginHandler_Success(t *testing.T) {
	mockService := &mockAuthService{
		loginFunc: func(ctx context.Context, username, password string) (string, error) {
			return "mocked-token", nil
		},
	}

	handler := NewAuthHandler(mockService)

	body := `{
		"username": "john",
		"password": "password123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Login Success")
	assert.Contains(t, rec.Body.String(), "mocked-token")
}

func TestMeHandler_UserNotFound(t *testing.T) {
	mockService := &mockAuthService{
		getMeFunc: func(ctx context.Context, userID string) (*users.User, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "USER-1")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMeHandler_Success(t *testing.T) {
	mockService := &mockAuthService{
		getMeFunc: func(ctx context.Context, userID string) (*users.User, error) {
			return &users.User{
				USER_ID:  userID,
				Username: "john",
				Email:    "john@mail.com",
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "USER-1")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "john")
}
