package cetegoriesbudget

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CategoriesBudgetService interface {
	Create(ctx context.Context, categoryID, allocatedAmount, period string) error
	GetByCategoryID(ctx context.Context, categoryID string) (*CategoriesBudget, error)
}

type categoriesBudgetService struct {
	categoriesBudgetRepo CategoriesBudgetRepository
}

func NewCategoriesBudgetService(categoriesBudgetRepo CategoriesBudgetRepository) CategoriesBudgetService {
	return &categoriesBudgetService{
		categoriesBudgetRepo: categoriesBudgetRepo,
	}
}

var (
	ErrInvalidCredentials     = errors.New("Invalid Credentials")
	ErrInvalidAllocatedAmount = errors.New("allocated_amount must be a valid integer")
	ErrCategoryBudgetNotFound = errors.New("category budget not found")
	ErrInternal               = errors.New("internal server error")
)

func (s *categoriesBudgetService) Create(ctx context.Context, categoryID, allocatedAmount, period string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	if _, err := strconv.ParseInt(allocatedAmount, 10, 64); err != nil {
		return ErrInvalidAllocatedAmount
	}

	if period == "" {
		period = "monthly"
	}

	budgetID, err := s.generateBudgetID(ctx)
	if err != nil {
		return err
	}

	categoryBudget := &CategoriesBudget{
		BUDGET_ID:       budgetID,
		USER_ID:         userID,
		CATEGORY_ID:     categoryID,
		AllocatedAmount: allocatedAmount,
		UsedAmount:      "0",
		Period:          period,
	}

	return s.categoriesBudgetRepo.Create(ctx, categoryBudget)
}

func (s *categoriesBudgetService) GetByCategoryID(ctx context.Context, categoryID string) (*CategoriesBudget, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	categoryBudget, err := s.categoriesBudgetRepo.GetByCategoryID(ctx, userID, categoryID)
	if err != nil {
		switch err {
		case ErrCategoryBudgetNotFound:
			return nil, ErrCategoryBudgetNotFound
		default:
			return nil, ErrInternal
		}
	}

	return categoryBudget, nil
}

func (s *categoriesBudgetService) generateBudgetID(ctx context.Context) (string, error) {
	dateStr := time.Now().Format("20060102")

	count, err := s.categoriesBudgetRepo.CountByDate(ctx, dateStr)
	if err != nil {
		return "", err
	}

	sequence := count + 1

	budgetID := fmt.Sprintf("BUD-%s-%06d", dateStr, sequence)

	return budgetID, nil
}
