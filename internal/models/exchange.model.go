package models

type ExchangeInput struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Source string `json:"source" binding:"required"`
}

type NBU struct {
	Rate   float64 `json:"rate"`
	Symbol string  `json:"cc"`
}

type ExchangeFrontendResponse struct {
	Change  float64 `json:"change"`
	EndRate float64 `json:"end_rate"`
	From    string  `json:"from"`
	To      string  `json:"to"`
	Source  string  `json:"source"`
}

type LayerResponse struct {
	Rates map[string]ExchangeFrontendResponse `json:"rates"`
}
