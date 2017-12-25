// Package bigone bigone rest api package
package bigone

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

const (
	RestHost = "api.big.one"
)

// Client BigONE client
type Client struct {
	URL        url.URL
	HTTPClient *http.Client
	ApiKey     string
}

// New creates a new BigONE Client
func New(apiKey string) *Client {
	u := url.URL{
		Scheme: "https",
		Host:   RestHost,
		Path:   "/",
	}

	c := Client{
		URL:        u,
		HTTPClient: &http.Client{},
		ApiKey:     apiKey,
	}

	return &c
}

func (c *Client) newRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader) (*http.Request, error) {
	u := c.URL
	u.Path = path.Join(c.URL.Path, spath)
	u.RawQuery = values.Encode()
	log.Debugf("Request URL: %#v", u.String())

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) getResponse(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		log.Errorf("body: %v", string(body))
		return nil, errors.New(fmt.Sprintf("status code: %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetTicker Get a Market, for GET https://api.big.one/markets/{symbol}
func (c *Client) GetTicker(ctx context.Context, quote string, base string) (*model.BigONETicker, error) {
	req, err := c.newRequest(ctx, "GET", "markets/"+strings.ToUpper(quote)+"-"+strings.ToUpper(base), nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	lastRes := gjson.GetBytes(body, "data.ticker.price").String()
	last, err := strconv.ParseFloat(lastRes, 64)
	if err != nil {
		return nil, err
	}

	openRes := gjson.GetBytes(body, "data.ticker.open").String()
	open, err := strconv.ParseFloat(openRes, 64)
	if err != nil {
		return nil, err
	}

	closeRes := gjson.GetBytes(body, "data.ticker.close").String()
	clos, err := strconv.ParseFloat(closeRes, 64)
	if err != nil {
		return nil, err
	}

	highRes := gjson.GetBytes(body, "data.ticker.high").String()
	high, err := strconv.ParseFloat(highRes, 64)
	if err != nil {
		return nil, err
	}

	lowRes := gjson.GetBytes(body, "data.ticker.low").String()
	low, err := strconv.ParseFloat(lowRes, 64)
	if err != nil {
		return nil, err
	}

	volRes := gjson.GetBytes(body, "data.ticker.volume").String()
	vol, err := strconv.ParseFloat(volRes, 64)
	if err != nil {
		return nil, err
	}

	return &model.BigONETicker{
		Last:  last,
		Open:  open,
		Close: clos,
		Low:   low,
		High:  high,
		Vol:   vol,
	}, nil
}

// GetDepth Order book, for GET https://api.big.one/markets/{symbol}/book
func (c *Client) GetDepth(ctx context.Context, quote string, base string) (*model.OrderBook, error) {
	req, err := c.newRequest(ctx, "GET", "markets/"+strings.ToUpper(quote)+"-"+strings.ToUpper(base)+"/book", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.MarketOrder
	var orderBook model.OrderBook
	gjson.GetBytes(body, "data.bids").ForEach(func(key, value gjson.Result) bool {
		order.Price = cast.ToFloat64(value.Get("price").String())
		order.Amount = cast.ToFloat64(value.Get("amount").String())
		orderBook.Bids = append(orderBook.Bids, order)
		return true // keep iterating
	})

	gjson.GetBytes(body, "data.asks").ForEach(func(key, value gjson.Result) bool {
		order.Price = cast.ToFloat64(value.Get("price").String())
		order.Amount = cast.ToFloat64(value.Get("amount").String())
		orderBook.Asks = append(orderBook.Asks, order)
		return true // keep iterating
	})

	return &orderBook, nil
}

// GetTrades List Trades in the Market, for GET https://api.big.one/markets/{symbol}/trades
func (c *Client) GetTrades(ctx context.Context, quote string, base string) ([]model.BigONETrade, error) {
	req, err := c.newRequest(ctx, "GET", "markets/"+strings.ToUpper(quote)+"-"+strings.ToUpper(base)+"/trades", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var trade model.BigONETrade
	var trades []model.BigONETrade
	gjson.GetBytes(body, "data").ForEach(func(key, value gjson.Result) bool {
		trade.ID = value.Get("trade_id").String()
		trade.Price = cast.ToFloat64(value.Get("price").String())
		trade.Amount = cast.ToFloat64(value.Get("amount").String())
		trade.Type = value.Get("trade_side").String()
		trade.Time = cast.ToTime(value.Get("created_at").String())
		trades = append(trades, trade)
		return true // keep iterating
	})

	return trades, nil
}
