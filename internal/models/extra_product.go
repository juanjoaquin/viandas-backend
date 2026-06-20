package models

type ExtraProduct struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Category  *ProductCategory `json:"category,omitempty"`
	Price     float64          `json:"price"`
	Active    bool             `json:"active"`
	CreatedAt string           `json:"created_at"`
}
