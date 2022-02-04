package ircbot

import (
	"fmt"
	ircevent "github.com/thoj/go-ircevent"
	"log"
	"regexp"
)

var (
	reValidString = regexp.MustCompile("^([a-zA-Z0-9\\-\\_]+)$")
)

func (bot *IrcBot) IsValidString(role string) bool {
	roleResult := reValidString.FindAllStringSubmatch(role, -1)
	if len(roleResult) != 1 {
		return false
	}
	return role == roleResult[0][1]
}

func (bot *IrcBot) AddHostmaskFor(nickname string) {
	bot.conn.AddCallback("311", func(e *ircevent.Event) {
		go func(e *ircevent.Event) {
			if len(e.Arguments) != 6 {
				log.Fatalf("whois reply: not enough arguments\n")
			}

			hostmask := fmt.Sprintf("%s!%s@%s", e.Arguments[1], e.Arguments[2], e.Arguments[3])
			if !bot.state.HasHostmask(hostmask) {
				bot.state.AddHostmaskToNickname(hostmask, nickname)
			}
			bot.conn.RemoveCallback("311", 0)
		}(e)
	})

	bot.conn.Whois(nickname)
}
