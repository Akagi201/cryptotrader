package xueqiu

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const (
	XueqiuAPI = "https://stock.xueqiu.com/v5/stock/chart/kline.json"
)

func GetXueqiuCookies() ([]*http.Cookie, error) {
	request, err := http.NewRequest("GET", "https://xueqiu.com", nil)
	if err != nil {
		return nil, errors.Wrap(err, "http new request")
	}

	client := http.DefaultClient
	for {
		response, err := client.Do(request)
		if err != nil {
			log.Errorf("http get cookie failed, err: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return response.Cookies(), nil
	}
}

// klineType: normal: 不复权, before: 前复权, after: 后复权
func GetXueqiuKline(symbol string, now int64, cookies []*http.Cookie, klineType string) ([]byte, error) {
	return GetContents(XueqiuAPI+"?symbol="+symbol+"&begin="+cast.ToString(now*1000)+"&period=1m&type="+klineType+"&count=-31&indicator=kline,volume,amount", cookies)
}

func GetContents(_url string, cookies []*http.Cookie) ([]byte, error) {
	request, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http new request")
	}

	// proxy := func(_ *http.Request) (*url.URL, error) {
	// 	return url.Parse(MayiURL)
	// }

	// transport := &http.Transport{
	// 	Proxy: proxy,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 		DualStack: true,
	// 	}).DialContext,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// 	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	// }

	client := &http.Client{
		Timeout: time.Second * 60,
		// Transport: transport,
	}

	if cookies != nil {
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	// loc, _ := tz.LoadLocation("Asia/Shanghai")
	// request.Header.Set("Proxy-Authorization", getSign(loc))
	ip, err := getMyIP("https://api.ipify.org?format=text", client)
	log.Debugf("my external IP: %v", ip)

	for {
		resp, err := client.Do(request)
		if err != nil {
			log.Errorf("http request failed, err: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if resp.Body != nil {
				resp.Body.Close()
			}
			time.Sleep(5 * time.Second)
			continue
		}

		return body, nil
	}
}

// https://api.ipify.org?format=text
// https://www.ipify.org
// http://myexternalip.com
// http://api.ident.me
// http://whatismyipaddress.com/api
func getMyIP(_url string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return "", errors.Wrap(err, "http new request")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	// loc, _ := tz.LoadLocation("Asia/Shanghai")
	// req.Header.Set("Proxy-Authorization", getSign(loc))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
