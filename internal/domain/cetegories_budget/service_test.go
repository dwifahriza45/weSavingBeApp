package cetegoriesbudget

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCategoriesBudgetRepo struct {
	countByDateFunc                        func(ctx context.Context, date string) (int, error)
	categoryExistsFunc                     func(ctx context.Context, userID, categoryID string) (bool, error)
	categoryBudgetExistsInCurrentMonthFunc func(ctx context.Context, userID, categoryID string) (bool, error)
	createFunc                             func(ctx context.Context, category *CategoriesBudget) error
	updateFunc                             func(ctx context.Context, category *CategoriesBudget) error
	deleteFunc                             func(ctx context.Context, category *CategoriesBudget) error
	getByCategoryIDFunc                    func(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error)
	getAllByCategoryIDFunc                 func(ctx context.Context, userID, categoryID string) ([]*CategoriesBudget, error)
	getByBudgetIDFunc                      func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error)
}

func (m *mockCategoriesBudgetRepo) CountByDate(ctx context.Context, date string) (int, error) {
	if m.countByDateFunc != nil {
		return m.countByDateFunc(ctx, date)
	}

	return 0, nil
}

func (m *mockCategoriesBudgetRepo) CategoryExists(ctx context.Context, userID, categoryID string) (bool, error) {
	if m.categoryExistsFunc != nil {
		return m.categoryExistsFunc(ctx, userID, categoryID)
	}

	return false, nil
}

func (m *mockCategoriesBudgetRepo) CategoryBudgetExistsInCurrentMonth(ctx context.Context, userID, categoryID string) (bool, error) {
	if m.categoryBudgetExistsInCurrentMonthFunc != nil {
		return m.categoryBudgetExistsInCurrentMonthFunc(ctx, userID, categoryID)
	}

	return false, nil
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

func (m *mockCategoriesBudgetRepo) Delete(ctx context.Context, category *CategoriesBudget) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, category)
	}

	return nil
}

func (m *mockCategoriesBudgetRepo) GetByCategoryID(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error) {
	if m.getByCategoryIDFunc != nil {
		return m.getByCategoryIDFunc(ctx, userID, categoryID)
	}

	return nil, nil
}

func (m *mockCategoriesBudgetRepo) GetAllByCategoryID(ctx context.Context, userID, categoryID string) ([]*CategoriesBudget, error) {
	if m.getAllByCategoryIDFunc != nil {
		return m.getAllByCategoryIDFunc(ctx, userID, categoryID)
	}

	return []*CategoriesBudget{}, nil
}

func (m *mockCategoriesBudgetRepo) GetByBudgetID(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
	if m.getByBudgetIDFunc != nil {
		return m.getByBudgetIDFunc(ctx, userID, budgetID)
	}

	return nil, nil
}

func TestCreateCategoriesBudget_Success(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, nil
		},
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)
			return true, nil
		},
		categoryBudgetExistsInCurrentMonthFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)
			return false, nil
		},
		createFunc: func(ctx context.Context, category *CategoriesBudget) error {
			assert.Equal(t, "USER-001", category.USER_ID)
			assert.Equal(t, "CAT-001", category.CATEGORY_ID)
			assert.Equal(t, "1000000", category.AllocatedAmount)
			assert.Equal(t, "0", category.UsedAmount)
			assert.Contains(t, category.BUDGET_ID, "BUD-")
			return nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.NoError(t, err)
}

func TestCreateCategoriesBudget_InvalidCredential(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	err := service.Create(context.Background(), "CAT-001", "1000000")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestCreateCategoriesBudget_InvalidAllocatedAmount(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "abc")

	assert.Equal(t, ErrInvalidAllocatedAmount, err)
}

func TestCreateCategoriesBudget_GenerateBudgetIDError(t *testing.T) {
	expectedErr := errors.New("db error")

	mockRepo := &mockCategoriesBudgetRepo{
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return true, nil
		},
		categoryBudgetExistsInCurrentMonthFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return false, nil
		},
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, expectedErr
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.Equal(t, expectedErr, err)
}

func TestCreateCategoriesBudget_RepoError(t *testing.T) {
	expectedErr := errors.New("insert failed")

	mockRepo := &mockCategoriesBudgetRepo{
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return true, nil
		},
		categoryBudgetExistsInCurrentMonthFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return false, nil
		},
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

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.Equal(t, expectedErr, err)
}

func TestCreateCategoriesBudget_CategoryNotFound(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)
			return false, nil
		},
		createFunc: func(ctx context.Context, category *CategoriesBudget) error {
			t.Fatal("Create should not be called when category does not exist")
			return nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.Equal(t, ErrCategoryNotFound, err)
}

func TestCreateCategoriesBudget_CategoryExistsError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.Equal(t, ErrInternal, err)
}

func TestCreateCategoriesBudget_AlreadyExistsInCurrentMonth(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)
			return true, nil
		},
		categoryBudgetExistsInCurrentMonthFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)
			return true, nil
		},
		createFunc: func(ctx context.Context, category *CategoriesBudget) error {
			t.Fatal("Create should not be called when budget for current month already exists")
			return nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.Equal(t, ErrCategoryBudgetAlreadyExists, err)
}

func TestCreateCategoriesBudget_ExistsInCurrentMonthError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		categoryExistsFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return true, nil
		},
		categoryBudgetExistsInCurrentMonthFunc: func(ctx context.Context, userID, categoryID string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "CAT-001", "1000000")

	assert.Equal(t, ErrInternal, err)
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

