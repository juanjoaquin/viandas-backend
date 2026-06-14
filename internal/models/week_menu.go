package models

type WeekMenu struct {
	ID            string         `json:"id"`
	WeekStartDate string         `json:"week_start_date"`
	CreatedBy     string         `json:"created_by"`
	Items         []WeekMenuItem `json:"items,omitempty"`
	CreatedAt     string         `json:"created_at"`
}

type WeekMenuItem struct {
	ID               string `json:"id"`
	WeekMenuID       string `json:"week_menu_id"`
	MenuDate         string `json:"menu_date"`
	TraditionalDish  *Dish  `json:"traditional_dish,omitempty"`
	HealthyDish      *Dish  `json:"healthy_dish,omitempty"`
	VegetarianDish   *Dish  `json:"vegetarian_dish,omitempty"`
}

type DayMenu struct {
	Date            string `json:"date"`
	TraditionalDish *Dish  `json:"traditional"`
	HealthyDish     *Dish  `json:"healthy"`
	VegetarianDish  *Dish  `json:"vegetarian"`
}
