package entity

import "time"

type WeekMenuItem struct {
	ID                 string    `db:"id"`
	WeekMenuID         string    `db:"week_menu_id"`
	MenuDate           time.Time `db:"menu_date"`
	TraditionalDishID  string    `db:"traditional_dish_id"`
	HealthyDishID      string    `db:"healthy_dish_id"`
	VegetarianDishID   string    `db:"vegetarian_dish_id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
