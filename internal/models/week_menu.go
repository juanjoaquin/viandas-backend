package models

type WeekMenu struct {
	ID            string         `json:"id"`
	WeekStartDate string         `json:"week_start_date"`
	WeekEndDate   string         `json:"week_end_date"`
	CreatedBy     string         `json:"created_by"`
	Items         []WeekMenuItem `json:"items,omitempty"`
	CreatedAt     string         `json:"created_at"`
}

type WeekMenuItem struct {
	ID         string    `json:"id"`
	WeekMenuID string    `json:"week_menu_id"`
	MenuDate   string    `json:"menu_date"`
	MenuType   *MenuType `json:"menu_type,omitempty"`
	Dish       *Dish     `json:"dish,omitempty"`
}

type DayMenu struct {
	Date  string        `json:"date"`
	Lines []DayMenuLine `json:"lines"`
}

type DayMenuLine struct {
	MenuType *MenuType `json:"menu_type"`
	Dish     *Dish     `json:"dish"`
}
