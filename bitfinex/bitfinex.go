package bitfinex

import (
	"strconv"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
)

// Bitfinex API data
type Bitfinex struct {
	bitfinex.Client
}

// New create new Allcoin API data
func New(accessKey string, secretKey string) *Bitfinex {
	var client *bitfinex.Client
	if accessKey != "" && secretKey != "" {
		client = bitfinex.NewClient().Auth(accessKey, secretKey)
	} else {
		client = bitfinex.NewClient()
	}

	return &Bitfinex{
		*client,
	}
}

// GetTicker 行情
func (bf *Bitfinex) GetTicker(base string, quote string) (*model.Ticker, error) {
	tick, err := bf.Ticker.Get(strings.ToUpper(quote) + strings.ToUpper(base))

	buy, err := strconv.ParseFloat(tick.Bid, 64)
	if err != nil {
		return nil, err
	}

	sell, err := strconv.ParseFloat(tick.Ask, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(tick.LastPrice, 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(tick.Low, 64)
	if err != nil {
		return nil, err
	}

	high, err := strconv.ParseFloat(tick.High, 64)
	if err != nil {
		return nil, err
	}

	vol, err := strconv.ParseFloat(tick.Volume, 64)
	if err != nil {
		return nil, err
	}

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
