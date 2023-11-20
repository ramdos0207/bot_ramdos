package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		cmd := strings.Split(p.Message.Text, " ")
		if cmd[1] == "check" {
			checkHandrer(bot, p)
		} else if cmd[1] == "checkuser" {
			checkUserHandrer(bot, p)
		} else if cmd[1] == "heatmap" {
			heatMapHandrer(bot, p)
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
