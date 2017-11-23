package model

import (
	"time"
)

// ZBOrder 获取委托买单或卖单
type ZBOrder struct {
	Currency    string    // 交易类型（目前仅支持btc_cny/ltc_cny/eth_cny/eth_btc/etc_cny）
	Fees        float64   // 交易手续费,卖单的话,显示的是收入货币(如人民币);买单的话,显示的是买入货币(如etc)
	ID          string    // 委托挂单号
	Price       float64   // 单价
	Status      int64     // 挂单状态(0: 待成交, 1: 取消, 2: 交易完成, 3: 待成交未交易部份)
	TotalAmount float64   // 挂单总数量
	TradeAmount float64   // 已成交数量
	TradePrice  float64   // 成交均价
	TradeDate   time.Time // 委托时间
	TradeMoney  float64   // 已成交总金额
	Type        int64     // 挂单类型 1/0[buy/sell]
}
