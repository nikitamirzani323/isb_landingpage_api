package entities

type Model_uom struct {
	Uom_id     int    `json:"uom_id"`
	Uom_name   string `json:"uom_name"`
	Uom_status string `json:"uom_status"`
	Uom_create string `json:"uom_create"`
	Uom_update string `json:"uom_update"`
}
type Controller_uom struct {
	Uom_search string `json:"uom_search"`
	Uom_page   int    `json:"uom_page"`
}
type Controller_uomsave struct {
	Page       string `json:"page" validate:"required"`
	Sdata      string `json:"sdata" validate:"required"`
	Uom_search string `json:"uom_search"`
	Uom_page   int    `json:"uom_page"`
	Uom_id     int    `json:"uom_id"`
	Uom_name   string `json:"uom_name" validate:"required"`
	Uom_status string `json:"uom_status" validate:"required"`
}
