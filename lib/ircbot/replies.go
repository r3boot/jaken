package ircbot

import (
	"time"
)

func (bot *IrcBot) replyWorker() {
	for {
		select {
		case msg := <-bot.privmsgChan:
			bot.conn.Privmsg(msg.Recipient, msg.Message)
		case msg := <-bot.noticeChan:
			bot.conn.Notice(msg.Recipient, msg.Message)
		case topic := <-bot.topicChan:
			bot.conn.SendRawf("TOPIC %s %s", topic.Channel, topic.Topic)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
