package models

type MenuType struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Price     *float64 `json:"price"`
	Active    bool     `json:"active"`
	CreatedAt string   `json:"created_at"`
}
