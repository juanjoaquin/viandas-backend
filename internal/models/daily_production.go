package models

type DailyProduction struct {
	ID             string                 `json:"id"`
	ProductionDate string                 `json:"production_date"`
	Customer       *Customer              `json:"customer,omitempty"`
	Delivery       *Delivery              `json:"delivery,omitempty"`
	Lines          []DailyProductionLine  `json:"lines,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
	Extras         []DailyProductionExtra `json:"extras,omitempty"`
	CreatedBy      string                 `json:"created_by"`
	CreatedAt      string                 `json:"created_at"`
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
}

type KitchenTotals struct {
	Date   string             `json:"date"`
	Totals []MenuTypeTotalQty `json:"totals"`
}

type MenuTypeTotalQty struct {
	MenuType *MenuType `json:"menu_type"`
	TotalQty int       `json:"total_qty"`
}

type ExtraTotals struct {
	Date   string       `json:"date"`
	Totals []ExtraTotal `json:"totals"`
}

type ExtraTotal struct {
	ExtraProduct *ExtraProduct `json:"extra_product"`
	TotalQty     int           `json:"total_qty"`
}
