package models

type ExtraProduct struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}
