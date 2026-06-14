package entity

import "time"

type ExtraProduct struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Category  string    `db:"category"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
