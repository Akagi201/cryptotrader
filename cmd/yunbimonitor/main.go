package main

import (
	"os"
	"time"

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

	for {
		_, err := yunbiApi.GetTicker("cny", "eth")
		channel := FindChannelByName(rtm, "devops")
		if err != nil {
			//rtm.SendMessage(rtm.NewOutgoingMessage("yunbi is down", channel.ID))
			log.Errorf("Get ETH ticker failed, err: %v", err)
		} else {
			rtm.SendMessage(rtm.NewOutgoingMessage("yunbi is up", channel.ID))
			log.Infoln("yunbi is up")
		}

		time.Sleep(5 * time.Second)
	}
}
