package models

type Dish struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MenuType    string `json:"menu_type"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
}
