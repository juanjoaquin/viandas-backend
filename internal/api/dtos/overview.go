package dtos

type ProductionOverviewFilters struct {
	From string `query:"from"`
	To   string `query:"to"`
}
