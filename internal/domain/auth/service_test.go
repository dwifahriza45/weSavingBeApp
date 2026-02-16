package auth

import (
	"BE_WE_SAVING/internal/domain/users"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	findByEmailFunc    func(ctx context.Context, email string) (*users.User, error)
	findByUsernameFunc func(ctx context.Context, username string) (*users.User, error)
	findByUserIDFunc   func(ctx context.Context, userID string) (*users.User, error)
	countByDateFunc    func(ctx context.Context, date string) (int, error)
	createFunc         func(ctx context.Context, user *users.User) error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	return m.findByEmailFunc(ctx, email)
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*users.User, error) {
	return m.findByUsernameFunc(ctx, username)
}

func (m *mockUserRepo) FindByUserID(ctx context.Context, userID string) (*users.User, error) {
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserRepo) CountByDate(ctx context.Context, date string) (int, error) {
	if m.countByDateFunc != nil {
		return m.countByDateFunc(ctx, date)
	}
	return 0, nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *users.User) error {
	return m.createFunc(ctx, user)
}

func TestRegister_EmailNotFound_SQLNoRows(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, sql.ErrNoRows
		},
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, sql.ErrNoRows
		},
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, nil
		},
		createFunc: func(ctx context.Context, user *users.User) error {
			return nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.NoError(t, err)
}

func TestRegister_ErrorFromGenerateUserID(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, nil
		},
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, nil
		},
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, errors.New("db error")
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.Error(t, err)
}

func TestRegister_DBErrorOnFindByEmail(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, errors.New("db error")
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.Error(t, err)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return &users.User{Email: email}, nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.Equal(t, ErrEmailAlreadyExists, err)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, nil
		},
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return &users.User{Username: username}, nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.Equal(t, ErrUsernameAlreadyExists, err)
}

func TestGenerateUserID_Error(t *testing.T) {
	mockRepo := &mockUserRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, errors.New("db error")
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	userID, err := service.generateUserID(context.Background())

	assert.Empty(t, userID)
	assert.Error(t, err)
}

func TestRegister_DBErrorOnFindByUsername(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, nil
		},
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, errors.New("db error")
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.Error(t, err)
}

func TestRegister_DBErrorOnCreate(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, nil
		},
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, nil
		},
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 1, nil
		},
		createFunc: func(ctx context.Context, user *users.User) error {
			return errors.New("insert failed")
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.Error(t, err)
}

func TestGenerateUserID_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 5, nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	userID, err := service.generateUserID(context.Background())

	assert.NoError(t, err)
	assert.Contains(t, userID, "USER-")
	assert.Contains(t, userID, "-000006")
}

func TestRegister_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByEmailFunc: func(ctx context.Context, email string) (*users.User, error) {
			return nil, nil
		},
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, nil
		},
		createFunc: func(ctx context.Context, user *users.User) error {
			return nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	err := service.Register(context.Background(), "username", "Full Name", "test@mail.com", "password")

	assert.NoError(t, err)
}

func TestGenerateUserID_FormatValidation(t *testing.T) {
	mockRepo := &mockUserRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 5, nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	userID, err := service.generateUserID(context.Background())

	assert.NoError(t, err)
	assert.Regexp(t, `^USER-\d{8}-\d{6}$`, userID)
}

func TestLogin_InvalidCredentials_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := &authService{
		userRepo:  mockRepo,
		jwtSecret: "secret",
	}

	token, err := service.Login(context.Background(), "john", "password123")

	assert.Empty(t, token)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestLogin_InvalidCredentials_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return &users.User{
				USER_ID:  "USER-20260216-000001",
				Username: username,
				Password: string(hashedPassword),
			}, nil
		},
	}

	service := &authService{
		userRepo:  mockRepo,
		jwtSecret: "secret",
	}

	token, err := service.Login(context.Background(), "john", "wrongpassword")

	assert.Empty(t, token)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestLogin_InvalidCredentials_DBError(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return nil, errors.New("db connection lost")
		},
	}

	service := &authService{
		userRepo:  mockRepo,
		jwtSecret: "secret",
	}

	token, err := service.Login(context.Background(), "john", "password")

	assert.Empty(t, token)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestLogin_InvalidCredentials_EmptyPasswordHash(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return &users.User{
				USER_ID:  "USER-20260216-000001",
				Username: username,
				Password: "", // invalid hash
			}, nil
		},
	}

	service := &authService{
		userRepo:  mockRepo,
		jwtSecret: "secret",
	}

	token, err := service.Login(context.Background(), "john", "password123")

	assert.Empty(t, token)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestLogin_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepo{
		findByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
			return &users.User{
				USER_ID:  "USER-20260216-000001",
				Username: username,
				Password: string(hashedPassword),
			}, nil
		},
	}

	service := &authService{
		userRepo:  mockRepo,
		jwtSecret: "secret",
	}

	token, err := service.Login(context.Background(), "john", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetMe_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByUserIDFunc: func(ctx context.Context, userID string) (*users.User, error) {
			return &users.User{
				USER_ID:  userID,
				Username: "john",
				Email:    "john@mail.com",
			}, nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	user, err := service.GetMe(context.Background(), "USER-1")

	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
}

func TestGetMe_NilUserNoError(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByUserIDFunc: func(ctx context.Context, userID string) (*users.User, error) {
			return nil, nil
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	user, err := service.GetMe(context.Background(), "USER-1")

	assert.Nil(t, user)
	assert.NoError(t, err)
}

func TestGetMe_Error(t *testing.T) {
	mockRepo := &mockUserRepo{
		findByUserIDFunc: func(ctx context.Context, userID string) (*users.User, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := &authService{
		userRepo: mockRepo,
	}

	user, err := service.GetMe(context.Background(), "USER-1")

	assert.Nil(t, user)
	assert.Error(t, err)
}
