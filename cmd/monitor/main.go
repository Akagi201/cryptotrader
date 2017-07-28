package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Akagi201/cryptotrader/viabtc"
	"github.com/Akagi201/cryptotrader/yunbi"
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

func main() {
	// ctx := context.Background()
	slackApi := slack.New("")
	slackApi.SetDebug(false)

	rtm := slackApi.NewRTM()
	go rtm.ManageConnection()

	func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
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

			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT Current: %v", sntTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ETH Current: %v", ethTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BCC Current: %v", bccTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC Current: %v", btcTicker.Last), channel.ID))
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

		if sntTicker.Last < 0.35 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT Low: %v", sntTicker.Last), channel.ID))
		}

		ethTicker, err := yunbiApi.GetTicker("cny", "eth")
		if err != nil {
			log.Error("Get ETH ticker failed, err: %v", err)
		}

		log.Infof("ETH latest: %+v", ethTicker.Last)
		if ethTicker.Last > 1350 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ETH High: %v", ethTicker.Last), channel.ID))
		}

		if ethTicker.Last < 1300 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ETH Low: %v", ethTicker.Last), channel.ID))
		}

		btcTicker, err := viabtcApi.GetTicker("cny", "btc")
		if err != nil {
			log.Error("Get BTC ticker failed, err: %v", err)
		}

		log.Infof("BTC Latest: %+v", btcTicker.Last)

		if btcTicker.Last > 17000 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC High: %v", btcTicker.Last), channel.ID))
		}

		if btcTicker.Last < 16800 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC Low: %v", btcTicker.Last), channel.ID))
		}

		time.Sleep(5 * time.Second)
	}
}
