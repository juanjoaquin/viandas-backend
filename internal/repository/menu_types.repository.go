package repository

import (
	"context"
	"fmt"
	"strings"

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

func buildMenuTypeWhere(nameQuery string, activeFilter *bool) (string, []interface{}) {
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

func (r *repo) CountMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool) (int, error) {
	where, args := buildMenuTypeWhere(nameQuery, activeFilter)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM menu_types`+where, args...)
	return count, err
}

func (r *repo) GetMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]entity.MenuType, error) {
	where, args := buildMenuTypeWhere(nameQuery, activeFilter)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM menu_types%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var types []entity.MenuType
	err := r.db.SelectContext(ctx, &types, query, args...)
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
