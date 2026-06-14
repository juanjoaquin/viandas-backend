package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveDelivery(ctx context.Context, name string) (*entity.Delivery, error) {
	var d entity.Delivery
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO deliveries (name) VALUES ($1) RETURNING *`,
		name,
	).StructScan(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) GetDeliveries(ctx context.Context) ([]entity.Delivery, error) {
	var deliveries []entity.Delivery
	err := r.db.SelectContext(ctx, &deliveries, `SELECT * FROM deliveries ORDER BY name`)
	return deliveries, err
}

func (r *repo) GetDeliveryByID(ctx context.Context, id string) (*entity.Delivery, error) {
	var d entity.Delivery
	err := r.db.GetContext(ctx, &d, `SELECT * FROM deliveries WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) UpdateDelivery(ctx context.Context, id, name string, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE deliveries SET name=$1, active=$2, updated_at=NOW() WHERE id=$3`,
		name, active, id,
	)
	return err
}

func (r *repo) DeleteDelivery(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM deliveries WHERE id = $1`, id)
	return err
}
