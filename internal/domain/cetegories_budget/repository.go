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
	Delete(ctx context.Context, category *CategoriesBudget) error
	GetByCategoryID(ctx context.Context, userID, categoryID string) (*CategoriesBudget, error)
	GetByBudgetID(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error)
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
		(budget_id, user_id, category_id, allocated_amount, used_amount, budget_date)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		category.BUDGET_ID,
		category.USER_ID,
		category.CATEGORY_ID,
		category.AllocatedAmount,
		category.UsedAmount,
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
			user_id = $3
			AND budget_id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		category.AllocatedAmount,
		category.UsedAmount,
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

func (r *categoriesBudgetRepository) Delete(ctx context.Context, category *CategoriesBudget) error {
	query := `
		DELETE FROM category_budgets
		WHERE user_id = $1
		  AND budget_id = $2
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
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
			used_amount
		FROM
			category_budgets
		WHERE
			user_id = $1
			AND category_id = $2
		ORDER BY budget_date DESC, budget_id DESC
		LIMIT 1
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

func (r *categoriesBudgetRepository) GetByBudgetID(ctx context.Context, userID, budgetID string) (*CategoriesBudget, error) {
	query := `
		SELECT
			id,
			budget_id,
			user_id,
			category_id,
			allocated_amount,
			used_amount
		FROM
			category_budgets
		WHERE
			user_id = $1
			AND budget_id = $2
		LIMIT 1
	`

	var categoryBudget CategoriesBudget

	err := r.db.GetContext(ctx, &categoryBudget, query, userID, budgetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryBudgetNotFound
		}

		return nil, err
	}

	return &categoryBudget, nil
}
