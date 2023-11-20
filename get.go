package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

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
