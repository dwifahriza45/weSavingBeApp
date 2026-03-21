package categories

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type CategoriesService interface {
	Create(ctx context.Context, name, description string) error
	GetAllByUserID(ctx context.Context) ([]*Categories, error)
	GetByCategoryID(ctx context.Context, categoryID string) (*Categories, error)
	Update(ctx context.Context, name, description, categoryID string) error
	Delete(ctx context.Context, categoryID string) error
}

type categoriesService struct {
	categoryRepo CategoriesRepository
}

func NewCategoriesService(categoryRepo CategoriesRepository) CategoriesService {
	return &categoriesService{
		categoryRepo: categoryRepo,
	}
}

var (
	ErrInvalidCredentials = errors.New("Invalid Credentials")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategoryHasBudget  = errors.New("category cannot be deleted because it is still used in category budgets")
	ErrInternal           = errors.New("internal server error")
)

func (s *categoriesService) Create(ctx context.Context, name, description string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	categoryID, err := s.generateCategoryID(ctx)
	if err != nil {
		return err
	}

	category := &Categories{
		CATEGORY_ID: categoryID,
		USER_ID:     userID,
		Name:        name,
		Description: description,
	}

	return s.categoryRepo.Create(ctx, category)
}

func (s *categoriesService) generateCategoryID(ctx context.Context) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	count, err := s.categoryRepo.CountByDate(ctx, dateStr)
	if err != nil {
		return "", err
	}

	sequence := count + 1

	categoryID := fmt.Sprintf("CAT-%s-%06d", dateStr, sequence)

	return categoryID, nil
}

func (s *categoriesService) GetAllByUserID(ctx context.Context) ([]*Categories, error) {

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	categories, err := s.categoryRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}

	return categories, nil
}

func (s *categoriesService) GetByCategoryID(ctx context.Context, categoryID string) (*Categories, error) {

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	category, err := s.categoryRepo.GetByCategoryID(ctx, userID, categoryID)
	if err != nil {
		switch err {
		case ErrCategoryNotFound:
			return nil, ErrCategoryNotFound
		default:
			return nil, ErrInternal
		}
	}

	return category, nil
}

func (s *categoriesService) Update(ctx context.Context, name, description, categoryID string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	category := &Categories{
		Name:        name,
		Description: description,
		CATEGORY_ID: categoryID,
		USER_ID:     userID,
	}

	err := s.categoryRepo.Update(ctx, category)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return ErrCategoryNotFound
		}

		return ErrInternal
	}

	return nil
}

func (s *categoriesService) Delete(ctx context.Context, categoryID string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	hasBudget, err := s.categoryRepo.HasCategoryBudget(ctx, categoryID, userID)
	if err != nil {
		return ErrInternal
	}

	if hasBudget {
		return ErrCategoryHasBudget
	}

	err = s.categoryRepo.Delete(ctx, categoryID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return ErrInternal
	}

	return nil
}
