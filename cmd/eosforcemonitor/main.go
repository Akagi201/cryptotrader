package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/Akagi201/cryptotrader/eosforce"
	log "github.com/sirupsen/logrus"
)

type AccountBalance struct {
	Staked    string  `json:"staked"`
	Available float64 `json:"avalable"`
	Created   string  `json:"created"`
	Unstaking string  `json:"unstaking"`
	Name      string  `json:"name"`
}

var balances []AccountBalance
var fuckingGuys []string

func main() {

	content, _ := ioutil.ReadFile(Opts.AccountFile)
	err := json.Unmarshal(content, &balances)
	if err != nil {
		log.Fatalf("json file parse json failed, err: %v", err)
	}

	c := eosforce.New([]string{"docker", "exec", "eosforce", "cleos"}, Opts.RpcScheme, Opts.RpcHost)

	for {
		for i, v := range balances {
			available, err := c.GetAvailable(context.Background(), v.Name)
			if err != nil {
				log.Fatalf("eosforce get available balance failed, err: %v", err)
			}
			if available != v.Available {
				log.Infof("Fucking Guy %v is dumping! before: %v, after: %v, diff: %v", v.Name, v.Available, available, v.Available-available)
				balances[i].Available = available
				fuckingGuys = append(fuckingGuys, v.Name)
			}
			time.Sleep(1 * time.Second)
		}
		log.Infof("Fucking Guys: %v", fuckingGuys)
		time.Sleep(2*time.Minute)
	}
}
