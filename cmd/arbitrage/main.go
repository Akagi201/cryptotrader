package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Akagi201/cryptotrader/huobi"
	"github.com/Akagi201/cryptotrader/okcoin"
	"github.com/Akagi201/cryptotrader/zb"
	"github.com/goSTL/sort"
	"github.com/nlopes/slack"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type ExchangePrice struct {
	Exchange string
	Price    float64
}

// GetPriceList get ticker between different exchanges.
func GetPriceList(base string, quote string, exchanges []string) []ExchangePrice {
	priceList := []ExchangePrice{}
	for _, ex := range exchanges {
		switch ex {
		case "huobi":
			api := huobi.New("", "")
			ticker, _ := api.GetTicker(base, quote)
			priceList = append(priceList, ExchangePrice{
				Exchange: ex,
				Price:    ticker.Last,
			})
		case "okcoin":
			api := okcoin.New("", "")
			ticker, _ := api.GetTicker(base, quote)
			priceList = append(priceList, ExchangePrice{
				Exchange: ex,
				Price:    ticker.Last,
			})
		case "zb":
			api := zb.New("", "")
			ticker, _ := api.GetTicker(base, quote)
			priceList = append(priceList, ExchangePrice{
				Exchange: ex,
				Price:    ticker.Last,
			})
		}
	}

	return priceList
}

func FindChannelByName(rtm *slack.RTM, name string) *slack.Channel {
	for _, ch := range rtm.GetInfo().Channels {
		if ch.Name == name {
			return &ch
		}
	}
	return nil
}

func cmp(i, j interface{}) bool {
	ii := i.(ExchangePrice)
	jj := j.(ExchangePrice)
	return ii.Price < jj.Price
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

	//go func() {
	//	for {
	//		time.Sleep(35 * time.Minute)
	//		os.Exit(0)
	//	}
	//}()

	for {
		l := GetPriceList(Opts.Base, Opts.Quote, Opts.Exchanges)
		log.Debugf("Price list: %v", l)

		out := sort.MergeSort(l, cmp)

		exchangeList := []string{"Exchange"}
		priceList := []string{Opts.Quote + "_" + Opts.Base}

		for _, v := range out {
			exchangeList = append(exchangeList, v.(ExchangePrice).Exchange)
			priceList = append(priceList, cast.ToString(v.(ExchangePrice).Price))
		}

		buf := bytes.NewBuffer([]byte{})
		table := tablewriter.NewWriter(buf)
		table.SetHeader(exchangeList)
		table.SetRowLine(true)
		table.Append(priceList)
		footer := make([]string, len(exchangeList))
		footer[0] = "diff"
		footer[1] = fmt.Sprintf("%.2f", out[len(out)-1].(ExchangePrice).Price-out[0].(ExchangePrice).Price)
		table.Append(footer)
		table.Render()

		channel := FindChannelByName(rtm, "arbitrage")
		rtm.SendMessage(rtm.NewOutgoingMessage("```"+buf.String()+"```", channel.ID))
		log.Infof("\n%v", buf.String())

		time.Sleep(5 * time.Second)
	}
}
