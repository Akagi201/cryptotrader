package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Akagi201/cryptotrader/binance"
	"github.com/Akagi201/cryptotrader/viabtc"
	"github.com/Akagi201/cryptotrader/yunbi"
	mapset "github.com/deckarep/golang-set"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

func FindChannelByName(rtm *slack.RTM, name string) *slack.Channel {
	for _, ch := range rtm.GetInfo().Channels {
		if ch.Name == name {
			return &ch
		}
	}
	return nil
}

func stringToInterfaceSlice(s []string) []interface{} {
	new := make([]interface{}, len(s))
	for i, v := range s {
		new[i] = v
	}
	return new
}

func main() {
	// ctx := context.Background()
	slackApi := slack.New(Opts.SlackKey)
	slackApi.SetDebug(false)

	rtm := slackApi.NewRTM()
	go rtm.ManageConnection()

	func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				_ = ev
				//log.Infof("Infos Channels: %+v", ev.Info.Channels[0].Name)
				//log.Infof("res: id %v, name %v", ev.Info.Channels[0].groupConversation.conversation.ID,
				//	ev.Info.Channels[0].groupConversation.Name)
				//channelStr := spew.Sdump(ev.Info.Channels)
				//log.Infof("channels: %v", channelStr)
				log.Infof("slack connected")
				return
			case *slack.InvalidAuthEvent:
				log.Errorf("Invalid credentials")
				return
			default:
			}
		}
	}()

	go func() {
		for {
			time.Sleep(35 * time.Minute)
			os.Exit(0)
		}
	}()

	yunbiApi := yunbi.New("", "")
	viabtcApi := viabtc.New("", "")
	binanceApi := binance.New("", "")

	go func() {
		for {
			sntTicker, err := yunbiApi.GetTicker("cny", "snt")
			if err != nil {
				log.Errorf("Get SNT ticker failed, err: %v", err)
			}
			ethTicker, err := yunbiApi.GetTicker("cny", "eth")
			if err != nil {
				log.Errorf("Get ETH ticker failed, err: %v", err)
			}

			btcTicker, err := viabtcApi.GetTicker("cny", "btc")
			if err != nil {
				log.Error("Get BTC ticker failed, err: %v", err)
			}

			bccTicker, err := viabtcApi.GetTicker("cny", "bcc")
			if err != nil {
				log.Errorf("Get BCC ticker failed, err: %v", err)
			}

			dntTicker, err := binanceApi.GetTicker("eth", "dnt")
			if err != nil {
				log.Errorf("Get DNT ticker failed, err: %v", err)
			}

			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT Current: %v", sntTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ETH Current: %v", ethTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BCC Current: %v", bccTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC Current: %v", btcTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("DNT Current: %v", dntTicker.Last), channel.ID))

			time.Sleep(30 * time.Minute)
		}
	}()

	var oldTickerList []string
	var newTickerList []string
	var err error
	oldTickerList, err = yunbiApi.GetTickerList()
	if err != nil {
		log.Fatalf("Yunbi get ticker list failed, err: %v", err)
	}
	newTickerList, err = yunbiApi.GetTickerList()
	if err != nil {
		log.Fatalf("Yunbi get ticker list failed, err: %v", err)
	}

	go func() {
		for {
			newTickerList, err = yunbiApi.GetTickerList()
			newListInterface := stringToInterfaceSlice(newTickerList)
			oldListInterface := stringToInterfaceSlice(oldTickerList)
			newSet := mapset.NewSetFromSlice(newListInterface)
			oldSet := mapset.NewSetFromSlice(oldListInterface)
			if len(newSet.Difference(oldSet).ToSlice()) > 0 {
				channel := FindChannelByName(rtm, "devops")
				rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Yunbi Got New Coin: %v", newSet.Difference(oldSet).ToSlice()), channel.ID))
			}
			//oldTickerList = newTickerList
			time.Sleep(30 * time.Minute)
		}
	}()

	for {
		sntTicker, err := yunbiApi.GetTicker("cny", "snt")
		if err != nil {
			log.Errorf("Get SNT ticker failed, err: %v", err)
		}

		log.Infof("SNT Latest: %+v", sntTicker.Last)
		if sntTicker.Last > 0.4 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT High: %v", sntTicker.Last), channel.ID))
		}

		if sntTicker.Last < 0.15 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT Low: %v", sntTicker.Last), channel.ID))
		}

		ethTicker, err := yunbiApi.GetTicker("cny", "eth")
		if err != nil {
			log.Error("Get ETH ticker failed, err: %v", err)
		}

		log.Infof("ETH latest: %+v", ethTicker.Last)
		if ethTicker.Last > 2500 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ETH High: %v", ethTicker.Last), channel.ID))
		}

		if ethTicker.Last < 1500 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ETH Low: %v", ethTicker.Last), channel.ID))
		}

		btcTicker, err := viabtcApi.GetTicker("cny", "btc")
		if err != nil {
			log.Error("Get BTC ticker failed, err: %v", err)
		}

		log.Infof("BTC Latest: %+v", btcTicker.Last)

		if btcTicker.Last > 31000 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC High: %v", btcTicker.Last), channel.ID))
		}

		if btcTicker.Last < 20000 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC Low: %v", btcTicker.Last), channel.ID))
		}

		dntTicker, err := binanceApi.GetTicker("eth", "dnt")
		if err != nil {
			log.Error("Get DNT ticker failed, err: %v", err)
		}

		log.Infof("DNT Latest: %+v", dntTicker.Last)

		if dntTicker.Last > 0.000285 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("DNT High: %v", dntTicker.Last), channel.ID))
		}

		if dntTicker.Last < 0.0002 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("DNT Low: %v", dntTicker.Last), channel.ID))
		}

		time.Sleep(5 * time.Second)
	}
}
