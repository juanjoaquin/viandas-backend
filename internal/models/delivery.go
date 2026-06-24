package models

type Delivery struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Phone     *string `json:"phone"`
	Active    bool    `json:"active"`
	CreatedAt string  `json:"created_at"`
}
