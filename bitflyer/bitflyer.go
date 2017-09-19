// Package bitflyer bitflyer rest api
package bitflyer

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	JPAPI = "https://api.bitflyer.jp/v1/"
	USAPI = "https://api.bitflyer.com/v1/"
)

// Bitflyer API data
type Bitflyer struct {
	AccessKey string
	SecretKey string
}

// New create new Bitflyer API data
func New(accessKey string, secretKey string) *Bitflyer {
	return &Bitflyer{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (bf *Bitflyer) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := JPAPI + "ticker?product_code=" + strings.ToUpper(quote) + "_" + strings.ToUpper(base)

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

	buy := gjson.GetBytes(body, "best_bid").Float()
	sell := gjson.GetBytes(body, "best_ask").Float()
	last := gjson.GetBytes(body, "ltp").Float()
	low := gjson.GetBytes(body, "best_bid").Float()
	high := gjson.GetBytes(body, "best_ask").Float()
	vol := gjson.GetBytes(body, "volume").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
