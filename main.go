package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		fmt.Println(p.Message.Text)
		if p.Message.Text[len(p.Message.Text)-5:] == "check" {
			checkHandrer(bot, p.Message.ChannelID)
		} else {
			_, _, err := bot.API().
				MessageApi.
				PostMessage(context.Background(), p.Message.ChannelID).
				PostMessageRequest(traq.PostMessageRequest{
					Content: "No such command",
				}).
				Execute()
			if err != nil {
				log.Println(err)
			}
		}

	})

	if err := bot.Start(); err != nil {
		panic(err)
	}
}
func checkHandrer(bot *traqwsbot.Bot, c string) {
	stamplist, r, err := bot.API().StampApi.GetStamps(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	channelId := c     // string | チャンネルUUID
	limit := int32(50) // int32 | 取得する件数 (optional)
	offset := int32(0) // int32 | 取得するオフセット (optional) (default to 0)
	//since := time.Now()              // time.Time | 取得する時間範囲の開始日時 (optional) (default to "0000-01-01T00:00Z")
	//until := time.Now()              // time.Time | 取得する時間範囲の終了日時 (optional)
	inclusive := true        // bool | 範囲の端を含めるかどうか (optional) (default to false)
	order := "order_example" // string | 昇順か降順か (optional) (default to "desc")
	log.Println(channelId)
	resp, r, err := bot.API().ChannelApi.GetMessages(context.Background(), channelId).Limit(limit).Offset(offset). /*.Since(since).Until(until)*/ Inclusive(inclusive).Order(order).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	s := ""
	for _, v := range resp {
		if len(v.Content) > 50 {
			s += v.Content[:50] + "... : "
		} else {
			s += v.Content + " : "
		}
		for _, w := range v.Stamps {
			for _, stamp := range stamplist {
				if w.StampId == stamp.Id {
					s += ":" + stamp.Name + ":"
				}
			}
		}
		s += "\n"
	}
	_, _, err = bot.API().
		MessageApi.
		PostMessage(context.Background(), c).
		PostMessageRequest(traq.PostMessageRequest{
			Content: s,
		}).
		Execute()
	if err != nil {
		log.Println(err)
	}
}
