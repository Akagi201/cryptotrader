package model

import (
	"time"
)

// Trades 多个历史成交
type Trades []*Trade

// Trade  单个历史成交
type Trade struct {
	Price     float64   // 交易价格
	Amount    float64   // 交易数量
	Tid       int64     // 交易生成ID
	Type      string    // 交易类型, buy(买)/sell(卖)
	TradeType string    // 委托类型, ask(卖)/bid(买)
	Date      time.Time // 交易时间
}
