package models

type ProductCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}
