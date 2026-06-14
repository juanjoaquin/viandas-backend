package dtos

type CreateDelivery struct {
	Name string `json:"name" validate:"required,min=2"`
}

type UpdateDelivery struct {
	Name   string `json:"name"   validate:"required,min=2"`
	Active bool   `json:"active"`
}
