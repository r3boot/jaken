package ircbot

func (bot *IrcBot) WhoAmI(channel, caller, nickname string) {
	if bot.state.HasNickname(nickname) {
		bot.conn.Privmsgf(channel, "You are %s (%s)", nickname, caller)
	} else {
		bot.conn.Privmsgf(channel, "I dont know you")
	}
}

func (bot *IrcBot) Meet(channel, caller, nickname string) {
	if !bot.IsOwner(caller) {
		bot.conn.Privmsgf(channel, "I dont know your hostmask %s", caller)
		return
	}

	if nickname == "" {
		bot.conn.Privmsgf(channel, "Need a nickname")
		return
	}

	if bot.state.HasNickname(nickname) {
		bot.conn.Privmsg(channel, "I already know that nickname")
		return
	}

	bot.AddHostmask(nickname)
	bot.conn.Privmsgf(channel, "Pleased to meet you %s", nickname)
}

func (bot *IrcBot) Forget(channel, caller, nickname string) {
	if !bot.IsOwner(caller) {
		bot.conn.Privmsgf(channel, "I dont know your hostmask %s", caller)
		return
	}

	if nickname == "" {
		bot.conn.Privmsgf(channel, "Need a nickname")
		return
	}

	if !bot.state.HasNickname(nickname) {
		bot.conn.Privmsg(channel, "Unknown nickname")
		return
	}

	bot.state.RemoveUser(nickname)
	bot.conn.Privmsgf(channel, "Forgot about %s", nickname)
}
