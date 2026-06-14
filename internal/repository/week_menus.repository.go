package repository

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveWeekMenu(ctx context.Context, weekStartDate time.Time, createdBy string) (*entity.WeekMenu, error) {
	var m entity.WeekMenu
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO week_menus (week_start_date, created_by) VALUES ($1, $2) RETURNING *`,
		weekStartDate, createdBy,
	).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repo) GetWeekMenus(ctx context.Context) ([]entity.WeekMenu, error) {
	var menus []entity.WeekMenu
	err := r.db.SelectContext(ctx, &menus, `SELECT * FROM week_menus ORDER BY week_start_date DESC`)
	return menus, err
}

func (r *repo) GetWeekMenuByID(ctx context.Context, id string) (*entity.WeekMenu, error) {
	var m entity.WeekMenu
	err := r.db.GetContext(ctx, &m, `SELECT * FROM week_menus WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repo) GetWeekMenuByDate(ctx context.Context, date time.Time) (*entity.WeekMenu, error) {
	var m entity.WeekMenu
	err := r.db.GetContext(ctx, &m,
		`SELECT * FROM week_menus WHERE week_start_date <= $1 AND week_start_date + INTERVAL '6 days' >= $1`,
		date,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repo) DeleteWeekMenu(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM week_menus WHERE id = $1`, id)
	return err
}

func (r *repo) SaveWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, traditionalDishID, healthyDishID, vegetarianDishID string) (*entity.WeekMenuItem, error) {
	var item entity.WeekMenuItem
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO week_menu_items (week_menu_id, menu_date, traditional_dish_id, healthy_dish_id, vegetarian_dish_id)
		 VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		weekMenuID, menuDate, traditionalDishID, healthyDishID, vegetarianDishID,
	).StructScan(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repo) GetWeekMenuItems(ctx context.Context, weekMenuID string) ([]entity.WeekMenuItem, error) {
	var items []entity.WeekMenuItem
	err := r.db.SelectContext(ctx, &items,
		`SELECT * FROM week_menu_items WHERE week_menu_id = $1 ORDER BY menu_date`,
		weekMenuID,
	)
	return items, err
}

func (r *repo) GetWeekMenuItemByDate(ctx context.Context, date time.Time) (*entity.WeekMenuItem, error) {
	var item entity.WeekMenuItem
	err := r.db.GetContext(ctx, &item,
		`SELECT * FROM week_menu_items WHERE menu_date = $1`,
		date,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repo) UpdateWeekMenuItem(ctx context.Context, id, traditionalDishID, healthyDishID, vegetarianDishID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE week_menu_items SET traditional_dish_id=$1, healthy_dish_id=$2, vegetarian_dish_id=$3, updated_at=NOW() WHERE id=$4`,
		traditionalDishID, healthyDishID, vegetarianDishID, id,
	)
	return err
}

func (r *repo) DeleteWeekMenuItem(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM week_menu_items WHERE id = $1`, id)
	return err
}
