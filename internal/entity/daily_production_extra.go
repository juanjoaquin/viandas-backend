package entity

import "time"

type DailyProductionExtra struct {
	ID                string    `db:"id"`
	DailyProductionID string    `db:"daily_production_id"`
	ExtraProductID    string    `db:"extra_product_id"`
	Quantity          int       `db:"quantity"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
