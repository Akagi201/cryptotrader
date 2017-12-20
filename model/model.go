package model

import "time"

// Ticker 行情数据
type Ticker struct {
	Buy  float64 // 买一价
	Sell float64 // 卖一价
	Last float64 // 最新成交价
	Low  float64 // 最低价
	High float64 // 最高价
	Vol  float64 // 成交量(最近 24 小时)
	Time time.Time
	Raw  string // exchange original info
}

type SimpleTicker struct {
	Price  float64
	Symbol string
}

type BookTicker struct {
	Symbol    string
	BidPrice  float64
	BidAmount float64
	AskPrice  float64
	AskAmount float64
}

// Trades 多个历史成交
type Trades []*Trade

// 历史成交
type Trade struct {
	ID     int64 // trade id
	Type   string
	Price  float64
	Amount float64
	Time   time.Time
	Raw    string // exchange original info
}

// Kline OHLC struct
type Record struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
	Vol   float64
	Time  time.Time
	Raw   string
}

type MarketOrder struct {
	Price  float64
	Amount float64
}

type OrderBook struct {
	Asks []MarketOrder // 卖方深度
	Bids []MarketOrder // 买方深度
	Time time.Time
	Raw  string // exchange original info
}

type Order struct {
	ID         int64
	Amount     float64
	DealAmount float64
	Price      float64
	Status     string
	Type       string
	Side       string
	Raw        string
}

type Balance struct {
	Currency string
	Free     float64
	Frozen   float64
}

type MarketInfo struct {
	Symbol        string
	DecimalPlaces int64
	MinAmount     float64
	Fee           float64
}
