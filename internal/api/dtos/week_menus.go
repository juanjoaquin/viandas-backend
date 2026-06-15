package dtos

type CreateWeekMenu struct {
	WeekStartDate string `json:"week_start_date" validate:"required"`
	WeekEndDate   string `json:"week_end_date"   validate:"required"`
}

type AddWeekMenuItem struct {
	MenuDate    string `json:"menu_date"    validate:"required"`
	MenuTypeID  string `json:"menu_type_id" validate:"required"`
	DishID      string `json:"dish_id"      validate:"required"`
}

type UpdateWeekMenuItem struct {
	DishID string `json:"dish_id" validate:"required"`
}
