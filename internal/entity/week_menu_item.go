package entity

import "time"

type WeekMenuItem struct {
	ID          string    `db:"id"`
	WeekMenuID  string    `db:"week_menu_id"`
	MenuDate    time.Time `db:"menu_date"`
	MenuTypeID  string    `db:"menu_type_id"`
	DishID      string    `db:"dish_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
