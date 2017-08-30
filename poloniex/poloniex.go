// Package poloniex poloniex rest api package
package poloniex

import (
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://poloniex.com/public"
)

// Poloniex API data
type Poloniex struct {
	AccessKey string
	SecretKey string
}

// New create new Poloniex API data
func New(accessKey string, secretKey string) *Poloniex {
	return &Poloniex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (pl *Poloniex) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "?command=returnTicker"

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

	tickers := gjson.ParseBytes(body).Map()
	key := strings.ToUpper(base) + "_" + strings.ToUpper(quote)
	if !tickers[key].Exists() {
		return nil, errors.New("The ticker not exists")
	}

	v := tickers[key]
	return &model.Ticker{
		Buy:  v.Get("highestBid").Float(),
		Sell: v.Get("lowestAsk").Float(),
		Last: v.Get("last").Float(),
		Low:  v.Get("low24hr").Float(),
		High: v.Get("high24hr").Float(),
		Vol:  v.Get("baseVolume").Float(),
	}, nil
}
