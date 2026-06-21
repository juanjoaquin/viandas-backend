package repository

import (
	"context"
	"fmt"
	"strings"

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

func buildDishWhere(nameQuery, menuTypeID string) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if menuTypeID != "" {
		args = append(args, menuTypeID)
		conditions = append(conditions, fmt.Sprintf("menu_type_id = $%d", len(args)))
		conditions = append(conditions, "active = TRUE")
	}
	if nameQuery != "" {
		args = append(args, "%"+nameQuery+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *repo) CountDishes(ctx context.Context, nameQuery, menuTypeID string) (int, error) {
	where, args := buildDishWhere(nameQuery, menuTypeID)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM dishes`+where, args...)
	return count, err
}

func (r *repo) GetDishes(ctx context.Context, nameQuery, menuTypeID string, offset, limit int) ([]entity.Dish, error) {
	where, args := buildDishWhere(nameQuery, menuTypeID)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM dishes%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var dishes []entity.Dish
	err := r.db.SelectContext(ctx, &dishes, query, args...)
	return dishes, err
}

func (r *repo) GetDishesByMenuTypeID(ctx context.Context, menuTypeID string, offset, limit int) ([]entity.Dish, error) {
	return r.GetDishes(ctx, "", menuTypeID, offset, limit)
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
