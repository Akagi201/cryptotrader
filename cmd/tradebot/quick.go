package main

import (
	"context"
	"strings"
	"time"

	"github.com/Akagi201/cryptotrader/bigone"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/Akagi201/cryptotrader/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type QuickCmd struct {
	Pair     string  `short:"r" long:"pair" description:"exchange pair to trade"`
	Side     string  `short:"s" long:"side" default:"" description:"trade side: buy or sell"`
	Limit    float64 `short:"l" long:"limit" description:"the highest price or lowest price"`
	Amount   float64 `long:"amount" description:"the amount coins you want to buy"`
	Step     float64 `short:"p" long:"step" description:"step price"`
	All      bool    `short:"a" long:"all" description:"trade all my balance"`
	Interval int     `short:"i" long:"interval" default:"3" description:"tick interval in second"`
	ApiKey   string  `short:"k" long:"apikey" default:"" description:"API key"`
	quote    string
	base     string
	client   *bigone.Client
}

func (qc *QuickCmd) Execute(args []string) error {
	log.Infof("trade command side: %v, limit: %v, all:%v, interval: %v, args: %v", qc.Side, qc.Limit, qc.All, qc.Interval, args)

	var exchange string
	var pair string
	words := strings.Split(qc.Pair, ".")
	if len(words) == 2 {
		exchange = words[0]
		pair = words[1]
	} else {
		return errors.New("Wrong pair format, use [exchange].<quote-base> instead")
	}
	pairWords := strings.Split(pair, "-")
	if len(pairWords) != 2 {
		return errors.New("Wrong pair format, use [exchange].<quote-base> instead")
	}
	qc.quote = pairWords[0]
	qc.base = pairWords[1]

	switch exchange {
	case "bigone":
		qc.BigoneTrade()
	default:
		return errors.New("exchange not supported")
	}

	return nil
}

func (qc *QuickCmd) BigoneTrade() error {
	qc.client = bigone.New(qc.ApiKey)

	log.Infoln("bigone trade")

	for {
		qc.onTick()
		time.Sleep(time.Duration(qc.Interval) * time.Second)
	}
}

func (qc *QuickCmd) onTick() {
	log.Infof("on tick, interval: %v", qc.Interval)

	var depth *model.OrderBook
	{
		// get order book status
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		depth, err = qc.client.GetDepth(ctx, qc.quote, qc.base)
		if err != nil {
			log.Fatalf("bigone get depth failed: %v", err)
		}

		log.Infof("bigone sell one: %+v", depth.Asks[0])
		log.Infof("bigone buy one: %+v", depth.Bids[0])
	}

	needNewOrder := false
	{
		// get my orders
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		orders, err := qc.client.GetOrders(ctx, qc.quote, qc.base, 20)
		if err != nil {
			log.Fatalf("bigone get my orders failed: %v", err)
		}
		log.Infof("orders: %+v", orders)

		if len(orders) == 0 {
			needNewOrder = true
		}

		for _, v := range orders {
			log.Infof("order: %+v", v)
			switch qc.Side {
			case "buy":
				log.Infof("price: %v, buy one: %v", v.Price, depth.Bids[0].Price)
				if v.Price < depth.Bids[0].Price {
					// cancel order below buy one
					qc.client.CancelOrder(ctx, qc.quote, qc.base, v.ID)
					log.Infof("Cancel order price: %v, amount: %v", v.Price, v.Amount)
					time.Sleep(2 * time.Second)
				}
				needNewOrder = true
			case "sell":
				if v.Price > depth.Asks[0].Price {
					// cancel order above sell one
					qc.client.CancelOrder(ctx, qc.quote, qc.base, v.ID)
					log.Infof("Cancel order price: %v, amount: %v", v.Price, v.Amount)
					time.Sleep(2 * time.Second)
				}
				needNewOrder = true
			}
		}
	}

	noZeroBalance := []model.Balance{}
	{
		// get balance
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		balance, err := qc.client.GetAccount(ctx)
		if err != nil {
			log.Fatalf("bigone get balance failed: %v", err)
		}

		noZeroBalance = util.GetNonZeroBalance(balance)

		log.Infof("bigone balance: %+v", noZeroBalance)
	}

	// new order
	if needNewOrder {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var freeBTC float64 = 0
		for _, v := range noZeroBalance {
			if v.Currency == "BTC" {
				freeBTC = v.Free
			}
		}
		switch qc.Side {
		case "buy":
			id, err := qc.client.Trade(ctx, qc.quote, qc.base, "BID", freeBTC, depth.Bids[0].Price+qc.Step)
			if err != nil {
				log.Fatalf("client trade error: %v", err)
			}
			log.Infof("New Buy Order %v-%v: amount: %v, price:%v, order_id: %v", qc.quote, qc.base, freeBTC, depth.Bids[0].Price+qc.Step, id)
		case "sell":
			id, err := qc.client.Trade(ctx, qc.quote, qc.base, "ASK", freeBTC, depth.Asks[0].Price-qc.Step)
			if err != nil {
				log.Fatalf("client trade error: %v", err)
			}
			log.Infof("New Sell Order %v-%v: amount: %v, price:%v, order_id: %v", qc.quote, qc.base, freeBTC, depth.Asks[0].Price-qc.Step, id)
		}
	} else {
		log.Infof("No new order needed")
	}
}

var quickCmd QuickCmd

func init() {
	parser.AddCommand("quick", "quick trade against an exchange pair", "The quick command quick trade against an exchange pair", &quickCmd)
}
