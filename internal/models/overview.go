package models

type OverviewPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ProductionOverviewSummary struct {
	TotalMenusQty       int     `json:"total_menus_qty"`
	TotalMenusAmount    float64 `json:"total_menus_amount"`
	TotalProductsQty    int     `json:"total_products_qty"`
	TotalProductsAmount float64 `json:"total_products_amount"`
	GrandTotalAmount    float64 `json:"grand_total_amount"`
}

type ProductionOverview struct {
	Period      OverviewPeriod              `json:"period"`
	Summary     ProductionOverviewSummary   `json:"summary"`
	MenusByType []ProductionOverviewMenu    `json:"menus_by_type"`
	Products    []ProductionOverviewProduct `json:"products"`
}

type ProductionOverviewMenu struct {
	MenuType    *MenuType `json:"menu_type"`
	TotalQty    int       `json:"total_qty"`
	TotalAmount *float64  `json:"total_amount"`
}

type ProductionOverviewProduct struct {
	ExtraProduct *ExtraProduct `json:"extra_product"`
	TotalQty     int           `json:"total_qty"`
	TotalAmount  float64       `json:"total_amount"`
}
