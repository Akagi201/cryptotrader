// Package eosforce eosforce RPC and transaction package

package eosforce

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// EosForce Client data
type EosForce struct {
	URL        url.URL
	CmdPrefix  []string
	HTTPClient *http.Client
}

// New create new EosForce Client
func New(cmdPrefix []string, rpcScheme string, rpcHost string) *EosForce {
	u := url.URL{
		Scheme: rpcScheme,
		Host:   rpcHost,
		Path:   "/v1",
	}

	return &EosForce{
		URL:        u,
		CmdPrefix:  cmdPrefix,
		HTTPClient: &http.Client{},
	}
}

func (ef *EosForce) newRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader) (*http.Request, error) {
	u := ef.URL
	u.Path = path.Join(ef.URL.Path, spath)
	u.RawQuery = values.Encode()
	log.Debugf("Request URL: %#v", u.String())

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (ef *EosForce) getResponse(req *http.Request) ([]byte, error) {
	res, err := ef.HTTPClient.Do(req)

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

// GetAvailable 获取指定账户可用余额
func (ef *EosForce) GetAvailable(ctx context.Context, account string) (float64, error) {
	reqBody := `{
		"scope": "eosio",
		"code": "eosio",
		"table": "accounts",
        "json": true,
        "table_key": "%s",
		"limit": 1
	}`

	reqBody, _ = sjson.Set(reqBody, "table_key", account)

	reqBuf := bytes.NewBufferString(reqBody)

	req, err := ef.newRequest(ctx, "POST", "chain/get_table_rows", nil, reqBuf)
	if err != nil {
		return -1, err
	}

	body, err := ef.getResponse(req)
	if err != nil {
		return -1, err
	}

	if gjson.GetBytes(body, "rows.#").Int() == 0 {
		return -1, errors.New("account not found")
	}

	available := gjson.GetBytes(body, "rows.0").Get("available").String()

	return cast.ToFloat64(strings.Split(available, " ")[0]), nil
}

// GetStaked 获取指定账户投票金额
func (ef *EosForce) GetStaked(ctx context.Context, account string, bp string) (float64, error) {
	reqBody := `{
		"scope": "%s",
		"code": "eosio",
		"table": "votes",
        "json": true,
        "table_key": "%s",
		"limit": 1
	}`

	reqBody, _ = sjson.Set(reqBody, "scope", account)
	reqBody, _ = sjson.Set(reqBody, "table_key", bp)

	reqBuf := bytes.NewBufferString(reqBody)

	req, err := ef.newRequest(ctx, "POST", "chain/get_table_rows", nil, reqBuf)
	if err != nil {
		return -1, err
	}

	body, err := ef.getResponse(req)
	if err != nil {
		return -1, err
	}

	if gjson.GetBytes(body, "rows.#").Int() == 0 {
		return -1, errors.New("account not found")
	}

	available := gjson.GetBytes(body, "rows.0").Get("staked").String()

	return cast.ToFloat64(strings.Split(available, " ")[0]), nil
}
