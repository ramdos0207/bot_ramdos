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
	if len(cmd) >= 3 {
		for _, v := range userlist {
			if v.Name == cmd[2] {
				userID = v.Id
			}
		}
	}
	resp, err := getUserMessages(bot, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	atleast := 1
	if len(cmd) >= 4 {
		atleast, _ = strconv.Atoi(cmd[3])
	}
	PostMessagesWithStamp(bot, resp, p.Message.ChannelID, atleast)
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

// https://git.trap.jp/pikachu/traQ-BOT-pika-test/src/branch/main/commands/stamps.go
func getUserMessages(bot *traqwsbot.Bot, userID string /*, progressMessageID string*/) ([]traq.Message, error) {
	var messages []traq.Message
	var before = time.Now()
	for {
		t1 := time.Now()

		res, r, err := bot.API().MessageApi.SearchMessages(context.Background()).From(userID).Limit(int32(100)).Offset(int32(0)).Before(before).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}

		fmt.Println(time.Since(t1))
		if err != nil {
			return nil, err
		}
		if len(res.Hits) == 0 {
			break
		}
		messages = append(messages, res.Hits...)
		// for i := range res.Hits {
		// 	messages = append(messages, res.Hits[i])
		// }
		time.Sleep(time.Millisecond * 100)
		before = messages[len(messages)-1].CreatedAt
		fmt.Println(len(messages))
		/*Bot.API().
			MessageApi.EditMessage(context.Background(), progressMessageID).PostMessageRequest(traq.PostMessageRequest{
			Content: fmt.Sprintf("Searching...(%d):loading:", len(messages)),
		}).Execute()*/
	}

	return messages, nil
}
func getChannelMessages(bot *traqwsbot.Bot, channelID string /*, progressMessageID string*/) ([]traq.Message, error) {
	var messages []traq.Message
	var before = time.Now()
	for {
		t1 := time.Now()

		res, r, err := bot.API().MessageApi.SearchMessages(context.Background()).In(channelID).Limit(int32(100)).Offset(int32(0)).Before(before).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}

		fmt.Println(time.Since(t1))
		if err != nil {
			return nil, err
		}
		if len(res.Hits) == 0 {
			break
		}

		messages = append(messages, res.Hits...)
		// for i := range res.Hits {
		// 	messages = append(messages, res.Hits[i])
		// }
		time.Sleep(time.Millisecond * 100)
		before = messages[len(messages)-1].CreatedAt
		fmt.Println(len(messages))
		/*Bot.API().
			MessageApi.EditMessage(context.Background(), progressMessageID).PostMessageRequest(traq.PostMessageRequest{
			Content: fmt.Sprintf("Searching...(%d):loading:", len(messages)),
		}).Execute()*/
	}

	return messages, nil
}
