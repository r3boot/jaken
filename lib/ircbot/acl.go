package ircbot

func (bot *IrcBot) IsOwner(hostmask string) bool {
	return bot.params.Owner == hostmask
}
