package models

type Dish struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MenuType    *MenuType `json:"menu_type,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   string    `json:"created_at"`
}
