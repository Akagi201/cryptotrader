package main

import (
	"github.com/Akagi201/cryptotrader/chbtc"
	log "github.com/sirupsen/logrus"
)

func main() {
	api := chbtc.New("", "")

	res, err := api.GetTicker("cny", "eth")
	if err != nil {
		log.Errorf("Get ticker failed, err: %v", err)
	}
	log.Infof("Get ticker: %+v", res)
}
