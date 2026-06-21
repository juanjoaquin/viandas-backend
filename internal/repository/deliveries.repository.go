package repository

import (
	"context"
	"fmt"

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

func buildDeliveryWhere(nameQuery string) (string, []interface{}) {
	if nameQuery == "" {
		return "", nil
	}
	return " WHERE name ILIKE $1", []interface{}{"%" + nameQuery + "%"}
}

func (r *repo) CountDeliveries(ctx context.Context, nameQuery string) (int, error) {
	where, args := buildDeliveryWhere(nameQuery)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM deliveries`+where, args...)
	return count, err
}

func (r *repo) GetDeliveries(ctx context.Context, nameQuery string, offset, limit int) ([]entity.Delivery, error) {
	where, args := buildDeliveryWhere(nameQuery)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM deliveries%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var deliveries []entity.Delivery
	err := r.db.SelectContext(ctx, &deliveries, query, args...)
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
