// Package binance binance rest api
package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Akagi201/binancego/model"
	"github.com/Akagi201/utilgo/enums"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

const (
	RestHost = "www.binance.com"
	ApiV1    = "v1"
	ApiV3    = "v3"
)

// Kline intervals
// m -> minutes; h -> hours; d -> days; w -> weeks; M -> months
var (
	Interval    enums.Enum
	Interval1m  = Interval.Iota("1m")
	Interval3m  = Interval.Iota("3m")
	Interval5m  = Interval.Iota("5m")
	Interval15m = Interval.Iota("15m")
	Interval30m = Interval.Iota("30m")
	Interval1h  = Interval.Iota("1h")
	Interval2h  = Interval.Iota("2h")
	Interval4h  = Interval.Iota("4h")
	Interval6h  = Interval.Iota("6h")
	Interval8h  = Interval.Iota("8h")
	Interval12h = Interval.Iota("12h")
	Interval1d  = Interval.Iota("1d")
	Interval3d  = Interval.Iota("3d")
	Interval1w  = Interval.Iota("1w")
	Interval1M  = Interval.Iota("1M")
)

// Order types
var (
	OrderType   enums.Enum
	OrderLimit  = OrderType.Iota("LIMIT")
	OrderMarket = OrderType.Iota("MARKET")
)

// Order Side
var (
	OrderSide enums.Enum
	OrderBuy  = OrderSide.Iota("BUY")
	OrderSell = OrderSide.Iota("SELL")
)

// Time in force
var (
	TimeInForce enums.Enum
	GTC         = TimeInForce.Iota("GTC")
	IOC         = TimeInForce.Iota("IOC")
)

// Client Binance client
type Client struct {
	URL        url.URL
	HTTPClient *http.Client
	AccessKey  string
	SecretKey  string
}

// New creates a new Binance Client
func New(accessKey string, secretKey string) *Client {
	u := url.URL{
		Scheme: "https",
		Host:   RestHost,
		Path:   "/api/",
	}

	c := Client{
		URL:        u,
		HTTPClient: &http.Client{},
		AccessKey:  accessKey,
		SecretKey:  secretKey,
	}

	return &c
}

func (c *Client) newRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader, version string) (*http.Request, error) {
	u := c.URL
	u.Path = path.Join(c.URL.Path, version, spath)
	u.RawQuery = values.Encode()
	log.Debugf("Request URL: %#v", u.String())

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

// Sign sign the params with a secret key
func (c *Client) Sign(secretKey string, totalParams string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(totalParams))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Client) newPrivateRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader, recvWindow int64) (*http.Request, error) {
	req, err := c.newRequest(ctx, method, spath, values, body, ApiV3)
	if err != nil {
		return nil, err
	}

	var bodyText string
	if body != nil {
		bodyBytes, _ := ioutil.ReadAll(body)
		bodyText = string(bodyBytes)
	} else {
		bodyText = ""
	}

	//TODO: support body encode
	_ = bodyText

	if values == nil {
		values = url.Values{}
	}
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	values.Set("timestamp", timestamp)

	if recvWindow != 0 {
		values.Set("recvWindow", strconv.FormatInt(recvWindow, 10))
	} else {
		values.Set("recvWindow", "5000")
	}

	sign := c.Sign(c.SecretKey, values.Encode())
	values.Add("signature", sign)

	req.Header.Set("X-MBX-APIKEY", c.AccessKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

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

// GetTicker 24hr ticker price change statistics, for GET /api/v1/ticker/24hr
func (c *Client) GetTicker(ctx context.Context, quote string, base string) (*model.Ticker, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))

	req, err := c.newRequest(ctx, "GET", "ticker/24hr", v, nil, ApiV1)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	buyRes := gjson.GetBytes(body, "bidPrice").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "askPrice").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "lastPrice").String()
	last, err := strconv.ParseFloat(lastRes, 64)
	if err != nil {
		return nil, err
	}

	lowRes := gjson.GetBytes(body, "lowPrice").String()
	low, err := strconv.ParseFloat(lowRes, 64)
	if err != nil {
		return nil, err
	}

	highRes := gjson.GetBytes(body, "highPrice").String()
	high, err := strconv.ParseFloat(highRes, 64)
	if err != nil {
		return nil, err
	}

	volRes := gjson.GetBytes(body, "volume").String()
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
		Raw:  string(body),
	}, nil
}

