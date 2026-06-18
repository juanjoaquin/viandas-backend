package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveDelivery(ctx context.Context, name string, phone *string) (*entity.Delivery, error) {
	var d entity.Delivery
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO deliveries (name, phone) VALUES ($1, $2) RETURNING *`,
		name, phone,
	).StructScan(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) GetDeliveries(ctx context.Context, nameQuery string) ([]entity.Delivery, error) {
	var deliveries []entity.Delivery
	var err error

	if nameQuery == "" {
		err = r.db.SelectContext(ctx, &deliveries, `SELECT * FROM deliveries ORDER BY name`)
	} else {
		err = r.db.SelectContext(ctx, &deliveries,
			`SELECT * FROM deliveries WHERE name ILIKE $1 ORDER BY name`,
			"%"+nameQuery+"%",
		)
	}
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

func (r *repo) UpdateDelivery(ctx context.Context, id, name string, phone *string, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE deliveries SET name=$1, phone=$2, active=$3, updated_at=NOW() WHERE id=$4`,
		name, phone, active, id,
	)
	return err
}

func (r *repo) DeleteDelivery(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM deliveries WHERE id = $1`, id)
	return err
}
