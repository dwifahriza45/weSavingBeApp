package categories

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"fmt"
	"time"
)

type CategoriesService interface {
	Create(ctx context.Context, name, description string) error
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
)

func (s *categoriesService) Create(
	ctx context.Context,
	name, description string,
) error {
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
