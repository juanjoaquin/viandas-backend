package dtos

type CreateDish struct {
	Name        string `json:"name"         validate:"required,min=2"`
	Description string `json:"description"`
	MenuTypeID  string `json:"menu_type_id" validate:"required"`
}

type UpdateDish struct {
	ID          string `json:"id"           validate:"required,uuid"`
	Name        string `json:"name"         validate:"required,min=2"`
	Description string `json:"description"`
	MenuTypeID  string `json:"menu_type_id" validate:"required"`
	Active      bool   `json:"active"`
}

type DeleteDish struct {
	ID string `json:"id" validate:"required,uuid"`
}
