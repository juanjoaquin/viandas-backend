package models

type MenuType struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}
