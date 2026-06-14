package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveExtraProduct(ctx context.Context, name, category string) (*entity.ExtraProduct, error) {
	var e entity.ExtraProduct
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO extra_products (name, category) VALUES ($1, $2) RETURNING *`,
		name, category,
	).StructScan(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repo) GetExtraProducts(ctx context.Context) ([]entity.ExtraProduct, error) {
	var products []entity.ExtraProduct
	err := r.db.SelectContext(ctx, &products, `SELECT * FROM extra_products ORDER BY category, name`)
	return products, err
}

func (r *repo) GetExtraProductByID(ctx context.Context, id string) (*entity.ExtraProduct, error) {
	var e entity.ExtraProduct
	err := r.db.GetContext(ctx, &e, `SELECT * FROM extra_products WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repo) UpdateExtraProduct(ctx context.Context, id, name string, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE extra_products SET name=$1, active=$2, updated_at=NOW() WHERE id=$3`,
		name, active, id,
	)
	return err
}

func (r *repo) DeleteExtraProduct(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM extra_products WHERE id = $1`, id)
	return err
}
