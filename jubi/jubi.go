// Package jubi jubi rest api package
package jubi

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
	API = "https://www.jubi.com/"
)

// Jubi API data
type Jubi struct {
	AccessKey string
	SecretKey string
}

// New create new Jubi API data
func New(accessKey string, secretKey string) *Jubi {
	return &Jubi{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (jb *Jubi) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "api/v1/ticker" + "?coin=" + strings.ToLower(quote)

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

	buyRes := gjson.GetBytes(body, "buy").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "sell").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "last").String()
	last, err := strconv.ParseFloat(lastRes, 64)
	if err != nil {
		return nil, err
	}

	lowRes := gjson.GetBytes(body, "low").String()
	low, err := strconv.ParseFloat(lowRes, 64)
	if err != nil {
		return nil, err
	}

	highRes := gjson.GetBytes(body, "high").String()
	high, err := strconv.ParseFloat(highRes, 64)
	if err != nil {
		return nil, err
	}

	vol := gjson.GetBytes(body, "vol").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
