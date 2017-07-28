// Package yunbi yunbi rest api package
package yunbi

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://yunbi.com/api/v2/"
)

// Yunbi API data
type Yunbi struct {
	AccessKey string
	SecretKey string
}

// New create new Yunbi API data
func New(accessKey string, secretKey string) *Yunbi {
	return &Yunbi{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (yb *Yunbi) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "tickers/" + strings.ToLower(quote) + strings.ToLower(base)

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

// GetTickerList 获取当前支持币种列表
func (yb *Yunbi) GetTickerList() ([]string, error) {
	url := API + "tickers.json"

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

	var tickerList []string
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		tickerList = append(tickerList, key.String())
		return true
	})

	return tickerList, nil
}
