package cetegoriesbudget

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCategoriesBudgetRepo struct {
	countByDateFunc     func(ctx context.Context, date string) (int, error)
	createFunc          func(ctx context.Context, category *CategoriesBudget) error
	updateFunc          func(ctx context.Context, category *CategoriesBudget) error
	getByCategoryIDFunc func(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error)
}

func (m *mockCategoriesBudgetRepo) CountByDate(ctx context.Context, date string) (int, error) {
	if m.countByDateFunc != nil {
		return m.countByDateFunc(ctx, date)
	}

	return 0, nil
}

func (m *mockCategoriesBudgetRepo) Create(ctx context.Context, category *CategoriesBudget) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, category)
	}

	return nil
}

func (m *mockCategoriesBudgetRepo) Update(ctx context.Context, category *CategoriesBudget) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, category)
	}

	return nil
}

func (m *mockCategoriesBudgetRepo) GetByCategoryID(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error) {
	if m.getByCategoryIDFunc != nil {
		return m.getByCategoryIDFunc(ctx, userID, categoryID)
	}

	return nil, nil
}

func TestCreateCategoriesBudget_Success(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, nil
		},
		createFunc: func(ctx context.Context, category *CategoriesBudget) error {
			assert.Equal(t, "USER-001", category.USER_ID)
			assert.Equal(t, "CAT-001", category.CATEGORY_ID)
			assert.Equal(t, "1000000", category.AllocatedAmount)
			assert.Equal(t, "0", category.UsedAmount)
			assert.Equal(t, "monthly", category.Period)
			assert.Contains(t, category.BUDGET_ID, "BUD-")
			return nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000", "")

	assert.NoError(t, err)
}

func TestCreateCategoriesBudget_InvalidCredential(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	err := service.Create(context.Background(), "CAT-001", "1000000", "")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestCreateCategoriesBudget_InvalidAllocatedAmount(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "abc", "")

	assert.Equal(t, ErrInvalidAllocatedAmount, err)
}

func TestCreateCategoriesBudget_GenerateBudgetIDError(t *testing.T) {
	expectedErr := errors.New("db error")

	mockRepo := &mockCategoriesBudgetRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, expectedErr
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000", "")

	assert.Equal(t, expectedErr, err)
}

func TestCreateCategoriesBudget_RepoError(t *testing.T) {
	expectedErr := errors.New("insert failed")

	mockRepo := &mockCategoriesBudgetRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 1, nil
		},
		createFunc: func(ctx context.Context, category *CategoriesBudget) error {
			return expectedErr
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000", "monthly")

	assert.Equal(t, expectedErr, err)
}

func TestGetCategoriesBudgetByCategoryID_Success(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByCategoryIDFunc: func(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)

			return &CategoriesBudget{
				BUDGET_ID:       "BUD-001",
				USER_ID:         userID,
				CATEGORY_ID:     categoryID,
				AllocatedAmount: "1000000",
				UsedAmount:      "250000",
				Period:          "monthly",
			}, nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetByCategoryID(ctx, "CAT-001")

	assert.NoError(t, err)
	assert.Equal(t, "BUD-001", result.BUDGET_ID)
	assert.Equal(t, "CAT-001", result.CATEGORY_ID)
}

func TestGetCategoriesBudgetByCategoryID_InvalidCredential(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	result, err := service.GetByCategoryID(context.Background(), "CAT-001")

	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestGetCategoriesBudgetByCategoryID_NotFound(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByCategoryIDFunc: func(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error) {
			return nil, ErrCategoryBudgetNotFound
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetByCategoryID(ctx, "CAT-001")

	assert.Nil(t, result)
	assert.Equal(t, ErrCategoryBudgetNotFound, err)
}

func TestGetCategoriesBudgetByCategoryID_RepoError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByCategoryIDFunc: func(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error) {
			return nil, errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetByCategoryID(ctx, "CAT-001")

	assert.Nil(t, result)
	assert.Equal(t, ErrInternal, err)
}
