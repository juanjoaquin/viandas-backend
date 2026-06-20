package entity

import "time"

type ExtraProduct struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	CategoryID string    `db:"category_id"`
	Price      float64   `db:"price"`
	Active     bool      `db:"active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
