package repository

import (
	"context"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveMenuType(ctx context.Context, name string, sortOrder int) (*entity.MenuType, error) {
	var mt entity.MenuType
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO menu_types (name, sort_order) VALUES ($1, $2) RETURNING *`,
		name, sortOrder,
	).StructScan(&mt)
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repo) GetMenuTypes(ctx context.Context) ([]entity.MenuType, error) {
	var types []entity.MenuType
	err := r.db.SelectContext(ctx, &types, `SELECT * FROM menu_types ORDER BY sort_order, name`)
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

func (r *repo) UpdateMenuType(ctx context.Context, id, name string, sortOrder int, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu_types SET name=$1, sort_order=$2, active=$3, updated_at=NOW() WHERE id=$4`,
		name, sortOrder, active, id,
	)
	return err
}

func (r *repo) DeleteMenuType(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu_types WHERE id = $1`, id)
	return err
}
