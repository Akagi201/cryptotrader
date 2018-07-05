package main

import (
	"math"

	"github.com/Akagi201/binancego/model"
	"github.com/Akagi201/utilgo/slices"
	talib "github.com/markcheno/go-talib"
	log "github.com/sirupsen/logrus"
)

func priceTrendFactor(trades []model.Trade, orderBook *model.OrderBook, isSymmetric bool) float64 {
	var tradePrices []float64
	for _, v := range trades {
		tradePrices = append(tradePrices, v.Price)
	}
	latestTrades := tradePrices[len(tradePrices)-6:]
	midPrice := (orderBook.Bids[0].Price+orderBook.Asks[0].Price)/2*0.7 + (orderBook.Bids[1].Price+orderBook.Asks[1].Price)/2*0.2 + (orderBook.Bids[2].Price+orderBook.Asks[2].Price)/2*0.1
	latestTrades = append(latestTrades, midPrice)
	isBullTrend := false
	isBearTrend := false
	lastPriceTooFarFromLatest := false
	hasLargeVolTrade := false

	if latestTrades[len(latestTrades)-1] > slices.MaxFloat(latestTrades[:len(latestTrades)-1])+latestTrades[len(latestTrades)-1]*0.00005 || (latestTrades[len(latestTrades)-1] > slices.MaxFloat(latestTrades[:len(latestTrades)-2])+latestTrades[len(latestTrades)-1]*0.00005 && latestTrades[len(latestTrades)-1] > latestTrades[len(latestTrades)-2]) {
		isBullTrend = true
		log.Info("Bull Tending!!")
	} else if latestTrades[len(latestTrades)-1] < slices.MinFloat(latestTrades[:len(latestTrades)-1])-latestTrades[len(latestTrades)-1]*0.00005 || (latestTrades[len(latestTrades)-1] < slices.MinFloat(latestTrades[:len(latestTrades)-2])-latestTrades[len(latestTrades)-1]*0.00005 && latestTrades[len(latestTrades)-1] < latestTrades[len(latestTrades)-2]) {
		isBearTrend = true
		log.Info("Bear Tending!!")
	}

	if math.Abs(latestTrades[len(latestTrades)-1]-latestTrades[len(latestTrades)-2]*0.7-latestTrades[len(latestTrades)-3]*0.2-latestTrades[len(latestTrades)-4]*0.1) > latestTrades[len(latestTrades)-1]*0.01 {
		lastPriceTooFarFromLatest = true
		log.Info("Last price too far from latest!!")
	}

	var orderQty []float64
	maxLen := slices.MinInt([]int{15, len(orderBook.Bids), len(orderBook.Asks)})
	for _, v := range orderBook.Bids[:maxLen] {
		orderQty = append(orderQty, v.Price*v.Amount)
	}

	for _, v := range orderBook.Asks[:maxLen] {
		orderQty = append(orderQty, v.Price*v.Amount)
	}

	if slices.MaxFloat(orderQty) > 2 {
		log.Info("Has large vol trade")
		hasLargeVolTrade = true
	}

	if isBullTrend || isBearTrend || lastPriceTooFarFromLatest || hasLargeVolTrade {
		return 0
	}

	index := talib.Rsi(tradePrices, len(tradePrices)-1)
	if index[len(tradePrices)-1] <= 20 || index[len(tradePrices)-1] >= 80 {
		log.Infof("RSI index not good, index: %v", index[len(tradePrices)-1])
		return 0
	}

	var factor float64
	if isSymmetric {
		factor = 1 - math.Abs(index[len(tradePrices)-1]-50)/50.0
	} else {
		factor = index[len(tradePrices)-1] / 50.0
	}

	return factor
}
