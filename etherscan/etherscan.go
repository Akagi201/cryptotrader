package etherscan

import (
	"io/ioutil"
	"net/http"

	"math/big"

	"github.com/Akagi201/cryptotrader/util"
	"github.com/ethereum/go-ethereum/common/math"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://api.etherscan.io/api"
)

// EtherScan API data
type EtherScan struct {
	ApiKey string
}

// New create new EtherScan API data
func New(apiKey string) *EtherScan {
	return &EtherScan{
		ApiKey: apiKey,
	}
}

func (es *EtherScan) GetBalance(addr string) (*big.Float, error) {
	url := API + "?module=account&action=balance&address=" + addr + "&tag=latest&apikey=" + es.ApiKey

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

	balance := gjson.GetBytes(body, "result").String()

	balanceInWei := math.MustParseBig256(balance)
	balanceInEther := util.WeiToEther(balanceInWei)

	return balanceInEther, nil
}
