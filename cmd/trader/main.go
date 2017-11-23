package main

import (
	"github.com/Akagi201/cryptotrader/allcoin"
	"github.com/Akagi201/cryptotrader/binance"
	"github.com/Akagi201/cryptotrader/bitfinex"
	"github.com/Akagi201/cryptotrader/bitflyer"
	"github.com/Akagi201/cryptotrader/bittrex"
	"github.com/Akagi201/cryptotrader/cex"
	"github.com/Akagi201/cryptotrader/coincheck"
	"github.com/Akagi201/cryptotrader/coinegg"
	"github.com/Akagi201/cryptotrader/etherscan"
	"github.com/Akagi201/cryptotrader/fixer"
	"github.com/Akagi201/cryptotrader/huobi"
	"github.com/Akagi201/cryptotrader/lbank"
	"github.com/Akagi201/cryptotrader/liqui"
	"github.com/Akagi201/cryptotrader/okcoin"
	"github.com/Akagi201/cryptotrader/poloniex"
	"github.com/Akagi201/cryptotrader/zb"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {
	if false {
		// zb
		api := zb.New("", "")

		ticker, err := api.GetTicker("btc", "eth")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)

		orderBook, err := api.GetOrderBook("btc", "eth", 3, 0.1)
		if err != nil {
			log.Errorf("Get orderbook failed, err: %v", err)
		}

		orderBookStr := spew.Sdump(orderBook)
		log.Infof("Get orderbook: %+v", orderBookStr)

		trades, err := api.GetTrades("btc", "eth", 0)
		if err != nil {
			log.Errorf("Get trades failed, err: %v", err)
		}

		tradesStr := spew.Sdump(trades)
		log.Infof("Get trades: %v", tradesStr)

		kline, err := api.GetRecords("btc", "eth", "1min", 0, 0)
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

	if false {
		// huobi
		api := huobi.New("", "")
		ticker, err := api.GetTicker("btc", "eth")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// binance
		api := binance.New("", "")
		ticker, err := api.GetTicker("eth", "dnt")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// cex
		api := cex.New("", "")
		omgTicker, err := api.GetTicker("btc", "omg")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get OMG ticker: %+v", omgTicker)

		payTicker, err := api.GetTicker("btc", "pay")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get PAY ticker: %+v", payTicker)
	}

	if false {
		// okcoin
		api := okcoin.New("", "")

		ticker, err := api.GetTicker("btc", "eth")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)

	}

	if false {
		// bittrex
		api := bittrex.New("", "")

		ticker, err := api.GetTicker("eth", "storj")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// liqui
		api := liqui.New("", "")

		ticker, err := api.GetTicker("eth", "zrx")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// etherscan
		api := etherscan.New("")

		balance, err := api.GetBalance("0x258ce53268BEaA9BA97fA6b7790d7555ae4044fc")
		if err != nil {
			log.Errorf("Get balance failed, err: %v", err)
		}

		log.Infof("Get balance: %v", balance)

		block, err := api.GetBlockNumber()

		if err != nil {
			log.Errorf("Get block number failed, err: %v", err)
		}

		log.Infof("Get block number: %v", block)
	}

	if false {
		// fixer
		api := fixer.New()

		rate, err := api.GetRate("USD", "btc")
		if err != nil {
			log.Errorf("Get rate failed, err: %v", err)
		}

		log.Infof("Get rate: %v", rate)
	}

	if false {
		// poloniex
		api := poloniex.New("", "")

		ticker, err := api.GetTicker("btc", "eth")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// lbank
		api := lbank.New("", "")

		ticker, err := api.GetTicker("btc", "btc")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// coinegg
		api := coinegg.New("", "")

		ticker, err := api.GetTicker("btc", "btc")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// allcoin
		api := allcoin.New("", "")

		ticker, err := api.GetTicker("usd", "btc")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// bitfinex
		api := bitfinex.New("", "")

		ticker, err := api.GetTicker("usd", "btc")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// coincheck
		api := coincheck.New("", "")

		ticker, err := api.GetTicker("jpy", "btc")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}

	if false {
		// bitflyer
		api := bitflyer.New("", "")

		ticker, err := api.GetTicker("jpy", "btc")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get ticker: %+v", ticker)
	}
}
