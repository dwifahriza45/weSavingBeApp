package categories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CategoriesRepository interface {
	CountByDate(ctx context.Context, date string) (int, error)
	Create(ctx context.Context, categories *Categories) error
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
		WHERE user_id LIKE $1
	`

	err := r.db.GetContext(ctx, &count, query, pattern)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *categoriesRepository) Create(ctx context.Context, category *Categories) error {
	query := `INSERT INTO categories (category_id, user_id, name, description, updated_at) VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.db.ExecContext(ctx, query, category.CATEGORY_ID, category.USER_ID, category.Name, category.Description)
	return err
}
