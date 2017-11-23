// Package zb zb rest api package
package zb

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	MarketAPI = "http://api.zb.com/data/v1/"
	TradeAPI  = "https://trade.zb.com/api/"
)

// ZB API data
type ZB struct {
	AccessKey string
	SecretKey string
}

// New create new Zb API data
func New(accessKey string, secretKey string) *ZB {
	return &ZB{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (z *ZB) GetTicker(base string, quote string) (*model.Ticker, error) {
	log.Debugf("Currency base: %s, quote: %s", base, quote)

	url := MarketAPI + "ticker?market=" + quote + "_" + base

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

	buyRes := gjson.GetBytes(body, "ticker.buy").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "ticker.sell").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "ticker.last").String()
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

// GetOrderBook 市场深度
//
// * size: 档位 1-50, 如果有合并深度, 只能返回 5 档深度
// * merge:
//   * btc_cny: 可选 1, 0.1
//   * ltc_cny: 可选 0.5, 0.3, 0.1
//   * eth_cny: 可选 0.5, 0.3, 0.1
//   * etc_cny: 可选 0.3, 0.1
//   * bts_cny: 可选 1, 0.1
func (z *ZB) GetOrderBook(base string, quote string, size int, merge float64) (*model.OrderBook, error) {
	url := MarketAPI + "depth?currency=" + quote + "_" + base + "&size=" + strconv.Itoa(size) + "&merge=" + strconv.FormatFloat(merge, 'f', -1, 64)

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

	orderBook := &model.OrderBook{
		Time: time.Unix(gjson.GetBytes(body, "timestamp").Int(), 0),
	}

	gjson.GetBytes(body, "asks").ForEach(func(k, v gjson.Result) bool {
		orderBook.Asks = append(orderBook.Asks, model.MarketOrder{
			Price:  v.Array()[0].Float(),
			Amount: v.Array()[1].Float(),
		})

		return true
	})

	gjson.GetBytes(body, "bids").ForEach(func(k, v gjson.Result) bool {
		orderBook.Bids = append(orderBook.Bids, model.MarketOrder{
			Price:  v.Array()[0].Float(),
			Amount: v.Array()[1].Float(),
		})

		return true
	})

	return orderBook, nil
}

// GetTrades 获取历史成交
//
// * currency: quote_base
//   * btc_cny: 比特币/人民币
//   * ltc_cny: 莱特币/人民币
//   * eth_cny: 以太币/人民币
//   * etc_cny: ETC币/人民币
//   * bts_cny: BTS币/人民币
// * since: 从指定交易 ID 后 50 条数据
func (z *ZB) GetTrades(base string, quote string, since int) (*model.Trades, error) {
	url := MarketAPI + "trades?currency=" + quote + "_" + base
	if since != 0 {
		url += "&since=" + strconv.Itoa(since)
	}

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

	trades := new(model.Trades)

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		trade := &model.Trade{
			Amount:    v.Get("amount").Float(),
			Price:     v.Get("price").Float(),
			Tid:       v.Get("tid").Int(),
			TradeType: v.Get("trade_type").String(),
			Type:      v.Get("type").String(),
			Date:      time.Unix(v.Get("date").Int(), 0),
		}
		*trades = append(*trades, trade)
		return true
	})

	return trades, nil
}

// GetKline 获取 K 线数据
//
// * currency: quote_base
//   * btc_cny: 比特币/人民币
//   * ltc_cny: 莱特币/人民币
//   * eth_cny: 以太币/人民币
//   * etc_cny: ETC币/人民币
//   * bts_cny: BTS币/人民币
// * typ:
//   * 1min: 1 分钟
//   * 3min: 3 分钟
//   * 5min: 5 分钟
//   * 15min: 15 分钟
//   * 30min: 30 分钟
//   * 1day: 1 日
//   * 3day: 3 日
//   * 1week: 1 周
//   * 1hour: 1 小时
//   * 2hour: 2 小时
//   * 4hour: 4 小时
//   * 6hour: 6小时
//   * 12hour: 12 小时
// * since: 从这个时间戳之后的
// * size: 返回数据的条数限制(默认为 1000, 如果返回数据多于 1000 条, 那么只返回 1000 条)
func (z *ZB) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := MarketAPI + "kline?currency=" + quote + "_" + base

	if len(typ) != 0 {
		url += "&type=" + typ
	}

	if since != 0 {
		url += "&since=" + strconv.Itoa(since)
	}

	if size != 0 {
		url += "&size=" + strconv.Itoa(size)
	}

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

	var records []model.Record

	gjson.GetBytes(body, "data").ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
		}

		records = append(records, record)
		return true
	})

	return records, nil
}

// SecretDigest calc secert digest
func (z *ZB) SecretDigest() string {
	sha := sha1.New()
	sha.Write([]byte(z.SecretKey))
	return hex.EncodeToString(sha.Sum(nil))
}

