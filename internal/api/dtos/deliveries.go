package dtos

type CreateDelivery struct {
	Name  string  `json:"name" validate:"required,min=2"`
	Phone *string `json:"phone"`
}

type UpdateDelivery struct {
	Name   string  `json:"name"   validate:"required,min=2"`
	Phone  *string `json:"phone"`
	Active bool    `json:"active"`
}

type DeleteDelivery struct {
	ID string `json:"id" validate:"required,uuid"`
}
