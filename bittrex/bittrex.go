// Package bittrex bittrex rest api package
package bittrex

import (
	"io/ioutil"
	"net/http"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://bittrex.com/api/v1.1"
)

// Bittrex API data
type Bittrex struct {
	AccessKey string
	SecretKey string
}

// New create new Yunbi API data
func New(accessKey string, secretKey string) *Bittrex {
	return &Bittrex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (yb *Bittrex) GetTicker(base string, quote string) (ticker *model.Ticker, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			ticker = nil
			rerr = err.(error)
		}
	}()

	url := API + "/public/getmarketsummary?market=" + base + "-" + quote

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

	buy := gjson.GetBytes(body, "result.#.Ask").Array()[0].Float()
	sell := gjson.GetBytes(body, "result.#.Bid").Array()[0].Float()
	last := gjson.GetBytes(body, "result.#.Last").Array()[0].Float()
	low := gjson.GetBytes(body, "result.#.Low").Array()[0].Float()
	high := gjson.GetBytes(body, "result.#.High").Array()[0].Float()
	vol := gjson.GetBytes(body, "result.#.BaseVolume").Array()[0].Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
