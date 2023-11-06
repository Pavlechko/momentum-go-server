package models

type NBU struct {
	Rate   float64 `json:"rate"`
	Symbol string  `json:"cc"`
}

type ExchangeResponse struct {
	Rate   float64 `json:"rate"`
	Symbol string  `json:"symbol"`
}

type ExchangeFrontendResponse struct {
	Change  float64 `json:"change"`
	EndRate float64 `json:"end_rate"`
}

type ExchangeRatesResponse struct {
	NBU map[string]ExchangeFrontendResponse
}
