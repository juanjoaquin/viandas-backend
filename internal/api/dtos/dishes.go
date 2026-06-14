package dtos

type CreateDish struct {
	Name        string `json:"name"        validate:"required,min=2"`
	Description string `json:"description"`
	MenuType    string `json:"menu_type"   validate:"required,oneof=TRADITIONAL HEALTHY VEGETARIAN"`
}

type UpdateDish struct {
	Name        string `json:"name"        validate:"required,min=2"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}
