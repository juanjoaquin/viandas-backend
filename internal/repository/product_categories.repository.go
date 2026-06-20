package repository

import (
	"context"

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

func (r *repo) GetProductCategories(ctx context.Context, nameQuery string, activeFilter *bool) ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	var err error

	switch {
	case nameQuery == "" && activeFilter == nil:
		err = r.db.SelectContext(ctx, &categories, `SELECT * FROM product_categories ORDER BY name`)
	case nameQuery != "" && activeFilter == nil:
		err = r.db.SelectContext(ctx, &categories,
			`SELECT * FROM product_categories WHERE name ILIKE $1 ORDER BY name`,
			"%"+nameQuery+"%",
		)
	case nameQuery == "" && activeFilter != nil:
		err = r.db.SelectContext(ctx, &categories,
			`SELECT * FROM product_categories WHERE active = $1 ORDER BY name`,
			*activeFilter,
		)
	default:
		err = r.db.SelectContext(ctx, &categories,
			`SELECT * FROM product_categories WHERE name ILIKE $1 AND active = $2 ORDER BY name`,
			"%"+nameQuery+"%", *activeFilter,
		)
	}
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
