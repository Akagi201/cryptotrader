package main

import (
	"fmt"
	"time"

	"github.com/Akagi201/cryptotrader/viabtc"
	"github.com/Akagi201/cryptotrader/yunbi"
	"github.com/davecgh/go-spew/spew"
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
				channelStr := spew.Sdump(ev.Info.Channels)
				log.Infof("channels: %v", channelStr)
				log.Infof("slack connected")
				return
			case *slack.InvalidAuthEvent:
				log.Errorf("Invalid credentials")
				return
			default:
			}
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
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BCC Current: %v", bccTicker.Last), channel.ID))
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC Current: %v", btcTicker.Last), channel.ID))
			time.Sleep(10 * time.Minute)
		}
	}()

	for {
		sntTicker, err := yunbiApi.GetTicker("cny", "snt")
		if err != nil {
			log.Errorf("Get SNT ticker failed, err: %v", err)
		}

		log.Infof("SNT Latest: %+v", sntTicker.Last)
		if sntTicker.Last > 0.6 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT High: %v", sntTicker.Last), channel.ID))
		}

		if sntTicker.Last < 0.55 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("SNT Low: %v", sntTicker.Last), channel.ID))
		}

		btcTicker, err := viabtcApi.GetTicker("cny", "btc")
		if err != nil {
			log.Error("Get BTC ticker failed, err: %v", err)
		}

		log.Infof("BTC Latest: %+v", btcTicker.Last)

		if btcTicker.Last > 19000 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC High: %v", sntTicker.Last), channel.ID))
		}

		if btcTicker.Last < 18000 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("BTC Low: %v", sntTicker.Last), channel.ID))
		}

		time.Sleep(5 * time.Second)
	}
}
