package dtos

type CreateCustomer struct {
	Name    string `json:"name"    validate:"required,min=2"`
	Type    string `json:"type"    validate:"required,oneof=COMPANY PERSON"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}

type UpdateCustomer struct {
	Name    string `json:"name"    validate:"required,min=2"`
	Type    string `json:"type"    validate:"required,oneof=COMPANY PERSON"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}

type DeleteCustomer struct {
	ID string `json:"id" validate:"required,uuid"`
}
