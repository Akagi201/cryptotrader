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

// GetStaked 获取指定账户和 BP 的投票金额
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

// GetUnstaking 获取指定账户和 BP 的赎回金额
func (ef *EosForce) GetUnstaking(ctx context.Context, account string, bp string) (float64, error) {
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

	available := gjson.GetBytes(body, "rows.0").Get("unstaking").String()

	return cast.ToFloat64(strings.Split(available, " ")[0]), nil
}

// calcVoteage 计算最新票龄
// 票龄 + 投票金额 * （当前高度 - 票龄更新高度）
func (ef *EosForce) calcVoteage(voteage int64, staked int64, currentHeight int64, updateHeight int64) int64 {
	return voteage + staked*(currentHeight-updateHeight)
}

// calcReward 计算分红
func (ef *EosForce) calcReward(myVoteage int64, bpVoteage int64, rewardsPool float64) float64 {
	if bpVoteage == 0 {
		return 0
	}

	return float64(myVoteage) * rewardsPool / float64(bpVoteage)
}

// GetRewards 获取指定账户和 BP 的待领分红
func (ef *EosForce) GetRewards(ctx context.Context, account string, bp string) (float64, error) {
	// get info
	var currentHeight int64
	{
		req, err := ef.newRequest(ctx, "POST", "chain/get_info", nil, nil)
		if err != nil {
			return -1, err
		}

		body, err := ef.getResponse(req)
		if err != nil {
			return -1, err
		}

		currentHeight = gjson.GetBytes(body, "head_block_num").Int()
	}

	// get bps table
	var commissionRate int64
	var totalStaked int64
	var rewardsPool float64
	var totalVoteage int64
	var voteageUpdateHeight int64
	var bpVoteage int64
	{
		reqBody := `{
			"scope": "eosio",
			"code": "eosio",
			"table": "bps",
			"json": true,
			"table_key": "%s",
			"limit": 1
		}`

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

		commissionRate = gjson.GetBytes(body, "rows.0").Get("commission_rate").Int()
		totalStaked = gjson.GetBytes(body, "rows.0").Get("total_staked").Int()
		rewardsPoolStr := gjson.GetBytes(body, "rows.0").Get("rewards_pool").String()
		rewardsPool = cast.ToFloat64(strings.Split(rewardsPoolStr, " ")[0])
		totalVoteage = gjson.GetBytes(body, "rows.0").Get("total_voteage").Int()
		voteageUpdateHeight = gjson.GetBytes(body, "rows.0").Get("voteage_update_height").Int()

		bpVoteage = ef.calcVoteage(totalVoteage, totalStaked, currentHeight, voteageUpdateHeight)
	}

	_ = commissionRate

	// get votes table
	var myStaked int64
	var myVoteageUpdateHeight int64
	var myUnstaking float64
	var myUnstakeHeight int64
	var myVoteage int64
	{
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

		myStakedStr := gjson.GetBytes(body, "rows.0").Get("staked").String()
		myStaked = cast.ToInt64(cast.ToFloat64(strings.Split(myStakedStr, " ")[0]))
		myVoteage = gjson.GetBytes(body, "rows.0").Get("voteage").Int()
		myVoteageUpdateHeight = gjson.GetBytes(body, "rows.0").Get("voteage_update_height").Int()
		myUnstaking = gjson.GetBytes(body, "rows.0").Get("unstaking").Float()
		myUnstakeHeight = gjson.GetBytes(body, "rows.0").Get("unstake_height").Int()
	}

	_ = myUnstaking
	_ = myUnstakeHeight

	myFinalVoteage := ef.calcVoteage(myVoteage, myStaked, currentHeight, myVoteageUpdateHeight)

	reward := ef.calcReward(myFinalVoteage, bpVoteage, rewardsPool)

	return reward, nil
}