// Sign calc sign string
func (z *ZB) Sign(uri string) string {
	digest := z.SecretDigest()
	mac := hmac.New(md5.New, []byte(digest))
	mac.Write([]byte(uri))
	return hex.EncodeToString(mac.Sum(nil))
}

// GetUserAddress 获取用户充值地址
//
// * currency:
//   * btc: BTC
//   * ltc: LTC
//   * eth: 以太币
//   * etc: ETC币
func (z *ZB) GetUserAddress(currency string) (string, error) {
	url := "method=getUserAddress"
	url += "&accesskey=" + z.AccessKey
	url += "&currency=" + currency
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getUserAddress?" + url

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Debugf("Response body: %v", string(body))

	return gjson.GetBytes(body, "message.datas.key").String(), nil
}

// PlaceOrder 委托下单
//
// * price: 单价(cny 保留小数后 2 位, btc 保留小数后 6 位)
// * amount: 交易数量(btc, ltc, eth, etc保留小数后 3 位)
// * tradeType: 交易类型 1/0[buy / sell]
// * currency: quote_base
//   * btc_cny: 比特币/人民币
//   * ltc_cny: 莱特币/人民币
//   * eth_cny: 以太币/人民币
//   * etc_cny: ETC币/人民币
//   * bts_cny: BTS币/人民币
// return 委托挂单号
func (z *ZB) PlaceOrder(price float64, amount float64, tradeType int, base string, quote string) (string, error) {
	url := "method=order"
	url += "&accesskey=" + z.AccessKey
	url += "&price=" + strconv.FormatFloat(price, 'f', -1, 64)
	url += "&amount=" + strconv.FormatFloat(amount, 'f', -1, 64)
	url += "&tradeType=" + strconv.Itoa(tradeType)
	url += "&currency=" + quote + "_" + base
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "order?" + url

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Debugf("Response body: %v", string(body))
	code := gjson.GetBytes(body, "code").String()
	if code == "1000" {
		return gjson.GetBytes(body, "id").String(), nil
	}

	return "", errors.New(code)
}

// CancelOrder 取消委托
//
// * id: 委托挂单号
// * currency: quote_base
//   * btc_cny: 比特币/人民币
//   * ltc_cny: 莱特币/人民币
//   * eth_cny: 以太币/人民币
//   * etc_cny: ETC币/人民币
//   * bts_cny: BTS币/人民币
func (z *ZB) CancelOrder(id string, base string, quote string) error {
	url := "method=cancelOrder"
	url += "&accesskey=" + z.AccessKey
	url += "&id=" + id
	url += "&currency=" + quote + "_" + base
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "cancelOrder?" + url

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debugf("Response body: %v", string(body))

	code := gjson.GetBytes(body, "code").String()
	if code == "1000" {
		return nil
	}

	return errors.New(code)
}

// GetOrder 获取委托买单或卖单
// id: 委托挂单号
func (z *ZB) GetOrder(id string, base string, quote string) (*model.ZBOrder, error) {
	url := "method=getOrder"
	url += "&accesskey=" + z.AccessKey
	url += "&id=" + id
	url += "&currency=" + quote + "_" + base
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getOrder?" + url

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

	return &model.ZBOrder{
		Currency:    gjson.GetBytes(body, "currency").String(),
		Fees:        gjson.GetBytes(body, "fees").Float(),
		ID:          gjson.GetBytes(body, "id").String(),
		Price:       gjson.GetBytes(body, "price").Float(),
		Status:      gjson.GetBytes(body, "status").Int(),
		TotalAmount: gjson.GetBytes(body, "total_amount").Float(),
		TradeAmount: gjson.GetBytes(body, "trade_amount").Float(),
		TradePrice:  gjson.GetBytes(body, "trade_price").Float(),
		TradeDate:   time.Unix(gjson.GetBytes(body, "trade_date").Int(), 0),
		TradeMoney:  gjson.GetBytes(body, "trade_money").Float(),
		Type:        gjson.GetBytes(body, "type").Int(),
	}, nil
}

// GetOrders 获取多个委托买单或卖单, 每次请求返回 10 条记录
func (z *ZB) GetOrders(tradeType int, base string, quote string, pageIndex int) ([]*model.ZBOrder, error) {
	url := "method=getOrders"
	url += "&accesskey=" + z.AccessKey
	url += "&tradeType=" + strconv.Itoa(tradeType)
	url += "&currency=" + quote + "_" + base
	url += "&pageIndex=" + strconv.Itoa(pageIndex)
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getOrders?" + url

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

	var orders []*model.ZBOrder
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		orders = append(orders, &model.ZBOrder{
			Currency:    v.Get("currency").String(),
			Fees:        v.Get("fees").Float(),
			ID:          v.Get("id").String(),
			Price:       v.Get("price").Float(),
			Status:      v.Get("status").Int(),
			TotalAmount: v.Get("total_amount").Float(),
			TradeAmount: v.Get("trade_amount").Float(),
			TradePrice:  v.Get("trade_price").Float(),
			TradeDate:   time.Unix(v.Get("trade_date").Int(), 0),
			TradeMoney:  v.Get("trade_money").Float(),
			Type:        v.Get("type").Int(),
		})

		return true
	})

	return orders, nil
}

