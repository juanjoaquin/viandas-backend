package dtos

type CreateMenuType struct {
	Name  string   `json:"name"  validate:"required,min=2"`
	Price *float64 `json:"price"`
}

type UpdateMenuType struct {
	Name   string   `json:"name"   validate:"required,min=2"`
	Price  *float64 `json:"price"`
	Active bool     `json:"active"`
}
