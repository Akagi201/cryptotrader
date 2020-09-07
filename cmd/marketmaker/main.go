package main

import (
	"context"
	"time"

	"github.com/Akagi201/cryptotrader/binance"
	log "github.com/sirupsen/logrus"
)

var (
	bc  *binance.Client
	cnt int
)

func onTick() {
	log.Infof("on tick, interval: %v", opts.Interval)
	// get balance
	// get order book status
	// limit trade

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nulsTrades, err := bc.GetTrades(ctx, "nuls", "btc", 0, 0, 0, 0)
	if err != nil {
		log.Fatalf("Get NULS trades failed, err: %v", err)
	}

	nulsOrderBook, err := bc.GetDepth(ctx, "nuls", "btc", 0)
	if err != nil {
		log.Fatalf("Get NULS orderbook failed, err: %v", err)
	}

	factor := priceTrendFactor(nulsTrades, nulsOrderBook, true)
	log.Infof("factor: %v", factor)
}

func main() {
	bc = binance.New(opts.AccessKey, opts.SecretKey)

	for {
		onTick()
		time.Sleep(time.Duration(opts.Interval) * time.Second)
	}
}
