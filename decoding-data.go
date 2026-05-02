package main

// Simulating fyersws.DataResponse (since it's just JSON string)
type DataResponse string

// 🔹 Base message (for routing)
type BaseMessage struct {
	Type string `json:"type"`
}

// 🔹 Equity / Option Data (sf)
type MarketData struct {
	Symbol           string  `json:"symbol"`
	LTP              float64 `json:"ltp"`
	PrevClosePrice   float64 `json:"prev_close_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	OpenPrice        float64 `json:"open_price"`
	Change           float64 `json:"ch"`
	ChangePercent    float64 `json:"chp"`
	VolumeTraded     int64   `json:"vol_traded_today"`
	LastTradedTime   int64   `json:"last_traded_time"`
	ExchangeFeedTime int64   `json:"exch_feed_time"`
	BidSize          int64   `json:"bid_size"`
	AskSize          int64   `json:"ask_size"`
	BidPrice         float64 `json:"bid_price"`
	AskPrice         float64 `json:"ask_price"`
	LastTradedQty    int64   `json:"last_traded_qty"`
	TotalBuyQty      int64   `json:"tot_buy_qty"`
	TotalSellQty     int64   `json:"tot_sell_qty"`
	AvgTradePrice    float64 `json:"avg_trade_price"`
	LowerCircuit     float64 `json:"lower_ckt"`
	UpperCircuit     float64 `json:"upper_ckt"`
	Type             string  `json:"type"`
}

// 🔹 Index Data (if)
type IndexData struct {
	Symbol           string  `json:"symbol"`
	LTP              float64 `json:"ltp"`
	PrevClosePrice   float64 `json:"prev_close_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	OpenPrice        float64 `json:"open_price"`
	Change           float64 `json:"ch"`
	ChangePercent    float64 `json:"chp"`
	ExchangeFeedTime int64   `json:"exch_feed_time"`
	Type             string  `json:"type"`
}

// Depth Data (df)
type DepthData struct {
	Symbol string `json:"symbol"`
	Type   string `json:"type"`

	// Bid Prices
	BidPrice1 float64 `json:"bid_price1"`
	BidPrice2 float64 `json:"bid_price2"`
	BidPrice3 float64 `json:"bid_price3"`
	BidPrice4 float64 `json:"bid_price4"`
	BidPrice5 float64 `json:"bid_price5"`

	// Ask Prices
	AskPrice1 float64 `json:"ask_price1"`
	AskPrice2 float64 `json:"ask_price2"`
	AskPrice3 float64 `json:"ask_price3"`
	AskPrice4 float64 `json:"ask_price4"`
	AskPrice5 float64 `json:"ask_price5"`

	// Bid Sizes
	BidSize1 int64 `json:"bid_size1"`
	BidSize2 int64 `json:"bid_size2"`
	BidSize3 int64 `json:"bid_size3"`
	BidSize4 int64 `json:"bid_size4"`
	BidSize5 int64 `json:"bid_size5"`

	// Ask Sizes
	AskSize1 int64 `json:"ask_size1"`
	AskSize2 int64 `json:"ask_size2"`
	AskSize3 int64 `json:"ask_size3"`
	AskSize4 int64 `json:"ask_size4"`
	AskSize5 int64 `json:"ask_size5"`

	// Order Count
	BidOrder1 int64 `json:"bid_order1"`
	BidOrder2 int64 `json:"bid_order2"`
	BidOrder3 int64 `json:"bid_order3"`
	BidOrder4 int64 `json:"bid_order4"`
	BidOrder5 int64 `json:"bid_order5"`

	AskOrder1 int64 `json:"ask_order1"`
	AskOrder2 int64 `json:"ask_order2"`
	AskOrder3 int64 `json:"ask_order3"`
	AskOrder4 int64 `json:"ask_order4"`
	AskOrder5 int64 `json:"ask_order5"`
}
