package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveMenuType(ctx context.Context, name string, price *float64) (*entity.MenuType, error) {
	var mt entity.MenuType
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO menu_types (name, price) VALUES ($1, $2) RETURNING *`,
		name, price,
	).StructScan(&mt)
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repo) GetMenuTypes(ctx context.Context) ([]entity.MenuType, error) {
	var types []entity.MenuType
	err := r.db.SelectContext(ctx, &types, `SELECT * FROM menu_types ORDER BY name`)
	return types, err
}

func (r *repo) GetMenuTypeByID(ctx context.Context, id string) (*entity.MenuType, error) {
	var mt entity.MenuType
	err := r.db.GetContext(ctx, &mt, `SELECT * FROM menu_types WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repo) UpdateMenuType(ctx context.Context, id, name string, price *float64, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu_types SET name=$1, price=$2, active=$3, updated_at=NOW() WHERE id=$4`,
		name, price, active, id,
	)
	return err
}

func (r *repo) DeleteMenuType(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu_types WHERE id = $1`, id)
	return err
}
