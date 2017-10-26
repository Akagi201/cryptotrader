package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Akagi201/cryptotrader/binance"
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
	slackApi.SetDebug(true)

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

	binanceApi := binance.New("", "")

	go func() {
		for {
			zrxTicker, err := binanceApi.GetTicker("btc", "zrx")
			if err != nil {
				log.Errorf("Get ZRX ticker failed, err: %v", err)
			}

			channel := FindChannelByName(rtm, "monitor")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ZRX Current: %v", zrxTicker.Last), channel.ID))

			time.Sleep(30 * time.Minute)
		}
	}()

	for {
		zrxTicker, err := binanceApi.GetTicker("btc", "zrx")
		if err != nil {
			log.Error("Get ZRX ticker failed, err: %v", err)
		}

		log.Infof("ZRX Latest: %v", zrxTicker.Last)

		if zrxTicker.Last > 0.0001 {
			fmt.Println("High")
			channel := FindChannelByName(rtm, "monitor")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ZRX High: %v", zrxTicker.Last), channel.ID))
		}

		if zrxTicker.Last < 0.00009 {
			fmt.Println("Low")
			channel := FindChannelByName(rtm, "monitor")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("ZRX Low: %v", zrxTicker.Last), channel.ID))
		}

		time.Sleep(5 * time.Second)
	}
}
