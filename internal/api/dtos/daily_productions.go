package dtos

type CreateDailyProduction struct {
	ProductionDate string `json:"production_date" validate:"required"`
	CustomerID     string `json:"customer_id"     validate:"required"`
	DeliveryID     string `json:"delivery_id"`
	TraditionalQty int    `json:"traditional_qty" validate:"min=0"`
	HealthyQty     int    `json:"healthy_qty"     validate:"min=0"`
	VegetarianQty  int    `json:"vegetarian_qty"  validate:"min=0"`
	Notes          string `json:"notes"`
}

type UpdateDailyProduction struct {
	DeliveryID     string `json:"delivery_id"`
	TraditionalQty int    `json:"traditional_qty" validate:"min=0"`
	HealthyQty     int    `json:"healthy_qty"     validate:"min=0"`
	VegetarianQty  int    `json:"vegetarian_qty"  validate:"min=0"`
	Notes          string `json:"notes"`
}

type AddDailyProductionExtra struct {
	ExtraProductID string `json:"extra_product_id" validate:"required"`
	Quantity       int    `json:"quantity"         validate:"required,min=1"`
}