// Ping Test connectivity to the Rest API, for GET /api/v1/ping
func (c *Client) Ping(ctx context.Context) error {
	req, err := c.newRequest(ctx, "GET", "ping", nil, nil, ApiV1)
	if err != nil {
		return err
	}

	_, err = c.getResponse(req)
	if err != nil {
		return err
	}

	return nil
}

// GetTime Check server time, for GET /api/v1/time
func (c *Client) GetTime(ctx context.Context) (*time.Time, error) {
	req, err := c.newRequest(ctx, "GET", "time", nil, nil, ApiV1)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	serverTime := gjson.GetBytes(body, "serverTime").Int()
	t := time.Unix(0, serverTime*int64(time.Millisecond))
	return &t, nil
}

// GetDepth Order book, for GET /api/v1/depth
func (c *Client) GetDepth(ctx context.Context, quote string, base string, limit int64) (*model.OrderBook, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))

	if limit != 0 {
		v.Set("limit", strconv.FormatInt(limit, 10))
	}

	req, err := c.newRequest(ctx, "GET", "depth", v, nil, ApiV1)
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
	gjson.GetBytes(body, "bids").ForEach(func(key, value gjson.Result) bool {
		order.Price = cast.ToFloat64(value.Array()[0].String())
		order.Amount = cast.ToFloat64(value.Array()[1].String())
		orderBook.Bids = append(orderBook.Bids, order)
		return true // keep iterating
	})

	gjson.GetBytes(body, "asks").ForEach(func(key, value gjson.Result) bool {
		order.Price = cast.ToFloat64(value.Array()[0].String())
		order.Amount = cast.ToFloat64(value.Array()[1].String())
		orderBook.Asks = append(orderBook.Asks, order)
		return true // keep iterating
	})

	//t := gjson.GetBytes(body, "lastUpdateId").Int()
	//orderBook.Time = time.Unix(0, t*int64(time.Millisecond))
	orderBook.Raw = string(body)

	return &orderBook, nil
}

// GetTickers Symbols price ticker, for GET /api/v1/ticker/allPrices
func (c *Client) GetTickers(ctx context.Context) ([]model.SimpleTicker, error) {
	req, err := c.newRequest(ctx, "GET", "ticker/allPrices", nil, nil, ApiV1)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var simpleTicker model.SimpleTicker
	var simpleTickers []model.SimpleTicker
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		simpleTicker.Symbol = value.Get("symbol").String()
		simpleTicker.Price = cast.ToFloat64(value.Get("price").String())
		simpleTickers = append(simpleTickers, simpleTicker)
		return true // keep iterating
	})

	return simpleTickers, nil
}

// GetTrades Compressed/Aggregate trades list, for GET /api/v1/aggTrades
func (c *Client) GetTrades(ctx context.Context, quote string, base string, fromID int64, startTime int64, endTime int64, limit int64) ([]model.Trade, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))

	if fromID != 0 {
		v.Set("fromId", strconv.FormatInt(fromID, 10))
	}

	if startTime != 0 {
		v.Set("startTime", strconv.FormatInt(startTime, 10))
	}

	if endTime != 0 {
		v.Set("startTime", strconv.FormatInt(endTime, 10))
	}

	if limit != 0 {
		v.Set("limit", strconv.FormatInt(limit, 10))
	}

	req, err := c.newRequest(ctx, "GET", "aggTrades", v, nil, ApiV1)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var trade model.Trade
	var trades []model.Trade
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		trade.ID = value.Get("a").Int()
		trade.Price = cast.ToFloat64(value.Get("p").String())
		trade.Amount = cast.ToFloat64(value.Get("q").String())
		trade.Time = cast.ToTime(value.Get("T").Int() / 1000)
		trade.Raw = value.String()
		trades = append(trades, trade)
		return true // keep iterating
	})

	return trades, nil
}

