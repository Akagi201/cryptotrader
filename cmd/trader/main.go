package main

import (
	"context"
	"time"

	"github.com/Akagi201/cryptotrader/allcoin"
	"github.com/Akagi201/cryptotrader/bigone"
	"github.com/Akagi201/cryptotrader/binance"
	"github.com/Akagi201/cryptotrader/bitfinex"
	"github.com/Akagi201/cryptotrader/bitflyer"
	"github.com/Akagi201/cryptotrader/bittrex"
	"github.com/Akagi201/cryptotrader/cex"
	"github.com/Akagi201/cryptotrader/coincheck"
	"github.com/Akagi201/cryptotrader/coinegg"
	"github.com/Akagi201/cryptotrader/eosforce"
	"github.com/Akagi201/cryptotrader/etherscan"
	"github.com/Akagi201/cryptotrader/fixer"
	"github.com/Akagi201/cryptotrader/gateio"
	"github.com/Akagi201/cryptotrader/huobi"
	"github.com/Akagi201/cryptotrader/lbank"
	"github.com/Akagi201/cryptotrader/liqui"
	"github.com/Akagi201/cryptotrader/okcoin"
	"github.com/Akagi201/cryptotrader/okex"
	"github.com/Akagi201/cryptotrader/poloniex"
	"github.com/Akagi201/cryptotrader/util"
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

	if false {
		// binance
		rc := binance.New("xBhrsdymp92w3yTIf20x2TOs39fyyCM4TgeJCtbKuWQe1Rx2nCh2y6rDl1G5u5Th", "bC08CbIl5wBfVYJgrGkYgSl8dJ6JXVoqT57uLTstYyOj9ZUwo8r3dejHdonToiVw")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if false {
			ethTicker, err := rc.GetTicker(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("Binance get ETH-BTC ticker failed, err: %v", err)
			}

			log.Infof("ETH-BTC Ticker: %+v", ethTicker)
		}

		if false {
			tickers, err := rc.GetTickers(ctx)
			if err != nil {
				log.Fatalf("Get tickers failed, err: %v", err)
			}

			log.Infof("Get tickers: %+v", tickers)
		}

		if false {
			err := rc.Ping(ctx)
			if err != nil {
				log.Fatalf("Binance ping failed, err: %v", err)
			}

			log.Infoln("Binance Ping success")
		}

		if false {
			serverTime, err := rc.GetTime(ctx)
			if err != nil {
				log.Fatalf("Binance get time failed, err: %v", err)
			}

			log.Infof("Binance Server time: %v", serverTime)
		}

		if false {
			nulsDepth, err := rc.GetDepth(ctx, "nuls", "btc", 0)
			if err != nil {
				log.Fatalf("Get NULS depth failed, err: %v", err)
			}

			log.Infof("Get NULS Depth: %+v", nulsDepth)
			for i, v := range nulsDepth.Asks[:5] {
				log.Infof("ask orderbook, %v: %v: %v", i, v.Price, v.Amount)
			}
			for i, v := range nulsDepth.Bids[:5] {
				log.Infof("bids orderbook, %v: %v: %v", i, v.Price, v.Amount)
			}
		}

		if false {
			nulsTrades, err := rc.GetTrades(ctx, "nuls", "btc", 0, 0, 0, 0)
			if err != nil {
				log.Fatalf("Get NULS trades failed, err: %v", err)
			}

			log.Infof("Get NULS trades: %+v, len: %v", nulsTrades[:5], len(nulsTrades))
			for i, v := range nulsTrades[len(nulsTrades)-5:] {
				log.Infof("index: %v, price: %v, amount: %v", i, v.Price, v.Amount)
			}
		}

		if false {
			zrxRecords, err := rc.GetRecords(ctx, "zrx", "btc", "1m", 0, 0, 0)
			if err != nil {
				log.Fatalf("Get ZRX records failed, err: %v", err)
			}

			log.Infof("Get ZRX records: %+v", zrxRecords)
		}

		if false {
			bookTickers, err := rc.GetBookTickers(ctx)
			if err != nil {
				log.Fatalf("Get book tickers failed, err: %v", err)
			}

			log.Infof("Get book tickers: %+v", bookTickers)
		}

		{
			account, err := rc.GetAccount(ctx, 0)
			if err != nil {
				log.Fatalf("Get account failed, err: %v", err)
			}
			account = util.GetNonZeroBalance(account)

			log.Infof("Get account: %+v", account)
		}

		if false {
			tradeID, err := rc.Trade(ctx, "zrx", "btc", "BUY", "LIMIT", "GTC", 1, 0.02, 0, 0, 0)
			if err != nil {
				log.Fatalf("Trade ZRX_BTC failed, err: %v", err)
			}

			log.Infof("Trade ZRX_BTC ID: %v", tradeID)
		}

		if false {
			order, err := rc.GetOrder(ctx, "zrx", "btc", 0, 0)
			if err != nil {
				log.Fatalf("Get ZRX_BTC order failed, err: %v", err)
			}

			log.Infof("Get ZRX_BTC order: %+v", order)
		}

		if false {
			err := rc.CancelOrder(ctx, "zrx", "btc", 0, 0)
			if err != nil {
				log.Fatalf("Cancel ZRX_BTC order failed, err: %v", err)
			}

			log.Infof("Cancel ZRX_BTC order success")
		}

		if false {
			orders, err := rc.GetOrders(ctx, "zrx", "btc", 0)
			if err != nil {
				log.Fatalf("Get ZRX_BTC orders failed, err: %v", err)
			}

			log.Infof("Get ZRX_BTC orders: %+v", orders)
		}

		if false {
			orders, err := rc.GetAllOrders(ctx, "zrx", "btc", 0, 0, 0)
			if err != nil {
				log.Fatalf("Get all ZRX_BTC orders failed, err: %v", err)
			}

			log.Infof("Get all ZRX_BTC orders: %+v", orders)
		}

		if false {
			zrxTrades, err := rc.GetMyTrades(ctx, "dnt", "btc", 0, 0, 0)
			if err != nil {
				log.Fatalf("Get my ZRX trades failed, err: %v", err)
			}

			log.Infof("Get my ZRX trades: %+v", zrxTrades)
		}
	}

	if false {
		// OKEX
		c := okex.New("", "")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if false {
			ethTicker, err := c.GetTicker(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("OKEX get eth-btc ticker failed, err: %v", err)
			}

			log.Infof("OKEX ETH-BTC Ticker: %+v", ethTicker)
		}

		if false {
			ethDepth, err := c.GetDepth(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("OKEX get eth-btc depth failed, err: %v", err)
			}

			log.Infof("OKEX ETH-BTC depth: %+v", ethDepth)
		}

		if false {
			ethTrades, err := c.GetTrades(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("OKEX get eth-btc trades failed, err: %v", err)
			}

			log.Infof("OKEX ETH-BTC trades: %+v", ethTrades)
		}

		if false {
			ethKline, err := c.GetRecords(ctx, "eth", "btc", "1min", 0, 0)
			if err != nil {
				log.Fatalf("OKEX get eth-btc kline failed, err: %v", err)
			}

			log.Infof("OKEX ETH-BTC kline: %+v", ethKline)
		}

		{
			balance, err := c.GetAccount(ctx)
			if err != nil {
				log.Fatalf("OKEX get balance failed, err: %v", err)
			}
			log.Infof("OKEX balance: %+v", balance)
		}
	}

	if false {
		// gate.io
		c := gateio.New("", "")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if false {
			pairs, err := c.GetPairs(ctx)
			if err != nil {
				log.Fatalf("gate.io get pairs failed, err: %v", err)
			}

			log.Infof("gate.io pairs: %+v", pairs)
		}

		if false {
			marketInfo, err := c.GetMarketInfo(ctx)
			if err != nil {
				log.Fatalf("gate.io get market_info failed, err: %v", err)
			}

			log.Infof("gate.io market_info: %+v", marketInfo)
		}

		{
			ethTicker, err := c.GetTicker(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("gate.io get ticker failed, err: %v", err)
			}

			log.Infof("gate.io get ticker eth-btc: %+v", ethTicker)
		}
	}

	if false {
		// big.one
		c := bigone.New("")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if false {
			ethTicker, err := c.GetTicker(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("big.one get ticker failed, err: %v", err)
			}

			log.Infof("big.one get ticker eth-btc: %+v", ethTicker)
		}

		if false {
			ethBook, err := c.GetDepth(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("big.one get eth-btc depth failed, err: %v", err)
			}

			log.Infof("big.one get depth eth-btc: %+v", ethBook)
		}

		if false {
			ethTrade, err := c.GetTrades(ctx, "eth", "btc")
			if err != nil {
				log.Fatalf("big.one get eth-btc trade failed, err: %v", err)
			}

			log.Infof("big.one get trade eth-btc: %+v", ethTrade)
		}

		if false {
			id, err := c.Trade(ctx, "eth", "btc", "BID", 0.1, 0.04)
			if err != nil {
				log.Fatalf("big.one trade eth-btc failed, err: %v", err)
			}
			log.Infof("big.one trade eth-btc success, order id: %v", id)

			order, err := c.GetOrder(ctx, "eth", "btc", id)
			if err != nil {
				log.Fatalf("big.one get order eth-btc failed, err: %v", err)
			}
			log.Infof("big.one get order eth-btc success, order id: %v", order)

			orders, err := c.GetOrders(ctx, "eth", "btc", 10)
			if err != nil {
				log.Fatalf("big.one get orders eth-btc failed, err: %v", err)
			}
			log.Infof("big.one get orders eth-btc success, orders: %+v", orders)

			err = c.CancelOrder(ctx, "eth", "btc", id)
			if err != nil {
				log.Fatalf("big.one cancel order eth-btc failed, err: %v", err)
			}
			log.Info("big.one cancel order eth-btc success")
		}

		if false {
			balance, err := c.GetAccount(ctx)
			if err != nil {
				log.Fatalf("big.one get balance failed, err: %v", err)
			}
			log.Infof("big.one get balance: %+v", balance)
		}
	}

	{
		c := eosforce.New([]string{"docker", "exec", "eosforce", "cleos"}, "https", "w2.eosforce.cn")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		{
			available, err := c.GetAvailable(ctx, "blockgw")
			if err != nil {
				log.Fatalf("eosforce get available balance failed, err: %v", err)
			}

			log.Infof("eosforce available balance: %+v", available)
		}

		{
			staked, err := c.GetStaked(ctx, "kuso", "blockgw")
			if err != nil {
				log.Fatalf("eosforce get staked balance failed, err: %v", err)
			}

			log.Infof("eosforce staked balance: %+v", staked)
		}

		{
			staked, err := c.GetUnstaking(ctx, "kuso", "blockgw")
			if err != nil {
				log.Fatalf("eosforce get unstaking balance failed, err: %v", err)
			}

			log.Infof("eosforce unstaking balance: %+v", staked)
		}

		{
			reward, err := c.GetRewards(ctx, "kuso", "blockgw")
			if err != nil {
				log.Fatalf("eosforce get reward balance failed, err: %v", err)
			}

			log.Infof("eosforce reward balance: %+v", reward)
		}
	}
}
