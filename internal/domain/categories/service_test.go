package categories

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCategoriesRepo struct {
	countByDateFunc       func(ctx context.Context, date string) (int, error)
	createFunc            func(ctx context.Context, categories *Categories) error
	getAllByUserIDFunc    func(ctx context.Context, userID string) ([]*Categories, error)
	getByCategoryIDFunc   func(ctx context.Context, userID, categoryID string) (*Categories, error)
	hasCategoryBudgetFunc func(ctx context.Context, categoryID string, userID string) (bool, error)
	updateFunc            func(ctx context.Context, categories *Categories) error
	deleteFunc            func(ctx context.Context, categoryID string, userID string) error
}

func (m *mockCategoriesRepo) CountByDate(ctx context.Context, date string) (int, error) {
	if m.countByDateFunc != nil {
		return m.countByDateFunc(ctx, date)
	}
	return 0, nil
}

func (m *mockCategoriesRepo) Create(ctx context.Context, categories *Categories) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, categories)
	}
	return nil
}

func (m *mockCategoriesRepo) GetAllByUserID(ctx context.Context, userID string) ([]*Categories, error) {
	if m.getAllByUserIDFunc != nil {
		return m.getAllByUserIDFunc(ctx, userID)
	}
	return []*Categories{}, nil
}

func (m *mockCategoriesRepo) GetByCategoryID(ctx context.Context, userID, categoryID string) (*Categories, error) {
	if m.getByCategoryIDFunc != nil {
		return m.getByCategoryIDFunc(ctx, userID, categoryID)
	}
	return nil, nil
}

func (m *mockCategoriesRepo) HasCategoryBudget(ctx context.Context, categoryID string, userID string) (bool, error) {
	if m.hasCategoryBudgetFunc != nil {
		return m.hasCategoryBudgetFunc(ctx, categoryID, userID)
	}
	return false, nil
}

func (m *mockCategoriesRepo) Update(ctx context.Context, categories *Categories) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, categories)
	}
	return nil
}

func (m *mockCategoriesRepo) Delete(ctx context.Context, categoryID string, userID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, categoryID, userID)
	}
	return nil
}

func TestCreate_Success(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, nil
		},
		createFunc: func(ctx context.Context, categories *Categories) error {
			assert.Equal(t, "USER-001", categories.USER_ID)
			assert.Equal(t, "test", categories.Name)
			assert.Equal(t, "test desc", categories.Description)
			return nil
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "test", "test desc")

	assert.NoError(t, err)
}

func TestCreate_ReturnsExactRepoError(t *testing.T) {
	expectedErr := errors.New("insert failed")

	mockRepo := &mockCategoriesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, nil
		},
		createFunc: func(ctx context.Context, categories *Categories) error {
			return expectedErr
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "test", "desc")

	assert.Equal(t, expectedErr, err)
}

func TestCreate_DBError(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		createFunc: func(ctx context.Context, category *Categories) error {
			return errors.New("internal server error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "test", "desc")

	assert.Error(t, err)
	assert.Equal(t, ErrInternal, err)
}

func TestGenerateCategoryID_Success(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 5, nil
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	categoryID, err := service.generateCategoryID(context.Background())

	assert.NoError(t, err)
	assert.Contains(t, categoryID, "CAT-")
	assert.Contains(t, categoryID, "-000006")
}

func TestCreate_ErrorFromGenerateCategoryID(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, errors.New("internal server error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "test", "test desc")

	assert.Error(t, err)
}

func TestGenerateCategoriesID_Error(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, errors.New("internal server error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	userID, err := service.generateCategoryID(context.Background())

	assert.Empty(t, userID)
	assert.Error(t, err)
}

func TestCreate_CategoryIDFormat(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, nil
		},
		createFunc: func(ctx context.Context, categories *Categories) error {
			assert.Regexp(t, `^CAT-\d{8}-\d{6}$`, categories.CATEGORY_ID)
			return nil
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "Food", "Desc")
	assert.NoError(t, err)
}

func TestGetAllByUserID_Success(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		getAllByUserIDFunc: func(ctx context.Context, userID string) ([]*Categories, error) {
			assert.Equal(t, "USER-001", userID)

			return []*Categories{
				{
					CATEGORY_ID: "CAT-20250101-000001",
					USER_ID:     "USER-001",
					Name:        "Food",
					Description: "Food category",
				},
			}, nil
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetAllByUserID(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Food", result[0].Name)
}

func TestGetAllByUserID_RepoError(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		getAllByUserIDFunc: func(ctx context.Context, userID string) ([]*Categories, error) {
			return nil, errors.New("database error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.GetAllByUserID(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreate_InvalidCredential(t *testing.T) {
	mockRepo := &mockCategoriesRepo{}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	err := service.Create(context.Background(), "test", "test desc")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		hasCategoryBudgetFunc: func(ctx context.Context, categoryID string, userID string) (bool, error) {
			assert.Equal(t, "CAT-20250101-000001", categoryID)
			assert.Equal(t, "USER-001", userID)
			return false, nil
		},
		deleteFunc: func(ctx context.Context, categoryID string, userID string) error {
			assert.Equal(t, "CAT-20250101-000001", categoryID)
			assert.Equal(t, "USER-001", userID)
			return nil
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "CAT-20250101-000001")

	assert.NoError(t, err)
}

func TestDelete_InvalidCredential(t *testing.T) {
	mockRepo := &mockCategoriesRepo{}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	err := service.Delete(context.Background(), "CAT-20250101-000001")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestDelete_CategoryNotFound(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		hasCategoryBudgetFunc: func(ctx context.Context, categoryID string, userID string) (bool, error) {
			return false, nil
		},
		deleteFunc: func(ctx context.Context, categoryID string, userID string) error {
			return sql.ErrNoRows
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "CAT-001")

	assert.Equal(t, ErrCategoryNotFound, err)
}

func TestDelete_CategoryHasBudget(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		hasCategoryBudgetFunc: func(ctx context.Context, categoryID string, userID string) (bool, error) {
			assert.Equal(t, "CAT-001", categoryID)
			assert.Equal(t, "USER-001", userID)
			return true, nil
		},
		deleteFunc: func(ctx context.Context, categoryID string, userID string) error {
			t.Fatal("Delete should not be called when category budget still exists")
			return nil
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "CAT-001")

	assert.Equal(t, ErrCategoryHasBudget, err)
}

func TestDelete_RepoError(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		hasCategoryBudgetFunc: func(ctx context.Context, categoryID string, userID string) (bool, error) {
			return false, nil
		},
		deleteFunc: func(ctx context.Context, categoryID string, userID string) error {
			return errors.New("internal server error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "CAT-20250101-000001")

	assert.Equal(t, ErrInternal, err)
}

func TestDelete_HasCategoryBudgetError(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		hasCategoryBudgetFunc: func(ctx context.Context, categoryID string, userID string) (bool, error) {
			return false, errors.New("database error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Delete(ctx, "CAT-20250101-000001")

	assert.Equal(t, ErrInternal, err)
}
