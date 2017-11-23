// Package cex cex rest api package
package cex

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://cex.com/Jsons"
)

// Cex API data
type Cex struct {
	AccessKey string
	SecretKey string
}

// New create new Cex API data
func New(accessKey string, secretKey string) *Cex {
	return &Cex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (bt *Cex) GetTicker(base string, quote string) (ticker *model.Ticker, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			ticker = nil
			rerr = err.(error)
		}
	}()

	var url string
	if quote == "pay" {
		url = API + "/trade_43.js?v=" + strconv.FormatFloat(rand.Float64(), 'g', 1, 64)
	} else if quote == "omg" {
		url = API + "/trade_41.js?v=" + strconv.FormatFloat(rand.Float64(), 'g', 1, 64)
	}

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

	buy := gjson.GetBytes(body, "depth.1.#.price").Array()[0].Float()
	sell := gjson.GetBytes(body, "depth.2.#.price").Array()[0].Float()
	last := gjson.GetBytes(body, "cmark.new_price").Float()
	low := gjson.GetBytes(body, "cmark.min_price").Float()
	high := gjson.GetBytes(body, "cmark.max_price").Float()
	vol := gjson.GetBytes(body, "cmark.H24_done_num").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
