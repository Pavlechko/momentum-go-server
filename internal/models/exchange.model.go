package models

type NBU struct {
	Currency string  `json:"txt"`
	Rate     float64 `json:"rate"`
	Symbol   string  `json:"cc"`
}

type ExchangeResponse struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
	Symbol   string  `json:"symbol"`
}
