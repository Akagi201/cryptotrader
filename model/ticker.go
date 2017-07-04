package model

// Ticker 行情数据
type Ticker struct {
	Buy  float64 // 买一价
	Sell float64 // 卖一价
	Last float64 // 最新成交价
	Low  float64 // 最低价
	High float64 // 最高价
	Vol  float64 // 成交量(最近 24 小时)
}
