package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func checkHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	resp, err := getChannelMessages(bot, p.Message.ChannelID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	PostMessagesWithStamp(bot, resp, p.Message.ChannelID, 1)
}
func checkUserHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	userlist, _, _ := bot.API().UserApi.GetUsers(context.Background()).Execute()
	cmd := strings.Split(p.Message.Text, " ")
	userID := p.Message.User.ID
	if len(cmd) >= 5 {
		for _, v := range userlist {
			if v.Name == cmd[4] {
				userID = v.Id
			}
		}
	}
	resp, err := getUserMessages(bot, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	atleast := 1
	if len(cmd) >= 6 {
		atleast, _ = strconv.Atoi(cmd[5])
	}
	PostMessagesWithStamp(bot, resp, p.Message.ChannelID, atleast)
}
func heatMapHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	resp, err := getUserMessages(bot, p.Message.User.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	daily := map[string]int{}
	for _, v := range resp {
		daily[v.CreatedAt.Format("2006-01-02")] += 1
	}
	s := ""
	for i := 0; i < 7; i++ {
		s += fmt.Sprintf("%s : %d\n", time.Now().AddDate(0, 0, -i).Format("2006-01-02"), daily[time.Now().AddDate(0, 0, -i).Format("2006-01-02")])
	}
	if len(s) > 3000 {
		s = s[:3000] + "\n(snip)"
	}
	_, _, err = bot.API().
		MessageApi.PostMessage(context.Background(), p.Message.ChannelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: s,
		}).
		Execute()
	if err != nil {
		log.Println(err)
	}
}
func PostMessagesWithStamp(bot *traqwsbot.Bot, resp []traq.Message, c string, atleast int) {
	stamplist, r, err := bot.API().StampApi.GetStamps(context.Background()).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	s := ""
	for _, v := range resp {
		c := ""
		f := 0
		q := []rune(v.Content)
		if len(q) > 50 {
			c += string(q[0:7]) + "... : "
		} else {
			c += string(q) + " : "
		}
		for _, w := range v.Stamps {
			for _, stamp := range stamplist {
				if w.StampId == stamp.Id {
					var i int32
					for i = 0; i < w.Count; i++ {
						c += ":" + stamp.Name + ":"
					}
					f += 1
				}
			}
		}
		c += "\n"
		if f >= atleast {
			s += c
		}
	}
	if len(s) > 3000 {
		s = s[:3000] + "\n(snip)"
	}
	simplePost(bot, c, s)
}

func simplePost(bot *traqwsbot.Bot, c string, s string) {
	_, r, err := bot.API().
		MessageApi.
		PostMessage(context.Background(), c).
		PostMessageRequest(traq.PostMessageRequest{
			Content: s,
		}).
		Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
