package models

type DailyProduction struct {
	ID             string                  `json:"id"`
	ProductionDate string                  `json:"production_date"`
	Customer       *Customer               `json:"customer,omitempty"`
	Delivery       *Delivery               `json:"delivery,omitempty"`
	TraditionalQty int                     `json:"traditional_qty"`
	HealthyQty     int                     `json:"healthy_qty"`
	VegetarianQty  int                     `json:"vegetarian_qty"`
	Notes          string                  `json:"notes,omitempty"`
	Extras         []DailyProductionExtra  `json:"extras,omitempty"`
	CreatedBy      string                  `json:"created_by"`
	CreatedAt      string                  `json:"created_at"`
}

type DailyProductionExtra struct {
	ID           string        `json:"id"`
	ExtraProduct *ExtraProduct `json:"extra_product,omitempty"`
	Quantity     int           `json:"quantity"`
}

type KitchenTotals struct {
	Date           string `json:"date"`
	TraditionalQty int    `json:"traditional_qty"`
	HealthyQty     int    `json:"healthy_qty"`
	VegetarianQty  int    `json:"vegetarian_qty"`
}

type ExtraTotals struct {
	Date   string       `json:"date"`
	Totals []ExtraTotal `json:"totals"`
}

type ExtraTotal struct {
	ExtraProduct *ExtraProduct `json:"extra_product"`
	TotalQty     int           `json:"total_qty"`
}