// GetRecords for Kline/candlesticks, for GET /api/v1/ticker/allPrices
func (c *Client) GetRecords(ctx context.Context, quote string, base string, interval string, startTime int64, endTime int64, limit int64) ([]model.Record, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))
	v.Set("interval", interval)

	if startTime != 0 {
		v.Set("startTime", strconv.FormatInt(startTime, 10))
	}

	if endTime != 0 {
		v.Set("startTime", strconv.FormatInt(endTime, 10))
	}

	if limit != 0 {
		v.Set("limit", strconv.FormatInt(limit, 10))
	}

	req, err := c.newRequest(ctx, "GET", "klines", v, nil, ApiV1)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var record model.Record
	var records []model.Record

	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		record.Open = cast.ToFloat64(value.Array()[1].String())
		record.High = cast.ToFloat64(value.Array()[2].String())
		record.Low = cast.ToFloat64(value.Array()[3].String())
		record.Close = cast.ToFloat64(value.Array()[4].String())
		record.Vol = cast.ToFloat64(value.Array()[5].String())
		record.Time = cast.ToTime(value.Array()[0].Int() / 1000)
		record.Raw = value.String()
		records = append(records, record)
		return true // keep iterating
	})
	return records, nil
}

// GetBookTickers Symbols order book ticker, for GET /api/v1/ticker/allBookTickers
func (c *Client) GetBookTickers(ctx context.Context) ([]model.BookTicker, error) {
	req, err := c.newRequest(ctx, "GET", "ticker/allBookTickers", nil, nil, ApiV1)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var bookTicker model.BookTicker
	var bookTickers []model.BookTicker
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		bookTicker.Symbol = value.Get("symbol").String()
		bookTicker.BidPrice = cast.ToFloat64(value.Get("bidPrice").String())
		bookTicker.BidAmount = cast.ToFloat64(value.Get("bidQty").String())
		bookTicker.AskPrice = cast.ToFloat64(value.Get("askPrice").String())
		bookTicker.AskAmount = cast.ToFloat64(value.Get("askQty").String())
		bookTickers = append(bookTickers, bookTicker)
		return true // keep iterating
	})

	return bookTickers, nil
}

// GetAccount Account information (SIGNED), for GET /api/v3/account
func (c *Client) GetAccount(ctx context.Context, recvWindow int64) ([]model.Balance, error) {
	req, err := c.newPrivateRequest(ctx, "GET", "account", nil, nil, recvWindow)
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
	gjson.GetBytes(body, "balances").ForEach(func(key, value gjson.Result) bool {
		balance.Currency = value.Get("asset").String()
		balance.Free = cast.ToFloat64(value.Get("free").String())
		balance.Frozen = cast.ToFloat64(value.Get("locked").String())

		balances = append(balances, balance)
		return true // keep iterating
	})

	return balances, nil
}

// Trade Send in a new order, for POST /api/v3/order
func (c *Client) Trade(ctx context.Context, quote string, base string, side string, typ string, timeInForce string, quantity float64, price float64, stopPrice float64, icebergQty float64, recvWindow int64) (int64, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))
	v.Set("side", side)
	v.Set("type", typ)
	v.Set("timeInForce", timeInForce)
	v.Set("quantity", cast.ToString(quantity))
	v.Set("price", cast.ToString(price))

	if stopPrice != 0 {
		v.Set("stopPrice", cast.ToString(stopPrice))
	}
	if icebergQty != 0 {
		v.Set("icebergQty", cast.ToString(icebergQty))
	}

	req, err := c.newPrivateRequest(ctx, "POST", "order/test", v, nil, recvWindow)
	if err != nil {
		return 0, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return 0, err
	}

	log.Debugf("Response body: %v", string(body))

	return gjson.GetBytes(body, "orderId").Int(), nil
}

