package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveExtraProduct(ctx context.Context, name, categoryID string, price float64) (*entity.ExtraProduct, error) {
	var e entity.ExtraProduct
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO extra_products (name, category_id, price) VALUES ($1, $2, $3) RETURNING *`,
		name, categoryID, price,
	).StructScan(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repo) GetExtraProducts(ctx context.Context, nameQuery string) ([]entity.ExtraProduct, error) {
	var products []entity.ExtraProduct
	var err error

	if nameQuery == "" {
		err = r.db.SelectContext(ctx, &products,
			`SELECT ep.* FROM extra_products ep
			 JOIN product_categories pc ON pc.id = ep.category_id
			 ORDER BY pc.name, ep.name`,
		)
	} else {
		err = r.db.SelectContext(ctx, &products,
			`SELECT ep.* FROM extra_products ep
			 JOIN product_categories pc ON pc.id = ep.category_id
			 WHERE ep.name ILIKE $1
			 ORDER BY pc.name, ep.name`,
			"%"+nameQuery+"%",
		)
	}
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

func (r *repo) UpdateExtraProduct(ctx context.Context, id, name, categoryID string, price float64, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE extra_products SET name=$1, category_id=$2, price=$3, active=$4, updated_at=NOW() WHERE id=$5`,
		name, categoryID, price, active, id,
	)
	return err
}

func (r *repo) DeleteExtraProduct(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM extra_products WHERE id = $1`, id)
	return err
}
