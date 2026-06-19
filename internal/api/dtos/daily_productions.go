package dtos

type CreateDailyProductionLine struct {
	MenuTypeID string `json:"menu_type_id" validate:"required"`
	Quantity   int    `json:"quantity"     validate:"min=0"`
}

type CreateDailyProduction struct {
	ProductionDate  string                      `json:"production_date"  validate:"required"`
	CustomerID      string                      `json:"customer_id"      validate:"required"`
	FulfillmentType string                      `json:"fulfillment_type"`
	DeliveryID      string                      `json:"delivery_id"`
	Notes           string                      `json:"notes"`
	Lines           []CreateDailyProductionLine `json:"lines"`
}

type UpdateDailyProduction struct {
	ID              string  `json:"id" validate:"required"`
	FulfillmentType *string `json:"fulfillment_type"`
	DeliveryID      *string `json:"delivery_id"`
	Notes           *string `json:"notes"`
}

type DeleteDailyProduction struct {
	ID string `json:"id" validate:"required"`
}

type UpsertDailyProductionLine struct {
	MenuTypeID string `json:"menu_type_id" validate:"required"`
	Quantity   int    `json:"quantity"     validate:"min=0"`
}

type AddDailyProductionExtra struct {
	ExtraProductID string `json:"extra_product_id" validate:"required"`
	Quantity       int    `json:"quantity"         validate:"required,min=1"`
}
