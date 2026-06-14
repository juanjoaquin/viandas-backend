package entity

import "time"

type DailyProduction struct {
	ID             string    `db:"id"`
	ProductionDate time.Time `db:"production_date"`
	CustomerID     string    `db:"customer_id"`
	DeliveryID     string    `db:"delivery_id"`
	TraditionalQty int       `db:"traditional_qty"`
	HealthyQty     int       `db:"healthy_qty"`
	VegetarianQty  int       `db:"vegetarian_qty"`
	Notes          string    `db:"notes"`
	CreatedBy      string    `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
