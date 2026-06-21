package entity

type OverviewMenuTotal struct {
	MenuTypeID    string   `db:"menu_type_id"`
	MenuTypeName  string   `db:"menu_type_name"`
	MenuTypePrice *float64 `db:"menu_type_price"`
	TotalQty      int      `db:"total_qty"`
	TotalAmount   *float64 `db:"total_amount"`
}

type OverviewProductTotal struct {
	ExtraProductID    string  `db:"extra_product_id"`
	ExtraProductName  string  `db:"extra_product_name"`
	ExtraProductPrice float64 `db:"extra_product_price"`
	CategoryID        string  `db:"category_id"`
	CategoryName      string  `db:"category_name"`
	TotalQty          int     `db:"total_qty"`
	TotalAmount       float64 `db:"total_amount"`
}
