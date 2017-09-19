// Package coincheck coincheck rest api
package coincheck

import (
	coincheck "github.com/Akagi201/coincheckgo"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/tidwall/gjson"
)

// Coincheck API data
type Coincheck struct {
	coincheck.CoinCheck
}

// New create new Allcoin API data
func New(accessKey string, secretKey string) *Coincheck {
	client := new(coincheck.CoinCheck).NewClient(accessKey, secretKey)

	return &Coincheck{
		client,
	}
}

// GetTicker 行情
func (cc *Coincheck) GetTicker(base string, quote string) (*model.Ticker, error) {

	resp := cc.Ticker.All()

	buy := gjson.Get(resp, "bid").Float()
	sell := gjson.Get(resp, "ask").Float()
	last := gjson.Get(resp, "last").Float()
	low := gjson.Get(resp, "low").Float()
	high := gjson.Get(resp, "high").Float()
	vol := gjson.Get(resp, "volume").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
