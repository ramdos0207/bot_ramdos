package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func checkHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	resp, err := getChannelMessages(bot, p.Message.ChannelID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	PostMessagesWithStamp(bot, resp, c, 1)
}
func checkUserHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
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
	resp, err := getUserMessages(bot, userID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	atleast := 1
	if len(cmd) >= 6 {
		atleast, _ = strconv.Atoi(cmd[5])
	}
	PostMessagesWithStamp(bot, resp, c, atleast)
}

func heatMapHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	resp, err := getUserMessages(bot, p.Message.User.ID, c)
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
	simpleEdit(bot, c, s)
}
func stampCountHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	userlist, _, _ := bot.API().UserApi.GetUsers(context.Background()).Execute()
	cmd := strings.Split(p.Message.Text, " ")
	userID := p.Message.User.ID
	if len(cmd) >= 4 {
		for _, v := range userlist {
			if v.Name == cmd[3] {
				userID = v.Id
			}
		}
	}
	resp, err := getUserMessages(bot, userID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	//-------------------------------------
	target := "w"
	if len(cmd) >= 5 {
		target = cmd[4]
	}
	stamplist, r, err := bot.API().StampApi.GetStamps(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", r)
	}
	for _, stamp := range stamplist {
		if target == stamp.Name || target == ":"+stamp.Name+":" {
			target = stamp.Id
		}
	}
	fmt.Println(target)
	count := 0
	total := 0
	for _, v := range resp {
		for _, w := range v.Stamps {
			if w.StampId == target {
				count += 1
				total += int(w.Count)
			}
		}
	}
	userstat, r, err := bot.API().UserApi.GetUserStats(context.Background(), userID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", r)
	}
	sendcount, sendtotal := 0, 0
	for _, v := range userstat.Stamps {
		if v.Id == target {
			sendcount = int(v.Count)
			sendtotal = int(v.Total)
		}
	}
	s := "received: " + strconv.Itoa(count) + "(total: " + strconv.Itoa(total) + ")\n"
	s += "sent: " + strconv.Itoa(sendcount) + "(total: " + strconv.Itoa(sendtotal) + ")\n"
	simpleEdit(bot, c, s)
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
			c += string(q[0:50]) + "... : "
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
	simpleEdit(bot, c, s)
}
func lengthHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	userId := p.Message.User.ID
	userlist, _, _ := bot.API().UserApi.GetUsers(context.Background()).Execute()
	cmd := strings.Split(p.Message.Text, " ")
	if len(cmd) >= 3 {
		for _, v := range userlist {
			if v.Name == cmd[2] {
				userId = v.Id
			}
		}
	}
	resp, err := getUserMessages(bot, userId, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	threshold := 200
	if len(cmd) >= 4 {
		threshold, _ = strconv.Atoi(cmd[3])
	}
	sum := 0
	over := 0
	runesum := 0
	runeover := 0
	for _, v := range resp {
		sum += len(v.Content)
		if len(v.Content) > threshold {
			over += 1
		}
		runesum += len([]rune(v.Content))
		if len([]rune(v.Content)) > threshold {
			runeover += 1
		}
	}
	average := float64(sum) / float64(len(resp))
	runeaverage := float64(runesum) / float64(len(resp))
	simpleEdit(bot, c, fmt.Sprintf("sum(byte): %d\naverage(byte): %f\nover %d(byte): %d\nsum(rune):%d\naverage(rune):%f\nover %d(rune):%d", sum, average, threshold, over, runesum, runeaverage, threshold,runeover))
}

type lenstr struct {
	sum int
	str string
}

func lengthgroupHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {

	groupId := ""
	grouplist, _, _ := bot.API().GroupApi.GetUserGroups(context.Background()).Execute()
	cmd := strings.Split(p.Message.Text, " ")
	if len(cmd) >= 3 {
		for _, v := range grouplist {
			if v.Name == cmd[2] {
				groupId = v.Id
			}
		}
	}
	userlist, _, _ := bot.API().GroupApi.GetUserGroupMembers(context.Background(), groupId).Execute()

	username, _, _ := bot.API().UserApi.GetUsers(context.Background()).Execute()
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	responsetext := "|user|sum|count|average|>200|sum(rune)|average(rune)|>200(rune)|\n|---|---|---|---|---|---|---|---|\n"
	x := simplePost(bot, p.Message.ChannelID, responsetext)
	sq := []lenstr{}
	for _, v := range userlist {
		resp, err := getUserMessages(bot, v.Id, c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		threshold := 200
		sum := 0
		over := 0
		runesum := 0
		runeover := 0
		for _, v := range resp {
			sum += len(v.Content)
			runesum += len([]rune(v.Content))
			if len(v.Content) > threshold {
				over += 1
			}
			if len([]rune(v.Content)) > threshold {
				runeover += 1
			}
		}
		average := float64(sum) / float64(len(resp))
		runeaverage := float64(runesum) / float64(len(resp))
		name := "?"
		for _, q := range username {
			if q.Id == v.Id {
				name = q.Name
			}
		}
		sortkey:="sum"
		if len(cmd) >= 4 {
			sortkey = cmd[3]
		}
		minimum := 0
		if len(cmd) >= 5 {
			minimum, _ = strconv.Atoi(cmd[4])
		}
		if name != "?" && len(resp) >= minimum{
			c := fmt.Sprintf("|:@%s: %s|%d|%d|%f|%d|%d|%f|%d|\n", name, name, sum, len(resp), average, over,runesum,runeaverage,runeover)
			responsetext += c
			if sortkey=="sum"{
				sq = append(sq, lenstr{sum,c})
			}else if sortkey=="count"{
				sq = append(sq, lenstr{len(resp), c})
			}else if sortkey=="average"{
				sq = append(sq, lenstr{int(100.0*average),c})
			}else if sortkey==">200"{
				sq = append(sq, lenstr{over,c})
			}
		}
		if len(responsetext) > 9900 {
			simpleEdit(bot, x, responsetext+"\n(snip)")
		} else {
			simpleEdit(bot, x, responsetext)
		}

	}
	sort.Slice(sq, func(i, j int) bool { return sq[i].sum > sq[j].sum })
	responsetext = "|user|sum|count|average|>200|sum(rune)|average(rune)|>200(rune)|\n|---|---|---|---|---|---|---|---|\n"
	for _, v := range sq {
		responsetext += v.str
	}
	if len(responsetext) > 9900 {
		simpleEdit(bot, x, responsetext+"\n(snip)")
	} else {
		simpleEdit(bot, x, responsetext)
	}

}
