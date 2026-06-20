package models

type DailyProduction struct {
	ID              string                 `json:"id"`
	ProductionDate  string                 `json:"production_date"`
	FulfillmentType string                 `json:"fulfillment_type"`
	Customer        *Customer              `json:"customer,omitempty"`
	Delivery        *Delivery              `json:"delivery,omitempty"`
	Lines           []DailyProductionLine  `json:"lines,omitempty"`
	Notes           *string                `json:"notes"`
	Extras          []DailyProductionExtra `json:"extras,omitempty"`
	CreatedBy       string                 `json:"created_by"`
	CreatedAt       string                 `json:"created_at"`
}

type DailyProductionLine struct {
	ID          string    `json:"id"`
	MenuType    *MenuType `json:"menu_type,omitempty"`
	Quantity    int       `json:"quantity"`
}

type DailyProductionExtra struct {
	ID           string        `json:"id"`
	ExtraProduct *ExtraProduct `json:"extra_product,omitempty"`
	Quantity     int           `json:"quantity"`
	TotalAmount  *float64      `json:"total_amount,omitempty"`
}

type KitchenTotals struct {
	Date       string             `json:"date"`
	Totals     []MenuTypeTotalQty `json:"totals"`
	GrandTotal *float64           `json:"grand_total"`
}

type MenuTypeTotalQty struct {
	MenuType    *MenuType `json:"menu_type"`
	TotalQty    int       `json:"total_qty"`
	TotalAmount *float64  `json:"total_amount"`
}

type ExtraTotals struct {
	Date       string       `json:"date"`
	Totals     []ExtraTotal `json:"totals"`
	GrandTotal *float64     `json:"grand_total"`
}

type ExtraTotal struct {
	ExtraProduct *ExtraProduct `json:"extra_product"`
	TotalQty     int           `json:"total_qty"`
	TotalAmount  *float64      `json:"total_amount"`
}
