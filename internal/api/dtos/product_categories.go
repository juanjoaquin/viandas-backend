package dtos

type CreateProductCategory struct {
	Name string `json:"name" validate:"required,min=2"`
}

type UpdateProductCategory struct {
	ID     string `json:"id"     validate:"required,uuid"`
	Name   string `json:"name"   validate:"required,min=2"`
	Active bool   `json:"active"`
}

type DeleteProductCategory struct {
	ID string `json:"id" validate:"required,uuid"`
}