// GetOrdersNew (新)获取多个委托买单或卖单，每次请求返回pageSize<100条记录
func (z *ZB) GetOrdersNew(tradeType int, base string, quote string, pageIndex int, pageSize int) ([]*model.ZBOrder, error) {
	url := "method=getOrdersNew"
	url += "&accesskey=" + z.AccessKey
	url += "&tradeType=" + strconv.Itoa(tradeType)
	url += "&currency=" + quote + "_" + base
	url += "&pageIndex=" + strconv.Itoa(pageIndex)
	url += "&pageSize=" + strconv.Itoa(pageSize)
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getOrdersNew?" + url

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

	var orders []*model.ZBOrder
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		orders = append(orders, &model.ZBOrder{
			Currency:    v.Get("currency").String(),
			Fees:        v.Get("fees").Float(),
			ID:          v.Get("id").String(),
			Price:       v.Get("price").Float(),
			Status:      v.Get("status").Int(),
			TotalAmount: v.Get("total_amount").Float(),
			TradeAmount: v.Get("trade_amount").Float(),
			TradePrice:  v.Get("trade_price").Float(),
			TradeDate:   time.Unix(v.Get("trade_date").Int(), 0),
			TradeMoney:  v.Get("trade_money").Float(),
			Type:        v.Get("type").Int(),
		})

		return true
	})

	return orders, nil
}

// GetOrdersIgnoreTradeType 与getOrdersNew的区别是取消tradeType字段过滤，可同时获取买单和卖单，每次请求返回pageSize<100条记录
func (z *ZB) GetOrdersIgnoreTradeType(base string, quote string, pageIndex int, pageSize int) ([]*model.ZBOrder, error) {
	url := "method=getOrdersIgnoreTradeType"
	url += "&accesskey=" + z.AccessKey
	url += "&currency=" + quote + "_" + base
	url += "&pageIndex=" + strconv.Itoa(pageIndex)
	url += "&pageSize=" + strconv.Itoa(pageSize)
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getOrdersIgnoreTradeType?" + url

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

	var orders []*model.ZBOrder
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		orders = append(orders, &model.ZBOrder{
			Currency:    v.Get("currency").String(),
			Fees:        v.Get("fees").Float(),
			ID:          v.Get("id").String(),
			Price:       v.Get("price").Float(),
			Status:      v.Get("status").Int(),
			TotalAmount: v.Get("total_amount").Float(),
			TradeAmount: v.Get("trade_amount").Float(),
			TradePrice:  v.Get("trade_price").Float(),
			TradeDate:   time.Unix(v.Get("trade_date").Int(), 0),
			TradeMoney:  v.Get("trade_money").Float(),
			Type:        v.Get("type").Int(),
		})

		return true
	})

	return orders, nil
}

// GetUnfinishedOrdersIgnoreTradeType 获取未成交或部份成交的买单和卖单，每次请求返回pageSize<=100条记录
func (z *ZB) GetUnfinishedOrdersIgnoreTradeType(base string, quote string, pageIndex int, pageSize int) ([]*model.ZBOrder, error) {
	url := "method=getUnfinishedOrdersIgnoreTradeType"
	url += "&accesskey=" + z.AccessKey
	url += "&currency=" + quote + "_" + base
	url += "&pageIndex=" + strconv.Itoa(pageIndex)
	url += "&pageSize=" + strconv.Itoa(pageSize)
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getUnfinishedOrdersIgnoreTradeType?" + url

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

	var orders []*model.ZBOrder
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		orders = append(orders, &model.ZBOrder{
			Currency:    v.Get("currency").String(),
			Fees:        v.Get("fees").Float(),
			ID:          v.Get("id").String(),
			Price:       v.Get("price").Float(),
			Status:      v.Get("status").Int(),
			TotalAmount: v.Get("total_amount").Float(),
			TradeAmount: v.Get("trade_amount").Float(),
			TradePrice:  v.Get("trade_price").Float(),
			TradeDate:   time.Unix(v.Get("trade_date").Int(), 0),
			TradeMoney:  v.Get("trade_money").Float(),
			Type:        v.Get("type").Int(),
		})

		return true
	})

	return orders, nil
}

// GetWithdrawAddress 获取用户认证的提现地址
//
// * currency:
//   * btc: BTC
//   * ltc: LTC
//   * eth: 以太币
//   * etc: ETC币
func (z *ZB) GetWithdrawAddress(currency string) (string, error) {
	url := "method=getWithdrawAddress"
	url += "&accesskey=" + z.AccessKey
	url += "&currency=" + currency
	sign := z.Sign(url)
	url += "&sign=" + sign
	url += "&reqTime=" + strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)

	log.Debugf("Request url: %v", url)

	url = TradeAPI + "getWithdrawAddress?" + url

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Debugf("Response body: %v", string(body))

	return gjson.GetBytes(body, "message.datas.key").String(), nil
}
