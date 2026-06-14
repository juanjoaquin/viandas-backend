package dtos

type CreateWeekMenu struct {
	WeekStartDate string `json:"week_start_date" validate:"required"`
}

type AddWeekMenuItem struct {
	MenuDate          string `json:"menu_date"           validate:"required"`
	TraditionalDishID string `json:"traditional_dish_id" validate:"required"`
	HealthyDishID     string `json:"healthy_dish_id"     validate:"required"`
	VegetarianDishID  string `json:"vegetarian_dish_id"  validate:"required"`
}

type UpdateWeekMenuItem struct {
	TraditionalDishID string `json:"traditional_dish_id" validate:"required"`
	HealthyDishID     string `json:"healthy_dish_id"     validate:"required"`
	VegetarianDishID  string `json:"vegetarian_dish_id"  validate:"required"`
}
