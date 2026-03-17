package salaries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SalariesRepository interface {
	CountByDate(ctx context.Context, date string) (int, error)
	Create(ctx context.Context, salary *Salaries) error
	CheckSalary(ctx context.Context, salary *Salaries) (int, error)
	GetTotalSalary(ctx context.Context, salary *Salaries) (int64, error)
	GetAllByUserID(ctx context.Context, salary *Salaries) ([]*Salaries, error)
	GetBySalaryID(ctx context.Context, salary *Salaries) (*Salaries, error)
	Update(ctx context.Context, salary *Salaries) error
	Delete(ctx context.Context, salaryID, userID string) error
}

type salariesRepository struct {
	db *sqlx.DB
}

func NewSalariesRepository(db *sqlx.DB) SalariesRepository {
	return &salariesRepository{db: db}
}

func NewCategoriesBudgetRepository(db *sqlx.DB) SalariesRepository {
	return NewSalariesRepository(db)
}

func (r *salariesRepository) CountByDate(ctx context.Context, date string) (int, error) {
	var count int

	pattern := fmt.Sprintf("SAL-%s-%%", date)

	query := `
		SELECT COUNT(*) 
		FROM salaries 
		WHERE salary_id LIKE $1
	`

	err := r.db.GetContext(ctx, &count, query, pattern)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *salariesRepository) Create(ctx context.Context, salary *Salaries) error {
	query := `
		INSERT INTO salaries
		(salary_id, user_id, amount, source, description, received_at)
		VALUES ($1, $2, $3::bigint, $4, $5, $6::timestamptz)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		salary.SalaryID,
		salary.UserID,
		salary.Amount,
		salary.Source,
		salary.Description,
		salary.ReceivedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *salariesRepository) CheckSalary(ctx context.Context, salary *Salaries) (int, error) {
	if salary == nil {
		return 0, errors.New("salary is nil")
	}

	if salary.UserID == "" {
		return 0, errors.New("user_id is required")
	}

	if salary.ReceivedAt == "" {
		return 0, errors.New("received_at is required")
	}

	var exists int

	query := `
		SELECT CASE
			WHEN EXISTS (
				SELECT 1
				FROM salaries
				WHERE user_id = $1
				  AND received_at >= DATE_TRUNC('month', $2::timestamptz)
				  AND received_at < DATE_TRUNC('month', $2::timestamptz) + INTERVAL '1 month'
			) THEN 1
			ELSE 0
		END
	`

	err := r.db.GetContext(ctx, &exists, query, salary.UserID, salary.ReceivedAt)
	if err != nil {
		return 0, err
	}

	return exists, nil
}

func (r *salariesRepository) GetTotalSalary(ctx context.Context, salary *Salaries) (int64, error) {
	if salary == nil {
		return 0, errors.New("salary is nil")
	}

	if salary.UserID == "" {
		return 0, errors.New("user_id is required")
	}

	if salary.ReceivedAt == "" {
		return 0, errors.New("received_at is required")
	}

	var total int64

	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM salaries
		WHERE user_id = $1
		  AND received_at >= DATE_TRUNC('month', $2::timestamptz)
		  AND received_at < DATE_TRUNC('month', $2::timestamptz) + INTERVAL '1 month'
	`

	err := r.db.GetContext(ctx, &total, query, salary.UserID, salary.ReceivedAt)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *salariesRepository) GetAllByUserID(ctx context.Context, salary *Salaries) ([]*Salaries, error) {
	if salary == nil {
		return nil, errors.New("salary is nil")
	}

	if salary.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	if salary.ReceivedAt == "" {
		return nil, errors.New("received_at is required")
	}

	var salaries []*Salaries

	query := `
		SELECT
			id,
			salary_id,
			user_id,
			amount,
			source,
			description,
			received_at
		FROM salaries
		WHERE user_id = $1
		  AND received_at >= DATE_TRUNC('month', $2::timestamptz)
		  AND received_at < DATE_TRUNC('month', $2::timestamptz) + INTERVAL '1 month'
		ORDER BY received_at DESC, salary_id DESC
	`

	err := r.db.SelectContext(ctx, &salaries, query, salary.UserID, salary.ReceivedAt)
	if err != nil {
		return nil, err
	}

	return salaries, nil
}

func (r *salariesRepository) GetBySalaryID(ctx context.Context, salary *Salaries) (*Salaries, error) {
	if salary == nil {
		return nil, errors.New("salary is nil")
	}

	if salary.SalaryID == "" {
		return nil, errors.New("salary_id is required")
	}

	if salary.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	var salaryData Salaries

	query := `
		SELECT
			id,
			salary_id,
			user_id,
			amount,
			source,
			description,
			received_at
		FROM salaries
		WHERE salary_id = $1
		  AND user_id = $2
	`

	err := r.db.GetContext(ctx, &salaryData, query, salary.SalaryID, salary.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSalaryNotFound
		}

		return nil, err
	}

	return &salaryData, nil
}

func (r *salariesRepository) Update(ctx context.Context, salary *Salaries) error {
	if salary == nil {
		return errors.New("salary is nil")
	}

	if salary.SalaryID == "" {
		return errors.New("salary_id is required")
	}

	if salary.UserID == "" {
		return errors.New("user_id is required")
	}

	query := `
		UPDATE salaries
		SET amount = $1::bigint,
		    source = $2,
		    description = $3,
		    updated_at = NOW()
		WHERE salary_id = $4
		  AND user_id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		salary.Amount,
		salary.Source,
		salary.Description,
		salary.SalaryID,
		salary.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSalaryNotFound
	}

	return nil
}

func (r *salariesRepository) Delete(ctx context.Context, salaryID, userID string) error {
	if salaryID == "" {
		return errors.New("salary_id is required")
	}

	if userID == "" {
		return errors.New("user_id is required")
	}

	query := `
		DELETE FROM salaries
		WHERE salary_id = $1
		  AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, salaryID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSalaryNotFound
	}

	return nil
}
