package cetegoriesbudget

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CategoriesBudgetRepository interface {
}

type categoriesBudgetRepository struct {
	db *sqlx.DB
}

func NewCategoriesBudgetRepository(db *sqlx.DB) CategoriesBudgetRepository {
	return &categoriesBudgetRepository{db: db}
}

func (r *categoriesBudgetRepository) CountByDate(ctx context.Context, date string) (int, error) {
	var count int

	pattern := fmt.Sprintf("BUD-%s-%%", date)

	query := `
		SELECT COUNT(*) 
		FROM category_budgets 
		WHERE budget_id LIKE $1
	`

	err := r.db.GetContext(ctx, &count, query, pattern)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *categoriesBudgetRepository) Create(ctx context.Context, category *CategoriesBudget) error {
	query := `
		INSERT INTO category_budgets
		(budget_id, user_id, category_id, allocated_amount, used_amount, period)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		category.BUDGET_ID,
		category.USER_ID,
		category.CATEGORY_ID,
		category.AllocatedAmount,
		category.UsedAmount,
		category.Period,
	)

	return err
}

func (r *categoriesBudgetRepository) Update(ctx context.Context, category *CategoriesBudget) error {
	query := `
		UPDATE
			category_budgets
		SET
			allocated_amount = $1, 
		    used_amount = $2, 
		    updated_at = NOW() 
		WHERE
			category_id = $3 
			AND user_id = $4
			AND budget_id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		category.AllocatedAmount,
		category.UsedAmount,
		category.CATEGORY_ID,
		category.USER_ID,
		category.BUDGET_ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category budget not found")
	}

	return nil
}
