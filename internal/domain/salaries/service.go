package salaries

import (
	"BE_WE_SAVING/internal/app/middleware"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type SalariesService interface {
	Create(ctx context.Context, amount, source, description string) error
	CheckSalary(ctx context.Context) (int, error)
	GetTotalSalary(ctx context.Context) (int64, error)
	GetAllByUserID(ctx context.Context) ([]*Salaries, error)
	GetBySalaryID(ctx context.Context, salaryID string) (*Salaries, error)
	Update(ctx context.Context, salaryID, amount, source, description string) error
	Delete(ctx context.Context, salaryID string) error
}

type salariesService struct {
	salaryRepo SalariesRepository
	now        func() time.Time
}

func NewSalariesService(salaryRepo SalariesRepository) SalariesService {
	return &salariesService{
		salaryRepo: salaryRepo,
		now:        time.Now,
	}
}

var (
	ErrInvalidCredentials = errors.New("Invalid Credentials")
	ErrInvalidAmount      = errors.New("amount must be a valid integer")
	ErrInternal           = errors.New("internal server error")
	ErrSalaryNotFound     = errors.New("salary not found")
)

func (s *salariesService) Create(ctx context.Context, amount, source, description string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	parsedAmount, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return ErrInvalidAmount
	}

	receivedAtTime := s.currentTime()

	salaryID, err := s.generateSalaryID(ctx, receivedAtTime)
	if err != nil {
		return err
	}

	salary := &Salaries{
		SalaryID:    salaryID,
		UserID:      userID,
		Amount:      strconv.FormatInt(parsedAmount, 10),
		Source:      source,
		Description: description,
		ReceivedAt:  receivedAtTime.Format(time.RFC3339),
	}

	return s.salaryRepo.Create(ctx, salary)
}

func (s *salariesService) CheckSalary(ctx context.Context) (int, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return 0, ErrInvalidCredentials
	}

	referenceTime := s.currentTime().Format(time.RFC3339)

	exists, err := s.salaryRepo.CheckSalary(ctx, &Salaries{
		UserID:     userID,
		ReceivedAt: referenceTime,
	})
	if err != nil {
		return 0, ErrInternal
	}

	return exists, nil
}

func (s *salariesService) GetTotalSalary(ctx context.Context) (int64, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return 0, ErrInvalidCredentials
	}

	referenceSalary := &Salaries{
		UserID:     userID,
		ReceivedAt: s.currentTime().Format(time.RFC3339),
	}

	exists, err := s.salaryRepo.CheckSalary(ctx, referenceSalary)
	if err != nil {
		return 0, ErrInternal
	}

	if exists == 0 {
		return 0, nil
	}

	total, err := s.salaryRepo.GetTotalSalary(ctx, referenceSalary)
	if err != nil {
		return 0, ErrInternal
	}

	return total, nil
}

func (s *salariesService) GetAllByUserID(ctx context.Context) ([]*Salaries, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	referenceSalary := &Salaries{
		UserID:     userID,
		ReceivedAt: s.currentTime().Format(time.RFC3339),
	}

	salaries, err := s.salaryRepo.GetAllByUserID(ctx, referenceSalary)
	if err != nil {
		return nil, ErrInternal
	}

	return salaries, nil
}

func (s *salariesService) GetBySalaryID(ctx context.Context, salaryID string) (*Salaries, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	salary, err := s.salaryRepo.GetBySalaryID(ctx, &Salaries{
		SalaryID: salaryID,
		UserID:   userID,
	})
	if err != nil {
		if errors.Is(err, ErrSalaryNotFound) {
			return nil, ErrSalaryNotFound
		}

		return nil, ErrInternal
	}

	return salary, nil
}

func (s *salariesService) Update(ctx context.Context, salaryID, amount, source, description string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	parsedAmount, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return ErrInvalidAmount
	}

	salary := &Salaries{
		SalaryID:    salaryID,
		UserID:      userID,
		Amount:      strconv.FormatInt(parsedAmount, 10),
		Source:      source,
		Description: description,
	}

	err = s.salaryRepo.Update(ctx, salary)
	if err != nil {
		if errors.Is(err, ErrSalaryNotFound) {
			return ErrSalaryNotFound
		}

		return ErrInternal
	}

	return nil
}

func (s *salariesService) Delete(ctx context.Context, salaryID string) error {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return ErrInvalidCredentials
	}

	err := s.salaryRepo.Delete(ctx, salaryID, userID)
	if err != nil {
		if errors.Is(err, ErrSalaryNotFound) {
			return ErrSalaryNotFound
		}

		return ErrInternal
	}

	return nil
}

func (s *salariesService) generateSalaryID(ctx context.Context, receivedAt time.Time) (string, error) {
	dateStr := receivedAt.Format("20060102")

	count, err := s.salaryRepo.CountByDate(ctx, dateStr)
	if err != nil {
		return "", err
	}

	sequence := count + 1

	salaryID := fmt.Sprintf("SAL-%s-%06d", dateStr, sequence)

	return salaryID, nil
}

func (s *salariesService) currentTime() time.Time {
	if s.now != nil {
		return s.now()
	}

	return time.Now()
}
