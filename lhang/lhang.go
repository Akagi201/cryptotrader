// Package lhang lhang rest api package
package lhang

import (
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://api.lhang.com/v1"
)

// LHang API data
type LHang struct {
	AccessKey string
	SecretKey string
}

// New create new LHang API data
func New(accessKey string, secretKey string) *LHang {
	return &LHang{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (lh *LHang) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "/ticker.do" + "?symbol=" + quote + "_" + base

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

	// TODO: use market depth API.
	buyRes := gjson.GetBytes(body, "ticker.latest").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "ticker.latest").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "ticker.latest").String()
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
