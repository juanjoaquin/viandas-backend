package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveDish(ctx context.Context, name, description, menuTypeID string) (*entity.Dish, error) {
	var d entity.Dish
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO dishes (name, description, menu_type_id) VALUES ($1, $2, $3) RETURNING *`,
		name, description, menuTypeID,
	).StructScan(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) GetDishes(ctx context.Context, nameQuery string) ([]entity.Dish, error) {
	var dishes []entity.Dish
	var err error

	if nameQuery == "" {
		err = r.db.SelectContext(ctx, &dishes, `SELECT * FROM dishes ORDER BY name`)
	} else {
		err = r.db.SelectContext(ctx, &dishes,
			`SELECT * FROM dishes WHERE name ILIKE $1 ORDER BY name`,
			"%"+nameQuery+"%",
		)
	}
	return dishes, err
}

func (r *repo) GetDishesByMenuTypeID(ctx context.Context, menuTypeID string) ([]entity.Dish, error) {
	var dishes []entity.Dish
	err := r.db.SelectContext(ctx, &dishes,
		`SELECT * FROM dishes WHERE menu_type_id = $1 AND active = TRUE ORDER BY name`,
		menuTypeID,
	)
	return dishes, err
}

func (r *repo) GetDishByID(ctx context.Context, id string) (*entity.Dish, error) {
	var d entity.Dish
	err := r.db.GetContext(ctx, &d, `SELECT * FROM dishes WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) UpdateDish(ctx context.Context, id, name, description, menuTypeID string, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE dishes SET name=$1, description=$2, menu_type_id=$3, active=$4, updated_at=NOW() WHERE id=$5`,
		name, description, menuTypeID, active, id,
	)
	return err
}

func (r *repo) DeleteDish(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM dishes WHERE id = $1`, id)
	return err
}
