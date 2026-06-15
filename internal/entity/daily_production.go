package entity

import "time"

type DailyProduction struct {
	ID              string    `db:"id"`
	ProductionDate  time.Time `db:"production_date"`
	CustomerID      string    `db:"customer_id"`
	FulfillmentType string    `db:"fulfillment_type"`
	DeliveryID      *string   `db:"delivery_id"`
	Notes           *string   `db:"notes"`
	CreatedBy       string    `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
