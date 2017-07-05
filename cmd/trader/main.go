package main

import (
	"github.com/Akagi201/cryptotrader/chbtc"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {
	api := chbtc.New("", "")

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
	log.Infof("Get trades: %+v", tradesStr)
}
