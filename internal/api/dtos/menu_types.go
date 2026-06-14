package dtos

type CreateMenuType struct {
	Name      string `json:"name"       validate:"required,min=2"`
	SortOrder int    `json:"sort_order"`
}

type UpdateMenuType struct {
	Name      string `json:"name"       validate:"required,min=2"`
	SortOrder int    `json:"sort_order"`
	Active    bool   `json:"active"`
}
