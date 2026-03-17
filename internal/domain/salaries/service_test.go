package salaries

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockSalariesRepo struct {
	countByDateFunc func(ctx context.Context, date string) (int, error)
	createFunc      func(ctx context.Context, salary *Salaries) error
	checkSalaryFunc func(ctx context.Context, salary *Salaries) (int, error)
	getTotalFunc    func(ctx context.Context, salary *Salaries) (int64, error)
}

func fixedNow() time.Time {
	return time.Date(2026, time.March, 17, 9, 30, 0, 0, time.FixedZone("WIB", 7*60*60))
}

func (m *mockSalariesRepo) CountByDate(ctx context.Context, date string) (int, error) {
	if m.countByDateFunc != nil {
		return m.countByDateFunc(ctx, date)
	}

	return 0, nil
}

func (m *mockSalariesRepo) Create(ctx context.Context, salary *Salaries) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, salary)
	}

	return nil
}

func (m *mockSalariesRepo) CheckSalary(ctx context.Context, salary *Salaries) (int, error) {
	if m.checkSalaryFunc != nil {
		return m.checkSalaryFunc(ctx, salary)
	}

	return 0, nil
}

func (m *mockSalariesRepo) GetTotalSalary(ctx context.Context, salary *Salaries) (int64, error) {
	if m.getTotalFunc != nil {
		return m.getTotalFunc(ctx, salary)
	}

	return 0, nil
}

func TestCreateSalary_Success(t *testing.T) {
	now := fixedNow()
	expectedDate := now.Format("20060102")
	expectedDatePrefix := now.Format("2006-01-02")

	mockRepo := &mockSalariesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			assert.Equal(t, expectedDate, date)
			return 0, nil
		},
		createFunc: func(ctx context.Context, salary *Salaries) error {
			assert.Equal(t, "USER-001", salary.UserID)
			assert.Equal(t, "SAL-"+expectedDate+"-000001", salary.SalaryID)
			assert.Equal(t, "5000000", salary.Amount)
			assert.Equal(t, "main job", salary.Source)
			assert.Equal(t, "monthly salary", salary.Description)
			assert.Contains(t, salary.ReceivedAt, expectedDatePrefix)
			return nil
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
		now:        func() time.Time { return now },
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "5000000", "main job", "monthly salary")

	assert.NoError(t, err)
}

func TestCreateSalary_InvalidCredential(t *testing.T) {
	service := &salariesService{
		salaryRepo: &mockSalariesRepo{},
	}

	err := service.Create(context.Background(), "5000000", "", "")

	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestCreateSalary_InvalidAmount(t *testing.T) {
	service := &salariesService{
		salaryRepo: &mockSalariesRepo{},
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "abc", "", "")

	assert.Equal(t, ErrInvalidAmount, err)
}

func TestCreateSalary_GenerateSalaryIDError(t *testing.T) {
	expectedErr := errors.New("database error")

	mockRepo := &mockSalariesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 0, expectedErr
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
		now:        fixedNow,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "5000000", "", "")

	assert.Equal(t, expectedErr, err)
}

func TestCreateSalary_RepoCreateError(t *testing.T) {
	expectedErr := errors.New("insert failed")

	mockRepo := &mockSalariesRepo{
		countByDateFunc: func(ctx context.Context, date string) (int, error) {
			return 1, nil
		},
		createFunc: func(ctx context.Context, salary *Salaries) error {
			return expectedErr
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	err := service.Create(ctx, "5000000", "", "")

	assert.Equal(t, expectedErr, err)
}

func TestCheckSalary_Success(t *testing.T) {
	now := fixedNow()

	mockRepo := &mockSalariesRepo{
		checkSalaryFunc: func(ctx context.Context, salary *Salaries) (int, error) {
			assert.Equal(t, "USER-001", salary.UserID)
			assert.Equal(t, now.Format(time.RFC3339), salary.ReceivedAt)
			return 1, nil
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
		now:        func() time.Time { return now },
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.CheckSalary(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestCheckSalary_InvalidCredential(t *testing.T) {
	service := &salariesService{
		salaryRepo: &mockSalariesRepo{},
	}

	result, err := service.CheckSalary(context.Background())

	assert.Equal(t, 0, result)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestCheckSalary_RepoError(t *testing.T) {
	mockRepo := &mockSalariesRepo{
		checkSalaryFunc: func(ctx context.Context, salary *Salaries) (int, error) {
			return 0, errors.New("db error")
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	result, err := service.CheckSalary(ctx)

	assert.Equal(t, 0, result)
	assert.Equal(t, ErrInternal, err)
}

func TestGetTotalSalary_Success(t *testing.T) {
	now := fixedNow()

	mockRepo := &mockSalariesRepo{
		checkSalaryFunc: func(ctx context.Context, salary *Salaries) (int, error) {
			assert.Equal(t, "USER-001", salary.UserID)
			assert.Equal(t, now.Format(time.RFC3339), salary.ReceivedAt)
			return 1, nil
		},
		getTotalFunc: func(ctx context.Context, salary *Salaries) (int64, error) {
			assert.Equal(t, "USER-001", salary.UserID)
			assert.Equal(t, now.Format(time.RFC3339), salary.ReceivedAt)
			return 7500000, nil
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
		now:        func() time.Time { return now },
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	total, err := service.GetTotalSalary(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(7500000), total)
}

func TestGetTotalSalary_InvalidCredential(t *testing.T) {
	service := &salariesService{
		salaryRepo: &mockSalariesRepo{},
	}

	total, err := service.GetTotalSalary(context.Background())

	assert.Equal(t, int64(0), total)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestGetTotalSalary_RepoError(t *testing.T) {
	mockRepo := &mockSalariesRepo{
		checkSalaryFunc: func(ctx context.Context, salary *Salaries) (int, error) {
			return 1, nil
		},
		getTotalFunc: func(ctx context.Context, salary *Salaries) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	total, err := service.GetTotalSalary(ctx)

	assert.Equal(t, int64(0), total)
	assert.Equal(t, ErrInternal, err)
}

func TestGetTotalSalary_CheckSalaryError(t *testing.T) {
	mockRepo := &mockSalariesRepo{
		checkSalaryFunc: func(ctx context.Context, salary *Salaries) (int, error) {
			return 0, errors.New("db error")
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
		now:        fixedNow,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	total, err := service.GetTotalSalary(ctx)

	assert.Equal(t, int64(0), total)
	assert.Equal(t, ErrInternal, err)
}

func TestGetTotalSalary_NoSalaryInCurrentMonth(t *testing.T) {
	mockRepo := &mockSalariesRepo{
		checkSalaryFunc: func(ctx context.Context, salary *Salaries) (int, error) {
			return 0, nil
		},
		getTotalFunc: func(ctx context.Context, salary *Salaries) (int64, error) {
			t.Fatal("GetTotalSalary should not be called when current month has no salary")
			return 0, nil
		},
	}

	service := &salariesService{
		salaryRepo: mockRepo,
		now:        fixedNow,
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "USER-001")

	total, err := service.GetTotalSalary(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
}
