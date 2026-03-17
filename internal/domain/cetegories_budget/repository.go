package cetegoriesbudget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CategoriesBudgetRepository interface {
	CountByDate(ctx context.Context, date string) (int, error)
	Create(ctx context.Context, category *CategoriesBudget) error
	Update(ctx context.Context, category *CategoriesBudget) error
	GetByCategoryID(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error)
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
		return ErrCategoryBudgetNotFound
	}

	return nil
}

func (r *categoriesBudgetRepository) GetByCategoryID(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error) {
	query := `
		SELECT
			id,
			budget_id,
			user_id,
			category_id,
			allocated_amount,
			used_amount,
			period
		FROM
			category_budgets
		WHERE
			user_id = $1
			AND category_id = $2
	`

	var categoryBudget CategoriesBudget

	err := r.db.GetContext(ctx, &categoryBudget, query, userID, categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryBudgetNotFound
		}

		return nil, err
	}

	return &categoryBudget, nil
}
