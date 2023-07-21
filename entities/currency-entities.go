package entities

type Model_currency struct {
	Currency_id     string `json:"currency_id"`
	Currency_name   string `json:"currency_name"`
	Currency_create string `json:"currency_create"`
	Currency_update string `json:"currency_update"`
}
type Controller_currency struct {
	Currency_search string `json:"currency_search"`
	Currency_page   int    `json:"currency_page"`
}
type Controller_currencysave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Currency_search string `json:"currency_search"`
	Currency_page   int    `json:"currency_page"`
	Currency_id     string `json:"currency_id"`
	Currency_name   string `json:"currency_name" validate:"required"`
}
