package main

import (
	"fmt"
	"time"

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

	api := yunbi.New("", "")

	go func() {
		for {
			ticker, err := api.GetTicker("cny", "snt")
			if err != nil {
				log.Errorf("Get ticker failed, err: %v", err)
			}
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Current Price: %v", ticker.Last), channel.ID))
			time.Sleep(10 * time.Minute)
		}
	}()

	for {
		ticker, err := api.GetTicker("cny", "snt")
		if err != nil {
			log.Errorf("Get ticker failed, err: %v", err)
		}

		log.Infof("Get SNT Latest Price: %+v", ticker.Last)
		if ticker.Last > 0.6 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Price High: %v", ticker.Last), channel.ID))
		}

		if ticker.Last < 0.55 {
			channel := FindChannelByName(rtm, "devops")
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Price Low: %v", ticker.Last), channel.ID))
		}
		time.Sleep(2 * time.Second)
	}
}
