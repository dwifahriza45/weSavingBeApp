package auth

import (
	"BE_WE_SAVING/internal/domain/users"
	"BE_WE_SAVING/internal/shared/jwt"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, fullname, email, password string) error
	Login(ctx context.Context, username, password string) (string, error)
	GetMe(ctx context.Context, userID string) (*users.User, error)
}

type authService struct {
	userRepo  users.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo users.UserRepository, secret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: secret,
	}
}

var (
	ErrEmailAlreadyExists    = errors.New("Email already exists")
	ErrInvalidCredentials    = errors.New("Invalid Credentials")
	ErrUsernameAlreadyExists = errors.New("Username already exists")
	ErrUserIdAlreadyExists   = errors.New("User ID already exists")
)

func (s *authService) Register(
	ctx context.Context,
	username, fullname, email, password string,
) error {

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingUser != nil {
		return ErrEmailAlreadyExists
	}

	existingUsername, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingUsername != nil {
		return ErrUsernameAlreadyExists
	}

	userID, err := s.generateUserID(ctx)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &users.User{
		USER_ID:  userID,
		Username: username,
		Fullname: fullname,
		Email:    email,
		Password: string(hashed),
	}

	return s.userRepo.Create(ctx, user)
}

func (s *authService) generateUserID(ctx context.Context) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	count, err := s.userRepo.CountByDate(ctx, dateStr)
	if err != nil {
		return "", err
	}

	sequence := count + 1

	userID := fmt.Sprintf("USER-%s-%06d", dateStr, sequence)

	return userID, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateJWT(s.jwtSecret, user.USER_ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GetMe(ctx context.Context, userID string) (*users.User, error) {
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
