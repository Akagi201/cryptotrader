// Package bigone bigone rest api package
package bigone

import (
	"bytes"
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
	"github.com/golang-plus/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

func (c *Client) newPrivateRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader) (*http.Request, error) {
	req, err := c.newRequest(ctx, method, spath, values, body)
	if err != nil {
		return nil, err
	}

	if values == nil {
		values = url.Values{}
	}

	u, _ := uuid.NewTimeBased()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36")
	req.Header.Set("Big-Device-Id", u.String())
	req.Header.Set("Content-Type", "application/json")

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

// Trade Create an Order, side: BID or ASK, for POST https://api.big.one/orders
func (c *Client) Trade(ctx context.Context, quote string, base string, side string, quantity float64, price float64) (string, error) {
	reqBody := `{
		"order_market": "ETH-BTC",
		"order_side": "BID",
		"price": "0.04",
		"amount": "0.00001"
	}`

	reqBody, _ = sjson.Set(reqBody, "order_market", strings.ToUpper(quote)+"-"+strings.ToUpper(base))
	reqBody, _ = sjson.Set(reqBody, "order_side", side)
	reqBody, _ = sjson.Set(reqBody, "price", cast.ToString(price))
	reqBody, _ = sjson.Set(reqBody, "amount", cast.ToString(quantity))

	log.Debugf("reqBody: %v", reqBody)

	reqBuf := bytes.NewBufferString(reqBody)

	req, err := c.newPrivateRequest(ctx, "POST", "orders", nil, reqBuf)
	if err != nil {
		return "", err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return "", err
	}

	orderID := gjson.GetBytes(body, "data.order_id").String()

	return orderID, err
}

// GetOrder Check an order's status, for GET /api/v3/order
func (c *Client) GetOrder(ctx context.Context, quote string, base string, orderID string) (*model.BigONEOrder, error) {
	req, err := c.newPrivateRequest(ctx, "GET", "orders/"+orderID, nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.BigONEOrder
	order.ID = gjson.GetBytes(body, "data.order_id").String()
	order.Type = gjson.GetBytes(body, "data.order_type").String()
	order.Side = gjson.GetBytes(body, "data.order_side").String()
	order.Status = gjson.GetBytes(body, "data.order_state").String()
	order.Price = cast.ToFloat64(gjson.GetBytes(body, "data.price").String())
	order.Amount = cast.ToFloat64(gjson.GetBytes(body, "data.amount").String())
	order.DealAmount = cast.ToFloat64(gjson.GetBytes(body, "data.filled_amount").String())
	order.Time = cast.ToTime(gjson.GetBytes(body, "data.updated_at").String())
	order.Raw = string(body)

	return &order, nil
}

// GetOrders Get all open orders on a symbol, for GET /orders{?market,limit}
func (c *Client) GetOrders(ctx context.Context, quote string, base string, limit int64) ([]model.BigONEOrder, error) {
	v := url.Values{}
	v.Set("market", strings.ToUpper(quote)+"-"+strings.ToUpper(base))
	v.Set("limit", cast.ToString(limit))

	req, err := c.newPrivateRequest(ctx, "GET", "orders", v, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.BigONEOrder
	var orders []model.BigONEOrder

	gjson.GetBytes(body, "data").ForEach(func(key, value gjson.Result) bool {
		order.ID = value.Get("order_id").String()
		order.Type = value.Get("order_type").String()
		order.Side = value.Get("order_side").String()
		order.Status = value.Get("order_state").String()
		order.Price = cast.ToFloat64(value.Get("price").String())
		order.Amount = cast.ToFloat64(value.Get("amount").String())
		order.DealAmount = cast.ToFloat64(value.Get("filled_amount").String())
		order.Time = cast.ToTime(value.Get("updated_at").String())
		order.Raw = string(body)

		orders = append(orders, order)
		return true // keep iterating
	})
	return orders, nil
}

// CancelOrder Cancel an Order for DELETE https://api.big.one/orders/{id}
func (c *Client) CancelOrder(ctx context.Context, quote string, base string, orderID string) error {
	req, err := c.newPrivateRequest(ctx, "DELETE", "orders/"+strings.ToUpper(quote)+"-"+strings.ToUpper(base), nil, nil)
	if err != nil {
		return err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return err
	}

	log.Debugf("Response body: %v", string(body))

	return nil
}

// GetAccount List Accounts of Current User, for GET https://api.big.one/accounts
func (c *Client) GetAccount(ctx context.Context) ([]model.Balance, error) {
	req, err := c.newPrivateRequest(ctx, "GET", "accounts", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var balance model.Balance
	var balances []model.Balance
	gjson.GetBytes(body, "data").ForEach(func(key, value gjson.Result) bool {
		balance.Currency = value.Get("account_type").String()
		balance.Free = cast.ToFloat64(value.Get("active_balance").String())
		balance.Frozen = cast.ToFloat64(value.Get("frozen_balance").String())

		balances = append(balances, balance)
		return true // keep iterating
	})

	return balances, nil
}
