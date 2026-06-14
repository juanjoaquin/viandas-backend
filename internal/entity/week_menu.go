package entity

import "time"

type WeekMenu struct {
	ID            string    `db:"id"`
	WeekStartDate time.Time `db:"week_start_date"`
	CreatedBy     string    `db:"created_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
