package fixer

import (
	"io/ioutil"
	"net/http"

	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	API = "https://api.fixer.io/latest"
)

// Fixer API data
type Fixer struct {
}

// New create new Fixer API data
func New() *Fixer {
	return &Fixer{}
}

func (f *Fixer) GetRate(base string, quote string) (float64, error) {
	url := API + "?base=" + strings.ToUpper(base) + "&symbols=" + strings.ToUpper(quote)

	log.Debugf("Request url: %v", url)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	log.Debugf("Response body: %v", string(body))

	rateRes := gjson.GetBytes(body, "rates."+strings.ToUpper(quote)).String()
	rate, err := strconv.ParseFloat(rateRes, 64)
	if err != nil {
		return 0, err
	}

	return rate, nil
}
