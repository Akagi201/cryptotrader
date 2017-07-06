package main

import (
	"github.com/Akagi201/cryptotrader/chbtc"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {
	api := chbtc.New("c390fceb-cee2-44bd-980a-0662aed39142", "dfffc1e0-bab1-46ca-a947-49d814154836")

	ticker, err := api.GetTicker("cny", "eth")
	if err != nil {
		log.Errorf("Get ticker failed, err: %v", err)
	}

	log.Infof("Get ticker: %+v", ticker)

	orderBook, err := api.GetOrderBook("cny", "eth", 3, 0.1)
	if err != nil {
		log.Errorf("Get orderbook failed, err: %v", err)
	}

	orderBookStr := spew.Sdump(orderBook)
	log.Infof("Get orderbook: %+v", orderBookStr)

	trades, err := api.GetTrades("cny", "eth", 0)
	if err != nil {
		log.Errorf("Get trades failed, err: %v", err)
	}

	tradesStr := spew.Sdump(trades)
	log.Infof("Get trades: %v", tradesStr)

	kline, err := api.GetKline("cny", "eth", "1min", 0, 0)
	if err != nil {
		log.Errorf("Get kline failed, err: %v", err)
	}

	klineStr := spew.Sdump(kline)
	log.Infof("Get kline: %v", klineStr)

	ethAddr, err := api.GetUserAddress("eth")
	if err != nil {
		log.Errorf("Get UserAddress failed, err: %v", err)
	}

	log.Infof("Get eth addr: %v", ethAddr)
}
