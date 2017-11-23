// Package huobi huobi rest api package
package huobi

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://api.huobi.pro/"
)

// Huobi API data
type Huobi struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Huobi {
	return &Huobi{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (hb *Huobi) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "market/detail/merged?symbol=" + strings.ToLower(quote) + strings.ToLower(base)

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

	buy := gjson.GetBytes(body, "tick.bid").Array()[0].Float()
	sell := gjson.GetBytes(body, "tick.ask").Array()[0].Float()
	last := gjson.GetBytes(body, "tick.close").Float()
	low := gjson.GetBytes(body, "tick.low").Float()
	high := gjson.GetBytes(body, "tick.high").Float()
	vol := gjson.GetBytes(body, "tick.vol").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
