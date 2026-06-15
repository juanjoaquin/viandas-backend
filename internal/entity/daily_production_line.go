package entity

import "time"

type DailyProductionLine struct {
	ID                string    `db:"id"`
	DailyProductionID string    `db:"daily_production_id"`
	MenuTypeID        string    `db:"menu_type_id"`
	Quantity          int       `db:"quantity"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type ProductionLineInput struct {
	MenuTypeID string
	Quantity   int
}
