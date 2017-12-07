package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Akagi201/cryptotrader/binance"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type BalanceCmd struct {
	Exchange  string `short:"e" long:"exchange" description:"The exchange to query balance"`
	AccessKey string `long:"access" default:"" description:"access key for the exchange"`
	SecretKey string `long:"secret" default:"" description:"secret key for the exchange"`
}

func (bc *BalanceCmd) Execute(args []string) error {
	log.Infof("balance command exchange: %v, args: %v", bc.Exchange, args)
	switch bc.Exchange {
	case "binance":
		log.Info("binance")
		c := binance.New(bc.AccessKey, bc.SecretKey)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		b, _ := c.GetAccount(ctx, 0)
		fmt.Printf("%v: %+v\n", bc.Exchange, b)
	default:
		return errors.New("exchange not supported")
	}
	return nil
}

var balanceCmd BalanceCmd

func init() {
	parser.AddCommand("balance", "get the balance from an exchange", "The balance command get the balance from an exchange", &balanceCmd)
}
