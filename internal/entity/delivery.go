package entity

import "time"

type Delivery struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Phone     *string   `db:"phone"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
