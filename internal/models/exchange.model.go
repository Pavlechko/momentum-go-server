package models

type NBU struct {
	Rate   float64 `json:"rate"`
	Symbol string  `json:"cc"`
}

type ExchangeFrontendResponse struct {
	Change  float64 `json:"change"`
	EndRate float64 `json:"end_rate"`
}

type LayerResponse struct {
	Rates map[string]ExchangeFrontendResponse `json:"rates"`
}

type ExchangeRatesResponse struct {
	NBU   map[string]ExchangeFrontendResponse
	Layer map[string]ExchangeFrontendResponse
}
