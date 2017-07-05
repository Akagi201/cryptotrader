package model

import (
	"time"
)

// Order 深度
type Order struct {
	Price  float64
	Amount float64
}

// OrderBook 市场深度
type OrderBook struct {
	Base  string
	Quote string
	Asks  []*Order  // 卖方深度
	Bids  []*Order  // 买方深度
	Time  time.Time // 此次深度的产生时间戳
}