// GetOrder Check an order's status, for GET /api/v3/order
func (c *Client) GetOrder(ctx context.Context, quote string, base string, orderID int64, recvWindow int64) (*model.Order, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))
	v.Set("orderId", cast.ToString(orderID))

	req, err := c.newPrivateRequest(ctx, "GET", "order", v, nil, recvWindow)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.Order
	order.ID = gjson.GetBytes(body, "orderId").Int()
	order.Amount = cast.ToFloat64(gjson.GetBytes(body, "origQty").String())
	order.DealAmount = cast.ToFloat64(gjson.GetBytes(body, "executedQty").String())
	order.Price = cast.ToFloat64(gjson.GetBytes(body, "price").String())
	order.Status = gjson.GetBytes(body, "status").String()
	order.Type = gjson.GetBytes(body, "type").String()
	order.Side = gjson.GetBytes(body, "side").String()
	order.Raw = string(body)

	return &order, nil
}

// CancelOrder Cancel an active orderf for DELETE /api/v3/order
func (c *Client) CancelOrder(ctx context.Context, quote string, base string, orderID int64, recvWindow int64) error {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))
	v.Set("orderId", cast.ToString(orderID))

	req, err := c.newPrivateRequest(ctx, "DELETE", "order", v, nil, recvWindow)
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

// GetOrders Get all open orders on a symbol, for GET /api/v3/openOrders
func (c *Client) GetOrders(ctx context.Context, quote string, base string, recvWindow int64) ([]model.Order, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))

	req, err := c.newPrivateRequest(ctx, "GET", "openOrders", v, nil, recvWindow)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.Order
	var orders []model.Order

	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		order.ID = value.Get("orderId").Int()
		order.Amount = cast.ToFloat64(value.Get("origQty").String())
		order.DealAmount = cast.ToFloat64(value.Get("executedQty").String())
		order.Price = cast.ToFloat64(value.Get("price").String())
		order.Status = value.Get("status").String()
		order.Type = value.Get("type").String()
		order.Side = value.Get("side").String()
		order.Raw = value.String()

		orders = append(orders, order)
		return true // keep iterating
	})
	return orders, nil
}

// GetAllOrders Get all account orders; active, canceled, or filled, for GET /api/v3/allOrders
func (c *Client) GetAllOrders(ctx context.Context, quote string, base string, orderID int64, limit int64, recvWindow int64) ([]model.Order, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))
	if orderID != 0 {
		v.Set("orderId", cast.ToString(orderID))
	}
	if limit != 0 {
		v.Set("limit", cast.ToString(limit))
	}

	req, err := c.newPrivateRequest(ctx, "GET", "allOrders", v, nil, recvWindow)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.Order
	var orders []model.Order

	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		order.ID = value.Get("orderId").Int()
		order.Amount = cast.ToFloat64(value.Get("origQty").String())
		order.DealAmount = cast.ToFloat64(value.Get("executedQty").String())
		order.Price = cast.ToFloat64(value.Get("price").String())
		order.Status = value.Get("status").String()
		order.Type = value.Get("type").String()
		order.Side = value.Get("side").String()
		order.Raw = value.String()

		orders = append(orders, order)
		return true // keep iterating
	})
	return orders, nil
}

// GetMyTrades Get trades for a specific account and symbol, for GET /api/v3/myTrades
func (c *Client) GetMyTrades(ctx context.Context, quote string, base string, fromID int64, limit int64, recvWindow int64) ([]model.Trade, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToUpper(quote)+strings.ToUpper(base))

	if fromID != 0 {
		v.Set("fromId", strconv.FormatInt(fromID, 10))
	}

	if limit != 0 {
		v.Set("limit", strconv.FormatInt(limit, 10))
	}

	req, err := c.newPrivateRequest(ctx, "GET", "myTrades", v, nil, recvWindow)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var trade model.Trade
	var trades []model.Trade
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		trade.ID = value.Get("id").Int()
		trade.Price = cast.ToFloat64(value.Get("price").String())
		trade.Amount = cast.ToFloat64(value.Get("qty").String())
		trade.Time = cast.ToTime(value.Get("time").Int() / 1000)
		if value.Get("isBuyer").Bool() {
			trade.Type = "buy"
		} else {
			trade.Type = "sell"
		}
		trade.Raw = value.String()
		trades = append(trades, trade)
		return true // keep iterating
	})

	return trades, nil
}
