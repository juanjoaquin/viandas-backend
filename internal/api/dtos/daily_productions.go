package dtos

type CreateDailyProduction struct {
	ProductionDate string `json:"production_date" validate:"required"`
	CustomerID     string `json:"customer_id"     validate:"required"`
	DeliveryID     string `json:"delivery_id"`
	Notes          string `json:"notes"`
}

type UpdateDailyProduction struct {
	DeliveryID string `json:"delivery_id"`
	Notes      string `json:"notes"`
}

type UpsertDailyProductionLine struct {
	MenuTypeID string `json:"menu_type_id" validate:"required"`
	Quantity   int    `json:"quantity"     validate:"min=0"`
}

type AddDailyProductionExtra struct {
	ExtraProductID string `json:"extra_product_id" validate:"required"`
	Quantity       int    `json:"quantity"         validate:"required,min=1"`
}
