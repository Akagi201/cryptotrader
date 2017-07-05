// Package chbtc CHBTC rest api package
package chbtc

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	MarketAPI = "http://api.chbtc.com/data/v1/"
	TradeAPI  = "https://trade.chbtc.com/api/"
)

type CHBTC struct {
	AccessKey string
	SecretKey string
}

func New(accessKey string, secretKey string) *CHBTC {
	return &CHBTC{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (cb *CHBTC) GetTicker(base string, quote string) (*model.Ticker, error) {
	log.Debugf("Currency base: %s, quote: %s", base, quote)

	url := MarketAPI + "ticker?currency=" + quote + "_" + base

	log.Debugf("Request url: %v", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	buyRes := gjson.GetBytes(body, "ticker.buy").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "ticker.sell").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "ticker.last").String()
	last, err := strconv.ParseFloat(lastRes, 64)
	if err != nil {
		return nil, err
	}

	lowRes := gjson.GetBytes(body, "ticker.low").String()
	low, err := strconv.ParseFloat(lowRes, 64)
	if err != nil {
		return nil, err
	}

	highRes := gjson.GetBytes(body, "ticker.high").String()
	high, err := strconv.ParseFloat(highRes, 64)
	if err != nil {
		return nil, err
	}

	volRes := gjson.GetBytes(body, "ticker.vol").String()
	vol, err := strconv.ParseFloat(volRes, 64)
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

// GetOrderBook 市场深度
// size: 档位 1-50, 如果有合并深度, 只能返回 5 档深度
// merge:
// btc_cny: 可选 1, 0.1
// ltc_cny: 可选 0.5, 0.3, 0.1
// eth_cny: 可选 0.5, 0.3, 0.1
// etc_cny: 可选 0.3, 0.1
// bts_cny: 可选 1, 0.1
func (cb *CHBTC) GetOrderBook(base string, quote string, size int, merge float64) (*model.OrderBook, error) {
	url := MarketAPI + "depth?currency=" + quote + "_" + base + "&size=" + strconv.Itoa(size) + "&merge=" + strconv.FormatFloat(merge, 'f', -1, 64)

	log.Debugf("Request url: %v", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	orderBook := &model.OrderBook{
		Base:  base,
		Quote: quote,
		Time:  time.Unix(gjson.GetBytes(body, "timestamp").Int(), 0),
	}

	gjson.GetBytes(body, "asks").ForEach(func(k, v gjson.Result) bool {
		price := v.Array()[0].Float()
		amount := v.Array()[1].Float()

		orderBook.Asks = append(orderBook.Asks, &model.Order{
			Price:  price,
			Amount: amount,
		})

		return true
	})

	gjson.GetBytes(body, "bids").ForEach(func(k, v gjson.Result) bool {
		price := v.Array()[0].Float()
		amount := v.Array()[1].Float()

		orderBook.Bids = append(orderBook.Bids, &model.Order{
			Price:  price,
			Amount: amount,
		})

		return true
	})

	return orderBook, nil
}

// GetTrades 获取历史成交
// currency: quote_base
// btc_cny: 比特币/人民币
// ltc_cny: 莱特币/人民币
// eth_cny: 以太币/人民币
// etc_cny: ETC币/人民币
// bts_cny: BTS币/人民币
// since: 从指定交易 ID 后 50 条数据
func (cb *CHBTC) GetTrades(base string, quote string, since int) (*model.Trades, error) {
	url := MarketAPI + "trades?currency=" + quote + "_" + base
	if since != 0 {
		url += "&since=" + strconv.Itoa(since)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	trades := new(model.Trades)

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		trade := &model.Trade{
			Amount:    v.Get("amount").Float(),
			Price:     v.Get("price").Float(),
			Tid:       v.Get("tid").Int(),
			TradeType: v.Get("trade_type").String(),
			Type:      v.Get("type").String(),
			Date:      time.Unix(v.Get("date").Int(), 0),
		}
		*trades = append(*trades, trade)
		return true
	})

	return trades, nil
}
