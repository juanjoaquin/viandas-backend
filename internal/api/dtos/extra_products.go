package dtos

type CreateExtraProduct struct {
	Name     string `json:"name"     validate:"required,min=2"`
	Category string `json:"category" validate:"required,oneof=SALAD SANDWICH"`
}

type UpdateExtraProduct struct {
	Name   string `json:"name"   validate:"required,min=2"`
	Active bool   `json:"active"`
}
