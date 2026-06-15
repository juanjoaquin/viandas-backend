package entity

import "time"

type MenuType struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Price     *float64  `db:"price"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
