// Package liqui liqui rest api package
package liqui

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://api.liqui.io/api/3/"
)

// Liqui API data
type Liqui struct {
	AccessKey string
	SecretKey string
}

// New create new Liqui API data
func New(accessKey string, secretKey string) *Liqui {
	return &Liqui{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (lq *Liqui) GetTicker(base string, quote string) (*model.Ticker, error) {
	pair := strings.ToLower(quote) + "_" + strings.ToLower(base)
	url := API + "ticker/" + pair

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

	buy := gjson.GetBytes(body, pair+".buy").Float()
	sell := gjson.GetBytes(body, pair+".sell").Float()
	last := gjson.GetBytes(body, pair+".last").Float()
	low := gjson.GetBytes(body, pair+".low").Float()
	high := gjson.GetBytes(body, pair+".high").Float()
	vol := gjson.GetBytes(body, pair+".vol").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
