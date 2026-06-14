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

func (r *repo) SaveWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, menuTypeID, dishID string) (*entity.WeekMenuItem, error) {
	var item entity.WeekMenuItem
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO week_menu_items (week_menu_id, menu_date, menu_type_id, dish_id)
		 VALUES ($1, $2, $3, $4) RETURNING *`,
		weekMenuID, menuDate, menuTypeID, dishID,
	).StructScan(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repo) GetWeekMenuItems(ctx context.Context, weekMenuID string) ([]entity.WeekMenuItem, error) {
	var items []entity.WeekMenuItem
	err := r.db.SelectContext(ctx, &items,
		`SELECT wmi.* FROM week_menu_items wmi
		 JOIN menu_types mt ON mt.id = wmi.menu_type_id
		 WHERE wmi.week_menu_id = $1
		 ORDER BY wmi.menu_date, mt.sort_order`,
		weekMenuID,
	)
	return items, err
}

func (r *repo) GetWeekMenuItemsByDate(ctx context.Context, date time.Time) ([]entity.WeekMenuItem, error) {
	var items []entity.WeekMenuItem
	err := r.db.SelectContext(ctx, &items,
		`SELECT wmi.* FROM week_menu_items wmi
		 JOIN menu_types mt ON mt.id = wmi.menu_type_id
		 WHERE wmi.menu_date = $1
		 ORDER BY mt.sort_order`,
		date,
	)
	return items, err
}

func (r *repo) UpdateWeekMenuItem(ctx context.Context, id, dishID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE week_menu_items SET dish_id=$1, updated_at=NOW() WHERE id=$2`,
		dishID, id,
	)
	return err
}

func (r *repo) DeleteWeekMenuItem(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM week_menu_items WHERE id = $1`, id)
	return err
}