func TestGetAllCategoriesBudgetByCategoryID_Success(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getAllByCategoryIDFunc: func(ctx context.Context, userID, categoryID string) ([]*CategoriesBudget, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "CAT-001", categoryID)

			return []*CategoriesBudget{
				{
					BUDGET_ID:       "BUD-002",
					USER_ID:         userID,
					CATEGORY_ID:     categoryID,
					AllocatedAmount: "1500000",
					UsedAmount:      "500000",
				},
				{
					BUDGET_ID:       "BUD-001",
					USER_ID:         userID,
					CATEGORY_ID:     categoryID,
					AllocatedAmount: "1000000",
					UsedAmount:      "250000",
				},
			}, nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetAllByCategoryID(ctx, "CAT-001")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BUD-002", result[0].BUDGET_ID)
	assert.Equal(t, "BUD-001", result[1].BUDGET_ID)
}

func TestGetAllCategoriesBudgetByCategoryID_InvalidCredential(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	result, err := service.GetAllByCategoryID(context.Background(), "CAT-001")

	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestGetAllCategoriesBudgetByCategoryID_Empty(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getAllByCategoryIDFunc: func(ctx context.Context, userID, categoryID string) ([]*CategoriesBudget, error) {
			return []*CategoriesBudget{}, nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetAllByCategoryID(ctx, "CAT-001")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetAllCategoriesBudgetByCategoryID_RepoError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getAllByCategoryIDFunc: func(ctx context.Context, userID, categoryID string) ([]*CategoriesBudget, error) {
			return nil, errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetAllByCategoryID(ctx, "CAT-001")

	assert.Nil(t, result)
	assert.Equal(t, ErrInternal, err)
}

func TestUpdateCategoriesBudget_Success(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "BUD-001", budgetID)

			return &CategoriesBudget{
				BUDGET_ID:       budgetID,
				USER_ID:         userID,
				CATEGORY_ID:     "CAT-001",
				AllocatedAmount: "1000000",
				UsedAmount:      "250000",
			}, nil
		},
		updateFunc: func(ctx context.Context, category *CategoriesBudget) error {
			assert.Equal(t, "BUD-001", category.BUDGET_ID)
			assert.Equal(t, "USER-001", category.USER_ID)
			assert.Equal(t, "CAT-001", category.CATEGORY_ID)
			assert.Equal(t, "1500000", category.AllocatedAmount)
			assert.Equal(t, "250000", category.UsedAmount)
			return nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Update(ctx, "BUD-001", "1500000")

	assert.NoError(t, err)
}

func TestUpdateCategoriesBudget_InvalidCredential(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	err := service.Update(context.Background(), "BUD-001", "1500000")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestUpdateCategoriesBudget_InvalidAllocatedAmount(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Update(ctx, "BUD-001", "abc")

	assert.Equal(t, ErrInvalidAllocatedAmount, err)
}

func TestUpdateCategoriesBudget_NotFound(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			return nil, ErrCategoryBudgetNotFound
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Update(ctx, "BUD-001", "1500000")

	assert.Equal(t, ErrCategoryBudgetNotFound, err)
}

func TestUpdateCategoriesBudget_GetRepoError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			return nil, errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Update(ctx, "BUD-001", "1500000")

	assert.Equal(t, ErrInternal, err)
}

func TestUpdateCategoriesBudget_UpdateRepoError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			return &CategoriesBudget{
				BUDGET_ID:       budgetID,
				USER_ID:         userID,
				CATEGORY_ID:     "CAT-001",
				AllocatedAmount: "1000000",
				UsedAmount:      "250000",
			}, nil
		},
		updateFunc: func(ctx context.Context, category *CategoriesBudget) error {
			return errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Update(ctx, "BUD-001", "1500000")

	assert.Equal(t, ErrInternal, err)
}

func TestDeleteCategoriesBudget_Success(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			assert.Equal(t, "USER-001", userID)
			assert.Equal(t, "BUD-001", budgetID)

			return &CategoriesBudget{
				BUDGET_ID:       budgetID,
				USER_ID:         userID,
				CATEGORY_ID:     "CAT-001",
				AllocatedAmount: "1000000",
				UsedAmount:      "250000",
			}, nil
		},
		deleteFunc: func(ctx context.Context, category *CategoriesBudget) error {
			assert.Equal(t, "BUD-001", category.BUDGET_ID)
			assert.Equal(t, "USER-001", category.USER_ID)
			assert.Equal(t, "CAT-001", category.CATEGORY_ID)
			return nil
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "BUD-001")

	assert.NoError(t, err)
}

func TestDeleteCategoriesBudget_InvalidCredential(t *testing.T) {
	service := &categoriesBudgetService{
		categoriesBudgetRepo: &mockCategoriesBudgetRepo{},
	}

	err := service.Delete(context.Background(), "BUD-001")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestDeleteCategoriesBudget_NotFound(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			return nil, ErrCategoryBudgetNotFound
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "BUD-001")

	assert.Equal(t, ErrCategoryBudgetNotFound, err)
}

func TestDeleteCategoriesBudget_GetRepoError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			return nil, errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "BUD-001")

	assert.Equal(t, ErrInternal, err)
}

func TestDeleteCategoriesBudget_DeleteRepoError(t *testing.T) {
	mockRepo := &mockCategoriesBudgetRepo{
		getByBudgetIDFunc: func(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
			return &CategoriesBudget{
				BUDGET_ID:       budgetID,
				USER_ID:         userID,
				CATEGORY_ID:     "CAT-001",
				AllocatedAmount: "1000000",
				UsedAmount:      "250000",
			}, nil
		},
		deleteFunc: func(ctx context.Context, category *CategoriesBudget) error {
			return errors.New("db error")
		},
	}

	service := &categoriesBudgetService{
		categoriesBudgetRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "BUD-001")

	assert.Equal(t, ErrInternal, err)
}
