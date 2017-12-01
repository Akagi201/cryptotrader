package main

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TradeCmd struct {
	Strategy string `short:"s" long:"strategy" default:"market_maker" description:"strategy to use"`
	List     bool   `short:"l" long:"list" description:"list strategy"`
	Interval int    `long:"interval" default:"5" description:"tick interval in second"`
}

func (tc *TradeCmd) Execute(args []string) error {
	log.Infof("trade command strategy: %v, args: %v", tc.Strategy, args)
	switch tc.Strategy {
	case "market_maker":
		tc.marketMaker()
	default:
		return errors.New("strategy not supported")
	}
	return nil
}

func (tc *TradeCmd) onTick() {
	log.Infof("on tick, interval: %v", tc.Interval)

	// get balance
	// get order book status
	// limit trade
}

func (tc *TradeCmd) marketMaker() {
	log.Infoln("use strategy market maker")

	for {
		tc.onTick()
		time.Sleep(time.Duration(tc.Interval) * time.Second)
	}
}

var tradeCmd TradeCmd

func init() {
	parser.AddCommand("trade", "run trading bot against live market data", "The trade command runs trading bot against live market data", &tradeCmd)
}
