// Package gateio gate.io rest api package
package gateio

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	RestHost = "data.gate.io"
	ApiVer   = "api2/1"
)

// Client OkEx client
type Client struct {
	URL        url.URL
	HTTPClient *http.Client
	AccessKey  string
	SecretKey  string
}

// New creates a new OkEx Client
func New(accessKey string, secretKey string) *Client {
	u := url.URL{
		Scheme: "http",
		Host:   RestHost,
		Path:   "/",
	}

	c := Client{
		URL:        u,
		HTTPClient: &http.Client{},
		AccessKey:  accessKey,
		SecretKey:  secretKey,
	}

	return &c
}

func (c *Client) newRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader) (*http.Request, error) {
	u := c.URL
	u.Path = path.Join(c.URL.Path, ApiVer, spath)
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

// GetPairs 返回所有系统支持的交易对, for http://data.gate.io/api2/1/pairs
func (c *Client) GetPairs(ctx context.Context) ([]string, error) {
	req, err := c.newRequest(ctx, "GET", "pairs", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))
	var pairs []string
	for _, v := range gjson.ParseBytes(body).Array() {
		pairs = append(pairs, v.String())
	}
	return pairs, nil
}

// GetMarketInfo 交易市场订单参数, 返回所有系统支持的交易市场的参数信息，包括交易费，最小下单量，价格精度等。for http://data.gate.io/api2/1/marketinfo
func (c *Client) GetMarketInfo(ctx context.Context) ([]model.MarketInfo, error) {
	req, err := c.newRequest(ctx, "GET", "marketinfo", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var marketInfos []model.MarketInfo
	var marketInfo model.MarketInfo

	for _, v := range gjson.GetBytes(body, "pairs").Array() {
		v.ForEach(func(key, value gjson.Result) bool {
			marketInfo.Symbol = key.String()
			marketInfo.DecimalPlaces = value.Get("decimal_places").Int()
			marketInfo.MinAmount = value.Get("min_amount").Float()
			marketInfo.Fee = value.Get("fee").Float()
			marketInfos = append(marketInfos, marketInfo)
			return true // keep iterating
		})
	}

	return marketInfos, nil
}
