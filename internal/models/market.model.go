package models

type StockMarket struct {
	Market struct {
		Symbol        string `json:"01. symbol"`
		Price         string `json:"05. price"`
		Change        string `json:"09. change"`
		ChangePercent string `json:"10. change percent"`
	} `json:"Global Quote"`
}

type StockMarketResponse struct {
	Symbol        string `json:"symbol"`
	Price         string `json:"price"`
	Change        string `json:"change"`
	ChangePercent string `json:"change_percent"`
}
