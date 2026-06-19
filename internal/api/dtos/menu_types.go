package dtos

type CreateMenuType struct {
	Name  string   `json:"name"  validate:"required,min=2"`
	Price *float64 `json:"price"`
}

type UpdateMenuType struct {
	ID     string   `json:"id"     validate:"required,uuid"`
	Name   string   `json:"name"   validate:"required,min=2"`
	Price  *float64 `json:"price"`
	Active bool     `json:"active"`
}

type DeleteMenuType struct {
	ID string `json:"id" validate:"required,uuid"`
}
