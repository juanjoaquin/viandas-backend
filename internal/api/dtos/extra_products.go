package dtos

type CreateExtraProduct struct {
	Name       string  `json:"name"        validate:"required,min=2"`
	CategoryID string  `json:"category_id" validate:"required,uuid"`
	Price      float64 `json:"price"       validate:"required,gt=0"`
}

type UpdateExtraProduct struct {
	ID         string  `json:"id"          validate:"required,uuid"`
	Name       string  `json:"name"        validate:"required,min=2"`
	CategoryID string  `json:"category_id" validate:"required,uuid"`
	Price      float64 `json:"price"       validate:"required,gt=0"`
	Active     bool    `json:"active"`
}

type DeleteExtraProduct struct {
	ID string `json:"id" validate:"required,uuid"`
}
