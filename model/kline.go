package model

import (
	"time"
)

// Kline K 线返回结果
type Kline struct {
	Data      []*KlineData // K 线内容
	MoneyType string       // 买入货币
	Symbol    string       // 卖出货币
}

// KlineData K 线单个内容
type KlineData struct {
	Time   time.Time
	Open   float64 // 开
	High   float64 // 高
	Low    float64 // 低
	Close  float64 // 收
	Amount float64 // 交易量
}
