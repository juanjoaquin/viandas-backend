package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveProductCategory(ctx context.Context, name string) (*entity.ProductCategory, error) {
	var category entity.ProductCategory
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO product_categories (name) VALUES ($1) RETURNING *`,
		name,
	).StructScan(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func buildProductCategoryWhere(nameQuery string, activeFilter *bool) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if nameQuery != "" {
		args = append(args, "%"+nameQuery+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if activeFilter != nil {
		args = append(args, *activeFilter)
		conditions = append(conditions, fmt.Sprintf("active = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *repo) CountProductCategories(ctx context.Context, nameQuery string, activeFilter *bool) (int, error) {
	where, args := buildProductCategoryWhere(nameQuery, activeFilter)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM product_categories`+where, args...)
	return count, err
}

func (r *repo) GetProductCategories(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]entity.ProductCategory, error) {
	where, args := buildProductCategoryWhere(nameQuery, activeFilter)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM product_categories%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var categories []entity.ProductCategory
	err := r.db.SelectContext(ctx, &categories, query, args...)
	return categories, err
}

func (r *repo) GetProductCategoryByID(ctx context.Context, id string) (*entity.ProductCategory, error) {
	var category entity.ProductCategory
	err := r.db.GetContext(ctx, &category, `SELECT * FROM product_categories WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repo) UpdateProductCategory(ctx context.Context, id, name string, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE product_categories SET name=$1, active=$2, updated_at=NOW() WHERE id=$3`,
		name, active, id,
	)
	return err
}

func (r *repo) DeleteProductCategory(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_categories WHERE id = $1`, id)
	return err
}
