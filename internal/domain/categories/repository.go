package categories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CategoriesRepository interface {
	CountByDate(ctx context.Context, date string) (int, error)
	Create(ctx context.Context, categories *Categories) error
	GetAllByUserID(ctx context.Context, userID string) ([]*Categories, error)
	GetByCategoryID(ctx context.Context, userID, categoryID string) (*Categories, error)
	Update(ctx context.Context, categories *Categories) error
	Delete(ctx context.Context, categoryID string, userID string) error
}

type categoriesRepository struct {
	db *sqlx.DB
}

func NewCategoriesRepository(db *sqlx.DB) CategoriesRepository {
	return &categoriesRepository{db: db}
}

func (r *categoriesRepository) CountByDate(ctx context.Context, date string) (int, error) {
	var count int

	pattern := fmt.Sprintf("CAT-%s-%%", date)

	query := `
		SELECT COUNT(*) 
		FROM categories 
		WHERE category_id LIKE $1
	`

	err := r.db.GetContext(ctx, &count, query, pattern)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *categoriesRepository) Create(ctx context.Context, category *Categories) error {
	query := `INSERT INTO categories (category_id, user_id, name, description) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, category.CATEGORY_ID, category.USER_ID, category.Name, category.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoriesRepository) GetAllByUserID(ctx context.Context, userID string) ([]*Categories, error) {
	query := `
		SELECT 
			id,
			category_id, 
			user_id, 
			name, 
			description 
		FROM 
			categories 
		WHERE user_id = $1
	`

	var categories []*Categories

	err := r.db.SelectContext(ctx, &categories, query, userID)
	if err != nil {
		return nil, err
	}

	return categories, nil

}

func (r *categoriesRepository) GetByCategoryID(ctx context.Context, userID, categoryID string) (*Categories, error) {
	query := `
		SELECT 
			id,
			category_id, 
			user_id, 
			name, 
			description 
		FROM 
			categories 
		WHERE
			user_id = $1
			AND category_id = $2
	`

	var categories Categories

	err := r.db.GetContext(ctx, &categories, query, userID, categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &categories, nil

}

func (r *categoriesRepository) Update(ctx context.Context, categories *Categories) error {
	query := `
		UPDATE categories 
		SET name = $1, 
		    description = $2, 
		    updated_at = NOW() 
		WHERE category_id = $3 
		  AND user_id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		categories.Name,
		categories.Description,
		categories.CATEGORY_ID,
		categories.USER_ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *categoriesRepository) Delete(ctx context.Context, categoryID string, userID string) error {
	query := `
		DELETE FROM categories 
		WHERE category_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, categoryID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
