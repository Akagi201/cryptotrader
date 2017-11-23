package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Akagi201/binancego"
	"github.com/Akagi201/cryptotrader/zb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TickerCmd struct {
	Average bool `short:"a" long:"average" description:"Get average price"`
}

func (tc *TickerCmd) Execute(args []string) error {
	log.Debugf("ticker command average: %v, args: %v", tc.Average, args)
	if len(args) != 1 {
		return errors.New("Wrong args format, use tradebot ticker [exchange].<quote-base> instead")
	}
	var exchange string
	var pair string
	var quote string
	var base string
	words := strings.Split(args[0], ".")
	if len(words) == 2 {
		exchange = words[0]
		pair = words[1]
	} else if len(words) == 1 {
		pair = args[0]
	} else {
		return errors.New("Wrong args format, use tradebot ticker [exchange].<quote-base> instead")
	}
	pairWords := strings.Split(pair, "-")
	if len(pairWords) != 2 {
		return errors.New("Wrong args format, use tradebot ticker [exchange].<quote-base> instead")
	}
	quote = pairWords[0]
	base = pairWords[1]

	prices := make(map[string]float64)
	if exchange == "" {
		for _, v := range opts.Exchanges {
			switch v {
			case "binance":
				c := binancego.NewClient("", "")
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				t, _ := c.GetTicker(ctx, quote, base)
				prices[v] = t.Last
				fmt.Printf("%v.%v: %v\n", v, pair, t.Last)
			case "zb":
				c := zb.New("", "")
				t, _ := c.GetTicker(base, quote)
				prices[v] = t.Last
				fmt.Printf("%v.%v: %v\n", v, pair, t.Last)
			}
		}
		if tc.Average {
			var sum float64 = 0
			for _, v := range prices {
				sum += v
			}
			avgPrice := sum / float64(len(prices))
			fmt.Printf("%v average price: %v\n", pair, avgPrice)
		}
	} else {
		switch exchange {
		case "binance":
			c := binancego.NewClient("", "")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			t, _ := c.GetTicker(ctx, quote, base)
			fmt.Printf("%v.%v: %v\n", exchange, pair, t.Last)
		case "zb":
			c := zb.New("", "")
			t, _ := c.GetTicker(base, quote)
			fmt.Printf("%v.%v: %v\n", exchange, pair, t.Last)
		}
	}

	return nil
}

var tickerCmd TickerCmd

func init() {
	parser.AddCommand("ticker", "get a trade pair ticker", "The ticker command get a trade pair ticker", &tickerCmd)
}
